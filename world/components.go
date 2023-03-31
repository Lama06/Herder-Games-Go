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
	Width, Height int
	Color         color.Color
}

type VelocityComponent struct {
	VelocityX, VelocityY float64
}

type KeyboardControllerComponent struct {
	Speed float64
}

type MoveToCoordinateComponent struct {
	Coordinate Coordinate
	Speed      float64
}

type MoveToCoordinatesComponent struct {
	Coordinates       []Coordinate
	CurrentCoordinate int
}

func (m *MoveToCoordinatesComponent) SetCoordinates(coordinates []Coordinate) {
	m.Coordinates = coordinates
	m.CurrentCoordinate = 0
}

type PathfinderComponentState byte

const (
	PathfinderComponentStateNotStarted PathfinderComponentState = iota
	PathfinderComponentStateToPortal
	PathfinderComponentStateToDestination
	PathfinderComponentStateFinished
)

type PathfinderComponent struct {
	Destination option.Option[Position]
	State       PathfinderComponentState

	Portal *Entity
}

func (p *PathfinderComponent) SetDestination(destination option.Option[Position]) {
	p.Destination = destination
	p.State = PathfinderComponentStateNotStarted
}

type RectColliderComponent struct {
	Width, Height float64
	Trigger       bool
}

type ImageBoundsColliderComponent struct{}

type PortalComponent struct {
	Destination Position
}
