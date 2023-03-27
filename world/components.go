package world

import (
	"image/color"

	"github.com/Lama06/Herder-Games/option"
	"github.com/hajimehoshi/ebiten/v2"
)

type BackgroundComponent struct {
	Color color.Color
}

type ImageComponent struct {
	Layer int
	Image *ebiten.Image
}

type RectColliderComponent struct {
	Width, Height int
}

type ImageBoundsColliderComponent struct{}

type CollisionsComponent struct {
	Collisions []*Entity
}

type PreventCollisionsComponent struct {
	LastLegalPosition option.Option[Position]
}

type KeyboardControllerComponent struct {
	Speed int
}
