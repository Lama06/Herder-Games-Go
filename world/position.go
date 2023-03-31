package world

const TileSize = 32

type Level int

type TilePosition struct {
	Level Level
	Tile  TileCoordinates
}

type Position struct {
	Level    Level
	Position Coordinates
}

type Coordinates interface {
	WorldCoordinates() WorldCoordinates
}

type WorldCoordinates struct {
	WorldX, WorldY float64
}

func (w WorldCoordinates) WorldCoordinates() WorldCoordinates {
	return w
}

func (w WorldCoordinates) Add(x, y float64) WorldCoordinates {
	return WorldCoordinates{
		WorldX: w.WorldX + x,
		WorldY: w.WorldY + y,
	}
}

type TileCoordinates struct {
	TileX, TileY int
}

func TileCoordinatesFromWorldCoordinates(coordinates WorldCoordinates) TileCoordinates {
	return TileCoordinates{
		TileX: int(coordinates.WorldX / TileSize),
		TileY: int(coordinates.WorldY / TileSize),
	}
}

func (t TileCoordinates) WorldCoordinates() WorldCoordinates {
	return WorldCoordinates{
		WorldX: float64(t.TileX) * TileSize,
		WorldY: float64(t.TileY) * TileSize,
	}
}
