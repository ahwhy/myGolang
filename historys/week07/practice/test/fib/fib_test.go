package main

import (
	"testing"
	"time"
)

func BenchmarkFib(b *testing.B) {
	time.Sleep(3 * time.Second)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		fib(30)
	}
}
