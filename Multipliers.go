package main

import (
	"gonum.org/v1/gonum/mat"
)

type Multipliers struct {
	inner       *mat.Dense
	outer       *mat.Dense
	outerRadius float64
	annulus     *mat.Dense
	M           *mat.CDense
	N           *mat.CDense
}

func ConstructMultipliers(
	innerRadius float64,
	width int,
	height int,
	logres float64,
) *Multipliers {
	outerRadius := 3 * innerRadius
	inner := AntialiasedCircle(width, height, innerRadius, true, logres)
	outer := AntialiasedCircle(width, height, outerRadius, true, logres)
	annulus := mat.NewDense(height, width, nil)
	annulus.Sub(outer, inner)

	// Scale each kernel so the sum is 1
	inner_magnitude := SumDenseMatrix(inner)
	annulus_magnitude := SumDenseMatrix(annulus)

	inner = DivideDenseMatrix(inner, inner_magnitude)
	annulus = DivideDenseMatrix(annulus, annulus_magnitude)

	// Precompute the FFT's
	M := Fft2RealIn(inner)
	N := Fft2RealIn(annulus)

	return &Multipliers{
		inner:       inner,
		outer:       outer,
		outerRadius: outerRadius,
		annulus:     annulus,
		M:           M,
		N:           N,
	}
}
