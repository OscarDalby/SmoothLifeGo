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

type Game struct {
	img    *image.RGBA
	matrix *mat.Dense
}

func NewGame(screenWidth int, screenHeight int, matrix *mat.Dense) *Game {
	return &Game{
		img:    image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
		matrix: matrix,
	}
}

func (g *Game) Update() error {
	rows, cols := g.matrix.Dims()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			val := g.matrix.At(y, x)
			intensity := uint8(math.Round(val * 255))
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
	cm := CellMath{}
	radius := 7.0
	// mp := ConstructMultipliers(cm, radius)
	// width := 3
	// height := 3
	// br := BasicRules{B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
	// sl := ConstructSmoothLife(cm, mp, br, width, height)
	logres := 0.0
	matrix := cm.AntialiasedCircle(screenWidth, screenHeight, radius, true, logres)
	game := NewGame(screenWidth, screenHeight, matrix)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("AntialiasedCircle Viz")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
