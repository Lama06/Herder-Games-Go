package world

const TileSize = 32

type Level int

type TilePosition struct {
	Level          Level
	TileCoordinate TileCoordinate
}

type Position struct {
	Level      Level
	Coordinate Coordinate
}

type Coordinate interface {
	WorldCoordinate() WorldCoordinate
}

type WorldCoordinate struct {
	WorldX, WorldY float64
}

func (w WorldCoordinate) WorldCoordinate() WorldCoordinate {
	return w
}

func (w WorldCoordinate) Add(x, y float64) WorldCoordinate {
	return WorldCoordinate{
		WorldX: w.WorldX + x,
		WorldY: w.WorldY + y,
	}
}

type TileCoordinate struct {
	TileX, TileY int
}

func TileCoordinateFromCoordinate(coordinate Coordinate) TileCoordinate {
	return TileCoordinate{
		TileX: int(coordinate.WorldCoordinate().WorldX / TileSize),
		TileY: int(coordinate.WorldCoordinate().WorldY / TileSize),
	}
}

func (t TileCoordinate) WorldCoordinate() WorldCoordinate {
	return WorldCoordinate{
		WorldX: float64(t.TileX) * TileSize,
		WorldY: float64(t.TileY) * TileSize,
	}
}
