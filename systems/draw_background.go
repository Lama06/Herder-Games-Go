package systems

import (
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

func drawBackground(w *world.World, screen *ebiten.Image) {
	for entity := range w.Entities {
		if !entity.Background.Present {
			continue
		}
		background := entity.Background.Data

		screen.Fill(background.Color)
	}
}
