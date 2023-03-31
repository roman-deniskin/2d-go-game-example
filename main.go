package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	scale float64
	x, y  float64
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.scale += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.scale -= 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {

	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.x -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.x += 1
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	p1 := NewEdge(color.RGBA64{uint16(65535), uint16(0), uint16(0), 1}, 5, g.scale, screen)
	p1.DrawImage(10+g.x, 40+g.y)

	p2 := NewEdge(color.RGBA64{uint16(65535), uint16(65535), uint16(0), 0}, 5, g.scale, screen)
	p2.DrawImage(30+g.x, 40+g.y)

	p3 := NewEdge(color.RGBA64{uint16(65535), uint16(0), uint16(0), 1}, 5, g.scale, screen)
	p3.DrawImage(170+g.x, 170+g.y)

	vector.StrokeLine(screen, p1.x, p1.y, p2.x, p2.y, float32(1+g.scale), color.White, true)
	vector.StrokeLine(screen, p2.x, p2.y, p3.x, p3.y, float32(1+g.scale), color.White, true)
	vector.StrokeLine(screen, p3.x, p3.y, p1.x, p1.y, float32(1+g.scale), color.White, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func createCircle(radius int) *ebiten.Image {
	img := ebiten.NewImage(radius*2, radius*2)
	img.Fill(color.Transparent)

	for y := -radius; y < radius; y++ {
		for x := -radius; x < radius; x++ {
			if x*x+y*y < radius*radius {
				img.Set(x+radius, y+radius, color.White)
			}
		}
	}

	return img
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ebiten Example")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{scale: 0}); err != nil {
		log.Fatal(err)
	}
}
