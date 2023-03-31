package systems

import (
	"errors"
	"sort"

	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func addRectImages(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.RectComponent.Present {
			continue
		}
		rectComponent := entity.RectComponent.Data

		if !entity.ImageComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "image"))
			continue
		}
		imageComponent := &entity.ImageComponent.Data

		image := ebiten.NewImage(rectComponent.Width, rectComponent.Height)
		image.Fill(rectComponent.Color)

		imageComponent.Image = image
	}
	return errors.Join(errs...)
}

func drawBackground(w *world.World, screen *ebiten.Image) {
	for entity := range w.Entities {
		if !entity.BackgroundComponent.Present {
			continue
		}
		backgroundComponent := entity.BackgroundComponent.Data

		if entity.Level != w.Player.Level {
			continue
		}

		screen.Fill(backgroundComponent.Color)
	}
}

func worldPositionToScreenPosition(w *world.World, position world.Coordinates) (x, y float64) {
	coordinates := position.Coordinates()
	playerCoordinates := w.Player.Position.Data.Coordinates()
	playerXOffset := coordinates.WorldX - playerCoordinates.WorldX
	playerYOffset := coordinates.WorldY - playerCoordinates.WorldY
	return graphics.PlayerX + playerXOffset, graphics.PlayerY + playerYOffset
}

func findAllImageLayers(w *world.World) []int {
	layersSet := make(map[int]struct{})
	for entity := range w.Entities {
		if !entity.ImageComponent.Present {
			continue
		}
		imageComponent := entity.ImageComponent.Data
		layersSet[imageComponent.Layer] = struct{}{}
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
	for _, layer := range findAllImageLayers(w) {
		for entity := range w.Entities {
			if entity.Level != w.Player.Level {
				continue
			}

			if !entity.ImageComponent.Present {
				continue
			}
			imageComponent := entity.ImageComponent.Data

			if imageComponent.Layer != layer {
				continue
			}

			if !entity.Position.Present {
				errs = append(errs, newRequireComponentError(entity, "position"))
				continue
			}
			position := entity.Position.Data

			screenPositionX, screenPositionY := worldPositionToScreenPosition(w, position)

			var drawOptions ebiten.DrawImageOptions
			drawOptions.GeoM.Translate(screenPositionX, screenPositionY)
			screen.DrawImage(imageComponent.Image, &drawOptions)
		}
	}
	return errors.Join(errs...)
}
