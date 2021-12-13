package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"fmt"

	"github.com/Yi-Jiahe/planet-harvester/src/game"
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
	port := os.Getenv("PORT")
	log.Println(fmt.Sprintf("Listening on %s...", port))
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	token := r.FormValue("credential")

	email, err := validatejwt(r.Context(), token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusForbidden)
	}

	if email != "" {
		userId := game.GetUser(email)
		if userId == "" {
			userId = game.NewUser(email)
		}
		w.WriteHeader(http.StatusOK)
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
				return
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
			if strings.HasPrefix(string(message), "jwt") {
				token := strings.Split(string(message), ":")[1]
				email, err := validatejwt(r.Context(), token)
				if err != nil {
					log.Println(err)
				}
				if email != "" {
					userId = game.GetUser(email)
				}
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

func validatejwt(ctx context.Context, token string) (string, error) {
	audience := os.Getenv("AUDIENCE")

	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return "", err
	}

	payload, err := validator.Validate(ctx, token, audience)
	if err != nil {
		return "", err
	}

	if payload != nil {
		email := payload.Claims["email"].(string)
		return email, nil
	}

	return "", errors.New("Something went wrong")
}
