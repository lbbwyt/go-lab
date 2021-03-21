package go_concurrent

import "sync"

type Singleton struct {
}

// sync 实现单例类
var (
	instance *Singleton
	once     sync.Once
)

func GetSingletonInstacne() *Singleton {
	once.Do(func() {
		instance = &Singleton{}
	})
	return instance
}
