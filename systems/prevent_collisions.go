package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/world"
)

func preventCollisions(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.PreventCollisions.Present {
			continue
		}
		preventCollisions := &entity.PreventCollisions.Data

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
			continue
		}
		position := &entity.Position.Data

		if !entity.Collissions.Present {
			errs = append(errs, newRequireComponentError(entity, "collisions"))
			continue
		}
		collisions := entity.Collissions.Data

		if len(collisions.Collisions) != 0 {
			if !preventCollisions.LastLegalPosition.Present {
				continue
			}
			lastLegalPosition := preventCollisions.LastLegalPosition.Data
			*position = lastLegalPosition
		} else {
			preventCollisions.LastLegalPosition = option.Some(*position)
		}
	}
	return errors.Join(errs...)
}
