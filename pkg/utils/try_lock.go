package utils

import "time"

type Mutex struct {
	c chan struct{}
}

func NewMutex() *Mutex {
	return &Mutex{make(chan struct{}, 1)}
}

func (m *Mutex) Lock() {
	m.c <- struct{}{}
}

func (m *Mutex) Unlock() {
	<-m.c
}

//延迟等待锁
func (m *Mutex) TryLock(timeout time.Duration) bool {

	select {
	case m.c <- struct{}{}:
		return true
	default:
		timer := time.NewTimer(timeout)
		select {
		case m.c <- struct{}{}:
			timer.Stop()
			return true
		case <-timer.C:
			return false
		}

	}
	return false
}
