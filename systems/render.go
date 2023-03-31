package systems

import (
	"errors"
	"sort"

	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/option"
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
		entity.RectComponent = option.None[world.RectComponent]()
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

func worldCoordinateToScreenPosition(w *world.World, coordinate world.Coordinate) (x, y float64) {
	worldCoordinate := coordinate.WorldCoordinate()
	playerCoordinates := w.Player.Coordinate.Data.WorldCoordinate()
	playerXOffset := worldCoordinate.WorldX - playerCoordinates.WorldX
	playerYOffset := worldCoordinate.WorldY - playerCoordinates.WorldY
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

			if !entity.Coordinate.Present {
				errs = append(errs, newRequireComponentError(entity, "coordinate"))
				continue
			}
			coordinate := entity.Coordinate.Data

			screenPositionX, screenPositionY := worldCoordinateToScreenPosition(w, coordinate)

			var drawOptions ebiten.DrawImageOptions
			drawOptions.GeoM.Translate(screenPositionX, screenPositionY)
			screen.DrawImage(imageComponent.Image, &drawOptions)
		}
	}
	return errors.Join(errs...)
}
