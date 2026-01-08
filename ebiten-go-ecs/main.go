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
	"github.com/mlange-42/ark/ecs"
)

const SCREEN_WIDTH int = 1920
const SCREEN_HEIGHT int = 1080

type Position struct {
	X, Y float64
}

type Velocity struct {
	DX, DY float64
}

type Game struct {
	gopherImg      *ebiten.Image
	gopherAddCount int
	fontFace       *text.GoXFace
	world          *ecs.World
	mapper         *ecs.Map2[Position, Velocity]
	filter         *ecs.Filter2[Position, Velocity]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})

	op := &ebiten.DrawImageOptions{}
	query := g.filter.Query()
	for query.Next() {
		pos, _ := query.Get()
		op.GeoM.Reset()
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(g.gopherImg, op)
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("Gophers: %d", len(g.gophers)))
	textOp := &text.DrawOptions{}
	textOp.ColorScale.ScaleWithColor(color.Black) // Make it explicitly black
	textOp.GeoM.Translate(10, 10)
	text.Draw(screen, fmt.Sprintf("Gophers: %d\n", g.world.Stats().Entities.Total), g.fontFace, textOp)
	textOp.GeoM.Translate(0, 20)
	text.Draw(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()), g.fontFace, textOp)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func (g *Game) Update() error {
	query := g.filter.Query()
	for query.Next() {
		pos, vel := query.Get()
		if pos.X+vel.DX > float64(SCREEN_WIDTH-g.gopherImg.Bounds().Dx()) || pos.X < 0 {
			vel.DX *= -1
		}

		if pos.Y+vel.DY > float64(SCREEN_HEIGHT-g.gopherImg.Bounds().Dy()) || pos.Y < 0 {
			vel.DY *= -1
		}
		pos.X += vel.DX
		pos.Y += vel.DY
	}
	if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		for range g.gopherAddCount {
			_ = g.mapper.NewEntity(
				&Position{X: float64(rand.Intn(SCREEN_WIDTH - g.gopherImg.Bounds().Dx())), Y: float64(rand.Intn(SCREEN_HEIGHT - g.gopherImg.Bounds().Dy()))},
				&Velocity{DX: float64(1 + rand.Int31n(10)), DY: float64(1 + rand.Int31n(10))},
			)
		}
	}
	return nil
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("gopher.png")
	if err != nil {
		log.Fatal(err)
	}

	world := ecs.NewWorld()
	mapper := ecs.NewMap2[Position, Velocity](world)
	filter := ecs.NewFilter2[Position, Velocity](world)

	g := &Game{
		gopherImg:      img,
		gopherAddCount: 1000,
		fontFace:       text.NewGoXFace(bitmapfont.Face),
		world:          world,
		mapper:         mapper,
		filter:         filter,
	}

	for range g.gopherAddCount {
		_ = mapper.NewEntity(
			&Position{X: float64(rand.Intn(SCREEN_WIDTH - img.Bounds().Dx())), Y: float64(rand.Intn(SCREEN_HEIGHT - img.Bounds().Dy()))},
			&Velocity{DX: float64(1 + rand.Int31n(10)), DY: float64(1 + rand.Int31n(10))},
		)
	}

	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("gopher mark")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
