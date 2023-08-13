package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Player struct {
	Id       int      `json:"id"`
	Position Coordinates `json:"position"`
}

type Game struct {
	Players []Player `json:"players"`
}

var mainGame = Game{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		log.Print("Websocket connection initiated")

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed:", err)
				break
			}

			log.Print(string(message))
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./client/public/")))

	log.Print("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
