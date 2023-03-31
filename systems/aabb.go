package systems

import (
	"github.com/Lama06/Herder-Games/world"
)

type aabb struct {
	x, y, width, height float64
}

func aabbFromEntity(entity *world.Entity) (bounds aabb, trigger bool, err error) {
	if !entity.Position.Present {
		return aabb{}, false, newRequireComponentError(entity, "position")
	}
	position := entity.Position.Data.WorldCoordinates()

	if !entity.RectColliderComponent.Present {
		return aabb{}, false, newRequireComponentError(entity, "rect collider")
	}
	rectColliderComponent := entity.RectColliderComponent.Data

	return aabb{
		x:      position.WorldX,
		y:      position.WorldY,
		width:  rectColliderComponent.Width,
		height: rectColliderComponent.Height,
	}, rectColliderComponent.Trigger, nil
}

func (first aabb) collidesWith(second aabb) bool {
	firstMinX := first.x
	firstMaxX := firstMinX + first.width
	firstMinY := first.y
	firstMaxY := firstMinY + first.height

	secondMinX := second.x
	secondMaxX := secondMinX + second.width
	secondMinY := second.y
	secondMaxY := secondMinY + second.height

	if firstMinX >= secondMaxX {
		return false
	}

	if firstMaxX <= secondMinX {
		return false
	}

	if firstMinY >= secondMaxY {
		return false
	}

	if firstMaxY <= secondMinY {
		return false
	}

	return true
}

func (bounds aabb) blockedTiles() map[world.TileCoordinates]struct{} {
	result := make(map[world.TileCoordinates]struct{})

	numTilesX := int(bounds.width/world.TileSize) + 1
	numTilesY := int(bounds.height/world.TileSize) + 1

	for tileOffsetX := 0; tileOffsetX < numTilesX; tileOffsetX++ {
		for tileOffsetY := 0; tileOffsetY < numTilesY; tileOffsetY++ {
			coordinates := world.WorldCoordinates{
				WorldX: bounds.x + float64(tileOffsetX*world.TileSize),
				WorldY: bounds.y + float64(tileOffsetY*world.TileSize),
			}
			result[world.TileCoordinatesFromWorldCoordinates(coordinates)] = struct{}{}
		}
	}

	return result
}
