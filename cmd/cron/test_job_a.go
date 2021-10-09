package main

import (
	"fmt"
)

type TestBJob struct {
}

func NewTestBJob() *TestBJob {
	return &TestBJob{}
}

//实现接口的Run方法
func (a *TestBJob) Run() {

	fmt.Println("TestB")
	//runtime_utils.ExecWithRecover(TestPanic, recover)
	TestPanic()
}

func TestPanic() {
	panic("panic form TestB")
}

func recover(r interface{}) {
	fmt.Println(r)
}
