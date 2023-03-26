package check_collisions

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func checkCollision(
	firstPosition world.WorldCoordinates,
	firstCollider world.RectColliderComponent,
	secondPosition world.WorldCoordinates,
	secondCollider world.RectColliderComponent,
) bool {
	firstMinX := firstPosition.WorldX
	firstMaxX := firstMinX + firstCollider.Width
	firstMinY := firstPosition.WorldY
	firstMaxY := firstMinY + firstCollider.Height

	secondMinX := secondPosition.WorldX
	secondMaxX := secondMinX + secondCollider.Width
	secondMinY := secondPosition.WorldY
	secondMaxY := secondMinY + secondCollider.Height

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

type checkCollisionsSystem struct{}

func New() system.System {
	return checkCollisionsSystem{}
}

func (c checkCollisionsSystem) Update(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.Collissions.Present {
			continue
		}
		entityCollisions := &entity.Collissions.Data

		if !entity.Position.Present {
			errs = append(errs, fmt.Errorf("position component required: %v", entity))
			continue
		}
		entityPosition := entity.Position.Data.Coordinates()

		if !entity.RectCollider.Present {
			errs = append(errs, fmt.Errorf("rect collider component required: %v", entity))
			continue
		}
		entityCollider := entity.RectCollider.Data

		var collisions []*world.Entity
		for other := range w.Entities {
			if entity == other {
				continue
			}

			if !other.Position.Present {
				continue
			}
			otherPosition := other.Position.Data.Coordinates()

			if !other.RectCollider.Present {
				continue
			}
			otherCollider := other.RectCollider.Data

			if checkCollision(entityPosition, entityCollider, otherPosition, otherCollider) {
				collisions = append(collisions, other)
			}
		}

		entityCollisions.Collisions = collisions

	}
	return errors.Join(errs...)
}

func (c checkCollisionsSystem) Draw(w *world.World, image *ebiten.Image) error {
	return nil
}
