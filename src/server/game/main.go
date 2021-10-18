package game

import (
	"fmt"

	"github.com/Yi-Jiahe/planet-harvester/src/server/models"
)

var (
	player = models.Player{
		Storage: map[models.Resource]int{
			models.Wood: 0,
			models.Iron: 0,
			models.Coal: 0,
		},
	}
)

func ChopWood() {
	player.Storage[models.Wood] += 1
}

func ShowResources() string {
	s := ""
	for resource, amount := range player.Storage {
		s += fmt.Sprintf("%s: %d\n", resource.Name, amount)
	}
	return s
}
