package kmutex

import (
	"context"
	"sync"
)

// New return new instance of KMutex
func New() *KMutex {
	lock := new(sync.Mutex)
	return &KMutex{
		mem:  make(map[interface{}]struct{}, 2<<5),
		lock: lock,
		cond: sync.NewCond(lock),
	}
}

// KMutex represent mutex with lock by key
type KMutex struct {
	mem  map[interface{}]struct{}
	lock sync.Locker
	cond *sync.Cond
}

func (km *KMutex) hasKey(key interface{}) bool {
	_, ok := km.mem[key]
	return ok
}

// Lock represent blocking lock by key function, can be canceled by context
func (km *KMutex) Lock(ctx context.Context, key interface{}) error {
	km.lock.Lock()
	defer km.lock.Unlock()

	for km.hasKey(key) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			km.cond.Wait()
		}
	}

	km.mem[key] = struct{}{}
	return nil
}

// Unlock represent unlock function
func (km *KMutex) Unlock(key interface{}) {
	km.lock.Lock()
	delete(km.mem, key)
	km.cond.Broadcast()
	km.lock.Unlock()
}

// WithLock represent wrapper for synchronized call given function by key
func (km *KMutex) WithLock(ctx context.Context, key interface{}, app func() error) error {
	if err := km.Lock(ctx, key); err != nil {
		return err
	}
	defer km.Unlock(key)
	return app()
}
