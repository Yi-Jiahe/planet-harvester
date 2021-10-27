package models

import (
	"errors"
	"log"
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
	"logger": {
		Rate: 0.1,
		Cost: map[Resource]float64{
			Wood: 10,
		},
	},
}

func NewExtractor(name string, player *Player, node *ResourceNode) (Extractor, error) {
	extractor, ok := extractors[name]
	if !ok {
		return extractor, errors.New("Name not in extractors")
	}

	log.Println(extractor)
	for resource, cost := range extractor.Cost {
		log.Println(resource)
		if player.Storage[resource] < cost {
			return extractor, errors.New("Insufficent materials")
		}
	}
	for resource, cost := range extractor.Cost {
		player.Storage[resource] -= cost
	}
	extractor.Player = player
	extractor.ResourceNode = node
	extractor.Resource = node.Resource

	return extractor, nil
}
