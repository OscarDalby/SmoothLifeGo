package main

import (
	"fmt"
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

func (cm CellMath) Fft2(input *mat.CDense, rows, cols int) *mat.Dense {
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

// LogisticThreshold computes the sigmoid curve of a provided slice, with alpha adjusting
// the steepness and direction of the transition.
// Used to smoothly transition from near 0 -> 1 as x moves from left to right past x0, with alpha
// a parameter to determine how abruptly this change occurs
func (cm CellMath) LogisticThresholdElementWise(x []float64, x0 float64, alpha float64) []float64 {
	result := make([]float64, len(x))
	for i, val := range x {
		result[i] = 1.0 / (1.0 + math.Exp(-4.0/alpha*(val-x0)))
	}
	return result
}

// LogisticThreshold computes the sigmoid curve of a provided x, with alpha adjusting
// the steepness and direction of the transition.
// Used to smoothly transition from near 0 -> 1 as x moves from left to right past x0, with alpha
// a parameter to determine how abruptly this change occurs
func (CellMath CellMath) LogisticThreshold(x float64, x0 float64, alpha float64) float64 {
	return 1.0 / (1.0 + math.Exp(-4.0/alpha*(x-x0)))
}

// Logistic function on x between a and b with transition width alpha
// ~:
// x < a		: 0
// a < x < b 	: 1
// a > b		: 0
func (cm CellMath) LogisticInterval(x float64, a float64, b float64, alpha float64) float64 {
	return cm.LogisticThreshold(x, a, alpha) * (1.0 - cm.LogisticThreshold(x, b, alpha))
}

func (cm CellMath) HardThreshold(x1 []int, x2 []int) []bool {
	if len(x1) != len(x2) {
		panic("slices are not of the same length and broadcasting is not implemented")
	}
	result := make([]bool, len(x1))
	for i := range x1 {
		result[i] = x1[i] > x2[i]
	}
	return result
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
	for i, val := range x {
		result[i] = cm.Clamp((val-x0)/alpha+0.5, 0, 1)
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

func secondary() {
	CellMath := CellMath{}
	dataToExp := []float64{1, 2, 3}
	exp_result := CellMath.Exp(dataToExp)
	fmt.Println(exp_result)

	a := []int{1, 3, 5}
	b := []int{2, 3, 2}
	greater_result := CellMath.Greater(a, b)
	fmt.Println(greater_result)

	dataToClip := []int{1, 2, 3, 4, 5, 6}
	clipResult := CellMath.Clip(dataToClip, 2, 4)
	fmt.Println(clipResult)

}
