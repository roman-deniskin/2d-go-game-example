package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Color struct {
	r, g, b, a float32
}

type Edge struct {
	x, y   float32
	c      color.RGBA64
	r      int
	op     *ebiten.DrawImageOptions
	i      *ebiten.Image
	scale  float64
	screen *ebiten.Image
}

func NewEdge(c color.RGBA64, r int, s float64, screen *ebiten.Image) *Edge {
	op := &ebiten.DrawImageOptions{}
	rgbaColor := color.RGBA64{R: 65535, G: 32768, B: 0, A: 65535}
	op.ColorScale.ScaleWithColor(rgbaColor)
	op.GeoM.Scale(1+s, 1+s)
	return &Edge{
		c:      c,
		r:      r,
		op:     op,
		i:      initCircle(r),
		screen: screen,
	}
}

func initCircle(radius int) *ebiten.Image {
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

func (e *Edge) DrawImage(x, y float64) {
	e.op.GeoM.Translate(x, y)
	e.x = float32(x)
	e.y = float32(y)
	e.screen.DrawImage(e.i, e.op)
}
