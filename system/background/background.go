package background

import (
	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type backgroundSystem struct{}

func New() system.System {
	return backgroundSystem{}
}

func (backgroundSystem) Update(w *world.World) error {
	return nil
}

func (b backgroundSystem) Draw(w *world.World, screen *ebiten.Image) error {
	for entity := range w.Entities {
		if !entity.BackgroundRenderer.Present {
			continue
		}
		backgroundRenderer := entity.BackgroundRenderer.Data

		screen.Fill(backgroundRenderer.Color)
	}

	return nil
}
