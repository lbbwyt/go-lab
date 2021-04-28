package utils

import (
	"testing"
	"time"
)

func TestTrylock(t *testing.T) {
	m := NewMutex()
	ok := m.TryLock(time.Second)
	if !ok {
		t.Error("it should be lock suc but failed!")
	}

	ok = m.TryLock(time.Second * 3)
	if ok {
		t.Error("it should be lock failed but suc")
	}

	m.Unlock()

	ok = m.TryLock(time.Second)
	if !ok {
		t.Error("it should be lock suc but failed!")
	}

}
