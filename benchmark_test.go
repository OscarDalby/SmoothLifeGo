package main

// run benchmarks with: go test -bench=. -benchmem

import "testing"

func BenchmarkAntialiasedCircle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AntialiasedCircle(3, 3, 4, true, 0)
	}
}
