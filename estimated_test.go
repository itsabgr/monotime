package monotime

import (
	"testing"
	"time"
)

func TestMonotime_Now(t *testing.T) {
	mt := NewEstimated(time.Microsecond*1, time.Now())
	defer mt.Stop()
	for i := 0; i < 10; i++ {
		if time.Since(time.UnixMicro(mt.Now().Microseconds())) > time.Second {
			t.Fail()
		}
		time.Sleep(time.Second / 4)
	}
}

func TestMonotime_Stop(t *testing.T) {
	mt := NewEstimated(time.Microsecond*1, time.Now())
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
	mt := NewEstimated(time.Microsecond*1, time.Now())
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
