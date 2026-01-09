package main

import (
	"github.com/r3labs/sse/v2"
	"net/http"
	"time"
)

const StreamTest = "test"

func main() {
	server := sse.New()
	server.CreateStream(StreamTest)

	go producer(server)

	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			println("The client is disconnected here")
			return
		}()

		server.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", mux)
}

func producer(server *sse.Server) {

	for {
		server.Publish("test", &sse.Event{
			Data: []byte("ping"),
		})
		time.Sleep(time.Second)
	}

}
