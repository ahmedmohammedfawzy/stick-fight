package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

    "github.com/gorilla/websocket"
)

type Position struct {
    X int `json:"x"`
    Y int `json:"y"`
}

type Player struct {
    Id       int      `json:"id"`
    Position Position `json:"position"`
}

type Game struct {
    Players []Player `json:"players"`
}

var mainGame = Game{}

// todo should switch to uid
var idCounter = 0

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
    http.HandleFunc("/ws", func (w http.ResponseWriter, r *http.Request) {
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

    http.Handle("/", http.FileServer(http.Dir("./public")))

    log.Print("Listening on port 8080")
    http.ListenAndServe(":8080", nil)
}

func registerPlayer(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Print("Could not parse request body")
        return
    }

    var player Player
    err = json.Unmarshal(body, &player)
    if err != nil {
        log.Print("Could not parse request body")
        return
    }

    player.Id = idCounter
    mainGame.Players = append(mainGame.Players, player)
    idCounter++

    jData, err := json.Marshal(player)
    if err != nil {
        log.Print("Could not encode json response")
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jData)
}

func getGameStatus(w http.ResponseWriter, r *http.Request) {
    jData, err := json.Marshal(mainGame)
    if err != nil {
        log.Print("Could not encode json response")
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jData)
}

func updatePlayer(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Print("Could not parse request body")
        return
    }

    var updatedPlayer Player
    err = json.Unmarshal(body, &updatedPlayer)
    if err != nil {
        log.Print("Could not parse request body")
        return
    }

    for updateIndex, player := range mainGame.Players {
        if player.Id == updatedPlayer.Id {
            mainGame.Players[updateIndex].Position = updatedPlayer.Position
            break
        }
    }

    log.Printf("%+v\n", mainGame)

    jData, err := json.Marshal(updatedPlayer)
    if err != nil {
        log.Print("Could not encode json response")
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jData)
}

