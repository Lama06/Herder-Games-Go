package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/world"
)

func addImageBoundsColliders(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.ImageBoundsColliderComponent.Present {
			continue
		}

		if !entity.ImageComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "image"))
			continue
		}
		imageComponent := entity.ImageComponent.Data

		imageWidth := imageComponent.Image.Bounds().Dx()
		imageHeight := imageComponent.Image.Bounds().Dy()

		if !entity.RectColliderComponent.Present {
			entity.RectColliderComponent = option.Some(world.RectColliderComponent{
				Width:  float64(imageWidth),
				Height: float64(imageHeight),
			})
		}
	}
	return errors.Join(errs...)
}

func getCollidingEntities(entity *world.Entity, w *world.World) ([]*world.Entity, error) {
	entityAabb, err := aabbFromEntity(entity)
	if err != nil {
		return nil, err
	}

	var collisions []*world.Entity
	for other := range w.Entities {
		if entity == other {
			continue
		}

		if entity.Level != other.Level {
			continue
		}

		otherAabb, err := aabbFromEntity(other)
		if err != nil {
			continue
		}

		if entityAabb.collidesWith(otherAabb) {
			collisions = append(collisions, other)
		}
	}
	return collisions, nil

}

func preventCollisions(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.PreventCollisionsComponent.Present {
			continue
		}
		preventCollisionsComponent := &entity.PreventCollisionsComponent.Data

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
			continue
		}
		position := &entity.Position.Data

		collisions, err := getCollidingEntities(entity, w)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if len(collisions) != 0 {
			if !preventCollisionsComponent.LastLegalPosition.Present {
				continue
			}
			lastLegalPosition := preventCollisionsComponent.LastLegalPosition.Data
			if lastLegalPosition.Level != entity.Level {
				continue
			}

			*position = lastLegalPosition.Position
		} else {
			preventCollisionsComponent.LastLegalPosition = option.Some(world.Position{
				Level:    entity.Level,
				Position: *position,
			})
		}
	}
	return errors.Join(errs...)
}
