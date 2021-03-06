package game

import (
	"fmt"
	"time"

	"github.com/Yi-Jiahe/planet-harvester/src/models"

	"github.com/rs/xid"
)

var (
	users         = map[string]string{}
	players       = map[string]*models.Player{}
	resourceNodes = []models.ResourceNode{
		models.NewTree(),
	}
	extractors = []models.Extractor{}
)

func init() {
	// Start game loop
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			<-ticker.C
			for _, extractor := range extractors {
				extractor.Extract()
			}
		}
	}()
}

func NewPlayer() string {
	userId := xid.New().String()

	players[userId] = &models.Player{
		Storage: map[models.Resource]float64{
			models.Wood: 0,
			models.Iron: 0,
			models.Coal: 0,
		},
	}

	return userId
}

func PlayerExists(userId string) bool {
	_, exists := players[userId]

	return exists
}

func GetUser(email string) string {
	if userId, exists := users[email]; exists {
		return userId
	}

	return ""
}

func NewUser(email string) string {
	userId := NewPlayer()
	users[email] = userId

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

func PlaceLogger(userId string) error {
	logger, err := models.NewExtractor("logger", players[userId], &resourceNodes[0])
	if err != nil {
		return err
	}

	extractors = append(extractors, logger)
	return nil
}

func ShowResources(userId string) string {
	s := ""
	for resource, amount := range players[userId].Storage {
		s += fmt.Sprintf("%s: %.0f\n", resource.Name, amount)
	}
	return s
}
