package systems

import (
	"errors"
	"log"
	"math"
	"sort"

	"github.com/Lama06/Herder-Games/astar"
	"github.com/Lama06/Herder-Games/world"
)

func floatSign(number float64) float64 {
	switch {
	case number > 0:
		return 1
	case number < 0:
		return -1
	case number == 0:
		return 0
	default:
		panic("unreachable")
	}
}

func moveToCoordinate(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.MoveToCoordinateComponent.Present {
			continue
		}
		moveToCoordinateComponent := &entity.MoveToCoordinateComponent.Data

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
			continue
		}
		position := entity.Position.Data.WorldCoordinates()

		if !entity.VelocityComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "velocity"))
			continue
		}
		velocity := &entity.VelocityComponent.Data

		if moveToCoordinateComponent.Coordinate == nil {
			continue
		}

		//log.Println(moveToCoordinateComponent.Coordinate)

		destinationX := moveToCoordinateComponent.Coordinate.WorldCoordinates().WorldX
		destinationY := moveToCoordinateComponent.Coordinate.WorldCoordinates().WorldY

		xDistance := destinationX - position.WorldX
		yDistance := destinationY - position.WorldY
		log.Println(xDistance, yDistance)

		if xDistance == 0 && yDistance == 0 {
			continue
		}

		absXDistance := math.Abs(xDistance)
		absYDistance := math.Abs(yDistance)

		var distanceToVelocityMultiplier float64
		if absXDistance > absYDistance {
			distanceToVelocityMultiplier = moveToCoordinateComponent.Speed / absXDistance
		} else {
			distanceToVelocityMultiplier = moveToCoordinateComponent.Speed / absYDistance
		}

		absXVelocity := absXDistance * distanceToVelocityMultiplier
		absYVelocity := absYDistance * distanceToVelocityMultiplier

		xVelocity := floatSign(xDistance) * absXVelocity
		yVelocity := floatSign(yDistance) * absYVelocity

		velocity.VelocityX += xVelocity
		velocity.VelocityY += yVelocity
	}
	return errors.Join(errs...)
}

func moveToCoordinates(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.MoveToCoordinatesComponent.Present {
			continue
		}
		moveToCoordinatesComponent := &entity.MoveToCoordinatesComponent.Data

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
			continue
		}
		position := entity.Position.Data.WorldCoordinates()

		if !entity.MoveToCoordinateComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "move to coordinate"))
			continue
		}
		moveToCoordinateComponent := &entity.MoveToCoordinateComponent.Data

		if len(moveToCoordinatesComponent.Coordinates) == 0 {
			moveToCoordinateComponent.Coordinate = nil
			continue
		}

		currentCoordinate := moveToCoordinatesComponent.Coordinates[moveToCoordinatesComponent.CurrentCoordinate]

		if currentCoordinate.WorldCoordinates() == position && moveToCoordinatesComponent.CurrentCoordinate+1 < len(moveToCoordinatesComponent.Coordinates) {
			moveToCoordinatesComponent.CurrentCoordinate++
			currentCoordinate = moveToCoordinatesComponent.Coordinates[moveToCoordinatesComponent.CurrentCoordinate]
		}

		moveToCoordinateComponent.Coordinate = currentCoordinate
	}
	return errors.Join(errs...)
}

func initialiseBlockedPathfindingTiles(w *world.World) {
	if w.BlockedPathfindingTiles != nil {
		return
	}
	w.BlockedPathfindingTiles = make(map[world.TilePosition]struct{})

	for entity := range w.Entities {
		if !entity.Static {
			continue
		}

		aabb, trigger, err := aabbFromEntity(entity)
		if err != nil {
			continue
		}
		if trigger {
			continue
		}

		for blockedTile := range aabb.blockedTiles() {
			blockedTilePosition := world.TilePosition{
				Level: entity.Level,
				Tile:  blockedTile,
			}
			w.BlockedPathfindingTiles[blockedTilePosition] = struct{}{}
		}
	}
}

func tilePositionNeighbours(w *world.World, position world.TilePosition) []astar.Neighbour[world.TilePosition] {
	var result []astar.Neighbour[world.TilePosition]
	for _, offset := range [...]struct{ x, y int }{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		neighbourTile := world.TilePosition{
			Level: position.Level,
			Tile: world.TileCoordinates{
				TileX: position.Tile.TileX + offset.x,
				TileY: position.Tile.TileY + offset.y,
			},
		}
		if _, isBlocked := w.BlockedPathfindingTiles[neighbourTile]; isBlocked {
			continue
		}
		result = append(result, astar.Neighbour[world.TilePosition]{
			Node: neighbourTile,
			Cost: 1,
		})
	}
	return result
}

func tilePositionEstimateCost(from world.TilePosition, to world.TilePosition) astar.Cost {
	xDiff := float64(from.Tile.TileX - to.Tile.TileX)
	yDiff := float64(from.Tile.TileX - to.Tile.TileX)
	return astar.Cost(math.Sqrt(math.Pow(xDiff, 2)+math.Pow(yDiff, 2)) * 10)
}

func tilePositionPathToCoordinatesPath(path astar.Path[world.TilePosition]) astar.Path[world.Coordinates] {
	result := make(astar.Path[world.Coordinates], len(path))
	for i := range path {
		result[i] = path[i].Tile
	}
	return result
}

func findShortestPath(w *world.World, from world.TilePosition, to world.TilePosition) astar.Path[world.Coordinates] {
	return tilePositionPathToCoordinatesPath(astar.FindPath(astar.Options[world.TilePosition]{
		Start: from,
		End:   to,
		NeighboursFunc: func(position world.TilePosition) []astar.Neighbour[world.TilePosition] {
			return tilePositionNeighbours(w, position)
		},
		EstimateCostFunc: tilePositionEstimateCost,
	}))
}

func findShortestPathToPortal(w *world.World, from world.TilePosition) (shortestPath astar.Path[world.Coordinates], portal *world.Entity) {
	type path struct {
		path   astar.Path[world.Coordinates]
		portal *world.Entity
	}

	var paths []path

	for portal := range w.Entities {
		if !portal.PortalComponent.Present {
			continue
		}

		if portal.Level != from.Level {
			continue
		}

		if !portal.Position.Present {
			continue
		}
		portalPosition := portal.Position.Data.WorldCoordinates()

		pathToPortal := findShortestPath(
			w,
			from,
			world.TilePosition{
				Level: portal.Level,
				Tile:  world.TileCoordinatesFromWorldCoordinates(portalPosition),
			},
		)
		if pathToPortal == nil {
			continue
		}

		paths = append(paths, path{
			path:   pathToPortal,
			portal: portal,
		})
	}

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].path) < len(paths[j].path)
	})

	return paths[0].path, paths[0].portal
}

func findShortestPathFromPortal(w *world.World, to world.TilePosition) (shortestPath astar.Path[world.Coordinates], portal *world.Entity) {
	type path struct {
		path   astar.Path[world.Coordinates]
		portal *world.Entity
	}

	var paths []path

	for portal := range w.Entities {
		if !portal.PortalComponent.Present {
			continue
		}

		if portal.Level != to.Level {
			continue
		}

		if !portal.Position.Present {
			continue
		}
		position := portal.Position.Data.WorldCoordinates()

		pathFromPortal := findShortestPath(
			w,
			world.TilePosition{
				Level: portal.Level,
				Tile:  world.TileCoordinatesFromWorldCoordinates(position),
			},
			to,
		)
		if pathFromPortal == nil {
			continue
		}

		paths = append(paths, path{
			path:   pathFromPortal,
			portal: portal,
		})
	}

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].path) < len(paths[j].path)
	})

	return paths[0].path, paths[0].portal
}

func pathfind(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.PathfinderComponent.Present {
			continue
		}
		pathfinderComponent := &entity.PathfinderComponent.Data

		if !entity.MoveToCoordinatesComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "move to coodinates"))
			continue
		}
		moveToCoodinatesComponent := &entity.MoveToCoordinatesComponent.Data

		if !entity.Position.Present {
			errs = append(errs, newRequireComponentError(entity, "position"))
			continue
		}
		position := &entity.Position.Data

		if !pathfinderComponent.Destination.Present {
			moveToCoodinatesComponent.SetCoordinates(nil)
			continue
		}
		destination := pathfinderComponent.Destination.Data

		switch pathfinderComponent.State {
		case world.PathfinderComponentStateNotStarted:
			if entity.Level == destination.Level {
				path := findShortestPath(
					w,
					world.TilePosition{
						Level: entity.Level,
						Tile:  world.TileCoordinatesFromWorldCoordinates((*position).WorldCoordinates()),
					},
					world.TilePosition{
						Level: destination.Level,
						Tile:  world.TileCoordinatesFromWorldCoordinates(destination.Position.WorldCoordinates()),
					},
				)

				moveToCoodinatesComponent.SetCoordinates(path)
				pathfinderComponent.State = world.PathfinderComponentStateToDestination
				continue
			}

			path, portal := findShortestPathToPortal(w, world.TilePosition{
				Level: entity.Level,
				Tile:  world.TileCoordinatesFromWorldCoordinates((*position).WorldCoordinates()),
			})
			moveToCoodinatesComponent.SetCoordinates(path)
			pathfinderComponent.State = world.PathfinderComponentStateToPortal
			pathfinderComponent.Portal = portal
			continue
		case world.PathfinderComponentStateToPortal:
			if !isCollision(entity, pathfinderComponent.Portal, true) {
				continue
			}

			path, portal := findShortestPathFromPortal(w, world.TilePosition{
				Level: destination.Level,
				Tile:  world.TileCoordinatesFromWorldCoordinates(destination.Position.WorldCoordinates()),
			})
			*position = portal.Position.Data
			pathfinderComponent.State = world.PathfinderComponentStateToDestination
			moveToCoodinatesComponent.SetCoordinates(path)
		case world.PathfinderComponentStateToDestination:
			if (*position).WorldCoordinates() != destination.Position.WorldCoordinates() {
				continue
			}

			pathfinderComponent.State = world.PathfinderComponentStateFinished
		}
	}
	return errors.Join(errs...)
}
