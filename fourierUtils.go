package main

import (
	"github.com/mjibson/go-dsp/fft"
	"gonum.org/v1/gonum/mat"
)

// fourier transforms of real matrices yield complex matrices output
// fourier inverse transforms of complex matrices yield complex matrices output
// fourier inverse transforms of real matrices yield complex matrices with ~=0 imag part output

// ifft2 gets the inverse fast fourier transform of a CDense
func ifft2cdense(input *mat.CDense) *mat.CDense {
	inputComplexArr := cdenseToSlice(input)
	outputComplexArr := fft.IFFT2(inputComplexArr)
	output := sliceToCDense(outputComplexArr)
	return output
}

// Fft2 performs a 2D FFT on a given matrix and returns the result in a new matrix
func fft2cdense(input *mat.CDense) *mat.CDense {
	r, c := input.Dims()

	data := make([][]complex128, r)
	for i := range data {
		data[i] = make([]complex128, c)
		for j := range data[i] {
			data[i][j] = input.At(i, j)
		}
	}
	outputData := fft.FFT2(data)
	output := mat.NewCDense(r, c, nil)
	for i := range outputData {
		for j := range outputData[i] {
			output.Set(i, j, outputData[i][j])
		}
	}

	return output
}

func fft2dense(input *mat.Dense) *mat.CDense {
	r, c := input.Dims()

	data := make([][]complex128, r)
	for i := range data {
		data[i] = make([]complex128, c)
		for j := range data[i] {
			data[i][j] = complex((input.At(i, j)), 0)
		}
	}
	outputData := fft.FFT2(data)
	output := mat.NewCDense(r, c, nil)
	for i := range outputData {
		for j := range outputData[i] {
			output.Set(i, j, outputData[i][j])
		}
	}

	return output
}
