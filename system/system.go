package system

import (
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
)

type System interface {
	Update(w *world.World) error
	Draw(w *world.World, screen *ebiten.Image) error
}
