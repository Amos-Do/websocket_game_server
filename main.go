package main

import (
	"demo/websocket_game_demo/app"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/game", app.HandleConnection)
	fmt.Println("Server started on :8898")
	http.ListenAndServe(":8898", nil)
}
