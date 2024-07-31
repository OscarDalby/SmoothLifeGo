package main

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

func MeshGrid(sizeY, sizeX int) (*mat.Dense, *mat.Dense) {
	yy := mat.NewDense(sizeY, sizeX, nil)
	xx := mat.NewDense(sizeY, sizeX, nil)

	for i := 0; i < sizeY; i++ {
		for j := 0; j < sizeX; j++ {
			yy.Set(i, j, float64(i))
			xx.Set(i, j, float64(j))
		}
	}

	return yy, xx
}

func AntialiasedCircle(sizeX, sizeY int, radius float64, roll bool, logres float64) *mat.Dense {
	yy, xx := MeshGrid(sizeY, sizeX)
	logistic := mat.NewDense(sizeY, sizeX, nil)

	if logres == 0 {
		logres = math.Log2(math.Min(float64(sizeX), float64(sizeY)))
	}

	for i := 0; i < sizeY; i++ {
		for j := 0; j < sizeX; j++ {
			x := float64(xx.At(i, j)) - float64(sizeX)/2
			y := float64(yy.At(i, j)) - float64(sizeY)/2
			r := math.Sqrt(x*x + y*y)
			expValue := logres * (r - radius)
			logistic.Set(i, j, 1/(1+math.Exp(expValue)))
		}
	}

	if roll {
		logistic = rollMatrix(logistic, sizeY/2, sizeX/2)
	}

	return logistic
}

func rollMatrix(input *mat.Dense, shiftY, shiftX int) *mat.Dense {
	r, c := input.Dims()
	output := mat.NewDense(r, c, nil)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			newI := (i + shiftY) % r
			newJ := (j + shiftX) % c
			output.Set(newI, newJ, input.At(i, j))
		}
	}

	return output
}
