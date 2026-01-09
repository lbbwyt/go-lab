package scheduler

import (
	"context"
	"go-lab/app/go_dag/model"
	"go-lab/app/go_dag/store"
	"sync"
	"time"
)

type Scheduler struct {
	handlers map[string]TaskHandler
	store    store.Store
}

func NewScheduler(st store.Store) *Scheduler {
	return &Scheduler{
		handlers: make(map[string]TaskHandler, 0),
		store:    st,
	}
}

// 无状态的任务处理器
type TaskHandler func(ctx context.Context, params map[string]interface{}) error

func (s *Scheduler) Register(name string, h TaskHandler) {
	if s.handlers == nil {
		s.handlers = make(map[string]TaskHandler)
	}
	s.handlers[name] = h
}

func (s *Scheduler) Run(ctx context.Context, wf *model.Workflow) error {
	d := NewDAG()
	taskMap := make(map[string]*model.Task)
	for _, t := range wf.Tasks {
		taskMap[t.ID] = &t
		d.AddVertex(t.ID)
		for _, dep := range t.Deps {
			d.AddEdge(dep, t.ID)
		}
	}

	order, err := d.Topological()
	if err != nil {
		return err
	}

	// 剩余入度 & 完成标记
	remain := make(map[string]int)
	for _, id := range order {
		remain[id] = d.indeg[id]
	}
	done := make(chan string, len(order)) // 带缓冲，防止阻塞

	var wg sync.WaitGroup
	cond := sync.NewCond(&sync.Mutex{})

	// 1. 先把入度为 0 的任务放出去执行， 并放入done队列中
	for _, id := range order {
		if remain[id] == 0 {
			wg.Add(1)
			go func(tid string) {
				defer wg.Done()
				s.runTask(ctx, taskMap[tid])
				done <- tid
			}(id)
		}
	}

	// 2. 监听done队列中的完成事件，当前任务执行完毕， 所有临接的节点入度减1， 如果临接的节点入度为0后， 唤醒其执行任务
	go func() {
		for fin := range done {
			cond.L.Lock()
			for _, nxt := range d.adj[fin] {
				remain[nxt]--
				if remain[nxt] == 0 {
					wg.Add(1)
					go func(tid string) {
						defer wg.Done()
						s.runTask(ctx, taskMap[tid])
						done <- tid
					}(nxt)
				}
			}
			cond.L.Unlock()
		}
	}()

	wg.Wait()
	close(done)
	return nil
}

func (s *Scheduler) runTask(ctx context.Context, t *model.Task) {
	for attempt := 0; attempt <= t.Retry; attempt++ {
		c, cancel := context.WithTimeout(ctx, time.Duration(t.Timeout)*time.Second)
		err := s.handlers[t.Action](c, t.Params)
		cancel()
		if err == nil {
			break
		}
		if attempt == t.Retry {
			// 补偿逻辑
		}
	}
}
