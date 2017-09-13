package log

import "testing"

func BenchmarkProxy_D(b *testing.B) {
	logger := NewLogger(nil)
	for i := 0; i < b.N; i++ {
		logger.D("i = %d\n", i)
	}
}
