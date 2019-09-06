package rand

import (
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestStdConcurrency tests the Rand object is not thread safe.
func TestStdConcurrent(t *testing.T) {
	var (
		wg          sync.WaitGroup
		failedCount int32
		end         int64
		begin       = time.Now().UnixNano()
		threads     = 16
	)
	r := rand.New(rand.NewSource(1))
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				e := recover()
				if e == nil {
					t.Error("this should be panic")
					return
				}
				err := e.(error)
				if !strings.Contains(err.Error(), "index out of range") {
					t.Errorf("unexpected error. %v\n", err)
					return
				}
				atomic.AddInt32(&failedCount, 1)
			}()
			defer wg.Done()
			for {
				_ = r.Int63()
			}
		}()
	}
	for {
		time.Sleep(time.Microsecond)
		if atomic.LoadInt32(&failedCount) == int32(threads)-1 && end == 0 {
			end = time.Now().UnixNano()
		}
		// We get random numbers in several threads. If the Rand is non-thread
		// safe, the most threads will be panic but one. So the test will exit
		// successfully, otherwise will be timeout or never exit.
		if time.Now().UnixNano()-end > (end-begin)*10 {
			wg.Done()
			break
		}
	}
	wg.Wait()
}

// TestLockedConcurrent exercises the rand API concurrently, triggering situations
// where the race detector is likely to detect issues.
func TestLockedConcurrent(t *testing.T) {
	const (
		numRoutines = 10
		numCycles   = 10
	)

	r := NewLocked(time.Now().UnixNano())
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func(i int) {
			defer wg.Done()
			buf := make([]byte, 997)
			for j := 0; j < numCycles; j++ {
				var seed int64
				seed += int64(r.ExpFloat64())
				seed += int64(r.Float32())
				seed += int64(r.Float64())
				seed += int64(r.Intn(Int()))
				seed += int64(r.Int31n(Int31()))
				seed += int64(r.Int63n(Int63()))
				seed += int64(r.NormFloat64())
				seed += int64(r.Uint32())
				seed += int64(r.Uint64())
				for _, p := range Perm(10) {
					seed += int64(p)
				}
				_, err := r.Read(buf)
				if err != nil {
					t.Errorf("Read(): %v", err)
				}
				for _, b := range buf {
					seed += int64(b)
				}
				r.Seed(int64(i*j) * seed)
			}
		}(i)
	}
}
