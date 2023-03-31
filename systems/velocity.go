package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/world"
)

func clearVelocity(w *world.World) {
	for entity := range w.Entities {
		if !entity.VelocityComponent.Present {
			continue
		}
		velocity := &entity.VelocityComponent.Data

		velocity.VelocityX = 0
		velocity.VelocityY = 0
	}
}

func moveWithVelocity(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.VelocityComponent.Present {
			continue
		}
		velocityComponent := entity.VelocityComponent.Data

		if !entity.Coordinate.Present {
			errs = append(errs, newRequireComponentError(entity, "coordinate"))
			continue
		}
		coordinate := &entity.Coordinate.Data

		oldCoordinate := *coordinate

		collisionsBeforeMove, _ := getCollidingEntities(w, entity, false)
		*coordinate = oldCoordinate.WorldCoordinate().Add(velocityComponent.VelocityX, velocityComponent.VelocityY)
		collisionsAfterMove, _ := getCollidingEntities(w, entity, false)

		if len(collisionsBeforeMove) == 0 && len(collisionsAfterMove) != 0 {
			*coordinate = oldCoordinate
		}
	}
	return errors.Join(errs...)
}
