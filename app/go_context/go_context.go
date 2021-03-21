package go_context

//假设你开启了一个函数，你需要将一些常用的值传递给下游函数，但是不能通过函数参数传递，怎么办？
//假设你开启了一个协程 A，协程 A 衍生出很多子协程，这些子协程又衍生出子协程，如果协程 A 所完成的任务“成果”不再需要，那我们怎么通知衍生出的子协程及时退出并释放占用的系统资源呢？
//假设一个任务需要在 2s 内完成，如果超时，如何优雅地退出返回呢？
//假设一个任务需要在中午 12 点完成，如果到点没有完成，又该如何优雅地退出呢？

//WithCancel() 基于父 context，返回子 context 和取消函数；
//
//WithDeadline() 基于父 context，返回带截止时间的子 context 和取消函数；
//
//WithTimeout() 基于父 context，返回带超时时间的子 context 和取消函数；
//
//WithValue() 基于父 context，返回绑定键值对的子 context，没有取消函数；

//前三个函数都会返回取消函数，需要注意的是只有创建该 context 的协程才能调用取消函数，且不推荐将取消函数作为参数传递。
//
//我们可以调用取消函数取消一个 context，以及这个 context 下面所有的子 context。

//context 树
//rootCtx := context.Background()
//childCtx := context.WithValue(rootCtx, "request_Id", "seekload")
//childOfChildCtx, cancelFunc := context.WithCancel(childCtx)

//总结
//context 是线程安全的，可在多个 goroutine 中传递；
//使用 context 作为函数参数时，需作为第一个参数，并且命名为 ctx；
//不要把 context 放在结构体中，要以参数的方式传递；
//当不知道传递什么类型 context 时，可以使用 context.TODO()；
//context 只能被取消一次，应当避免从已取消的 context 衍生 context；
//只有父 context 和创建了该 context 的函数才能调用取消函数，避免传递取消函数 cancelFunc；
