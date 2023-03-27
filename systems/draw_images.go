package systems

import (
	"errors"
	"sort"

	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func worldPositionToScreenPosition(w *world.World, position world.Position) (x, y int) {
	worldCoordinates := position.Coordinates()
	playerWorldCoordinates := w.Player.Position.Data.Coordinates()
	playerXOffset := worldCoordinates.WorldX - playerWorldCoordinates.WorldX
	playerYOffset := worldCoordinates.WorldY - playerWorldCoordinates.WorldY
	return graphics.PlayerX + playerXOffset, graphics.PlayerY + playerYOffset
}

func findAllLayers(w *world.World) []int {
	layersSet := make(map[int]struct{})
	for entity := range w.Entities {
		if !entity.Image.Present {
			continue
		}
		image := entity.Image.Data
		layersSet[image.Layer] = struct{}{}
	}

	layers := make([]int, 0, len(layersSet))
	for layer := range layersSet {
		layers = append(layers, layer)
	}

	sort.Ints(layers)

	return layers
}

func drawImages(w *world.World, screen *ebiten.Image) error {
	var errs []error
	for _, layer := range findAllLayers(w) {
		for entity := range w.Entities {
			if !entity.Image.Present {
				continue
			}
			image := entity.Image.Data

			if !entity.Position.Present {
				errs = append(errs, newRequireComponentError(entity, "position"))
				continue
			}
			worldPosition := entity.Position.Data

			if image.Layer != layer {
				continue
			}

			screenPositionX, screenPositionY := worldPositionToScreenPosition(w, worldPosition)

			var drawOptions ebiten.DrawImageOptions
			if entity == w.Player {
				println(screenPositionX)
			}
			drawOptions.GeoM.Translate(float64(screenPositionX), float64(screenPositionY))
			screen.DrawImage(image.Image, &drawOptions)
		}
	}
	return errors.Join(errs...)
}
