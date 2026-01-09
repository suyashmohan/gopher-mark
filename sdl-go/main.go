package main

import (
	"log"
	"math/rand"

	"github.com/jupiterrider/purego-sdl3/img"
	"github.com/jupiterrider/purego-sdl3/sdl"
)

type Gopher struct {
	PosX float32
	PosY float32
	VelX float32
	VelY float32
}

func NewRandGopher(w, h, maxX, maxY float32) Gopher {
	x := rand.Float32() * (maxX - w)
	y := rand.Float32() * (maxY - h)
	vx := 1.0 + rand.Float32()*10
	vy := 1.0 + rand.Float32()*10
	return Gopher{
		PosX: float32(x),
		PosY: float32(y),
		VelX: float32(vx),
		VelY: float32(vy),
	}
}

func (g *Gopher) Move(w, h, maxX, maxY float32) {
	if g.PosX+g.VelX > maxX-w || g.PosX < 0 {
		g.VelX *= -1
	}

	if g.PosY+g.VelY > maxY-h || g.PosY < 0 {
		g.VelY *= -1
	}

	g.PosX += g.VelX
	g.PosY += g.VelY
}

func main() {
	const screenWidth float32 = 1920
	const screenHeight float32 = 1080

	if !sdl.SetHint(sdl.HintRenderVSync, "1") {
		log.Panic(sdl.GetError())
	}

	defer sdl.Quit()
	if !sdl.Init(sdl.InitVideo) {
		log.Panic(sdl.GetError())
	}

	var window *sdl.Window
	var renderer *sdl.Renderer
	if !sdl.CreateWindowAndRenderer("gopher mark", int32(screenWidth), int32(screenHeight), sdl.WindowResizable, &window, &renderer) {
		log.Panic(sdl.GetError())
	}
	defer sdl.DestroyWindow(window)
	defer sdl.DestroyRenderer(renderer)

	texture := img.LoadTexture(renderer, "gopher.png")
	if texture == nil {
		log.Panic(sdl.GetError())
	}
	textureWidth := float32(texture.W)
	textureHeight := float32(texture.H)
	defer sdl.DestroyTexture(texture)

	gopherAddCount := 1000
	gophers := make([]Gopher, gopherAddCount)
	for i := 0; i < gopherAddCount; i++ {
		gophers[i] = NewRandGopher(textureWidth, textureHeight, screenWidth, screenHeight)
	}

	var frameCount uint32
	var fpsUpdateTime uint64 = sdl.GetTicks()
	//var lastFrameTime uint64 = sdl.GetTicks()
	//var deltaTime float64
	running := true

	sdl.SetRenderDrawColor(renderer, 255, 255, 255, sdl.AlphaOpaque)
	for running {
		currentTime := sdl.GetTicks()
		//deltaTime = float64(currentTime-lastFrameTime) / 1000.0
		//lastFrameTime = currentTime

		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				running = false
			case sdl.EventKeyDown:
				if event.Key().Scancode == sdl.ScancodeEscape {
					running = false
				}
			}
		}

		frameCount++
		if currentTime-fpsUpdateTime >= 1000 {
			fps := float64(frameCount) / ((float64(currentTime) - float64(fpsUpdateTime)) / 1000.0)
			if frameCount > 58 {
				newBatch := make([]Gopher, gopherAddCount)
				for i := 0; i < gopherAddCount; i++ {
					newBatch[i] = NewRandGopher(textureWidth, textureHeight, screenWidth, screenHeight)
				}
				gophers = append(gophers, newBatch...)
			}
			log.Printf("FPS: %.2f --- Gophers: %d\n", fps, len(gophers))
			frameCount = 0
			fpsUpdateTime = currentTime
		}

		sdl.RenderClear(renderer)
		var dstRect sdl.FRect
		for idx := range gophers {
			gophers[idx].Move(textureWidth, textureHeight, screenWidth, screenHeight)
			dstRect.X = gophers[idx].PosX
			dstRect.Y = gophers[idx].PosY
			dstRect.W = textureWidth
			dstRect.H = textureHeight
			sdl.RenderTexture(renderer, texture, nil, &dstRect)
		}
		sdl.RenderPresent(renderer)
	}
}
