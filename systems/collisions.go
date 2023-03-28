package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/world"
)

func addImageBoundsColliders(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.ImageBoundsCollider.Present {
			continue
		}

		if !entity.Image.Present {
			errs = append(errs, newRequireComponentError(entity, "image"))
			continue
		}
		image := entity.Image.Data

		imageWidth := image.Image.Bounds().Dx()
		imageHeight := image.Image.Bounds().Dy()

		if !entity.RectCollider.Present {
			entity.RectCollider = option.Some(world.RectColliderComponent{
				Width:  imageWidth,
				Height: imageHeight,
			})
		}
	}
	return errors.Join(errs...)
}

func addCollisionsComponents(w *world.World) {
	for entity := range w.Entities {
		if !entity.PreventCollisions.Present {
			continue
		}

		entity.Collissions = option.Some(world.CollisionsComponent{})
	}
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
