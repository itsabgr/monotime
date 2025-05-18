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

func NewEstimated(accuracy time.Duration, timestamp time.Time) *Estimated {
	if accuracy <= 10*time.Nanosecond {
		panic(errors.New("accuracy should be more than 10 nanoseconds"))
	}

	est := &Estimated{
		closeChannel: make(chan struct{}),
		closeOnce:    sync.Once{},
		now:          nanotime(),
	}

	est.timestampOffset = -time.Duration(timestamp.UnixNano() - est.now)

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		monoTimer := time.NewTicker(accuracy)
		defer monoTimer.Stop()
		wg.Done()
		for {
			atomic.StoreInt64(&est.now, nanotime())
			select {
			case <-est.closeChannel:
				return
			case <-monoTimer.C:
			}
		}
	}()

	wg.Wait()

	return est
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
	return (time.Duration(atomic.LoadInt64(&est.now)) * time.Nanosecond) - est.timestampOffset
}
