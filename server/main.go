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
    players []Player
}

var mainGame = Game{}

// todo should switch to uid
var idCounter = 0

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/registerPlayer", registerPlayer).Methods(http.MethodPost, http.MethodOptions)
    r.Use(mux.CORSMethodMiddleware(r))
 
    log.Print("Listenning on port 8080")

    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

    log.Fatal(http.ListenAndServe(":8080", handlers.CORS(headersOk, originsOk, methodsOk)(r)))
}

func registerPlayer(w http.ResponseWriter, r *http.Request) {
    log.Print("registerPlayer")

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
    mainGame.players = append(mainGame.players, player)
    idCounter++
    log.Printf("%+v\n", mainGame)

    jData, err := json.Marshal(player)
    if err != nil {
        log.Print("Could not encode json response")
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jData)
}



