package systems

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type systems []system.System

func New(children ...system.System) system.System {
	return systems(children)
}

func (s systems) Update(w *world.World) error {
	var errs []error
	for _, system := range s {
		err := system.Update(w)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to update system: %w", err))
		}
	}
	return errors.Join(errs...)
}

func (s systems) Draw(w *world.World, screen *ebiten.Image) error {
	var errs []error
	for _, system := range s {
		err := system.Draw(w, screen)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to draw system: %w", err))
		}
	}
	return errors.Join(errs...)
}
