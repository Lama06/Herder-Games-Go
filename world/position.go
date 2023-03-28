package world

import "github.com/Lama06/Herder-Games/graphics"

type Level int

type Position struct {
	Level    Level
	Position Coordinates
}

type Coordinates interface {
	Coordinates() WorldCoordinates
}

type WorldCoordinates struct {
	WorldX, WorldY int
}

func (w WorldCoordinates) Coordinates() WorldCoordinates {
	return w
}

func (w WorldCoordinates) Add(x, y int) WorldCoordinates {
	return WorldCoordinates{
		WorldX: w.WorldX + x,
		WorldY: w.WorldY + y,
	}
}

type TileCoordinates struct {
	TileX, TileY int
}

func (t TileCoordinates) Coordinates() WorldCoordinates {
	return WorldCoordinates{
		WorldX: t.TileX * graphics.TileSize,
		WorldY: t.TileY * graphics.TileSize,
	}
}
