// https://www.youtube.com/watch?v=-L-WgKMFuhE

package astar

import (
	"github.com/Lama06/Herder-Games/option"
)

type Cost int

type Neighbour[Node comparable] struct {
	Node Node
	Cost Cost
}

type Path[Node comparable] []Node

func (first Path[Node]) Equals(second Path[Node]) bool {
	if len(first) != len(second) {
		return false
	}

	for i := range first {
		if first[i] != second[i] {
			return false
		}
	}

	return true
}

type Options[Node comparable] struct {
	Start            Node
	End              Node
	NeighboursFunc   func(Node) []Neighbour[Node]
	EstimateCostFunc func(from Node, to Node) Cost
}

func FindPath[Node comparable](options Options[Node]) Path[Node] {
	nodes := nodes[Node]{
		options.Start: {
			open:          true,
			parent:        option.None[Node](),
			costFromStart: 0,
			costToEnd:     options.EstimateCostFunc(options.Start, options.End),
		},
	}

	for {
		currentNode, ok := nodes.nextNode()
		if !ok {
			return nil
		}

		if currentNode == options.End {
			path := Path[Node]{currentNode}
			for {
				parent := nodes[path[0]].parent
				if !parent.Present {
					break
				}
				path = append(Path[Node]{parent.Data}, path...)
			}
			return path
		}

		currentNodeData := nodes[currentNode]

		currentNodeData.open = false

		for _, neighbour := range options.NeighboursFunc(currentNode) {
			neighbourNodeData, neighbourExists := nodes[neighbour.Node]

			if !neighbourExists {
				nodes[neighbour.Node] = nodeData[Node]{
					open:          true,
					parent:        option.Some(currentNode),
					costFromStart: currentNodeData.costFromStart + neighbour.Cost,
					costToEnd:     options.EstimateCostFunc(neighbour.Node, options.End),
				}
				continue
			}

			if !neighbourNodeData.open {
				continue
			}

			costFromStartToNeighbourComingFromCurrent := currentNodeData.costFromStart + neighbour.Cost
			if costFromStartToNeighbourComingFromCurrent < neighbourNodeData.costFromStart {
				neighbourNodeData.costFromStart = costFromStartToNeighbourComingFromCurrent
				neighbourNodeData.parent = option.Some(currentNode)
				nodes[neighbour.Node] = neighbourNodeData
				continue
			}
		}

		nodes[currentNode] = currentNodeData
	}
}

type nodeData[Node comparable] struct {
	open          bool
	parent        option.Option[Node]
	costFromStart Cost
	costToEnd     Cost
}

func (n nodeData[Node]) totalCost() Cost {
	return n.costFromStart + n.costToEnd
}

type nodes[Node comparable] map[Node]nodeData[Node]

func (n nodes[Node]) nextNode() (Node, bool) {
	var (
		resultFound bool
		result      Node
		resultData  nodeData[Node]
	)

	for node, data := range n {
		if !data.open {
			continue
		}

		if !resultFound {
			resultFound = true
			resultData = data
			result = node
			continue
		}

		if data.totalCost() < resultData.totalCost() {
			resultData = data
			result = node
			continue
		}

		if data.totalCost() == resultData.totalCost() && data.costFromStart < resultData.costFromStart {
			resultData = data
			result = node
			continue
		}
	}

	return result, resultFound
}
