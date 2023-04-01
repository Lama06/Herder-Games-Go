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

type MoveSpeedComponent struct {
	Speed float64
}

type MoveToCoordinateComponent struct {
	Coordinate Coordinate
	Arrived    bool
}

func NewMoveToCoordinateComponent(coordiante Coordinate) MoveToCoordinateComponent {
	if coordiante == nil {
		panic("cannot move to nil coordinate")
	}

	return MoveToCoordinateComponent{
		Coordinate: coordiante,
		Arrived:    false,
	}
}

type MoveToCoordinatesComponentState byte

const (
	MoveToCoordinatesComponentStateIdle MoveToCoordinatesComponentState = iota
	MoveToCoordinatesComponentStateMoving
	MoveToCoordinatesComponentStateFinished
)

type MoveToCoordinatesComponent struct {
	State             MoveToCoordinatesComponentState
	Coordinates       []Coordinate
	CurrentCoordinate int
}

func NewMoveToCoordinatesComponent(coordinates []Coordinate) MoveToCoordinatesComponent {
	if len(coordinates) == 0 {
		panic("coordinates are empty")
	}

	return MoveToCoordinatesComponent{
		State:             MoveToCoordinatesComponentStateIdle,
		Coordinates:       coordinates,
		CurrentCoordinate: 0,
	}
}

type PathfinderComponentState byte

const (
	PathfinderComponentStateIdle PathfinderComponentState = iota
	PathfinderComponentStateMovingToPortal
	PathfinderComponentStateMovingToDestination
	PathfinderComponentStateNoPath
	PathfinderComponentStateArrived
)

type PathfinderComponent struct {
	State       PathfinderComponentState
	Destination Position
}

func NewPathfinderComponent(destination Position) PathfinderComponent {
	return PathfinderComponent{
		State:       PathfinderComponentStateIdle,
		Destination: destination,
	}
}

type RectColliderComponent struct {
	Width, Height float64
	Trigger       bool
}

type ImageBoundsColliderComponent struct{}

type PortalComponent struct {
	Destination Position
}
