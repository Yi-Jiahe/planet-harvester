package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Yi-Jiahe/planet-harvester/src/server/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type app struct {
	c      *websocket.Conn
	userId string
}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// TODO: Implement an acutal check
		return true
	}
}

func main() {
	http.HandleFunc("/socket", handleSocket)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	userId := game.NewGame()

	h := http.Header{
		"user-id": {userId},
	}
	c, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

	a := app{
		userId: userId,
		c:      c,
	}

	err = c.WriteMessage(1, []byte(fmt.Sprintf("User Id: %s", userId)))
	if err != nil {
		log.Println("write:", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		log.Printf("recv: %s", message)

		switch strings.ToLower(string(message)) {
		case "chop wood":
			game.ChopWood(userId)
			a.returnResourceValues()
		case "mine iron":
			game.MineIron(userId)
			a.returnResourceValues()
		case "mine coal":
			game.MineCoal(userId)
			a.returnResourceValues()
		}
	}
}

func (a *app) returnResourceValues() {
	err := a.c.WriteMessage(1, []byte(game.ShowResources(a.userId)))
	if err != nil {
		log.Println("write:", err)
	}
}
