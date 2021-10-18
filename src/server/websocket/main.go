package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Yi-Jiahe/planet-harvester/src/server/game"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		// TODO: Implement an acutal check
		return true
	}
}

func main() {
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	userId := game.NewGame()

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

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
		game.ChopWood(userId)

		err = c.WriteMessage(1, []byte(game.ShowResources(userId)))
		if err != nil {
			log.Println("write:", err)
		}
	}
}
