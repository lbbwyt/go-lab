package main

import (
	"fmt"
	"time"
)

type Mutex struct {
	ch chan struct{}
}

func NewMutex() *Mutex {
	mu := &Mutex{make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *Mutex) Lock() {
	<-m.ch
}

func (m *Mutex) Unlock() {
	select {
	case m.ch <- struct{}{}:
	default:
		panic("unlock of unlocked mutex")
	}
}

func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *Mutex) TryLockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-time.After(timeout):
	}
	return false
}

func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}

func main() {
	m := NewMutex()
	ok := m.TryLock()
	fmt.Printf("locked v %v\n", ok)
	ok = m.TryLock()
	fmt.Printf("locked %v\n", ok)
	ok = m.TryLockTimeout(time.Second * 3)
	fmt.Printf("locked %v\n", ok)
}

//
//package main
//
//import (
//"fmt"
//)
//
//type Mutex struct {
//	ch chan struct{}
//}
//
//func NewMutex() *Mutex {
//	mu := &Mutex{make(chan struct{}, 1)}
//	return mu
//}
//
//func (m *Mutex) Lock() {
//	m.ch <- struct{}{}
//}
//
//func (m *Mutex) Unlock() {
//	select {
//	case <-m.ch:
//	default:
//		panic("unlock of unlocked mutex")
//	}
//}
//
//func (m *Mutex) TryLock() bool {
//	select {
//	case m.ch <- struct{}{}:
//		return true
//	default:
//	}
//	return false
//}
//
//func (m *Mutex) IsLocked() bool {
//	return len(m.ch) == 1
//}
//
//func main() {
//	m := NewMutex()
//	ok := m.TryLock()
//	fmt.Printf("locked v %v\n", ok)
//	ok = m.TryLock()
//	fmt.Printf("locked %v\n", ok)
//}
