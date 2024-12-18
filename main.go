package main

import (
	"image"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gonum.org/v1/gonum/mat"
)

var logger *log.Logger

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
// var br BasicRules = BasicRules{B1: 0.278, B2: 0.365, D1: 0.267, D2: 0.445, N: 0.028, M: 0.147}
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

var firstRun = true
var updateTimerStart = 5
var updateTimer = updateTimerStart

func (g *Game) Update() error {

	if updateTimer > 0 {
		updateTimer--
		return nil
	} else {
		updateTimer = updateTimerStart
	}

	if firstRun {
		sl.AddSpeckles()
		firstRun = false
	}

	newStep := sl.Step()

	pix := g.img.Pix
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			index := y*g.img.Stride + x*4
			val := newStep.At(y, x)
			r, i := real(val), imag(val)
			intensity := uint8(math.Round(r*8+i*8)) * 8
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
	var err error

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	defer logFile.Close()

	logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	go func() {
		log.Println("Starting server for profiling at http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()
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
