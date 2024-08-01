package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/mat"
)

const screenWidth = 600
const screenHeight = 600

var cm CellMath = CellMath{}
var logres float64 = 0.5
var radius float64 = 7.0
var width int = 512
var height int = 512
var mp *Multipliers = ConstructMultipliers(cm, radius, width, height, logres)
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

func (g *Game) Update() error {
	var newStep *mat.CDense = sl.Step()
	// matrix = cm.RealPartCDenseMatrix(newStep)
	rows, cols := newStep.Dims()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			val := newStep.At(y, x)
			r, i := real(val), imag(val)
			log.Printf("r: %v", r)
			log.Printf("i: %v", i)
			intensity := uint8(math.Round(r*255 + i*255))
			c := color.RGBA{
				R: intensity,
				G: intensity,
				B: intensity,
				A: 255,
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
