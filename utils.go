package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/davidkleiven/gosfft/sfft"
	"github.com/mjibson/go-dsp/fft"
	"gonum.org/v1/gonum/mat"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// almostEqual compares floats with a tolerance
func almostEqual(a float64, b float64, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// Greater does an element wise greater than and returns a slice of bools containing a[i] > b[i]
func Greater(a, b []int) []bool {
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

// Clamps every element in a slice
func ClampSlice(values []int, min int, max int) []int {
	clamped := make([]int, len(values))
	for i, v := range values {
		if v < min {
			clamped[i] = min
		} else if v > max {
			clamped[i] = max
		} else {
			clamped[i] = v
		}
	}
	return clamped
}

// Gets the sqrt of every element in an array
func SqrtSlice(slice []float64) []float64 {
	result := make([]float64, len(slice))
	for i, val := range slice {
		result[i] = math.Sqrt(val)
	}
	return result
}

func Roll(slice []int, shift int) []int {
	n := len(slice)
	if n == 0 {
		return slice
	}
	shift = ((shift % n) + n) % n
	return append(slice[n-shift:], slice[:n-shift]...)
}

func Sum(slice []int) int {
	total := 0
	for _, value := range slice {
		total += value
	}
	return total
}

func Real(complexSlice []complex128) []float64 {
	realParts := make([]float64, len(complexSlice))
	for i, c := range complexSlice {
		realParts[i] = real(c)
	}
	return realParts
}

func Randint(low int, high int) int {
	// return 3
	fmt.Printf("low: %v, high: %v\n", low, high)
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(high-low) + low
}

func Fft2(input *mat.CDense) *mat.CDense {
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

func Fft2RealIn(input *mat.Dense) *mat.CDense {
	r, c := input.Dims()
	data := make([]complex128, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			data[i*c+j] = complex(input.At(i, j), 0)
		}
	}
	output := sfft.NewFFT2(r, c).FFT(data)
	return mat.NewCDense(r, c, output)
}

// LogisticThreshold computes the sigmoid curve of a provided x, with alpha adjusting
// the steepness and direction of the transition.
// Used to smoothly transition from near 0 -> 1 as x moves from left to right past x0, with alpha
// a parameter to determine how abruptly this change occurs
func LogisticThreshold(x float64, x0 float64, alpha float64) float64 {
	return 1.0 / (1.0 + math.Exp(-4.0/alpha*(x-x0)))
}

// LogisticThreshold computes the sigmoid curve of a provided slice, with alpha adjusting
// the steepness and direction of the transition.
// Used to smoothly transition from near 0 -> 1 as x moves from left to right past x0, with alpha
// a parameter to determine how abruptly this change occurs
func LogisticThresholdElementWise(x []float64, x0 float64, alpha float64) []float64 {
	result := make([]float64, len(x))
	for i, x := range x {
		result[i] = LogisticThreshold(x, x0, alpha)
	}
	return result
}

// Logistic function on x between a and b with transition width alpha
// ~:
// x < a		: 0
// a < x < b 	: 1
// a > b		: 0
func LogisticInterval(x float64, a float64, b float64, alpha float64) float64 {
	return LogisticThreshold(x, a, alpha) * (1.0 - LogisticThreshold(x, b, alpha))
}

// LogisticIntervalElementWise computes the LogisticInterval function on a slice of x values between a and b with transition width alpha
// Element-wise version of LogisticInterval that handles slices of float64.
// This function is useful for performing vectorized-like operations in scenarios where input x is a sequence of values.
func LogisticIntervalElementWise(x []float64, a float64, b float64, alpha float64) []float64 {
	result := make([]float64, len(x))
	for i, xi := range x {
		result[i] = LogisticInterval(xi, a, b, alpha)
	}
	return result
}

// LogisticThresholdElementWise applies LogisticThreshold to each element of x in a mat.Dense matrix.
// func LogisticThresholdDenseElementWise(x *mat.Dense, x0 float64, alpha float64) *mat.Dense {
// 	rows, cols := x.Dims()
// 	result := mat.NewDense(rows, cols, nil)

// 	for i := 0; i < rows; i++ {
// 		for j := 0; j < cols; j++ {
// 			xi := x.At(i, j)
// 			result.Set(i, j, LogisticThreshold(xi, x0, alpha))
// 		}
// 	}
// 	return result
// }

func LogisticThresholdDenseElementWise(x *mat.Dense, x0 float64, alpha float64) *mat.Dense {
	rows, cols := x.Dims()
	result := mat.NewDense(rows, cols, nil)

	var wg sync.WaitGroup
	wg.Add(rows)
	for i := 0; i < rows; i++ {
		go func(row int) {
			defer wg.Done()
			for j := 0; j < cols; j++ {
				xi := x.At(row, j)
				result.Set(row, j, LogisticThreshold(xi, x0, alpha))
			}
		}(i)
	}
	wg.Wait()
	return result
}

// LogisticIntervalElementWise applies LogisticInterval to each element of x in a mat.Dense matrix.
func LogisticIntervalDenseElementWise(x *mat.Dense, a float64, b float64, alpha float64) *mat.Dense {
	rows, cols := x.Dims()
	result := mat.NewDense(rows, cols, nil)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			xi := x.At(i, j)
			result.Set(i, j, LogisticInterval(xi, a, b, alpha))
		}
	}
	return result
}

// Returns if arg1 is greater than arg2
func HardThreshold(arg1 float64, arg2 float64) bool {
	return arg1 > arg2
}

// Clamp keeps a float64 within a range
func Clamp(a float64, aMin float64, aMax float64) float64 {
	var output float64 = a
	if a < aMin {
		output = aMin
	} else if a > aMax {
		output = aMax
	}
	return output
}

// ClampDense applies the clamp operation to each element of the matrix a
func ClampDense(a *mat.Dense, aMin float64, aMax float64) *mat.Dense {
	rows, cols := a.Dims()
	result := mat.NewDense(rows, cols, nil)

	// Apply the clamp operation element-wise
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := a.At(i, j)
			clampedValue := Clamp(val, aMin, aMax)
			result.Set(i, j, clampedValue)
		}
	}
	return result
}

func LinearisedThresholdElementWise(x []float64, x0 float64, alpha float64) []float64 {
	result := make([]float64, len(x))
	for i, x := range x {
		result[i] = LinearisedThreshold(x, x0, alpha)
	}
	return result
}

// Threshold x around x0 with a linear transition region alpha
func LinearisedThreshold(x float64, x0 float64, alpha float64) float64 {
	return Clamp((x-x0)/alpha+0.5, 0, 1)
}

// a<x<b with linearised threshold regions
// ~:
// x < a		: 0
// a < x < b 	: 1
// a > b		: 0
func LinearisedInterval(x float64, a float64, b float64, alpha float64) float64 {
	return LinearisedThreshold(x, a, alpha) * (1.0 - LinearisedThreshold(x, b, alpha))
}

// Linear interpolate from a -> b with t in [0,1]
func Lerp(a, b, t float64) float64 {
	if t < 0 || t > 1 {
		panic("Lerp: t should be in [0,1]")
	}
	return (1.0-t)*a + t*b
}

// Lerp performs linear interpolation between a and b, where t is a matrix of values between [0,1].
func LerpDense(a, b float64, t *mat.Dense) *mat.Dense {
	r, c := t.Dims()
	result := mat.NewDense(r, c, nil)

	t.Apply(func(i, j int, v float64) float64 {
		if v < 0 || v > 1 {
			log.Panicf("LerpDense: t at position (%d, %d) should be in [0,1], got %f", i, j, v)
		}
		result.Set(i, j, (1.0-v)*a+v*b)
		return v
	}, t)

	return result
}

func AntialiasedCircle(sizeX int, sizeY int, radius float64, roll bool, logres float64) *mat.Dense {
	if logres == 0 {
		logres = math.Log2(math.Min(float64(sizeX), float64(sizeY)))
	}

	halfX := float64(sizeX) / 2
	halfY := float64(sizeY) / 2

	logistic := mat.NewDense(sizeY, sizeX, nil)

	for i := 0; i < sizeY; i++ {
		for j := 0; j < sizeX; j++ {
			x := float64(j) - halfX
			y := float64(i) - halfY
			r := math.Sqrt(x*x + y*y)
			expValue := logres * (r - radius)
			logistic.Set(i, j, 1/(1+math.Exp(expValue)))
		}
	}

	if roll {
		logistic = RollMatrix(logistic, sizeY/2, sizeX/2)
	}

	return logistic
}

func RollMatrix(input *mat.Dense, shiftY int, shiftX int) *mat.Dense {
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

func SumDenseMatrix(values *mat.Dense) float64 {
	var total float64
	r, c := values.Dims()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			total += values.At(i, j)
		}
	}
	return total
}

func DivideDenseMatrix(A *mat.Dense, divisor float64) *mat.Dense {
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
func ElementwiseMultiplyCDenseMatrices(A, B *mat.CDense) *mat.CDense {
	rA, cA := A.Dims()
	rB, cB := B.Dims()
	if rA != rB || cA != cB {
		panic("ElementwiseMultiplyCDenseMatrices: Matrices are of different dimensions")
	}

	result := mat.NewCDense(rA, cA, nil)

	for i := 0; i < rA; i++ {
		for j := 0; j < cA; j++ {
			valA := A.At(i, j)
			valB := B.At(i, j)
			result.Set(i, j, complex(real(valA)*real(valB)-imag(valA)*imag(valB), real(valA)*imag(valB)+imag(valA)*real(valB)))
		}
	}
	return result
}

// Helper function which
func sliceToCDense(slice [][]complex128) *mat.CDense {
	rows := len(slice)
	cols := len(slice[0])
	data := make([]complex128, 0, rows*cols)
	for _, row := range slice {
		data = append(data, row...)
	}
	B := mat.NewCDense(rows, cols, data)
	return B
}

func cdenseToSlice(A *mat.CDense) [][]complex128 {
	r, c := A.Dims()
	result := make([][]complex128, r)
	for i := range result {
		result[i] = make([]complex128, c)
		for j := range result[i] {
			result[i][j] = A.At(i, j)
		}
	}
	return result
}

// ifft2 gets the inverse fast fourier transform of a CDense
func ifft2(input *mat.CDense) *mat.CDense {
	inputComplexArr := cdenseToSlice(input)
	outputComplexArr := fft.IFFT2(inputComplexArr)
	output := sliceToCDense(outputComplexArr)
	return output
}

// RealPartCDenseMatrix gets the real part of a CDense and returns a Dense
func RealPartCDenseMatrix(cd *mat.CDense) *mat.Dense {
	r, c := cd.Dims()
	realParts := mat.NewDense(r, c, nil)

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			realParts.Set(i, j, real(cd.At(i, j)))
		}
	}

	return realParts
}

// LogisticThresholdDenseElementWise applies the LogisticThreshold to each element of x with corresponding x0 from matrix x0.
func LogisticThresholdDenseDoubleElementWise(x, x0 *mat.Dense, alpha float64) *mat.Dense {
	rows, cols := x.Dims()
	result := mat.NewDense(rows, cols, nil)

	if x0x, x0y := x0.Dims(); x0x != rows || x0y != cols {
		log.Panic("LogisticThresholdDenseDoubleElementWise: Dimensions of x and x0 must be the same")
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			xi := x.At(i, j)
			x0i := x0.At(i, j)
			result.Set(i, j, LogisticThreshold(xi, x0i, alpha))
		}
	}
	return result
}

// LogisticIntervalDense computes the logistic interval using matrices n, a, and b, with a uniform alpha.
func LogisticIntervalTripleDense(n, a, b *mat.Dense, alpha float64) *mat.Dense {
	nRows, nCols := n.Dims()
	aRows, aCols := a.Dims()
	bRows, bCols := b.Dims()
	if nRows != aRows || nCols != aCols {
		log.Panic("LogisticIntervalTripleDense: Dimensions of n and a must match")
	}
	if nRows != bRows || nCols != bCols {
		log.Panic("LogisticIntervalTripleDense: Dimensions of n and b must match")
	}
	thresholdA := LogisticThresholdDenseDoubleElementWise(n, a, alpha)
	thresholdB := LogisticThresholdDenseDoubleElementWise(n, b, alpha)
	rows, cols := thresholdB.Dims()
	invThresholdB := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			invThresholdB.Set(i, j, 1.0-thresholdB.At(i, j))
		}
	}
	result := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			result.Set(i, j, thresholdA.At(i, j)*invThresholdB.At(i, j))
		}
	}

	return result
}

// ConvertDenseToCDense takes a *mat.Dense matrix and converts it to a *mat.CDense matrix.
func ConvertDenseToCDense(input *mat.Dense) *mat.CDense {
	rows, cols := input.Dims()
	output := mat.NewCDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			realPart := input.At(i, j)
			output.Set(i, j, complex(realPart, 0))
		}
	}
	return output
}

// Add a constant to every ele in a Dense
func AddConstantDense(a *mat.Dense, constant float64) *mat.Dense {
	rows, cols := a.Dims()
	result := mat.NewDense(rows, cols, nil)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := a.At(i, j) + constant
			result.Set(i, j, val)
		}
	}
	return result
}

// NormaliseDense uses the Frobenius to normalise a dense matrix
func NormaliseDense(m *mat.Dense) {
	r, c := m.Dims()
	var sumSquares float64

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			sumSquares += math.Pow(m.At(i, j), 2)
		}
	}
	norm := math.Sqrt(sumSquares)
	if norm == 0 {
		return
	}
	normalized := mat.NewDense(r, c, nil)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			normalized.Set(i, j, m.At(i, j)/norm)
		}
	}
	m.Copy(normalized)
}
