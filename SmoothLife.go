package main

import (
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

func ConstructSmoothLife(cm CellMath, mp *Multipliers, br BasicRules, width int, height int) *SmoothLife {
	sl := &SmoothLife{
		width:  width,
		height: height,
		cm:     cm,
		mp:     mp,
		rules:  br,
	}
	sl.field = mat.NewCDense(sl.height, sl.width, nil)
	sl.field.Zero()
	if sl.field == nil {
		panic("here the field is nil!")
	}
	return sl
}

type SmoothLife struct {
	width  int
	height int
	cm     CellMath
	mp     *Multipliers
	rules  BasicRules
	field  *mat.CDense
}

func (sl SmoothLife) Clear() {
	sl.field = mat.NewCDense(sl.height, sl.width, nil)
	sl.field.Zero()
}

func (sl SmoothLife) Step() *mat.CDense {
	var newField *mat.CDense = sl.cm.Fft2(sl.field)

	var mBuffer = sl.cm.ElementwiseMultiplyCDenseMatrices(newField, sl.mp.M)
	var nBuffer = sl.cm.ElementwiseMultiplyCDenseMatrices(newField, sl.mp.N)

	var _mBuffer = sl.cm.ifft2(mBuffer)
	var _nBuffer = sl.cm.ifft2(nBuffer)

	var realMBuffer = sl.cm.RealPartCDenseMatrix(_mBuffer)
	var realNBuffer = sl.cm.RealPartCDenseMatrix(_nBuffer)

	sl.field = sl.cm.ConvertDenseToCDense(sl.rules.S(sl.cm, realNBuffer, realMBuffer))
	return sl.field
}

func (sl SmoothLife) AddSpeckles() {
	rand.New(rand.NewSource(time.Now().Unix()))
	var count int = int(float64(sl.width*sl.height) / math.Pow(float64(mp.outerRadius*2), 2))
	var intensity complex128 = 1.0 + 1i
	for i := 0; i < count; i++ {
		var radius int = int(mp.outerRadius)
		row := rand.Intn(sl.height - radius)
		col := rand.Intn(sl.width - radius)
		for dr := 0; dr < radius; dr++ {
			for dc := 0; dc < radius; dc++ {
				sl.field.Set(row+dr, col+dc, intensity)
			}
		}
	}
}
