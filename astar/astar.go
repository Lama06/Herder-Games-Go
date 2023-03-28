package astar

type Neighbour struct {
	Neighbour Node
	Cost      int
}

type Node interface {
	Neighbours(context any) []Neighbour
	EstimateCostTo(other Node, context any) int
}

type node struct {
	node          Node
	parent        *node
	costFromStart int
	costToEnd     int
	open          bool
}

func (n *node) totalCost() int {
	return n.costFromStart + n.costToEnd
}

type nodes map[Node]*node

func (n nodes) cheapestOpenNode() *node {
	var cheapest *node
	for _, node := range n {
		if !node.open {
			continue
		}

		if cheapest == nil || node.totalCost() < cheapest.totalCost() {
			cheapest = node
		}
	}
	return cheapest
}

func FindCheapestPath(from Node, to Node, context any) []Node {
	nodes := nodes{
		from: &node{
			node:          from,
			parent:        nil,
			costFromStart: 0,
			costToEnd:     from.EstimateCostTo(to, context),
			open:          true,
		},
	}

	for {
		current := nodes.cheapestOpenNode()
		current.open = false

		if current.node == to {
			path := []Node{current.node}
			for nodes[path[0]].parent != nil {
				path = append([]Node{nodes[path[0]].parent.node}, path...)
			}
			return path
		}

		for _, neighbour := range current.node.Neighbours(context) {
			neighbourNode, neighbourHasNode := nodes[neighbour.Neighbour]

			if !neighbourHasNode {
				nodes[neighbour.Neighbour] = &node{
					node:          neighbour.Neighbour,
					parent:        current,
					costFromStart: current.costFromStart + neighbour.Cost,
					costToEnd:     neighbour.Neighbour.EstimateCostTo(to, context),
					open:          true,
				}
				continue
			}

			if !neighbourNode.open {
				continue
			}

			if current.costFromStart+neighbour.Cost < neighbourNode.costFromStart {
				neighbourNode.costFromStart = current.costFromStart + neighbour.Cost
				neighbourNode.parent = current
				continue
			}
		}
	}
}
