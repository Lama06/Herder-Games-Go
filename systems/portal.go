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

		collisions, err := getCollidingEntities(portal, w)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, collision := range collisions {
			if !collision.Position.Present {
				continue
			}
			position := &collision.Position.Data

			collision.Level = portalComponent.Destination.Level
			*position = portalComponent.Destination.Position
		}
	}
	return errors.Join(errs...)
}
