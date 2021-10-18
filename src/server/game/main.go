package game

import (
	"fmt"

	"github.com/Yi-Jiahe/planet-harvester/src/server/models"

	"github.com/rs/xid"
)

var (
	players = map[string]models.Player{}
)

func NewGame() string {
	userId := xid.New().String()

	players[userId] = models.Player{
		Storage: map[models.Resource]int{
			models.Wood: 0,
			models.Iron: 0,
			models.Coal: 0,
		},
	}

	return userId
}

func ChopWood(userId string) {
	players[userId].Storage[models.Wood] += 1
}

func MineIron(userId string) {
	players[userId].Storage[models.Iron] += 1
}

func MineCoal(userId string) {
	players[userId].Storage[models.Coal] += 1
}

func ShowResources(userId string) string {
	s := ""
	for resource, amount := range players[userId].Storage {
		s += fmt.Sprintf("%s: %d\n", resource.Name, amount)
	}
	return s
}
