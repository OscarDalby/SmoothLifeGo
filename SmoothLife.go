package main

import (
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

func ConstructSmoothLife(mp *Multipliers, br BasicRules, width int, height int) *SmoothLife {
	sl := &SmoothLife{
		width:  width,
		height: height,
		mp:     mp,
		rules:  br,
	}
	sl.field = mat.NewCDense(height, width, nil)
	sl.field.Zero()
	if sl.field == nil {
		panic("here the field is nil!")
	}
	return sl
}

type SmoothLife struct {
	width  int
	height int
	mp     *Multipliers
	rules  BasicRules
	field  *mat.CDense
}

func (sl SmoothLife) Clear() {
	sl.field = mat.NewCDense(height, width, nil)
	sl.field.Zero()
}

func (sl SmoothLife) Step() *mat.CDense {
	var newField *mat.CDense = Fft2(sl.field)

	var mBuffer = ElementwiseMultiplyCDenseMatrices(newField, mp.M)
	var nBuffer = ElementwiseMultiplyCDenseMatrices(newField, mp.N)

	var _mBuffer = ifft2(mBuffer)
	var _nBuffer = ifft2(nBuffer)

	// saveMatrixAsImage(RealPartCDenseMatrix(mBuffer), "mBuffer.png")
	// saveMatrixAsImage(RealPartCDenseMatrix(nBuffer), "nBuffer.png")
	// saveMatrixAsImage(RealPartCDenseMatrix(_mBuffer), "_mBuffer.png")
	// saveMatrixAsImage(RealPartCDenseMatrix(_nBuffer), "_nBuffer.png")
	// fmt.Println("mBuffer")
	// printMatrix(RealPartCDenseMatrix(mBuffer))
	// fmt.Println("nBuffer")
	// printMatrix(RealPartCDenseMatrix(nBuffer))
	// fmt.Println("_mBuffer")
	// printMatrix(RealPartCDenseMatrix(_mBuffer))
	// fmt.Println("_nBuffer")
	// printMatrix(RealPartCDenseMatrix(_nBuffer))

	var realMBuffer = RealPartCDenseMatrix(_mBuffer)
	var realNBuffer = RealPartCDenseMatrix(_nBuffer)

	sl.field = ConvertDenseToCDense(sl.rules.S(realNBuffer, realMBuffer))
	return sl.field
}

func (sl SmoothLife) AddSpeckles() {
	rand.New(rand.NewSource(time.Now().Unix()))
	var count int = int(float64(width*height) / math.Pow(float64(mp.outerRadius*2), 2))
	var intensity complex128 = 1.0 + 1i
	for i := 0; i < count; i++ {
		var radius int = int(mp.outerRadius)
		row := rand.Intn(height - radius)
		col := rand.Intn(width - radius)
		for dr := 0; dr < radius; dr++ {
			for dc := 0; dc < radius; dc++ {
				sl.field.Set(row+dr, col+dc, intensity)
			}
		}
	}
}
