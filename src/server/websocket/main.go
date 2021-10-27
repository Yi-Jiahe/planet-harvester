package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Yi-Jiahe/planet-harvester/src/server/game"
	"github.com/gorilla/websocket"
	_ "github.com/joho/godotenv/autoload"
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
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/socket", handleSocket)
	hostname := os.Getenv("HOST")
	log.Println(hostname)
	log.Fatal(http.ListenAndServe(hostname, nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: Figure out what to put here
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "User-Id")
	w.Header().Set("Access-Control-Expose-Headers", "User-Id")

	userId := r.Header.Get("User-Id")
	if userId == "" {
		userId = game.NewPlayer()

		w.Header().Set("User-Id", userId)
		return
	}
	if game.PlayerExists(userId) {
		// Send a positive response?
		// TODO: Currently because of lack of persistence, existing player ids can be lost
		// Once persistence is up, there should be no problem
	} else {
		userId = game.NewPlayer()

		w.Header().Set("User-Id", userId)
		return
	}
}

func handleSocket(w http.ResponseWriter, r *http.Request) {
	var userId string

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer c.Close()

	err = c.WriteMessage(1, []byte("Connected"))
	if err != nil {
		log.Println("write:", err)
	}

	// Send updates to client
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			<-ticker.C
			if !game.PlayerExists(userId) {
				continue
			}
			err := c.WriteMessage(1, []byte(game.ShowResources(userId)))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
		}
		log.Printf("recv: %s", message)

		if !game.PlayerExists(userId) {
			if strings.HasPrefix(string(message), "userId") {
				userId = strings.Split(string(message), ":")[1]
			} else {
				err := c.WriteMessage(1, []byte("Please provide login"))
				if err != nil {
					log.Println("write:", err)
				}
			}
			continue
		}

		switch strings.ToLower(string(message)) {
		case "chop wood":
			game.ChopWood(userId)
		case "mine iron":
			game.MineIron(userId)
		case "mine coal":
			game.MineCoal(userId)
		case "place logger":
			err := game.PlaceLogger(userId)
			if err != nil {
				err = c.WriteMessage(1, []byte(err.Error()))
				if err != nil {
					log.Println("write:", err)
				}
				continue
			}
			err = c.WriteMessage(1, []byte("Logger Placed"))
			if err != nil {
				log.Println("write:", err)
			}
		}
	}
}
