package main

import (
	"testing"
)

func benchmarkRandomStringadd(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringadd(i)
	}
}

func benchmarkRandomStringfmts(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringfmts(i)
	}
}

func benchmarkRandomStringfmts2(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringfmts2(i)
	}
}

func benchmarkRandomStringfmtp(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringfmtp(i)
	}
}

func benchmarkRandomStringbuilder(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringbuilder(i)
	}
}

func benchmarkRandomStringbytes(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringbytes(i)
	}
}

func benchmarkRandomStringbyte(i int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		randomStringbyte(i)
	}
}

func BenchmarkRandomStringadd1000(b *testing.B) { benchmarkRandomStringadd(1000, b) }

func BenchmarkRandomStringfmts1000(b *testing.B) { benchmarkRandomStringfmts(1000, b) }

func BenchmarkRandomStringfmts21000(b *testing.B) { benchmarkRandomStringfmts2(1000, b) }

func BenchmarkRandomStringfmtp1000(b *testing.B) { benchmarkRandomStringfmtp(1000, b) }

func BenchmarkRandomStringbuilder1000(b *testing.B) { benchmarkRandomStringbuilder(1000, b) }

func BenchmarkRandomStringbytes1000(b *testing.B) { benchmarkRandomStringbytes(1000, b) }

func BenchmarkRandomStringbyte1000(b *testing.B) { benchmarkRandomStringbyte(1000, b) }
