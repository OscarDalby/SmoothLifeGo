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
	outerRadius float64,
	width int,
	height int,
	logres float64,
) *Multipliers {
	inner := AntialiasedCircle(width, height, innerRadius, true, logres)
	outer := AntialiasedCircle(width, height, outerRadius, true, logres)
	// saveMatrixAsImage(inner, "inner.png")
	// saveMatrixAsImage(outer, "outer.png")
	annulus := mat.NewDense(height, width, nil)
	annulus.Sub(outer, inner)
	// saveMatrixAsImage(annulus, "annulus.png")

	// Scale each kernel so the sum is 1
	inner_magnitude := SumDenseMatrix(inner)

	annulus_magnitude := SumDenseMatrix(annulus)

	inner = DivideDenseMatrix(inner, inner_magnitude)
	// saveMatrixAsImage(inner, "inner_scaled.png")
	annulus = DivideDenseMatrix(annulus, annulus_magnitude)
	// saveMatrixAsImage(annulus, "annulus_scaled.png")

	// Precompute the FFT's
	M := fft2dense(inner)
	N := fft2dense(annulus)

	return &Multipliers{
		inner:       inner,
		outer:       outer,
		outerRadius: outerRadius,
		annulus:     annulus,
		M:           M,
		N:           N,
	}
}
