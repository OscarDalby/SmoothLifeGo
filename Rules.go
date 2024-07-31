package main

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
func (br BasicRules) S(cm CellMath, n float64, m float64) float64 {
	// Convert the local cell average `m` to a metric of how alive the local cell is.
	// We transition around 0.5 (0 is fully dead and 1 is fully alive).
	// The transition width is set by `br.M`
	var aliveness = cm.LogisticThreshold(m, 0.5, br.M)
	// A fully dead cell will become alive if the neighbor density is between B1 and B2.
	// A fully alive cell will stay alive if the neighhbor density is between D1 and D2.
	// Interpolate between the two sets of thresholds depending on how alive/dead the cell is.
	var threshold1 float64 = cm.Lerp(br.B1, br.D1, aliveness)
	var threshold2 float64 = cm.Lerp(br.B2, br.D2, aliveness)
	// Now with the smoothness of `logisticInterval` determine if the neighbor density is
	// inside of the threshold to stay/become alive.
	var newAliveness = cm.LogisticInterval(n, threshold1, threshold2, br.N)
	return cm.Clamp(newAliveness, 0, 1)
}
