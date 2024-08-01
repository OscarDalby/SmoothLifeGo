package main

import (
	"github.com/mjibson/go-dsp/fft"
	"gonum.org/v1/gonum/mat"
)

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

func (cm CellMath) ifft2(input *mat.CDense) *mat.CDense {
	inputComplexArr := cdenseToSlice(input)
	outputComplexArr := fft.IFFT2(inputComplexArr)
	output := sliceToCDense(outputComplexArr)
	return output
}
