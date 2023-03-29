package astar_test

import (
	"math"
	"testing"

	"github.com/Lama06/Herder-Games/astar"
)

type position struct {
	x, y int
}

func (p position) estimateCostTo(other position) astar.Cost {
	xDiff := float64(p.x - other.x)
	yDiff := float64(p.y - other.y)
	return astar.Cost(math.Sqrt(math.Pow(xDiff, 2)+math.Pow(yDiff, 2)) * 10)
}

func (p position) neighbours() []astar.Neighbour[position] {
	neighbours := make([]astar.Neighbour[position], 0, 8)
	for _, xOffset := range [3]int{-1, 0, 1} {
		for _, yOffset := range [3]int{-1, 0, 1} {
			neighbour := position{
				x: p.x + xOffset,
				y: p.y + yOffset,
			}
			if neighbour == p {
				continue
			}

			cost := astar.Cost(10)
			if xOffset != 0 && yOffset != 0 {
				cost = 14
			}

			neighbours = append(neighbours, astar.Neighbour[position]{
				Node: neighbour,
				Cost: cost,
			})
		}
	}
	return neighbours
}

func TestFindStraightPath(t *testing.T) {
	path := astar.FindPath(astar.Options[position]{
		Start:            position{x: 0, y: 0},
		End:              position{x: 10, y: 0},
		NeighboursFunc:   position.neighbours,
		EstimateCostFunc: position.estimateCostTo,
	})

	expected := astar.Path[position]{
		{x: 0, y: 0},
		{x: 1, y: 0},
		{x: 2, y: 0},
		{x: 3, y: 0},
		{x: 4, y: 0},
		{x: 5, y: 0},
		{x: 6, y: 0},
		{x: 7, y: 0},
		{x: 8, y: 0},
		{x: 9, y: 0},
		{x: 10, y: 0},
	}

	if !path.Equals(expected) {
		t.FailNow()
	}
}

func TestFindDiagonalPath(t *testing.T) {
	path := astar.FindPath(astar.Options[position]{
		Start:            position{x: 0, y: 0},
		End:              position{x: -10, y: 10},
		NeighboursFunc:   position.neighbours,
		EstimateCostFunc: position.estimateCostTo,
	})

	expected := astar.Path[position]{
		{x: 0, y: 0},
		{x: -1, y: 1},
		{x: -2, y: 2},
		{x: -3, y: 3},
		{x: -4, y: 4},
		{x: -5, y: 5},
		{x: -6, y: 6},
		{x: -7, y: 7},
		{x: -8, y: 8},
		{x: -9, y: 9},
		{x: -10, y: 10},
	}

	if !path.Equals(expected) {
		t.FailNow()
	}
}

func TestFindImpossiblePath(t *testing.T) {
	path := astar.FindPath(astar.Options[position]{
		Start: position{x: 0, y: 0},
		End:   position{x: 10, y: 0},
		NeighboursFunc: func(p position) []astar.Neighbour[position] {
			return nil
		},
		EstimateCostFunc: position.estimateCostTo,
	})

	if path != nil {
		t.FailNow()
	}
}

func TestFindPathAroundBlockedPositions(t *testing.T) {
	blockedPositions := map[position]struct{}{
		{x: -1, y: 4}: {},
		{x: 0, y: 4}:  {},
		{x: 1, y: 4}:  {},
		{x: 2, y: 4}:  {},
	}

	path := astar.FindPath(astar.Options[position]{
		Start: position{x: 0, y: 0},
		End:   position{x: 0, y: 6},
		NeighboursFunc: func(p position) []astar.Neighbour[position] {
			neighbours := p.neighbours()
			neighboursExceptBlocked := make([]astar.Neighbour[position], 0, len(neighbours))
			for _, neighbour := range neighbours {
				if _, blocked := blockedPositions[p]; blocked {
					continue
				}

				neighboursExceptBlocked = append(neighboursExceptBlocked, neighbour)
			}
			return neighboursExceptBlocked
		},
		EstimateCostFunc: position.estimateCostTo,
	})

	expected := astar.Path[position]{
		{x: 0, y: 0},
		{x: 0, y: 1},
		{x: 0, y: 2},
		{x: -1, y: 3},
		{x: -2, y: 4},
		{x: -1, y: 5},
		{x: 0, y: 6},
	}

	if !path.Equals(expected) {
		t.FailNow()
	}
}

func TestFindPathWithPortals(t *testing.T) {
	portalLocation := position{x: -5, y: -5}
	portalDesintation := position{x: 105, y: 105}

	path := astar.FindPath(astar.Options[position]{
		Start: position{x: 0, y: 0},
		End:   position{x: 100, y: 100},
		NeighboursFunc: func(p position) []astar.Neighbour[position] {
			neighbours := p.neighbours()

			if p == portalLocation {
				neighbours = append(neighbours, astar.Neighbour[position]{
					Node: portalDesintation,
					Cost: 0,
				})
			}

			return neighbours
		},
		EstimateCostFunc: func(from, to position) astar.Cost {
			// Die Schätzfunktion sollte die Kosten nie überschätzen
			// Wenn es Portale gibt, kann es allerdings sein, dass es einen Weg gibt, der kürzer als die Luftlinie ist
			return 0
		},
	})

	expected := astar.Path[position]{
		{x: 0, y: 0},
		{x: -1, y: -1},
		{x: -2, y: -2},
		{x: -3, y: -3},
		{x: -4, y: -4},
		{x: -5, y: -5},
		{x: 105, y: 105},
		{x: 104, y: 104},
		{x: 103, y: 103},
		{x: 102, y: 102},
		{x: 101, y: 101},
		{x: 100, y: 100},
	}

	if !path.Equals(expected) {
		t.FailNow()
	}
}

func BenchmarkFindPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		astar.FindPath(astar.Options[position]{
			Start:            position{x: 0, y: 0},
			End:              position{x: 1000, y: 1000},
			NeighboursFunc:   position.neighbours,
			EstimateCostFunc: position.estimateCostTo,
		})
	}
}
