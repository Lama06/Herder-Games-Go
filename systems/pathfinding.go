package systems

import (
	"errors"
	"log"
	"math"

	"github.com/Lama06/Herder-Games/astar"
	"github.com/Lama06/Herder-Games/option"
	"github.com/Lama06/Herder-Games/world"
)

func compareFloatsWithTolerance(a, b, tolerance float64) bool {
	diff := math.Abs(a - b)
	return diff <= tolerance
}

func moveToCoordinate(w *world.World) error {
	const tolerance = 0.001

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
		coordinate := &entity.Coordinate.Data
		currentCoordiante := (*coordinate).WorldCoordinate()

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

		if moveToCoordinateComponent.Arrived {
			continue
		}

		destination := moveToCoordinateComponent.Coordinate.WorldCoordinate()

		xDistance := destination.WorldX - currentCoordiante.WorldX
		yDistance := destination.WorldY - currentCoordiante.WorldY

		if compareFloatsWithTolerance(xDistance, 0, tolerance) && compareFloatsWithTolerance(yDistance, 0, tolerance) {
			*coordinate = destination
			moveToCoordinateComponent.Arrived = true
			continue
		}

		var xVelocity, yVelocity float64
		if math.Abs(xDistance) > math.Abs(yDistance) {
			xSpeed := moveSpeedComponent.Speed
			if xSpeed > math.Abs(xDistance) {
				xSpeed = math.Abs(xDistance)
			}

			if xDistance > 0 {
				xVelocity = xSpeed
			} else {
				xVelocity = -xSpeed
			}

			xDistancePercentage := xDistance / xVelocity
			yVelocity = yDistance * xDistancePercentage
		} else {
			ySpeed := moveSpeedComponent.Speed
			if ySpeed > math.Abs(yDistance) {
				ySpeed = math.Abs(yDistance)
			}

			if yDistance > 0 {
				yVelocity = ySpeed
			} else {
				yVelocity = -ySpeed
			}

			yDistancePercentage := yDistance / yVelocity
			xVelocity = xDistance * yDistancePercentage
		}

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

		switch moveToCoordinatesComponent.State {
		case world.MoveToCoordinatesComponentStateIdle:
			entity.MoveToCoordinateComponent = option.Some(world.NewMoveToCoordinateComponent(moveToCoordinatesComponent.Coordinates[0]))
			moveToCoordinatesComponent.State = world.MoveToCoordinatesComponentStateMoving
		case world.MoveToCoordinatesComponentStateMoving:
			moveToCoordinateComponent := entity.MoveToCoordinateComponent.Data

			if !moveToCoordinateComponent.Arrived {
				continue
			}

			if moveToCoordinatesComponent.CurrentCoordinate+1 > len(moveToCoordinatesComponent.Coordinates)-1 {
				entity.MoveToCoordinateComponent = option.None[world.MoveToCoordinateComponent]()
				moveToCoordinatesComponent.State = world.MoveToCoordinatesComponentStateFinished
				continue
			}

			moveToCoordinatesComponent.CurrentCoordinate++
			currentCoordinate := moveToCoordinatesComponent.Coordinates[moveToCoordinatesComponent.CurrentCoordinate]
			entity.MoveToCoordinateComponent = option.Some(world.NewMoveToCoordinateComponent(currentCoordinate))
		case world.MoveToCoordinatesComponentStateFinished:
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

		switch pathfinderComponent.State {
		case world.PathfinderComponentStateNoPath, world.PathfinderComponentStateArrived:
			continue
		case world.PathfinderComponentStateIdle:
			if entity.Level == pathfinderComponent.Destination.Level {
				path := findShortestPath(
					w,
					world.TilePositionFromPosition(position),
					world.TilePositionFromPosition(pathfinderComponent.Destination),
				)
				if path == nil {
					entity.MoveToCoordinatesComponent = option.None[world.MoveToCoordinatesComponent]()
					pathfinderComponent.State = world.PathfinderComponentStateNoPath
					continue
				}
				entity.MoveToCoordinatesComponent = option.Some(world.NewMoveToCoordinatesComponent(path))
				pathfinderComponent.State = world.PathfinderComponentStateMovingToDestination
				continue
			}

			pathToPortal := findShortestPathToPortal(w, world.TilePositionFromPosition(position))
			if pathToPortal == nil {
				entity.MoveToCoordinatesComponent = option.None[world.MoveToCoordinatesComponent]()
				pathfinderComponent.State = world.PathfinderComponentStateNoPath
				continue
			}

			entity.MoveToCoordinatesComponent = option.Some(world.NewMoveToCoordinatesComponent(pathToPortal))
			pathfinderComponent.State = world.PathfinderComponentStateMovingToPortal
			continue
		case world.PathfinderComponentStateMovingToPortal:
			moveToCoordinatesComponent := entity.MoveToCoordinatesComponent.Data
			if moveToCoordinatesComponent.State != world.MoveToCoordinatesComponentStateFinished {
				continue
			}

			pathFromPortal := findShortestPathFromPortal(w, world.TilePositionFromPosition(pathfinderComponent.Destination))
			if pathFromPortal == nil {
				entity.MoveToCoordinatesComponent = option.None[world.MoveToCoordinatesComponent]()
				pathfinderComponent.State = world.PathfinderComponentStateNoPath
				continue
			}
			log.Println("From portal", pathFromPortal)

			entity.Level = pathfinderComponent.Destination.Level
			*coordinate = pathFromPortal[0]

			entity.MoveToCoordinatesComponent = option.Some(world.NewMoveToCoordinatesComponent(pathFromPortal))
			pathfinderComponent.State = world.PathfinderComponentStateMovingToDestination
		case world.PathfinderComponentStateMovingToDestination:
			moveToCoordinatesComponent := entity.MoveToCoordinatesComponent.Data
			if moveToCoordinatesComponent.State != world.MoveToCoordinatesComponentStateFinished {
				continue
			}

			entity.MoveToCoordinatesComponent = option.None[world.MoveToCoordinatesComponent]()
			pathfinderComponent.State = world.PathfinderComponentStateArrived
		}
	}
	return errors.Join(errs...)
}
