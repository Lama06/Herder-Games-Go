package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/world"
)

func addImageBoundsColliders(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.ImageBoundsColliderComponent.Present {
			continue
		}

		if !entity.RectColliderComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "rect collider"))
			continue
		}
		rectCollider := &entity.RectColliderComponent.Data

		if !entity.ImageComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "image"))
			continue
		}
		imageComponent := entity.ImageComponent.Data

		imageWidth := imageComponent.Image.Bounds().Dx()
		imageHeight := imageComponent.Image.Bounds().Dy()

		rectCollider.Width = float64(imageWidth)
		rectCollider.Height = float64(imageHeight)
	}
	return errors.Join(errs...)
}

func getCollidingEntities(w *world.World, entity *world.Entity, triggerCollisions bool) ([]*world.Entity, error) {
	entityAabb, entityTrigger, err := aabbFromEntity(entity)
	if err != nil {
		return nil, err
	}

	if entityTrigger && !triggerCollisions {
		return nil, nil
	}

	var collisions []*world.Entity
	for other := range w.Entities {
		if entity == other {
			continue
		}

		if entity.Level != other.Level {
			continue
		}

		otherAabb, otherTrigger, err := aabbFromEntity(other)
		if err != nil {
			continue
		}

		if otherTrigger && !triggerCollisions {
			continue
		}

		if entityAabb.collidesWith(otherAabb) {
			collisions = append(collisions, other)
		}
	}
	return collisions, nil
}
