package graphics

import (
	"image/color"
)

const (
	TileSize = 32

	WindowWidth  = 500
	WindowHeight = 500

	PixelBufferSize = WindowWidth * WindowHeight * 4

	PlayerX = WindowWidth / 2
	PlayerY = WindowHeight / 2
)

type ScreenPosition struct {
	ScreenX, ScreenY int
}

func (s ScreenPosition) Valid() bool {
	return s.ScreenX >= 0 && s.ScreenX < WindowWidth && s.ScreenY >= 0 && s.ScreenY < WindowHeight
}

type PixelBuffer [PixelBufferSize]byte

func (p *PixelBuffer) Set(position ScreenPosition, clr color.Color) {
	if !position.Valid() {
		return
	}

	rgba := color.RGBAModel.Convert(clr).(color.RGBA)

	pixelIndex := (position.ScreenY*WindowWidth + position.ScreenX) * 4

	p[pixelIndex] = rgba.R
	p[pixelIndex+1] = rgba.G
	p[pixelIndex+2] = rgba.B
	p[pixelIndex+3] = rgba.A
}

func (p *PixelBuffer) DrawRect(xStart, yStart, width, height int, clr color.Color) {
	for x := xStart; x < xStart+width; x++ {
		for y := yStart; y < yStart+height; y++ {
			position := ScreenPosition{
				ScreenX: x,
				ScreenY: y,
			}
			p.Set(position, clr)
		}
	}
}
