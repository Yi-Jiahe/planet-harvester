package models

type Extractor struct {
	Player       *Player
	Resource     Resource
	ResourceNode *ResourceNode
	Rate         float64
}

func (e *Extractor) Extract() {
	e.Player.Storage[e.Resource] += e.Rate
}

func NewLogger(player *Player, node *ResourceNode) Extractor {
	extractor := Extractor{
		Player:       player,
		Resource:     Wood,
		ResourceNode: node,
		Rate:         0.1,
	}

	return extractor
}
