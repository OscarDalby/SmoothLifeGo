package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/mat"
)

const logres float64 = 0.5
const radius float64 = 7.0
const width int = 1 << 8
const height int = 1 << 8
const screenWidth = width
const screenHeight = height

var cm CellMath = CellMath{}
var mp *Multipliers = ConstructMultipliers(cm, radius, width, height, logres)

// var br BasicRules = BasicRules{B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
// Birth range, survival range, sigmoid widths
var br BasicRules = BasicRules{B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
var sl *SmoothLife = ConstructSmoothLife(cm, mp, br, width, height)

var game *Game
var matrix *mat.Dense

type Game struct {
	img *image.RGBA
}

func NewGame(screenWidth int, screenHeight int, matrix *mat.Dense) *Game {
	return &Game{
		img: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
}

func init() {
	sl.Clear()
	game = NewGame(screenWidth, screenHeight, matrix)
}

var updateTimerStart = 10
var updateTimer = updateTimerStart

func (g *Game) Update() error {
	if updateTimer > 0 {
		updateTimer--
		return nil
	} else {
		updateTimer = updateTimerStart
	}
	sl.AddSpeckles()
	var newStep *mat.CDense = sl.Step()
	rows, cols := newStep.Dims()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			val := newStep.At(y, x)
			r, i := real(val), imag(val)
			// log.Printf("r: %v", r)
			// log.Printf("i: %v", i)
			intensity := uint8(math.Round(r*255 + i*255))
			c := color.RGBA{
				R: intensity + uint8(x*4*int(intensity)) + uint8(y*4*int(intensity)),
				G: intensity + uint8(x*int(intensity)) + uint8(y*int(intensity)),
				B: intensity + uint8(x*int(intensity)) + uint8(y*int(intensity)),
				A: intensity,
			}
			g.img.SetRGBA(x, y, c)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.img.Pix)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("SmoothLifeGo")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
