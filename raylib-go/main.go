package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Gopher struct {
	PosX int32
	PosY int32
	VelX int32
	VelY int32
}

func NewRandGopher(maxX int32, maxY int32) *Gopher {
	x := rand.Int31n(maxX)
	y := rand.Int31n(maxY)
	vx := 1 + rand.Int31n(10)
	vy := 1 + rand.Int31n(10)
	return &Gopher{
		PosX: x,
		PosY: y,
		VelX: vx,
		VelY: vy,
	}
}

func (g *Gopher) Move(w, h, maxX, maxY int32) {
	if g.PosX+g.VelX > maxX+w || g.PosX < 0 {
		g.VelX *= -1
	}

	if g.PosY+g.VelY > maxY+h || g.PosY < 0 {
		g.VelY *= -1
	}

	g.PosX += g.VelX
	g.PosY += g.VelY
}

func main() {
	screenWidth := int32(1920)
	screenHeight := int32(1080)

	rl.InitWindow(1920, 1080, "gopher mark")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	gopherTex := rl.LoadTexture("gopher.png")
	gopherCount := 1000
	gophers := []Gopher{}
	for i := 0; i < gopherCount; i++ {
		gophers = append(gophers, *NewRandGopher(screenWidth, screenHeight))
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyReleased(rl.KeySpace) {
			for i := 0; i < gopherCount; i++ {
				gophers = append(gophers, *NewRandGopher(screenWidth, screenHeight))
			}
		}

		for idx := range gophers {
			gophers[idx].Move(gopherTex.Width, gopherTex.Height, screenWidth-gopherTex.Width, screenHeight-gopherTex.Height)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		for idx := range gophers {
			rl.DrawTexture(gopherTex, gophers[idx].PosX, gophers[idx].PosY, rl.White)
		}
		rl.DrawText(fmt.Sprintf("Gophers: %d", len(gophers)), 8, 8, 14, rl.Black)
		rl.DrawFPS(8, 22)
		rl.EndDrawing()
	}
}
