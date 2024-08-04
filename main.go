package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/mat"
)

const logres float64 = 0.5
const innerRadius float64 = 20.0
const outerRadius float64 = 60.0
const width int = 1 << 9
const height int = 1 << 9
const screenWidth = width
const screenHeight = height

var mp *Multipliers = ConstructMultipliers(innerRadius, outerRadius, width, height, logres)

// var br BasicRules = BasicRules{B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
// Birth range, survival range, sigmoid widths
var br BasicRules = BasicRules{B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
var sl *SmoothLife = ConstructSmoothLife(mp, br, width, height)

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
	// g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	// key := g.keys[0]
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		sl.AddSpeckles()
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		fmt.Printf("W pressed\n")
	}

	if updateTimer > 0 {
		updateTimer--
		return nil
	} else {
		updateTimer = updateTimerStart
	}
	sl.AddSpeckles()
	newStep := sl.Step()
	rows, cols := newStep.Dims()

	real_sum := 0.0
	r, c := newStep.Dims()
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			real_sum += real((newStep.At(i, j)))
		}
	}
	fmt.Printf("real_sum: %v\n", int(real_sum))

	pix := g.img.Pix
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			index := y*g.img.Stride + x*4
			val := newStep.At(y, x)
			r, i := real(val), imag(val)
			intensity := uint8(math.Round(r*32+i*32)) * 2
			pix[index], pix[index+1], pix[index+2], pix[index+3] = intensity, intensity, intensity, intensity
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
	go func() {
		log.Println("Starting server for profiling at http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()
	var err error
	// endEarly, err := debug() // run any debug functions
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if endEarly {
	// 	return
	// }
	// main logic below
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("SmoothLifeGo")
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
