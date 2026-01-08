package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type Position struct {
	X, Y int32
}

type Velocity struct {
	DX, DY int32
}

func main() {
	screenWidth := int32(1920)
	screenHeight := int32(1080)

	rl.InitWindow(1920, 1080, "gopher mark")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	gopherTex := rl.LoadTexture("gopher.png")
	gopherCount := 1000

	world := ecs.NewWorld()
	mapper := ecs.NewMap2[Position, Velocity](world)
	filter := ecs.NewFilter2[Position, Velocity](world)

	for range gopherCount {
		_ = mapper.NewEntity(
			&Position{X: rand.Int31n(screenWidth - gopherTex.Width), Y: rand.Int31n(screenHeight - gopherTex.Height)},
			&Velocity{DX: 1 + rand.Int31n(10), DY: 1 + rand.Int31n(10)},
		)
	}

	for !rl.WindowShouldClose() {
		if rl.IsKeyReleased(rl.KeySpace) {
			for range gopherCount {
				_ = mapper.NewEntity(
					&Position{X: rand.Int31n(screenWidth - gopherTex.Width), Y: rand.Int31n(screenHeight - gopherTex.Height)},
					&Velocity{DX: 1 + rand.Int31n(10), DY: 1 + rand.Int31n(10)},
				)
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		query := filter.Query()
		for query.Next() {
			pos, vel := query.Get()
			if pos.X+vel.DX > screenWidth-gopherTex.Width || pos.X < 0 {
				vel.DX *= -1
			}

			if pos.Y+vel.DY > screenHeight-gopherTex.Height || pos.Y < 0 {
				vel.DY *= -1
			}
			pos.X += vel.DX
			pos.Y += vel.DY

			rl.DrawTexture(gopherTex, pos.X, pos.Y, rl.White)
		}
		rl.DrawText(fmt.Sprintf("Gophers: %d", world.Stats().Entities.Total), 8, 8, 14, rl.Black)
		rl.DrawFPS(8, 22)

		rl.EndDrawing()
	}
}
