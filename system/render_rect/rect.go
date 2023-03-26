package render_rect

import (
	"errors"
	"fmt"

	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type renderRectSystem struct {
	layer  int
	pixels graphics.PixelBuffer
}

func New(layer int) system.System {
	return &renderRectSystem{
		layer: layer,
	}
}

func (*renderRectSystem) Update(w *world.World) error {
	return nil
}

func (r *renderRectSystem) Draw(w *world.World, screen *ebiten.Image) error {
	var errs []error

	screen.ReadPixels(r.pixels[:])

	for entity := range w.Entities {
		if !entity.RectRenderer.Present {
			continue
		}
		rectRenderer := entity.RectRenderer.Data

		if !entity.Position.Present {
			errs = append(errs, fmt.Errorf("entity has RectRenderer component but no position: %v", entity))
			continue
		}
		position := entity.Position.Data
		screenPosition := world.WorldPositionToScreenPosition(w, position)

		if !entity.Renderer.Present {
			errs = append(errs, fmt.Errorf("renderer component required: %v", entity))
			continue
		}
		renderer := entity.Renderer.Data

		if renderer.Layer != r.layer {
			continue
		}

		r.pixels.DrawRect(screenPosition.ScreenX, screenPosition.ScreenY, rectRenderer.Width, rectRenderer.Height, rectRenderer.Color)
	}

	screen.WritePixels(r.pixels[:])

	return errors.Join(errs...)
}
