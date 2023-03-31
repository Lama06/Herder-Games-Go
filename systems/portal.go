package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/world"
)

func teleportEntitiesTouchingPortal(w *world.World) error {
	var errs []error
	for portal := range w.Entities {
		if !portal.PortalComponent.Present {
			continue
		}
		portalComponent := portal.PortalComponent.Data

		collisions, err := getCollidingEntities(w, portal, true)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, collision := range collisions {
			if collision.PathfinderComponent.Present {
				continue
			}

			if !collision.Position.Present {
				continue
			}
			collision.Level = portalComponent.Destination.Level
			collision.Position.Data = portalComponent.Destination.Position
		}
	}
	return errors.Join(errs...)
}
