package models

import (
	"errors"
	"fmt"
)

type Extractor struct {
	Player       *Player
	Resource     Resource
	ResourceNode *ResourceNode
	Rate         float64
	Cost         map[Resource]float64
}

func (e *Extractor) Extract() {
	e.Player.Storage[e.Resource] += e.Rate
}

var extractors = map[string]Extractor{
	"Logger": {
		Resource: Wood,
		Rate:     0.1,
		Cost: map[Resource]float64{
			Wood: 10,
		},
	},
}

func NewExtractor(name string, player *Player, node *ResourceNode) (Extractor, error) {
	extractor := extractors[name]
	for resource, cost := range extractor.Cost {
		if player.Storage[resource] < cost {
			return extractor, errors.New(fmt.Sprintf("Insufficent materials"))
		}
	}

	extractor.Player = player
	extractor.ResourceNode = node

	return extractor, nil
}
