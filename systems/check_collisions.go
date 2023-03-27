package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/world"
)

type aabb struct {
	x, y, width, height int
}

func aabbFromEntity(entity *world.Entity) (aabb, error) {
	if !entity.Position.Present {
		return aabb{}, newRequireComponentError(entity, "position")
	}
	position := entity.Position.Data.Coordinates()

	if entity.RectCollider.Present {
		rectCollider := entity.RectCollider.Data
		return aabb{
			x:      position.WorldX,
			y:      position.WorldY,
			width:  rectCollider.Width,
			height: rectCollider.Height,
		}, nil
	}

	if entity.ImageBoundsCollider.Present {
		if !entity.Image.Present {
			return aabb{}, newRequireComponentError(entity, "image")
		}
		image := entity.Image.Data
		imageWidth := image.Image.Bounds().Dx()
		imageHeight := image.Image.Bounds().Dy()
		return aabb{
			x:      position.WorldX,
			y:      position.WorldY,
			width:  imageWidth,
			height: imageHeight,
		}, nil
	}

	return aabb{}, newRequireComponentError(entity, "rect collider or image bounds collider")
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

func checkCollisions(w *world.World) error {
	var errs []error
	for first := range w.Entities {
		if !first.Collissions.Present {
			continue
		}
		firstCollisions := &first.Collissions.Data

		firstAabb, err := aabbFromEntity(first)
		if err != nil {
			errs = append(errs, err)
		}

		firstCollisions.Collisions = nil
		for second := range w.Entities {
			if first == second {
				continue
			}

			secondAabb, err := aabbFromEntity(second)
			if err != nil {
				continue
			}

			if firstAabb.collidesWith(secondAabb) {
				firstCollisions.Collisions = append(firstCollisions.Collisions, second)
			}
		}
	}
	return errors.Join(errs...)
}
