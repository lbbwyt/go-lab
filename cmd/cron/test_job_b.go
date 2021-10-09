package main

import "fmt"

type TestAJob struct {
}

func NewTestAJob() *TestAJob {
	return &TestAJob{}
}

//实现接口的Run方法
func (a *TestAJob) Run() {
	fmt.Println("TestA")
}
