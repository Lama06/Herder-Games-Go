package systems

import (
	"errors"
	"sort"

	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

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

func addRectImages(w *world.World) {
	for entity := range w.Entities {
		if !entity.RectComponent.Present {
			continue
		}
		rectComponent := entity.RectComponent.Data

		if !entity.ImageComponent.Present {
			image := ebiten.NewImage(rectComponent.Width, rectComponent.Height)
			image.Fill(rectComponent.Color)

			entity.ImageComponent = option.Some(world.ImageComponent{
				Image: image,
				Layer: rectComponent.Layer,
			})
		}
	}
}

func worldPositionToScreenPosition(w *world.World, position world.Coordinates) (x, y int) {
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
			drawOptions.GeoM.Translate(float64(screenPositionX), float64(screenPositionY))
			screen.DrawImage(imageComponent.Image, &drawOptions)
		}
	}
	return errors.Join(errs...)
}
