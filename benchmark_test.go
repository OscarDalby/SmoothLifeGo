package main

// run benchmarks with: go test -bench=. -benchmem

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func BenchmarkAntialiasedCircle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AntialiasedCircle(3, 3, 4, true, 0)
	}
}

func BenchmarkLogisticThresholdDenseElementWise(b *testing.B) {
	rows, cols := 512, 512
	data := make([]float64, rows*cols)
	for i := range data {
		data[i] = float64(i)
	}
	x := mat.NewDense(rows, cols, data)
	x0 := 0.5
	alpha := 1.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LogisticThresholdDenseElementWise(x, x0, alpha)
	}
}
