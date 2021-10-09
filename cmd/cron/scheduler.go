package main

import (
	cron "github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		c: cron.New(cron.WithSeconds(), cron.WithChain(
			cron.Recover(cron.DefaultLogger))),
	}
}

// 添加一个任务
func (s *Scheduler) AddJob(cronStr string, job cron.Job) (cron.EntryID, error) {
	return s.c.AddJob(cronStr, job)
}

func (s *Scheduler) RemoveJob(id cron.EntryID) {
	s.c.Remove(id)
}

//// 启动所有任务
func (s *Scheduler) StartJobs() {
	s.c.Start()
}

func (s *Scheduler) ListJobs() []cron.Entry {
	return s.c.Entries()
}
