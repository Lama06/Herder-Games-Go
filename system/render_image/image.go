package render_image

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type imageRenderSystem struct {
	layer int
}

func New(layer int) system.System {
	return imageRenderSystem{
		layer: layer,
	}
}

func (i imageRenderSystem) Update(w *world.World) error {
	return nil
}

func (i imageRenderSystem) Draw(w *world.World, screen *ebiten.Image) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.ImageRenderer.Present {
			continue
		}
		imageRenderer := entity.ImageRenderer.Data

		if !entity.Position.Present {
			errs = append(errs, fmt.Errorf("position component required: %v", entity))
			continue
		}
		worldPosition := entity.Position.Data
		screenPosition := world.WorldPositionToScreenPosition(w, worldPosition)

		if !entity.Renderer.Present {
			errs = append(errs, fmt.Errorf("renderer component required: %v", entity))
			continue
		}
		renderer := entity.Renderer.Data

		if renderer.Layer != i.layer {
			continue
		}

		var drawOptions ebiten.DrawImageOptions
		drawOptions.GeoM.Translate(float64(screenPosition.ScreenX), float64(screenPosition.ScreenY))
		screen.DrawImage(imageRenderer.Image, &drawOptions)
	}
	return errors.Join(errs...)
}
