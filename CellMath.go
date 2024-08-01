package main

import (
	"math"
	"math/cmplx"
	"math/rand"
	"time"

	"github.com/davidkleiven/gosfft/sfft"
	"gonum.org/v1/gonum/mat"
)

type CellMath struct {
}

func (cm CellMath) Exp(values []float64) []float64 {
	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = math.Exp(v)
	}
	return result
}

func (cm CellMath) Greater(a, b []int) []bool {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}
	result := make([]bool, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] > b[i]
	}
	return result
}

func (cm CellMath) Clip(values []int, min, max int) []int {
	clipped := make([]int, len(values))
	for i, v := range values {
		if v < min {
			clipped[i] = min
		} else if v > max {
			clipped[i] = max
		} else {
			clipped[i] = v
		}
	}
	return clipped
}

func (cm CellMath) Sqrt(slice []float64) []float64 {
	result := make([]float64, len(slice))
	for i, val := range slice {
		result[i] = math.Sqrt(val)
	}
	return result
}

func (cm CellMath) Roll(slice []int, shift int) []int {
	n := len(slice)
	if n == 0 {
		return slice
	}
	shift = ((shift % n) + n) % n
	return append(slice[n-shift:], slice[:n-shift]...)
}

func (cm CellMath) Sum(slice []int) int {
	total := 0
	for _, value := range slice {
		total += value
	}
	return total
}

func (cm CellMath) Zeros(length int) []float64 {
	return make([]float64, length)
}

func (cm CellMath) Real(complexSlice []complex128) []float64 {
	realParts := make([]float64, len(complexSlice))
	for i, c := range complexSlice {
		realParts[i] = real(c)
	}
	return realParts
}

func (cm CellMath) RandomRandint(low, high int) int {
	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	return rnd.Intn(high-low) + low
}

func (cm CellMath) Fft2RealOut(input *mat.CDense, rows, cols int) *mat.Dense {
	ft := sfft.NewFFT2(rows, cols)
	ftData := ft.FFT(input.RawCMatrix().Data)
	ftMat := mat.NewCDense(rows, cols, ftData)
	sfft.Center2(ftMat)

	amp := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			amp.Set(i, j, cmplx.Abs(ftMat.At(i, j)))
		}
	}
	return amp
}

func (cm CellMath) Fft2(input *mat.CDense) *mat.CDense {
	r, c := input.Dims()

	fft := sfft.NewFFT2(r, c)
	data := make([]complex128, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			var comp complex128 = input.At(i, j)
			realPart := real(comp)
			imagPart := imag(comp)
			data[i*c+j] = complex(realPart, imagPart)
		}
	}
	output := fft.FFT(data)
	return mat.NewCDense(r, c, output)
}

func (cm CellMath) Fft2RealIn(input *mat.Dense) *mat.CDense {
	r, c := input.Dims()
	fft := sfft.NewFFT2(r, c)
	data := make([]complex128, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = complex(input.At(i, j), 0)
		}
	}
	output := fft.FFT(data)
	return mat.NewCDense(r, c, output)
}

// LogisticThreshold computes the sigmoid curve of a provided x, with alpha adjusting
// the steepness and direction of the transition.
// Used to smoothly transition from near 0 -> 1 as x moves from left to right past x0, with alpha
// a parameter to determine how abruptly this change occurs
func (CellMath CellMath) LogisticThreshold(x float64, x0 float64, alpha float64) float64 {
	return 1.0 / (1.0 + math.Exp(-4.0/alpha*(x-x0)))
}

// LogisticThreshold computes the sigmoid curve of a provided slice, with alpha adjusting
// the steepness and direction of the transition.
// Used to smoothly transition from near 0 -> 1 as x moves from left to right past x0, with alpha
// a parameter to determine how abruptly this change occurs
func (cm CellMath) LogisticThresholdElementWise(x []float64, x0 float64, alpha float64) []float64 {
	result := make([]float64, len(x))
	for i, x := range x {
		result[i] = cm.LogisticThreshold(x, x0, alpha)
	}
	return result
}

// Logistic function on x between a and b with transition width alpha
// ~:
// x < a		: 0
// a < x < b 	: 1
// a > b		: 0
func (cm CellMath) LogisticInterval(x float64, a float64, b float64, alpha float64) float64 {
	return cm.LogisticThreshold(x, a, alpha) * (1.0 - cm.LogisticThreshold(x, b, alpha))
}

// Returns if arg1 is greater than arg2
func (cm CellMath) HardThreshold(arg1 float64, arg2 float64) bool {
	return arg1 > arg2
}

func (cm CellMath) Clamp(a float64, aMin float64, aMax float64) float64 {
	var output float64 = a
	if a < aMin {
		output = aMin
	} else if a > aMax {
		output = aMax
	}
	return output
}

func (cm CellMath) LinearisedThresholdElementWise(x []float64, x0 float64, alpha float64) []float64 {
	result := make([]float64, len(x))
	for i, x := range x {
		result[i] = cm.LinearisedThreshold(x, x0, alpha)
	}
	return result
}

// Threshold x around x0 with a linear transition region alpha
func (cm CellMath) LinearisedThreshold(x float64, x0 float64, alpha float64) float64 {
	return cm.Clamp((x-x0)/alpha+0.5, 0, 1)
}

// a<x<b with linearised threshold regions
// ~:
// x < a		: 0
// a < x < b 	: 1
// a > b		: 0
func (cm CellMath) LinearisedInterval(x float64, a float64, b float64, alpha float64) float64 {
	return cm.LinearisedThreshold(x, a, alpha) * (1.0 - cm.LinearisedThreshold(x, b, alpha))
}

// Linear interpolate from a -> b with t in [0,1]
func (cm CellMath) Lerp(a, b, t float64) float64 {
	if t < 0 || t > 1 {
		panic("t should be in [0,1]")
	}
	return (1.0-t)*a + t*b
}

func (cm CellMath) MeshGrid(sizeY, sizeX int) (*mat.Dense, *mat.Dense) {
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

func (cm CellMath) AntialiasedCircle(sizeX int, sizeY int, radius float64, roll bool, logres float64) *mat.Dense {
	yy, xx := cm.MeshGrid(sizeY, sizeX)
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
		logistic = cm.RollMatrix(logistic, sizeY/2, sizeX/2)
	}

	return logistic
}

func (cm CellMath) RollMatrix(input *mat.Dense, shiftY, shiftX int) *mat.Dense {
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

func (cm CellMath) SumDenseMatrix(values *mat.Dense) float64 {
	var total float64
	r, c := values.Dims()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			total += values.At(i, j)
		}
	}
	return total
}

func (cm CellMath) DivideDenseMatrix(A *mat.Dense, divisor float64) *mat.Dense {
	r, c := A.Dims()
	result := mat.NewDense(r, c, nil)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			val := A.At(i, j)
			result.Set(i, j, val/divisor)
		}
	}
	return result
}

// ElementwiseMultiplyCDenseMatrices multiplies two complex matrices element-wise.
func (cm CellMath) ElementwiseMultiplyCDenseMatrices(A, B *mat.CDense) *mat.CDense {
	r, c := A.Dims()

	rB, cB := B.Dims()
	if r != rB || c != cB {
		panic("matrices are of different dimensions")
	}

	result := mat.NewCDense(r, c, nil)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			valA := A.At(i, j)
			valB := B.At(i, j)
			result.Set(i, j, complex(real(valA)*real(valB)-imag(valA)*imag(valB), real(valA)*imag(valB)+imag(valA)*real(valB)))
		}
	}
	return result
}
