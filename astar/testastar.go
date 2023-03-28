//go:build ignore

package main

import (
	"log"
	"math"

	"github.com/Lama06/Herder-Games/astar"
)

type position struct {
	x, y int
}

func (p position) Neighbours(context any) []astar.Neighbour {
	var neighbours []astar.Neighbour
	for _, xOffset := range [...]int{-1, 0, 1} {
		for _, yOffset := range [...]int{-1, 0, 1} {
			if xOffset == 0 && yOffset == 0 {
				continue
			}

			neighbourPosition := position{
				x: p.x + xOffset,
				y: p.y + yOffset,
			}

			if neighbourPosition.x == 0 && neighbourPosition.y == 5 {
				continue
			}

			neighbourCost := 10
			if xOffset != 0 && yOffset != 0 {
				neighbourCost = 14
			}
			neighbours = append(neighbours, astar.Neighbour{
				Neighbour: neighbourPosition,
				Cost:      neighbourCost,
			})
		}
	}
	return neighbours
}

func (p position) EstimateCostTo(other astar.Node, context any) int {
	xDistance := p.x - other.(position).x
	yDistance := p.y - other.(position).y
	return int(math.Sqrt(float64(xDistance*xDistance)+float64(yDistance*yDistance)) * 10)
}

func main() {
	log.Println(astar.FindCheapestPath(position{0, 0}, position{0, 7}, nil))
}
