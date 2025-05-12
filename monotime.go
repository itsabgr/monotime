package monotime

import (
	"time"
	_ "unsafe"
)

//go:linkname nanotime runtime.nanotime
func nanotime() int64

var initial = Now()

// Now returns monotonic time duration
//
//go:nosplit
func Now() time.Duration {
	return time.Duration(nanotime()) * time.Nanosecond
}

// Initial returns the initial monotonic time
//
//go:nosplit
func Initial() time.Duration {
	return initial
}
