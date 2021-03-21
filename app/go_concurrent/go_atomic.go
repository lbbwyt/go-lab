package go_concurrent

//atomic.Value 原子对象 可实现对复杂对象的原子读写

//atomic.Value类型对外暴露的方法就两个：
//v.Store(c)- 写操作，将原始的变量c存放到一个atomic.Value类型的v里。
//c = v.Load()- 读操作，从线程安全的v中读取上一步存放的内容。
