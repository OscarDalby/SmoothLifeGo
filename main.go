package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const screenWidth = 600
const screenHeight = 600

type Game struct {
	img *image.RGBA
}

func NewGame(screenWidth int, screenHeight int) *Game {
	return &Game{
		img: image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight)),
	}
}

func (g *Game) Update() error {
	for x := 0; x < screenWidth; x++ {
		for y := 0; y < screenHeight; y++ {
			// this is just a static gradient at the mo
			c := color.RGBA{
				R: uint8((x + y) % 256),
				G: uint8((x * y) % 256),
				B: 128,
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
	fmt.Printf("basicRules: %v", basicRules)
	var cm = CellMath{}
	var inner = cm.AntialiasedCircle(512, 512, 7, true, 0)
	fmt.Printf("inner: %v", inner)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Continuous Domain Cellular Automaton")
	if err := ebiten.RunGame(NewGame(screenWidth, screenHeight)); err != nil {
		log.Fatal(err)
	}
}
