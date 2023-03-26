package keyboard_controller

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type offset struct{ x, y int }

var keyOffsets = map[ebiten.Key]offset{
	ebiten.KeyW: {x: 0, y: -1},
	ebiten.KeyA: {x: -1, y: 0},
	ebiten.KeyS: {x: 0, y: 1},
	ebiten.KeyD: {x: 1, y: 0},
}

type keyboardControllerSystem struct{}

func New() system.System {
	return keyboardControllerSystem{}
}

func (k keyboardControllerSystem) Update(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.KeyboardController.Present {
			continue
		}
		keyboardController := entity.KeyboardController.Data

		if !entity.Position.Present {
			errs = append(errs, fmt.Errorf("entity has keyboard controller component but no position component: %v", entity))
			continue
		}
		position := &entity.Position.Data

		for key, offset := range keyOffsets {
			if ebiten.IsKeyPressed(key) {
				*position = (*position).Coordinates().Add(offset.x*keyboardController.Speed, offset.y*keyboardController.Speed)
			}
		}
	}
	return errors.Join(errs...)
}

func (k keyboardControllerSystem) Draw(w *world.World, screen *ebiten.Image) error {
	return nil
}
