package main

func main() {
	// 调度任务集中管理
	sc := NewScheduler()
	sc.AddJob("*/5 * * * * *", NewTestAJob())
	sc.AddJob("*/15 * * * * *", NewTestBJob())
	sc.StartJobs()

	select {}
}
