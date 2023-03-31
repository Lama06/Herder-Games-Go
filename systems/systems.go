package systems

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type requireComponentError struct {
	entity        *world.Entity
	componentName string
}

func newRequireComponentError(entity *world.Entity, componentName string) error {
	return requireComponentError{
		entity:        entity,
		componentName: componentName,
	}
}

func (r requireComponentError) Error() string {
	return fmt.Sprintf("%s component not present on entity: %v", r.componentName, r.entity)
}

func Update(w *world.World) error {
	var errs []error

	err := addRectImages(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run add rect images system: %w", err))
	}

	err = addImageBoundsColliders(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run add image bounds colliders system: %w", err))
	}

	initialiseBlockedPathfindingTiles(w)

	clearVelocity(w)

	err = handleKeyboardInput(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run handle keyboard input system: %w", err))
	}

	err = pathfind(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run pathfind system: %w", err))
	}

	err = moveToCoordinates(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run move to coordinates system: %w", err))
	}

	err = moveToCoordinate(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run move to coordinate system: %w", err))
	}

	err = moveWithVelocity(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run move with velocity system: %w", err))
	}

	err = teleportEntitiesTouchingPortal(w)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run teleport entities touching portal system: %w", err))
	}

	return errors.Join(errs...)
}

func Draw(w *world.World, screen *ebiten.Image) error {
	var errs []error

	drawBackground(w, screen)

	err := drawImages(w, screen)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to run draw images system: %w", err))
	}

	return errors.Join(errs...)
}
