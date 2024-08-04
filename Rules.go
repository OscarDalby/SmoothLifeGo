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

	a, b := aliveness.Dims()

	var alivenessSum float64 = 0
	r, c := aliveness.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			alivenessSum += aliveness.At(i, j)
		}
	}

	percentageOfPotentialAliveness := alivenessSum / float64(a*b)
	fmt.Printf("alivenessSum: %v\n", alivenessSum)
	fmt.Printf("percentageOfPotentialAliveness: %v%%\n", percentageOfPotentialAliveness)

	// debug n and m values - buffers obtained from ifft2
	var nSum float64 = 0
	rn, cn := n.Dims()
	for i := 0; i < rn; i++ {
		for j := 0; j < cn; j++ {
			nSum += n.At(i, j)
		}
	}
	fmt.Printf("nSum: %v\n", nSum)
	var mSum float64 = 0
	rm, cm := m.Dims()
	for i := 0; i < rm; i++ {
		for j := 0; j < cm; j++ {
			mSum += m.At(i, j)
		}
	}
	fmt.Printf("mSum: %v\n", mSum)

	// A fully dead cell will become alive if the neighbor density is between B1 and B2.
	// A fully alive cell will stay alive if the neighhbor density is between D1 and D2.
	// Interpolate between the two sets of thresholds depending on how alive/dead the cell is.
	// {B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
	// threshold1 := LerpDense(br.B1, br.D1, aliveness)
	// threshold2 := LerpDense(br.B2, br.D2, aliveness)

	// t1 := threshold1.At(0, 0) // t1 is coming out as the same as D1 on every run?
	// t2 := threshold2.At(0, 0) // t2 is coming out as approximately the same as B2 on every run?

	// fmt.Printf("t1: %v\n", t1)
	// fmt.Printf("t2: %v\n", t2)

	// Now with the smoothness of `logisticInterval` determine if the neighbor density is
	// inside of the threshold to stay/become alive.
	// newAliveness := LogisticIntervalTripleDense(n, threshold1, threshold2, br.N)

	newAliveness := LogisticIntervalDenseElementWise(n, br.D1, br.B2, br.N)

	var newAlivenessSum float64 = 0
	nr, nc := newAliveness.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			newAlivenessSum += newAliveness.At(i, j)
		}
	}

	fmt.Printf("newAlivenessSum: %v\n", newAlivenessSum)
	// var alivenessDiff float64 = newAlivenessSum - alivenessSum

	// boostedAliveness := AddConstantDense(newAliveness, 0.6)
	// return ClampDense(boostedAliveness, 0, 1)
	var output *mat.Dense = ClampDense(newAliveness, 0, 1)
	return output
}
