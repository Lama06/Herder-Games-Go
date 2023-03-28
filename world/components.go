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

type RectComponent struct {
	Layer         int
	Width, Height int
	Color         color.Color
}

type RectColliderComponent struct {
	Width, Height int
}

type ImageBoundsColliderComponent struct{}

type PreventCollisionsComponent struct {
	LastLegalPosition option.Option[Position]
}

type KeyboardControllerComponent struct {
	Speed int
}

type PortalComponent struct {
	Destination Position
}
