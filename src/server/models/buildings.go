package models

type Extractor interface {
	Place(player *Player, node *ResourceNode)
	Extract()
}

type Factory interface {
}

type Logger struct {
	Player       *Player
	Resource     Resource
	ResourceNode *ResourceNode
	Rate         float64
}

func (l *Logger) Place(player *Player, node *ResourceNode) {
	l.Player = player
	l.Resource = Wood
	l.ResourceNode = node
	l.Rate = 1

	l.Player.Extractors = append(l.Player.Extractors, l)
}

func (l *Logger) Extract() {
	l.Player.Storage[l.Resource] += l.Rate
}
