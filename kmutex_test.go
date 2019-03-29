package kmutex

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randSleep(min, max time.Duration) {
	time.Sleep(time.Duration(rand.Int63n(int64(max-min)) + int64(min)))
}

func TestKMutex_WithLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	km := New()

	t.Run("sequence", func(t *testing.T) {
		keys := []string{"green", "green", "green", "green", "green"}
		waitExecution := new(sync.WaitGroup)
		waitExecution.Add(len(keys))
		var tmp int32 = 0
		wg := new(sync.WaitGroup)
		for _, key := range keys {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				waitExecution.Done()
				waitExecution.Wait()
				assert.Error(t,
					km.WithLock(ctx, key, func() error {
						atomic.AddInt32(&tmp, 1)
						assert.Equal(t, int32(1), atomic.LoadInt32(&tmp))
						atomic.AddInt32(&tmp, -1)
						return errors.New(key)
					}), key)
			}(key)
		}
		wg.Wait()
	})

	t.Run("parallel", func(t *testing.T) {
		keys := []string{"green", "red", "blue", "yellow", "white"}
		waitParallelLock := new(sync.WaitGroup)
		waitParallelLock.Add(len(keys))

		wg := new(sync.WaitGroup)
		for _, key := range keys {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				randSleep(20*time.Millisecond, 100*time.Millisecond)
				assert.Error(t,
					km.WithLock(ctx, key, func() error {
						waitParallelLock.Done()
						waitParallelLock.Wait()
						return errors.New(key)
					}), key)
			}(key)
		}
		wg.Wait()
	})
}
