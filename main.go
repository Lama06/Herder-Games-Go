package main

import (
	"log"
	"math/rand"

	"github.com/Lama06/Herder-Games/assets"
	"github.com/Lama06/Herder-Games/graphics"
	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/system"
	"github.com/Lama06/Herder-Games/system/background"
	"github.com/Lama06/Herder-Games/system/check_collisions"
	"github.com/Lama06/Herder-Games/system/keyboard_controller"
	"github.com/Lama06/Herder-Games/system/prevent_collisions"
	"github.com/Lama06/Herder-Games/system/render_image"
	"github.com/Lama06/Herder-Games/system/render_rect"
	"github.com/Lama06/Herder-Games/system/systems"
	"github.com/Lama06/Herder-Games/world"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type Game struct {
	world  *world.World
	system system.System
}

func (g *Game) Update() error {
	err := g.system.Update(g.world)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.White)
	err := g.system.Draw(g.world, screen)
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
					Position: option.Some[world.WorldPosition](world.TilePosition{
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
				Position: option.Some[world.WorldPosition](world.TilePosition{
					TileX: x,
					TileY: y,
				}),
				ImageRenderer: option.Some(world.ImageRendererComponent{
					Image: ebiten.NewImageFromImage(assets.BodenImg),
				}),
				Renderer: option.Some(world.RendererComponent{
					Layer: BodenLayer,
				}),
			}
			w.Entities[boden] = struct{}{}

			if rand.Float64() <= 0.03 {
				tisch := &world.Entity{
					Position: option.Some[world.WorldPosition](world.TilePosition{
						TileX: x,
						TileY: y,
					}),
					ImageRenderer: option.Some(world.ImageRendererComponent{
						Image: ebiten.NewImageFromImage(assets.TischImg),
					}),
					Renderer: option.Some(world.RendererComponent{
						Layer: TischLayer,
					}),
					RectCollider: option.Some(world.RectColliderComponent{
						Width:  graphics.TileSize * 2,
						Height: graphics.TileSize * 2,
					}),
				}
				w.Entities[tisch] = struct{}{}
			}
		}
	}
}

func main() {
	player := &world.Entity{
		Position: option.Some[world.WorldPosition](world.TilePosition{
			TileX: 15,
			TileY: 15,
		}),
		RectRenderer: option.Some(world.RectRendererComponent{
			Width:  20,
			Height: 20,
			Color:  colornames.Blue,
		}),
		Renderer: option.Some(world.RendererComponent{
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
		},
	}

	addBoden(world)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.RunGame(&Game{
		world: world,
		system: systems.New(
			background.New(),

			keyboard_controller.New(),

			check_collisions.New(),
			prevent_collisions.New(),

			background.New(),
			render_image.New(BodenLayer),
			render_image.New(TischLayer),
			render_rect.New(PlayerLayer),
		),
	})
}
