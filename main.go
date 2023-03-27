package main

import (
	"log"
	"math/rand"

	"github.com/Lama06/Herder-Games/assets"
	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/systems"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type Game struct {
	world *world.World
}

func (g *Game) Update() error {
	err := systems.Update(g.world)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	err := systems.Draw(g.world, screen)
	if err != nil {
		log.Println(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (width, height int) {
	return graphics.WindowWidth, graphics.WindowHeight
}

const (
	BodenLayer = iota
	TischLayer
	PlayerLayer
)

func addBoden(w *world.World) {
	for x := 0; x < 30; x++ {
		for y := 0; y < 30; y++ {
			if x == 0 || x == 29 || y == 0 || y == 29 {
				border := &world.Entity{
					Position: option.Some[world.Position](world.TilePosition{
						TileX: x,
						TileY: y,
					}),
					RectCollider: option.Some(world.RectColliderComponent{
						Width:  graphics.TileSize,
						Height: graphics.TileSize,
					}),
				}
				w.Entities[border] = struct{}{}
				continue
			}

			boden := &world.Entity{
				Position: option.Some[world.Position](world.TilePosition{
					TileX: x,
					TileY: y,
				}),
				Image: option.Some(world.ImageComponent{
					Layer: BodenLayer,
					Image: ebiten.NewImageFromImage(assets.BodenImg),
				}),
			}
			w.Entities[boden] = struct{}{}

			if rand.Float64() <= 0.03 {
				tisch := &world.Entity{
					Position: option.Some[world.Position](world.TilePosition{
						TileX: x,
						TileY: y,
					}),
					Image: option.Some(world.ImageComponent{
						Image: ebiten.NewImageFromImage(assets.TischImg),
						Layer: TischLayer,
					}),
					ImageBoundsCollider: option.Some(world.ImageBoundsColliderComponent{}),
				}
				w.Entities[tisch] = struct{}{}
			}
		}
	}
}

func main() {
	playerImg := ebiten.NewImage(20, 20)
	playerImg.Fill(colornames.Red)

	player := &world.Entity{
		Position: option.Some[world.Position](world.TilePosition{
			TileX: 15,
			TileY: 15,
		}),
		Image: option.Some(world.ImageComponent{
			Image: playerImg,
			Layer: PlayerLayer,
		}),
		KeyboardController: option.Some(world.KeyboardControllerComponent{Speed: 2}),
		RectCollider: option.Some(world.RectColliderComponent{
			Width:  20,
			Height: 20,
		}),
		Collissions:       option.Some(world.CollisionsComponent{}),
		PreventCollisions: option.Some(world.PreventCollisionsComponent{}),
	}

	world := &world.World{
		Player: player,
		Entities: map[*world.Entity]struct{}{
			player: {},
			{
				Background: option.Some(world.BackgroundComponent{
					Color: colornames.White,
				}),
			}: {},
		},
	}

	addBoden(world)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.RunGame(&Game{
		world: world,
	})
}
