package monotime

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Estimated struct {
	_               sync.Mutex
	closeChannel    chan struct{}
	closeOnce       sync.Once
	now             int64
	timestampOffset time.Duration
}

func NewEstimated(accuracy time.Duration) *Estimated {
	if accuracy <= time.Nanosecond {
		panic(errors.New("accuracy should be more than a nanosecond"))
	}

	est := &Estimated{
		closeChannel: make(chan struct{}),
		closeOnce:    sync.Once{},
		now:          nanotime(),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		timer := time.NewTicker(accuracy)
		defer timer.Stop()
		wg.Done()
		for {
			_ = est.Update()
			select {
			case <-est.closeChannel:
				return
			case <-timer.C:
			}
		}
	}()

	wg.Wait()

	return est
}

// Update force updates the time holder and returns current monotonic time
//
//go:nosplit
func (est *Estimated) Update() time.Duration {
	now := nanotime()
	atomic.StoreInt64(&est.now, now)
	return time.Duration(now)
}

// Stop stops the time updater goroutine
func (est *Estimated) Stop() (done bool) {
	est.closeOnce.Do(func() {
		close(est.closeChannel)
		done = true
	})
	return done
}

// Stopped returns true when time updater is stopped
func (est *Estimated) Stopped() bool {
	select {
	case <-est.closeChannel:
		return true
	default:
		return false
	}
}

// Now returns estimated monotonic time duration in given accuracy
//
//go:nosplit
func (est *Estimated) Now() time.Duration {
	return time.Duration(atomic.LoadInt64(&est.now)) * time.Nanosecond
}
