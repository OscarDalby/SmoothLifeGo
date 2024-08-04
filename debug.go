package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"gonum.org/v1/gonum/mat"
)

func printMatrix(m *mat.Dense) {
	r, c := m.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			fmt.Printf("%.2f ", m.At(i, j))
		}
		fmt.Println()
	}
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
		fmt.Println(err)
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
