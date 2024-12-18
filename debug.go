package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/mat"
)

func cdenseToGrayImage(m *mat.CDense) *image.Gray {
	r, c := m.Dims()

	img := image.NewGray(image.Rect(0, 0, c, r))

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			val := m.At(i, j)
			pixel := uint8(real(val))
			img.SetGray(j, i, color.Gray{Y: pixel})
		}
	}

	return img
}

func cdenseToEbitenImage(m *mat.CDense) *ebiten.Image {
	rows, cols := m.Dims()
	img := ebiten.NewImage(cols, rows)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := m.At(i, j)

			intensity := uint8(math.Abs(real(val)))
			fmt.Printf("val: %v\n", val)
			clr := color.RGBA{R: intensity, G: intensity, B: intensity, A: 255}
			img.Set(j, i, clr)
		}
	}

	return img
}

func saveMatrixAsImage(m *mat.Dense, filename string) {
	r, c := m.Dims()
	img := image.NewGray(image.Rect(0, 0, r, c))
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			val := m.At(i, j)
			gray := uint8(255 * val)
			img.SetGray(j, i, color.Gray{Y: gray})
		}
	}
	f, err := os.Create(fmt.Sprintf("debug/%v", filename))
	if err != nil {
		return
	}
	defer f.Close()
	png.Encode(f, img)
}

func debug() (bool, error) {
	var sizeX int = 200
	var sizeY int = 200
	var rad float64 = 50
	var roll bool = false
	var logres float64 = 0.5

	inner := AntialiasedCircle(sizeX, sizeY, rad, roll, logres)
	outer := AntialiasedCircle(sizeX, sizeY, rad*3, roll, logres)
	annulus := mat.NewDense(sizeX, sizeY, nil)
	annulus.Sub(outer, inner)
	saveMatrixAsImage(annulus, "og_circle.png")
	return true, nil
}

func printDenseSum(m *mat.Dense, matrixName string) {
	var mSum float64 = 0
	rm, cm := m.Dims()
	for i := 0; i < rm; i++ {
		for j := 0; j < cm; j++ {
			mSum += m.At(i, j)
		}
	}
	fmt.Printf("Dense sum of %v: %v\n", matrixName, mSum)
}

func cdenseRealSum(m *mat.CDense) float64 {
	var mSum float64 = 0
	rm, cm := m.Dims()
	for i := 0; i < rm; i++ {
		for j := 0; j < cm; j++ {
			mSum += real(m.At(i, j))
		}
	}
	return mSum
}

func printCDenseRealSum(m *mat.CDense, matrixName string) {
	var mSum float64 = 0
	rm, cm := m.Dims()
	for i := 0; i < rm; i++ {
		for j := 0; j < cm; j++ {
			mSum += real(m.At(i, j))
		}
	}
	fmt.Printf("CDense real sum of %v: %v\n", matrixName, mSum)
}
