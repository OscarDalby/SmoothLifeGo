package main

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type BasicRules struct {
	B1 float64
	B2 float64
	D1 float64
	D2 float64
	N  float64
	M  float64
}

func (BasicRules BasicRules) Clear() {
}

// State transition function
func (br BasicRules) S(n *mat.Dense, m *mat.Dense) *mat.Dense {
	// Convert the local cell average `m` to a metric of how alive the local cell is.
	// We transition around 0.5 (0 is fully dead and 1 is fully alive).
	// The transition width is set by `br.M`
	aliveness := LogisticThresholdDenseElementWise(m, 0.5, br.M)

	var alivenessSum float64 = 0
	r, c := aliveness.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			alivenessSum += aliveness.At(i, j)
		}
	}

	// A fully dead cell will become alive if the neighbor density is between B1 and B2.
	// A fully alive cell will stay alive if the neighhbor density is between D1 and D2.
	// Interpolate between the two sets of thresholds depending on how alive/dead the cell is.
	threshold1 := LerpDense(br.B1, br.D1, aliveness)
	threshold2 := LerpDense(br.B2, br.D2, aliveness)

	t1 := threshold1.At(0, 0)
	t2 := threshold2.At(0, 0)

	fmt.Printf("threshold1: %v\n", t1)
	fmt.Printf("threshold2: %v\n", t2)

	// Now with the smoothness of `logisticInterval` determine if the neighbor density is
	// inside of the threshold to stay/become alive.
	newAliveness := LogisticIntervalTripleDense(n, threshold1, threshold2, br.N)

	var newAlivenessSum float64 = 0
	nr, nc := newAliveness.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			newAlivenessSum += newAliveness.At(i, j)
		}
	}

	var alivenessDiff float64 = newAlivenessSum - alivenessSum
	fmt.Printf("alivenessSum: %v\n", int(alivenessSum))
	// fmt.Printf("newAliveness: %v\n", newAliveness)
	// fmt.Printf("newAlivenessSum: %v\n", int(newAlivenessSum))
	fmt.Printf("alivenessDiff: %v\n", int(alivenessDiff))

	// boostedAliveness := AddConstantDense(newAliveness, 0.6)
	// return ClampDense(boostedAliveness, 0, 1)
	var output *mat.Dense = ClampDense(newAliveness, 0, 1)
	// fmt.Printf("output: %v\n", output)
	return output
}
