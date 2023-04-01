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
	start bool
	world *world.World
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.start = true
	}

	if !g.start {
		return nil
	}

	err := systems.Update(g.world)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if !g.start {
		return
	}

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
					Coordinate: option.Some[world.Coordinate](world.TileCoordinate{
						TileX: x,
						TileY: y,
					}),
					Static: true,
					RectColliderComponent: option.Some(world.RectColliderComponent{
						Width:  world.TileSize,
						Height: world.TileSize,
					}),
				}
				w.Entities[border] = struct{}{}
				continue
			}

			boden := &world.Entity{
				Level: level,
				Coordinate: option.Some[world.Coordinate](world.TileCoordinate{
					TileX: x,
					TileY: y,
				}),
				ImageComponent: option.Some(world.ImageComponent{
					Layer: BodenLayer,
					Image: ebiten.NewImageFromImage(assets.BodenImg),
				}),
			}
			w.Entities[boden] = struct{}{}

			if rand.Float64() <= 0.05 {
				box := &world.Entity{
					Level: level,
					Coordinate: option.Some[world.Coordinate](world.TileCoordinate{
						TileX: x,
						TileY: y,
					}),
					Static: true,
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
					Level:  level,
					Static: true,
					Coordinate: option.Some[world.Coordinate](world.TileCoordinate{
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
							Coordinate: world.TileCoordinate{
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
		Coordinate: option.Some[world.Coordinate](world.TileCoordinate{
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
		VelocityComponent: option.Some(world.VelocityComponent{}),
		//KeyboardControllerComponent:  option.Some(world.KeyboardControllerComponent{Speed: 2}),
		MoveSpeedComponent: option.Some(world.MoveSpeedComponent{
			Speed: 1,
		}),
		MoveToCoordinateComponent:    option.Some(world.MoveToCoordinateComponent{}),
		MoveToCoordinatesComponent:   option.Some(world.MoveToCoordinatesComponent{}),
		PathfinderComponent:          option.Some(world.PathfinderComponent{}),
		RectColliderComponent:        option.Some(world.RectColliderComponent{}),
		ImageBoundsColliderComponent: option.Some(world.ImageBoundsColliderComponent{}),
	}

	player.PathfinderComponent.Data.SetDestination(world.Position{
		Level: 1,
		Coordinate: world.TileCoordinate{
			TileX: 25,
			TileY: 25,
		},
	})

	world := &world.World{
		Player: player,
		Entities: map[*world.Entity]struct{}{
			player: {},
		},
	}

	addBoden(world, 0)
	addBoden(world, 1)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(1000, 500)
	ebiten.RunGame(&Game{
		start: true,
		world: world,
	})
}
