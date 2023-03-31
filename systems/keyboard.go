package systems

import (
	"errors"

	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

var keyOffsets = map[ebiten.Key]struct{ x, y float64 }{
	ebiten.KeyW: {x: 0, y: -1},
	ebiten.KeyA: {x: -1, y: 0},
	ebiten.KeyS: {x: 0, y: 1},
	ebiten.KeyD: {x: 1, y: 0},
}

func handleKeyboardInput(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.KeyboardControllerComponent.Present {
			continue
		}
		keyboardControllerComponent := entity.KeyboardControllerComponent.Data

		if !entity.VelocityComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "velocity"))
			continue
		}
		velocity := &entity.VelocityComponent.Data

		for key, offset := range keyOffsets {
			if ebiten.IsKeyPressed(key) {
				velocity.VelocityX += offset.x * keyboardControllerComponent.Speed
				velocity.VelocityY += offset.y * keyboardControllerComponent.Speed
			}
		}
	}
	return errors.Join(errs...)
}
