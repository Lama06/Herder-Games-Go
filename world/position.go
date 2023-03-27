package world

import "github.com/Lama06/Herder-Games/graphics"

type Position interface {
	Coordinates() Coordinates
}

type Coordinates struct {
	WorldX, WorldY int
}

func (w Coordinates) Coordinates() Coordinates {
	return w
}

func (w Coordinates) Add(x, y int) Coordinates {
	return Coordinates{
		WorldX: w.WorldX + x,
		WorldY: w.WorldY + y,
	}
}

type TilePosition struct {
	TileX, TileY int
}

func (t TilePosition) Coordinates() Coordinates {
	return Coordinates{
		WorldX: t.TileX * graphics.TileSize,
		WorldY: t.TileY * graphics.TileSize,
	}
}
