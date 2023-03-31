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

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
			continue
		}
		position := &entity.Position.Data

		oldPosition := *position

		collisionsBeforeMove, _ := getCollidingEntities(w, entity, false)
		*position = oldPosition.WorldCoordinates().Add(velocityComponent.VelocityX, velocityComponent.VelocityY)
		collisionsAfterMove, _ := getCollidingEntities(w, entity, false)

		if len(collisionsBeforeMove) == 0 && len(collisionsAfterMove) != 0 {
			*position = oldPosition
		}
	}
	return errors.Join(errs...)
}
