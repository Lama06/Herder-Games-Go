package systems

import "github.com/Lama06/Herder-Games/world"

type aabb struct {
	x, y, width, height int
}

func aabbFromEntity(entity *world.Entity) (aabb, error) {
	if !entity.Position.Present {
		return aabb{}, newRequireComponentError(entity, "position")
	}
	position := entity.Position.Data.Coordinates()

	if !entity.RectCollider.Present {
		return aabb{}, newRequireComponentError(entity, "rect collider")
	}
	rectCollider := entity.RectCollider.Data

	return aabb{
		x:      position.WorldX,
		y:      position.WorldY,
		width:  rectCollider.Width,
		height: rectCollider.Height,
	}, nil

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
