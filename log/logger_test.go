package log

import (
	"testing"
)

func BenchmarkProxy_D(b *testing.B) {
	logger := NewLogger(nil)
	for i := 0; i < b.N; i++ {
		logger.D("i = %d\n", i)
	}
	logger.CloseSafely()
}

func BenchmarkProxy_I(b *testing.B) {
	setDefaultWriter(NewFileWriter("./log", 2))
	logger := NewLogger(nil)
	for i := 0; i < b.N; i++ {
		logger.I("i = %d\n", i)
	}
	logger.CloseSafely()
}
