package utils

import (
	"sync/atomic"
	"time"
)

type Timer struct {
	innerTimer *time.Timer
	stopped    int32
}

func NewTimer(d time.Duration, callback func()) *Timer {
	return &Timer{
		innerTimer: time.AfterFunc(d, callback),
	}
}

// Stop stops the timer.
func (t *Timer) Stop() {
	if t == nil {
		return
	}
	if !atomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		return
	}

	t.innerTimer.Stop()
}

func (t *Timer) Reset(d time.Duration) bool {
	if t == nil {
		return false
	}
	// already stopped
	if atomic.LoadInt32(&t.stopped) == 1 {
		return false
	}
	return t.innerTimer.Reset(d)
}
