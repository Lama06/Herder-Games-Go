package world

import "github.com/Lama06/Herder-Games/graphics"

type WorldPosition interface {
	Coordinates() WorldCoordinates
}

func WorldPositionToScreenPosition(world *World, position WorldPosition) graphics.ScreenPosition {
	worldCoordinates := position.Coordinates()
	playerWorldCoordinates := world.Player.Position.Data.Coordinates()
	playerXOffset := worldCoordinates.WorldX - playerWorldCoordinates.WorldX
	playerYOffset := worldCoordinates.WorldY - playerWorldCoordinates.WorldY
	return graphics.ScreenPosition{
		ScreenX: graphics.PlayerX + playerXOffset,
		ScreenY: graphics.PlayerY + playerYOffset,
	}
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

type TilePosition struct {
	TileX, TileY int
}

func (t TilePosition) Coordinates() WorldCoordinates {
	return WorldCoordinates{
		WorldX: t.TileX * graphics.TileSize,
		WorldY: t.TileY * graphics.TileSize,
	}
}
