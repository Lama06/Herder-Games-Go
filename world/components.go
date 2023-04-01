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

type MoveToCoordinateComponentState byte

const (
	MoveToCoordinateComponentStateDisabled MoveToCoordinateComponentState = iota
	MoveToCoordinateComponentStateMoving
	MoveToCoordinateComponentStateArrived
)

type MoveToCoordinateComponent struct {
	State      MoveToCoordinateComponentState
	Coordinate Coordinate
}

func (m *MoveToCoordinateComponent) Disable() {
	m.State = MoveToCoordinateComponentStateDisabled
	m.Coordinate = nil
}

func (m *MoveToCoordinateComponent) SetCoordinate(coordiante Coordinate) {
	if coordiante == nil {
		panic("cannot move to nil coordinate")
	}

	m.State = MoveToCoordinateComponentStateMoving
	m.Coordinate = coordiante
}

type MoveToCoordinatesComponentState byte

const (
	MoveToCoordinatesComponentStateDisabled MoveToCoordinatesComponentState = iota
	MoveToCoordinatesComponentStateIdle
	MoveToCoordinatesComponentStateMoving
	MoveToCoordinatesComponentStateFinished
)

type MoveToCoordinatesComponent struct {
	State             MoveToCoordinatesComponentState
	Coordinates       []Coordinate
	CurrentCoordinate int
}

func (m *MoveToCoordinatesComponent) Disable() {
	m.State = MoveToCoordinatesComponentStateDisabled
	m.Coordinates = nil
	m.CurrentCoordinate = 0
}

func (m *MoveToCoordinatesComponent) SetCoordinates(coordiantes []Coordinate) {
	if len(coordiantes) == 0 {
		panic("cannot move to no coordinates")
	}

	m.State = MoveToCoordinatesComponentStateIdle
	m.Coordinates = coordiantes
	m.CurrentCoordinate = 0
}

type PathfinderComponentState byte

const (
	PathfinderComponentStateDisabled PathfinderComponentState = iota
	PathfinderComponentStateIdle
	PathfinderComponentStateNoPath
	PathfinderComponentStateMovingToPortal
	PathfinderComponentStateMovingToDestination
	PathfinderComponentStateArrived
)

type PathfinderComponent struct {
	State       PathfinderComponentState
	Destination Position
}

func (p *PathfinderComponent) Disable() {
	p.State = PathfinderComponentStateDisabled
	p.Destination = Position{}
}

func (p *PathfinderComponent) SetDestination(destination Position) {
	p.State = PathfinderComponentStateIdle
	p.Destination = destination
}

type RectColliderComponent struct {
	Width, Height float64
	Trigger       bool
}

type ImageBoundsColliderComponent struct{}

type PortalComponent struct {
	Destination Position
}
