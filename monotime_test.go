package monotime

import (
	"testing"
	"time"
	_ "unsafe"
)

func TestInit(t *testing.T) {
	initialTime := Initial()
	time.Sleep(time.Millisecond * 100)
	if initialTime != Initial() {
		t.Fail()
	}
	if Now() <= Initial() {
		t.Fail()
	}
}

func BenchmarkInit(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if Initial() <= 0 {
			b.Fail()
		}
	}
}
