package world

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type BackgroundRendererComponent struct {
	Color color.Color
}

type RendererComponent struct {
	Layer int
}

type RectRendererComponent struct {
	Width, Height int
	Color         color.Color
}

type ImageRendererComponent struct {
	Image *ebiten.Image
}

type RectColliderComponent struct {
	Width, Height int
}

type CollisionsComponent struct {
	Collisions []*Entity
}

type PreventCollisionsComponent struct{}

type KeyboardControllerComponent struct {
	Speed int
}
