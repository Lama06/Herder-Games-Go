package prevent_collisions

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type preventCollisionsSystem struct {
	lastLegalPosition map[*world.Entity]world.WorldCoordinates
}

func New() system.System {
	return preventCollisionsSystem{
		lastLegalPosition: make(map[*world.Entity]world.WorldCoordinates),
	}
}

func (p preventCollisionsSystem) cleanupLastLegalPositions(w *world.World) {
	for entity := range p.lastLegalPosition {
		_, alive := w.Entities[entity]
		if !alive {
			delete(p.lastLegalPosition, entity)
		}
	}
}

func (p preventCollisionsSystem) Update(w *world.World) error {
	p.cleanupLastLegalPositions(w)

	var errs []error
	for entity := range w.Entities {
		if !entity.PreventCollisions.Present {
			continue
		}

		if !entity.Position.Present {
			errs = append(errs, fmt.Errorf("entity has prevent collisions component but no position component"))
			continue
		}
		position := &entity.Position.Data

		if !entity.Collissions.Present {
			errs = append(errs, fmt.Errorf("collisions component required: %v", entity))
			continue
		}
		collisions := entity.Collissions.Data

		if len(collisions.Collisions) != 0 {
			lastLegalPosition, hasLastLegalPosition := p.lastLegalPosition[entity]
			if !hasLastLegalPosition {
				continue
			}
			*position = lastLegalPosition
		} else {
			p.lastLegalPosition[entity] = (*position).Coordinates()
		}
	}
	return errors.Join(errs...)
}

func (p preventCollisionsSystem) Draw(w *world.World, screen *ebiten.Image) error {
	return nil
}
