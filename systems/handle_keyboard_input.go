package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

var keyOffsets = map[ebiten.Key]struct{ x, y int }{
	ebiten.KeyW: {x: 0, y: -1},
	ebiten.KeyA: {x: -1, y: 0},
	ebiten.KeyS: {x: 0, y: 1},
	ebiten.KeyD: {x: 1, y: 0},
}

func handleKeyboardInput(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.KeyboardController.Present {
			continue
		}
		keyboardController := entity.KeyboardController.Data

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
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
