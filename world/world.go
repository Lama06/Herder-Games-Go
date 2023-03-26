package world

type World struct {
	Player   *Entity
	Entities map[*Entity]struct{}
}
