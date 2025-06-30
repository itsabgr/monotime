package monotime

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Elapsed struct {
	_            sync.Mutex
	closeChannel chan struct{}
	closeOnce    sync.Once
	timestamp    int64
	now          int64
	start        int64
}

func NewElapsed(accuracy time.Duration, timestamp time.Time) *Elapsed {
	if accuracy < 10*time.Nanosecond {
		panic(errors.New("accuracy should be more than 10 nanoseconds"))
	}

	start := nanotime()

	ela := &Elapsed{
		closeChannel: make(chan struct{}),
		closeOnce:    sync.Once{},
		timestamp:    timestamp.UnixNano(),
		start:        start,
		now:          start,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		monoTimer := time.NewTicker(accuracy)
		defer monoTimer.Stop()
		wg.Done()
		var now, last int64 = 0, 0
		for {
			now = nanotime()
			if now < last {
				panic(ErrNegative)
			}
			atomic.StoreInt64(&ela.now, now)
			last = now
			select {
			case <-ela.closeChannel:
				return
			case <-monoTimer.C:
			}
		}
	}()

	wg.Wait()

	return ela
}

func (ela *Elapsed) Stop() (done bool) {
	ela.closeOnce.Do(func() {
		close(ela.closeChannel)
		done = true
	})
	return done
}

func (ela *Elapsed) Stopped() bool {
	select {
	case <-ela.closeChannel:
		return true
	default:
		return false
	}
}

//go:nosplit
func (ela *Elapsed) Now() time.Duration {
	return (time.Duration(atomic.LoadInt64(&ela.now)) * time.Nanosecond) - (time.Duration(atomic.LoadInt64(&ela.start)) * time.Nanosecond) + time.Duration(ela.timestamp)
}
