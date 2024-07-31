package main

import (
	"math"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestExp(t *testing.T) {
	cm := CellMath{}
	input := []float64{-1, 0, 1}
	expected := []float64{math.Exp(-1), 1, math.Exp(1)}

	result := cm.Exp(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Exp() = %v, want %v", result, expected)
	}
}

func TestGreater(t *testing.T) {
	cm := CellMath{}
	a := []int{1, 3, 5}
	b := []int{2, 2, 6}
	expected := []bool{false, true, false}

	result := cm.Greater(a, b)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Greater() = %v, want %v", result, expected)
	}
}

func TestClip(t *testing.T) {
	cm := CellMath{}
	input := []int{0, 2, 4, 6}
	min := 1
	max := 5
	expected := []int{1, 2, 4, 5}

	result := cm.Clip(input, min, max)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Clip() = %v, want %v", result, expected)
	}
}

func TestSqrt(t *testing.T) {
	cm := CellMath{}
	input := []float64{1, 4, 9}
	expected := []float64{1, 2, 3}

	result := cm.Sqrt(input)
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Sqrt()[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestRoll(t *testing.T) {
	cm := CellMath{}
	input := []int{1, 2, 3, 4, 5}
	shift := 2
	expected := []int{4, 5, 1, 2, 3}

	result := cm.Roll(input, shift)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Roll() = %v, want %v", result, expected)
	}
}

func TestSum(t *testing.T) {
	cm := CellMath{}
	input := []int{1, 2, 3}
	expected := 6

	result := cm.Sum(input)
	if result != expected {
		t.Errorf("Sum() = %d, want %d", result, expected)
	}
}

func TestZeros(t *testing.T) {
	cm := CellMath{}
	length := 3
	expected := []float64{0, 0, 0}

	result := cm.Zeros(length)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Zeros() = %v, want %v", result, expected)
	}
}

func TestReal(t *testing.T) {
	cm := CellMath{}
	input := []complex128{1 + 2i, 3 + 4i}
	expected := []float64{1, 3}

	result := cm.Real(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Real() = %v, want %v", result, expected)
	}
}

// Continuing from the previous set of tests, adding tests for the rest of the methods

func TestRandomRandint(t *testing.T) {
	cm := CellMath{}
	low, high := 5, 10

	// Since the output is random, we test it multiple times to ensure it stays within expected bounds
	for i := 0; i < 10; i++ {
		result := cm.RandomRandint(low, high)
		if result < low || result >= high {
			t.Errorf("RandomRandint() generated %d, which is outside the bounds [%d, %d)", result, low, high)
		}
	}
}

func TestFft2(t *testing.T) {
	// Example with a simplified case for a known FFT output
	cm := CellMath{}
	rows, cols := 2, 2
	data := []complex128{1, 1, 1, 1}
	matInput := mat.NewCDense(rows, cols, data)

	result := cm.Fft2(matInput, rows, cols)
	expectedValues := []float64{4, 0, 0, 0} // Simplified expected result for FFT of uniform data

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if !floatEquals(result.At(i, j), expectedValues[i*cols+j], 1e-9) {
				t.Errorf("Fft2() at position [%d, %d] = %v, want %v", i, j, result.At(i, j), expectedValues[i*cols+j])
			}
		}
	}
}

func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestLogisticThresholdElementWise(t *testing.T) {
	cm := CellMath{}
	x := []float64{-1, 0, 1}
	x0, alpha := 0.0, 1.0
	expected := []float64{1.0 / (1.0 + math.Exp(4)), 0.5, 1.0 / (1.0 + math.Exp(-4))}

	result := cm.LogisticThresholdElementWise(x, x0, alpha)
	for i, v := range result {
		if !floatEquals(v, expected[i], 1e-4) {
			t.Errorf("LogisticThresholdElementWise()[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestLogisticInterval(t *testing.T) {
	cm := CellMath{}
	x, a, b, alpha := 0.5, 0.0, 1.0, 0.1
	expected := cm.LogisticThreshold(x, a, alpha) * (1.0 - cm.LogisticThreshold(x, b, alpha))

	result := cm.LogisticInterval(x, a, b, alpha)
	if !floatEquals(result, expected, 1e-4) {
		t.Errorf("LogisticInterval() = %v, want %v", result, expected)
	}
}

func TestHardThreshold(t *testing.T) {
	cm := CellMath{}
	x1 := []int{1, 3, 5}
	x2 := []int{2, 2, 6}
	expected := []bool{false, true, false}

	result := cm.HardThreshold(x1, x2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("HardThreshold() = %v, want %v", result, expected)
	}
}

func TestClamp(t *testing.T) {
	cm := CellMath{}
	x, min, max := 0.5, 0.0, 1.0
	expected := 0.5

	result := cm.Clamp(x, min, max)
	if result != expected {
		t.Errorf("Clamp() = %v, want %v", result, expected)
	}
}

func TestLinearisedThresholdElementWise(t *testing.T) {
	cm := CellMath{}
	input := []float64{0, 0.5, 1, 1.5, 2}
	x0, alpha := 1.0, 1.0
	expected := []float64{0, 0.25, 0.5, 0.75, 1}

	result := cm.LinearisedThresholdElementWise(input, x0, alpha)
	for i, v := range result {
		if !floatEquals(v, expected[i], 1e-4) {
			t.Errorf("LinearisedThresholdElementWise()[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestLinearisedInterval(t *testing.T) {
	cm := CellMath{}
	var x float64 = 0.5
	var a float64 = 0
	var b float64 = 1
	var alpha float64 = 0.1
	expected := cm.LinearisedThreshold(x, a, alpha) * (1.0 - cm.LinearisedThreshold(x, b, alpha))

	result := cm.LinearisedInterval(x, a, b, alpha)
	if !floatEquals(result, expected, 1e-4) {
		t.Errorf("LinearisedInterval() = %v, want %v", result, expected)
	}
}
