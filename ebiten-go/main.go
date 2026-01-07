package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/bitmapfont/v4"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const SCREEN_WIDTH int = 1920
const SCREEN_HEIGHT int = 1080

type Gopher struct {
	PosX float64
	PosY float64
	VelX float64
	VelY float64
}

func NewRandGopher(maxX int, maxY int) *Gopher {
	x := rand.Intn(maxX)
	y := rand.Intn(maxY)
	vx := 1 + rand.Intn(10)
	vy := 1 + rand.Intn(10)
	return &Gopher{
		PosX: float64(x),
		PosY: float64(y),
		VelX: float64(vx),
		VelY: float64(vy),
	}
}

func (g *Gopher) Move(w, h, maxX, maxY int) {
	// Check and reverse X velocity
	if g.PosX+g.VelX > float64(maxX-w) || g.PosX+g.VelX < 0 {
		g.VelX *= -1
	}

	// Check and reverse Y velocity
	if g.PosY+g.VelY > float64(maxY-h) || g.PosY+g.VelY < 0 {
		g.VelY *= -1
	}

	// Update position
	g.PosX += g.VelX
	g.PosY += g.VelY

	// Clamp to bounds (prevents vibration)
	if g.PosX < 0 {
		g.PosX = 0
	} else if g.PosX > float64(maxX-w) {
		g.PosX = float64(maxX - w)
	}

	if g.PosY < 0 {
		g.PosY = 0
	} else if g.PosY > float64(maxY-h) {
		g.PosY = float64(maxY - h)
	}
}

type Game struct {
	gopherImg      *ebiten.Image
	gophers        []Gopher
	gopherAddCount int
	fontFace       *text.GoXFace
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	for idx := range g.gophers {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(g.gophers[idx].PosX, g.gophers[idx].PosY)
		screen.DrawImage(g.gopherImg, op)
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("Gophers: %d", len(g.gophers)))
	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(color.Black) // Make it explicitly black
	textOp.GeoM.Translate(10, 10)
	text.Draw(screen, fmt.Sprintf("Gophers: %d\n", len(g.gophers)), g.fontFace, textOp)
	textOp.GeoM.Translate(0, 20)
	text.Draw(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()), g.fontFace, textOp)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (g *Game) Update() error {
	for idx := range g.gophers {
		g.gophers[idx].Move(g.gopherImg.Bounds().Dx(), g.gopherImg.Bounds().Dy(), SCREEN_WIDTH, SCREEN_HEIGHT)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		for i := 0; i < g.gopherAddCount; i++ {
			g.gophers = append(g.gophers, *NewRandGopher(SCREEN_WIDTH, SCREEN_HEIGHT))
		}
	}
	return nil
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		gopherImg:      img,
		gophers:        []Gopher{},
		gopherAddCount: 1000,
		fontFace:       text.NewGoXFace(bitmapfont.Face),
	}

	for i := 0; i < g.gopherAddCount; i++ {
		g.gophers = append(g.gophers, *NewRandGopher(SCREEN_WIDTH, SCREEN_HEIGHT))
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("gopher mark")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
