package systems

import (
	"errors"
	"math"

	"github.com/Lama06/Herder-Games/astar"
	"github.com/Lama06/Herder-Games/world"
)

func moveToCoordinate(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.MoveToCoordinateComponent.Present {
			continue
		}
		moveToCoordinateComponent := &entity.MoveToCoordinateComponent.Data

		if !entity.Coordinate.Present {
			errs = append(errs, newRequireComponentError(entity, "coordinate"))
			continue
		}
		coordinate := entity.Coordinate.Data.WorldCoordinate()

		if !entity.VelocityComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "velocity"))
			continue
		}
		velocity := &entity.VelocityComponent.Data

		if !entity.MoveSpeedComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "move speed"))
			continue
		}
		moveSpeedComponent := entity.MoveSpeedComponent.Data

		switch moveToCoordinateComponent.State {
		case world.MoveToCoordinateComponentStateDisabled:
			continue
		case world.MoveToCoordinateComponentStateMoving:
			destination := moveToCoordinateComponent.Coordinate.WorldCoordinate()

			xDistance := destination.WorldX - coordinate.WorldX
			yDistance := destination.WorldY - coordinate.WorldY

			if xDistance == 0 && yDistance == 0 {
				moveToCoordinateComponent.State = world.MoveToCoordinateComponentStateArrived
				continue
			}

			absXDistance := math.Abs(xDistance)
			absYDistance := math.Abs(yDistance)

			var distanceToVelocityMultiplier float64
			if absXDistance > absYDistance {
				distanceToVelocityMultiplier = moveSpeedComponent.Speed / absXDistance
			} else {
				distanceToVelocityMultiplier = moveSpeedComponent.Speed / absYDistance
			}

			absXVelocity := absXDistance * distanceToVelocityMultiplier
			absYVelocity := absYDistance * distanceToVelocityMultiplier

			xVelocity := absXVelocity
			if xDistance < 0 {
				xVelocity *= -1
			}
			yVelocity := absYVelocity
			if yDistance < 0 {
				yVelocity *= -1
			}

			velocity.VelocityX += xVelocity
			velocity.VelocityY += yVelocity
		case world.MoveToCoordinateComponentStateArrived:
			continue
		}
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

		if !entity.MoveToCoordinateComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "move to coordinate"))
			continue
		}
		moveToCoordinateComponent := &entity.MoveToCoordinateComponent.Data

		switch moveToCoordinatesComponent.State {
		case world.MoveToCoordinatesComponentStateDisabled:
			moveToCoordinateComponent.Disable()
		case world.MoveToCoordinatesComponentStateIdle:
			moveToCoordinateComponent.SetCoordinate(moveToCoordinatesComponent.Coordinates[0])
			moveToCoordinatesComponent.State = world.MoveToCoordinatesComponentStateMoving
		case world.MoveToCoordinatesComponentStateMoving:
			if moveToCoordinateComponent.State != world.MoveToCoordinateComponentStateArrived {
				continue
			}

			if moveToCoordinatesComponent.CurrentCoordinate+1 > len(moveToCoordinatesComponent.Coordinates)-1 {
				moveToCoordinateComponent.Disable()
				moveToCoordinatesComponent.State = world.MoveToCoordinatesComponentStateFinished
				continue
			}

			moveToCoordinatesComponent.CurrentCoordinate++
			currentCoordinate := moveToCoordinatesComponent.Coordinates[moveToCoordinatesComponent.CurrentCoordinate]
			moveToCoordinateComponent.SetCoordinate(currentCoordinate)
		case world.MoveToCoordinatesComponentStateFinished:
			moveToCoordinateComponent.Disable()
		}
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
				Level:          entity.Level,
				TileCoordinate: blockedTile,
			}
			w.BlockedPathfindingTiles[blockedTilePosition] = struct{}{}
		}
	}
}

func tilePositionNeighbours(w *world.World, tilePosition world.TilePosition) []astar.Neighbour[world.TilePosition] {
	var result []astar.Neighbour[world.TilePosition]
	for _, neighbourOffset := range [...]struct{ x, y int }{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		neighbourTilePosition := world.TilePosition{
			Level: tilePosition.Level,
			TileCoordinate: world.TileCoordinate{
				TileX: tilePosition.TileCoordinate.TileX + neighbourOffset.x,
				TileY: tilePosition.TileCoordinate.TileY + neighbourOffset.y,
			},
		}
		if _, isBlocked := w.BlockedPathfindingTiles[neighbourTilePosition]; isBlocked {
			continue
		}
		result = append(result, astar.Neighbour[world.TilePosition]{
			Node: neighbourTilePosition,
			Cost: 10,
		})
	}
	return result
}

func tilePositionEstimateCost(from world.TilePosition, to world.TilePosition) astar.Cost {
	xDiff := float64(from.TileCoordinate.TileX - to.TileCoordinate.TileX)
	yDiff := float64(from.TileCoordinate.TileX - to.TileCoordinate.TileX)
	return astar.Cost(math.Sqrt(math.Pow(xDiff, 2)+math.Pow(yDiff, 2)) * 10)
}

func tilePositionPathToCoordinatePath(path astar.Path[world.TilePosition]) astar.Path[world.Coordinate] {
	result := make(astar.Path[world.Coordinate], len(path))
	for i := range path {
		result[i] = path[i].TileCoordinate
	}
	return result
}

func findShortestPath(w *world.World, from world.TilePosition, to world.TilePosition) astar.Path[world.Coordinate] {
	path := astar.FindPath(astar.Options[world.TilePosition]{
		Start: from,
		End:   to,
		NeighboursFunc: func(position world.TilePosition) []astar.Neighbour[world.TilePosition] {
			return tilePositionNeighbours(w, position)
		},
		EstimateCostFunc: tilePositionEstimateCost,
	})

	if path == nil {
		return nil
	}

	return tilePositionPathToCoordinatePath(path)
}

func getShortestPath(paths []astar.Path[world.Coordinate]) astar.Path[world.Coordinate] {
	if len(paths) == 0 {
		return nil
	}

	var shortestPath astar.Path[world.Coordinate]
	for i, path := range paths {
		if i == 0 {
			shortestPath = path
			continue
		}

		if len(path) < len(shortestPath) {
			shortestPath = path
		}
	}
	return shortestPath
}

func findShortestPathToPortal(w *world.World, from world.TilePosition) astar.Path[world.Coordinate] {
	var paths []astar.Path[world.Coordinate]

	for portal := range w.Entities {
		if !portal.PortalComponent.Present {
			continue
		}

		if portal.Level != from.Level {
			continue
		}

		if !portal.Coordinate.Present {
			continue
		}
		portalCoordinate := portal.Coordinate.Data.WorldCoordinate()

		pathToPortal := findShortestPath(
			w,
			from,
			world.TilePosition{
				Level:          portal.Level,
				TileCoordinate: world.TileCoordinateFromCoordinate(portalCoordinate),
			},
		)
		if pathToPortal == nil {
			continue
		}

		paths = append(paths, pathToPortal)
	}

	return getShortestPath(paths)
}

func findShortestPathFromPortal(w *world.World, to world.TilePosition) astar.Path[world.Coordinate] {
	var paths []astar.Path[world.Coordinate]

	for portal := range w.Entities {
		if !portal.PortalComponent.Present {
			continue
		}

		if portal.Level != to.Level {
			continue
		}

		if !portal.Coordinate.Present {
			continue
		}
		portalCoordinate := portal.Coordinate.Data.WorldCoordinate()

		pathFromPortal := findShortestPath(
			w,
			world.TilePosition{
				Level:          portal.Level,
				TileCoordinate: world.TileCoordinateFromCoordinate(portalCoordinate),
			},
			to,
		)
		if pathFromPortal == nil {
			continue
		}

		paths = append(paths, pathFromPortal)
	}

	return getShortestPath(paths)
}

func pathfind(w *world.World) error {
	var errs []error
	for entity := range w.Entities {
		if !entity.PathfinderComponent.Present {
			continue
		}
		pathfinderComponent := &entity.PathfinderComponent.Data

		if !entity.Coordinate.Present {
			errs = append(errs, newRequireComponentError(entity, "coordinate"))
			continue
		}
		coordinate := &entity.Coordinate.Data

		position := world.Position{
			Level:      entity.Level,
			Coordinate: (*coordinate).WorldCoordinate(),
		}

		if !entity.MoveToCoordinatesComponent.Present {
			errs = append(errs, newRequireComponentError(entity, "move to coodinates"))
			continue
		}
		moveToCoodinatesComponent := &entity.MoveToCoordinatesComponent.Data

		switch pathfinderComponent.State {
		case world.PathfinderComponentStateDisabled, world.PathfinderComponentStateNoPath, world.PathfinderComponentStateArrived:
			moveToCoodinatesComponent.Disable()
		case world.PathfinderComponentStateIdle:
			if entity.Level == pathfinderComponent.Destination.Level {
				path := findShortestPath(
					w,
					world.TilePositionFromPosition(position),
					world.TilePositionFromPosition(pathfinderComponent.Destination),
				)
				moveToCoodinatesComponent.SetCoordinates(path)
				pathfinderComponent.State = world.PathfinderComponentStateMovingToDestination
				continue
			}

			pathToPortal := findShortestPathToPortal(w, world.TilePositionFromPosition(position))
			if pathToPortal == nil {
				moveToCoodinatesComponent.Disable()
				pathfinderComponent.State = world.PathfinderComponentStateNoPath
				continue
			}

			moveToCoodinatesComponent.SetCoordinates(pathToPortal)
			pathfinderComponent.State = world.PathfinderComponentStateMovingToPortal
			continue
		case world.PathfinderComponentStateMovingToPortal:
			if moveToCoodinatesComponent.State != world.MoveToCoordinatesComponentStateFinished {
				continue
			}

			pathFromPortal := findShortestPathFromPortal(w, world.TilePositionFromPosition(pathfinderComponent.Destination))
			if pathFromPortal == nil {
				pathfinderComponent.State = world.PathfinderComponentStateNoPath
				moveToCoodinatesComponent.Disable()
				continue
			}

			entity.Level = pathfinderComponent.Destination.Level
			*coordinate = pathFromPortal[0]

			moveToCoodinatesComponent.SetCoordinates(pathFromPortal)
			pathfinderComponent.State = world.PathfinderComponentStateMovingToDestination
		case world.PathfinderComponentStateMovingToDestination:
			if moveToCoodinatesComponent.State != world.MoveToCoordinatesComponentStateFinished {
				continue
			}

			moveToCoodinatesComponent.Disable()
			pathfinderComponent.State = world.PathfinderComponentStateArrived
		}
	}
	return errors.Join(errs...)
}
