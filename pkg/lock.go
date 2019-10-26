package gql

import (
	"errors"
	"time"
)

type Mutex struct {
	lock bool
}
func WithLock(mutex *Mutex) (func(), error) {
	if mutex.lock {
		ch := make(chan bool)
		go func() {
			for i := 0; i < 10; i++ {
				time.Sleep(time.Millisecond * 100)
				if !mutex.lock {
					ch <- true
					return
				}
			}
			ch <- false
		}()
		data := <-ch
		if !data {
			return nil, errors.New("failed to acquire lock")
		}
	}
	mutex.lock = true
	return func() {
		if p := recover(); p != nil {
			mutex.lock = false
			panic(p)
		}
		mutex.lock = false
	}, nil
}
