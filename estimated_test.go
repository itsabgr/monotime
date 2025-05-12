package monotime

import (
	"testing"
	"time"
)

func TestMonotime_Update(t *testing.T) {
	mt := NewEstimated(time.Second * 1)
	defer mt.Stop()
	now := mt.Update()
	if mt.Now() != now {
		t.Fail()
	}
}

func TestMonotime_Stop(t *testing.T) {
	mt := NewEstimated(time.Microsecond * 1)
	if mt.Stopped() {
		t.Fail()
	}
	if !mt.Stop() {
		t.Fail()
	}
	if !mt.Stopped() {
		t.Fail()
	}
	now := mt.Now()
	time.Sleep(time.Millisecond * 100)
	if now != mt.Now() {
		t.Fail()
	}
}

func BenchmarkEstimated(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	mt := NewEstimated(time.Microsecond * 1)
	defer mt.Stop()
	for i := 0; i < b.N; i++ {
		if mt.Now() <= 0 {
			b.Fail()
		}
	}
}

func BenchmarkNow(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if Now() <= 0 {
			b.Fail()
		}
	}
}
