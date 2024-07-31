// class Multipliers:
//     """Kernel convulution for neighbor integral"""

//     INNER_RADIUS = 7.0
//     OUTER_RADIUS = INNER_RADIUS * 3.0

//     def __init__(self, size, inner_radius=INNER_RADIUS, outer_radius=OUTER_RADIUS):
//         inner = antialiased_circle(size, inner_radius)
//         outer = antialiased_circle(size, outer_radius)
//         annulus = outer - inner

//         # Scale each kernel so the sum is 1
//         inner /= np.sum(inner)
//         annulus /= np.sum(annulus)

//         # Precompute the FFT's
//         self.M = np.fft.fft2(inner)
//         self.N = np.fft.fft2(annulus)

package main

import "gonum.org/v1/gonum/mat"

// radius := 7.0
// logres := 0.0
// matrix := cm.AntialiasedCircle(screenWidth, screenHeight, radius, true, logres)

func ConstructMultipliers(
	cm CellMath,
	inner_radius float64,
	outer_radius float64,
	M []float64,
	N []float64,
) *Multipliers {
	logres := 0.5
	height, width := 1<<9, 1<<9 // 1<<9 == 512
	inner := cm.AntialiasedCircle(width, height, inner_radius, true, logres)
	outer := cm.AntialiasedCircle(width, height, outer_radius, true, logres)
	annulus := mat.NewDense(height, width, nil)
	annulus.Sub(outer, inner)

	// # Scale each kernel so the sum is 1
	// inner /= np.sum(inner)
	// annulus /= np.sum(annulus)

	// # Precompute the FFT's
	// self.M = np.fft.fft2(inner)
	// self.N = np.fft.fft2(annulus)

	return &Multipliers{
		cm:      CellMath{},
		inner:   inner,
		outer:   outer,
		annulus: annulus,
		M:       M,
		N:       N,
	}
}

type Multipliers struct {
	cm      CellMath
	inner   *mat.Dense
	outer   *mat.Dense
	annulus *mat.Dense
	M       []float64
	N       []float64
}
