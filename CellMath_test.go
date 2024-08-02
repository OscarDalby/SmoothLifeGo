package main

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// almostEqual compares floats with a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestLogisticThreshold(t *testing.T) {
	cases := []struct {
		name         string
		x, x0, alpha float64
		expected     float64
	}{
		{"Negative x", -1, 0.5, 0.1, 8.75651076269652e-27},
		{"x less than x0", 0.1, 0.5, 0.1, 1.12535162055095e-07},
		{"x close to x0 lower", 0.25, 0.5, 0.1, 4.5397868702434395e-05},
		{"x equal to x0", 0.5, 0.5, 0.1, 0.5},
		{"x close to x0 higher", 0.75, 0.5, 0.1, 0.9999546021312976},
		{"x greater than x0", 1, 0.5, 0.1, 0.9999999979388463},
		{"x much greater than x0", 2, 0.5, 0.1, 1.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := LogisticThreshold(tc.x, tc.x0, tc.alpha)
			if !almostEqual(result, tc.expected, 1e-5) {
				t.Errorf("LogisticThreshold(%v, %v, %v) = %v; want %v", tc.x, tc.x0, tc.alpha, result, tc.expected)
			}
		})
	}
}

func TestHardThreshold(t *testing.T) {
	cases := []struct {
		name     string
		x, x0    float64
		expected bool
	}{
		{"Negative x", -1, 0.5, false},
		{"x less than x0", 0.1, 0.5, false},
		{"x close to x0 lower", 0.25, 0.5, false},
		{"x equal to x0", 0.5, 0.5, false},
		{"x close to x0 higher", 0.75, 0.5, true},
		{"x greater than x0", 1, 0.5, true},
		{"x much greater than x0", 2, 0.5, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := HardThreshold(tc.x, tc.x0)
			if result != tc.expected {
				t.Errorf("HardThreshold(%v, %v) = %v; want %v", tc.x, tc.x0, result, tc.expected)
			}
		})
	}
}

func TestLinearisedThreshold(t *testing.T) {
	cases := []struct {
		name         string
		x, x0, alpha float64
		expected     float64
	}{
		{"Negative x", -1, 0.5, 0.1, 0.0},
		{"x less than x0", 0.1, 0.5, 0.1, 0.0},
		{"x close to x0 lower", 0.25, 0.5, 0.1, 0.0},
		{"x equal to x0", 0.5, 0.5, 0.1, 0.5},
		{"x close to x0 higher", 0.75, 0.5, 0.1, 1.0},
		{"x greater than x0", 1, 0.5, 0.1, 1.0},
		{"x much greater than x0", 2, 0.5, 0.1, 1.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := LinearisedThreshold(tc.x, tc.x0, tc.alpha)
			if !almostEqual(result, tc.expected, 1e-5) {
				t.Errorf("LinearisedThreshold(%v, %v, %v) = %v; want %v", tc.x, tc.x0, tc.alpha, result, tc.expected)
			}
		})
	}
}

func TestLogisticInterval(t *testing.T) {
	cases := []struct {
		name           string
		x, a, b, alpha float64
		expected       float64
	}{
		{"x at lower bound", 0.3, 0.3, 0.7, 0.1, 0.49999994373241896},
		{"x close to lower bound", 0.4, 0.3, 0.7, 0.1, 0.9820077563737207},
		{"x in the middle", 0.5, 0.3, 0.7, 0.1, 0.9993294121987771},
		{"x close to upper bound", 0.6, 0.3, 0.7, 0.1, 0.9820077563737207},
		{"x at upper bound", 0.7, 0.3, 0.7, 0.1, 0.49999994373241896},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := LogisticInterval(tc.x, tc.a, tc.b, tc.alpha)
			if !almostEqual(result, tc.expected, 1e-5) {
				t.Errorf("LogisticInterval(%v, %v, %v, %v) = %v; want %v", tc.x, tc.a, tc.b, tc.alpha, result, tc.expected)
			}
		})
	}
}

func TestLerp(t *testing.T) {
	cases := []struct {
		name     string
		a, b, t  float64
		expected float64
	}{
		{"Middle interpolation", 0, 1, 0.5, 0.5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Lerp(tc.a, tc.b, tc.t)
			if !almostEqual(result, tc.expected, 1e-5) {
				t.Errorf("Lerp(%v, %v, %v) = %v; want %v", tc.a, tc.b, tc.t, result, tc.expected)
			}
		})
	}
}

func TestSumDenseMatrix(t *testing.T) {
	cases := []struct {
		name            string
		testDenseMatrix *mat.Dense
		expected        float64
	}{
		{"Summing a dense matrix", mat.NewDense(1, 2, []float64{1, 2}), 3.0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := SumDenseMatrix(tc.testDenseMatrix)
			if !almostEqual(result, tc.expected, 1e-5) {
				t.Errorf("SumDenseMatrix(%v) = %v want %v", tc.testDenseMatrix, result, tc.expected)
			}
		})
	}
}

type DenseMatrixIterationTestCase struct {
	matrix        *mat.Dense
	operatorValue float64
}

func TestDivideDenseMatrix(t *testing.T) {
	cases := []struct {
		name     string
		testData DenseMatrixIterationTestCase
		expected *mat.Dense
	}{
		{
			"Dividing a dense matrix",
			DenseMatrixIterationTestCase{matrix: mat.NewDense(1, 2, []float64{1, 2}), operatorValue: 2.0},
			mat.NewDense(1, 2, []float64{0.5, 1}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := DivideDenseMatrix(tc.testData.matrix, tc.testData.operatorValue)
			if !mat.Equal(result, tc.expected) {
				t.Errorf("DivideDenseMatrix(%v, %v) = %v want %v", tc.testData.matrix, tc.testData.operatorValue, result, tc.expected)
			}
		})
	}
}
