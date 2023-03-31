package world

type World struct {
	Player   *Entity
	Entities map[*Entity]struct{}

	BlockedPathfindingTiles map[TilePosition]struct{}
}
