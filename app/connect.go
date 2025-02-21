package app

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Lobby struct {
	Players []*Player
	Mutex   sync.Mutex
}

var lobby = Lobby{}

// cross domain
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnection connect game
func HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	// create player
	lobby.Mutex.Lock()
	playerName := fmt.Sprintf("player-%d", len(lobby.Players)+1)
	player := NewPlayer(conn, playerName)
	lobby.Players = append(lobby.Players, player)
	lobby.Mutex.Unlock()

	player.Handler()
}
