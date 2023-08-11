package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
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

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/registerPlayer", registerPlayer).Methods(http.MethodPost, http.MethodOptions)
    r.HandleFunc("/getGameStatus", getGameStatus).Methods(http.MethodGet)
    r.HandleFunc("/updatePlayer", updatePlayer).Methods(http.MethodPut, http.MethodOptions)

    r.Use(mux.CORSMethodMiddleware(r))
 

    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    log.Print("Listenning on port 8080")
    log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
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

