package main

func main() {
	a := make(map[string]chan int, 0)
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2

	println(len(ch))

	a["lbb"] = ch

	for _, v := range a {
		<-v
		println(len(v))
	}

	println(len(ch))

}
