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
	"google.golang.org/api/idtoken"
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
	http.HandleFunc("/google-login", handleGoogleLogin)
	http.HandleFunc("/socket", handleSocket)
	hostname := os.Getenv("HOST")
	log.Println(hostname)
	log.Fatal(http.ListenAndServe(hostname, nil))
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "User-Id")

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	token := r.FormValue("credential")

	validator, err := idtoken.NewValidator(r.Context())
	if err != nil {
		log.Println(err)
	}

	payload, err := validator.Validate(r.Context(), token, "1089484973261-qsvvlihbqof12s2rgqdi6crtnk92svqi.apps.googleusercontent.com")
	if err != nil {
		log.Println(err)
	}

	if payload != nil {
		email := payload.Claims["email"].(string)
		userId := game.GetUser(email)
		if userId == "" {
			userId = game.NewUser(email)
		}
		w.Header().Set("User-Id", userId)
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
