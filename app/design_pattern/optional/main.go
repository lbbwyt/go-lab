package main

import "fmt"

//函数式选项模式

type Message struct {
	id   int
	name string
}

type Option func(msg *Message)

var DefaultMessage = &Message{
	id:   0,
	name: "默认",
}

func WithId(id int) Option {
	return func(msg *Message) {
		msg.id = id
	}
}

func WithName(name string) Option {
	return func(msg *Message) {
		msg.name = name
	}
}

func NewMessage(opts ...Option) *Message {
	msg := DefaultMessage
	for _, v := range opts {
		v(msg)
	}
	return msg
}

func main() {
	v := NewMessage(WithName("234"))
	fmt.Println(fmt.Sprintf("v : %v", v))
}
