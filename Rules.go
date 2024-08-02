package main

import "gonum.org/v1/gonum/mat"

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
	// A fully dead cell will become alive if the neighbor density is between B1 and B2.
	// A fully alive cell will stay alive if the neighhbor density is between D1 and D2.
	// Interpolate between the two sets of thresholds depending on how alive/dead the cell is.
	threshold1 := LerpDense(br.B1, br.D1, aliveness)
	threshold2 := LerpDense(br.B2, br.D2, aliveness)
	// Now with the smoothness of `logisticInterval` determine if the neighbor density is
	// inside of the threshold to stay/become alive.
	newAliveness := LogisticIntervalTripleDense(n, threshold1, threshold2, br.N)
	// boostedAliveness := AddConstantDense(newAliveness, 0.6)
	// return ClampDense(boostedAliveness, 0, 1)
	return ClampDense(newAliveness, 0, 1)
}
