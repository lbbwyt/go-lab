package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"go-lab/app/redis_terminal/internal"
	"net/http"
)

var addr = flag.String("addr", ":63342", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "/Users/mac/libaobao/github/private/go-lab/app/redis_terminal/frontend/home.html")
}

func main() {
	flag.Parse()
	hub := internal.NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		internal.ServeWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
