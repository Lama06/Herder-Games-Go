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

func addBoden(w *world.World, level world.Level) {
	background := &world.Entity{
		Level: level,
		BackgroundComponent: option.Some(world.BackgroundComponent{
			Color: colornames.White,
		}),
	}
	w.Entities[background] = struct{}{}

	for x := 0; x < 30; x++ {
		for y := 0; y < 30; y++ {
			if x == 0 || x == 29 || y == 0 || y == 29 {
				border := &world.Entity{
					Level: level,
					Position: option.Some[world.Coordinates](world.TileCoordinates{
						TileX: x,
						TileY: y,
					}),
					RectColliderComponent: option.Some(world.RectColliderComponent{
						Width:  graphics.TileSize,
						Height: graphics.TileSize,
					}),
				}
				w.Entities[border] = struct{}{}
				continue
			}

			boden := &world.Entity{
				Level: level,
				Position: option.Some[world.Coordinates](world.TileCoordinates{
					TileX: x,
					TileY: y,
				}),
				ImageComponent: option.Some(world.ImageComponent{
					Layer: BodenLayer,
					Image: ebiten.NewImageFromImage(assets.BodenImg),
				}),
			}
			w.Entities[boden] = struct{}{}

			if rand.Float64() <= 0.02 {
				box := &world.Entity{
					Level: level,
					Position: option.Some[world.Coordinates](world.TileCoordinates{
						TileX: x,
						TileY: y,
					}),
					ImageComponent: option.Some(world.ImageComponent{
						Layer: TischLayer,
					}),
					RectComponent: option.Some(world.RectComponent{
						Width:  20,
						Height: 20,
						Color:  colornames.Red,
					}),
					RectColliderComponent:        option.Some(world.RectColliderComponent{}),
					ImageBoundsColliderComponent: option.Some(world.ImageBoundsColliderComponent{}),
				}
				w.Entities[box] = struct{}{}
			}

			if rand.Float64() <= 0.01 {
				destinationLevel := world.Level(1)
				if level == 1 {
					destinationLevel = 0
				}
				tisch := &world.Entity{
					Level: level,
					Position: option.Some[world.Coordinates](world.TileCoordinates{
						TileX: x,
						TileY: y,
					}),
					ImageComponent: option.Some(world.ImageComponent{
						Image: ebiten.NewImageFromImage(assets.TischImg),
						Layer: TischLayer,
					}),
					RectColliderComponent: option.Some(world.RectColliderComponent{
						Trigger: true,
					}),
					ImageBoundsColliderComponent: option.Some(world.ImageBoundsColliderComponent{}),
					PortalComponent: option.Some(world.PortalComponent{
						Destination: world.Position{
							Level: destinationLevel,
							Position: world.TileCoordinates{
								TileX: 15,
								TileY: 15,
							},
						},
					}),
				}
				w.Entities[tisch] = struct{}{}
			}
		}
	}
}

func main() {
	player := &world.Entity{
		Position: option.Some[world.Coordinates](world.TileCoordinates{
			TileX: 15,
			TileY: 15,
		}),
		ImageComponent: option.Some(world.ImageComponent{
			Layer: PlayerLayer,
		}),
		RectComponent: option.Some(world.RectComponent{
			Width:  20,
			Height: 20,
			Color:  colornames.Red,
		}),
		VelocityComponent:            option.Some(world.VelocityComponent{}),
		KeyboardControllerComponent:  option.Some(world.KeyboardControllerComponent{Speed: 2}),
		RectColliderComponent:        option.Some(world.RectColliderComponent{}),
		ImageBoundsColliderComponent: option.Some(world.ImageBoundsColliderComponent{}),
	}

	world := &world.World{
		Player: player,
		Entities: map[*world.Entity]struct{}{
			player: {},
		},
	}

	addBoden(world, 0)
	addBoden(world, 1)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.RunGame(&Game{
		world: world,
	})
}
