package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type BackgroundComponent struct {
	Color color.Color
}

type ImageComponent struct {
	Layer int
	Image *ebiten.Image
}

type RectComponent struct {
	Width, Height int
	Color         color.Color
}

type VelocityComponent struct {
	VelocityX, VelocityY float64
}

type KeyboardControllerComponent struct {
	Speed float64
}

type RectColliderComponent struct {
	Width, Height float64
	Trigger       bool
}

type ImageBoundsColliderComponent struct{}

type PortalComponent struct {
	Destination Position
}
