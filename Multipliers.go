package main

import (
	"gonum.org/v1/gonum/mat"
)

type Multipliers struct {
	cm      CellMath
	inner   *mat.Dense
	outer   *mat.Dense
	annulus *mat.Dense
	M       *mat.CDense
	N       *mat.CDense
}

func ConstructMultipliers(
	cm CellMath,
	innerRadius float64,
	width int,
	height int,
	logres float64,
) *Multipliers {
	outerRadius := 3 * innerRadius
	inner := cm.AntialiasedCircle(width, height, innerRadius, true, logres)
	outer := cm.AntialiasedCircle(width, height, outerRadius, true, logres)
	annulus := mat.NewDense(height, width, nil)
	annulus.Sub(outer, inner)

	// Scale each kernel so the sum is 1
	inner_magnitude := cm.SumDenseMatrix(inner)
	annulus_magnitude := cm.SumDenseMatrix(annulus)

	inner = cm.DivideDenseMatrix(inner, inner_magnitude)
	annulus = cm.DivideDenseMatrix(annulus, annulus_magnitude)

	// Precompute the FFT's
	M := cm.Fft2RealIn(inner)
	N := cm.Fft2RealIn(annulus)

	return &Multipliers{
		cm:      CellMath{},
		inner:   inner,
		outer:   outer,
		annulus: annulus,
		M:       M,
		N:       N,
	}
}
