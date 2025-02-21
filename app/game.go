package app

import (
	"errors"
	"fmt"
	"log"
	"sync"
)

type GameStatus string

const (
	START   GameStatus = "start"
	DECIDED GameStatus = "decided"
)

type Game struct {
	Players []*Player
	Status  GameStatus
	Mutex   sync.Mutex
}

var games = make(map[string]*Game)
var gamesMutex sync.Mutex

// NewGame
func NewGame(player *Player) error {
	gamesMutex.Lock()
	gameID := fmt.Sprintf("game-%d", len(games)+1)
	if !CheckGameID(gameID) {
		games[gameID] = &Game{}
	}
	gamesMutex.Unlock()

	return AddGame(gameID, player)
}

// AddGame
func AddGame(gameID string, player *Player) error {
	if !CheckGameID(gameID) {
		return errors.New("gameID not exist")
	}
	if CheckGameStatus(gameID, START) {
		return errors.New("the gameID ready start! you can't add")
	}
	player.GameID = gameID

	gamesMutex.Lock()
	games[gameID].Mutex.Lock()
	games[gameID].Players = append(games[gameID].Players, player)
	games[gameID].Mutex.Unlock()
	gamesMutex.Unlock()

	// 滿二人就自動開始
	if CheckGamePlayersNum(gameID, 2) {
		return GameStart(gameID)
	}
	return nil
}

// GameStart
func GameStart(gameID string) error {
	if !CheckGameID(gameID) {
		return errors.New("gameID not exist")
	}

	games[gameID].Mutex.Lock()
	games[gameID].Status = START
	games[gameID].Mutex.Unlock()

	log.Printf("start %s ----", gameID)

	for _, player := range games[gameID].Players {
		msg := &Message{
			Key:  PUB_GAME_STAUTS,
			Data: fmt.Sprintf("Hi %s, you can choice rock, paper and scissors", player.Name),
		}
		player.WritMessage(msg)
	}
	return nil
}

// HandleGame
func HandleGame(player *Player, choice string) error {
	if !CheckGameID(player.GameID) {
		return errors.New("gameID not exist")
	}
	if !CheckGameStatus(player.GameID, START) {
		return errors.New("game not start")
	}
	if CheckGameStatus(player.GameID, DECIDED) {
		return nil
	}

	games[player.GameID].Mutex.Lock()
	player.Choice = choice
	games[player.GameID].Mutex.Unlock()

	if CheckGamePlayersNum(player.GameID, 2) &&
		games[player.GameID].Players[0].Choice != "" &&
		games[player.GameID].Players[1].Choice != "" {
		decideWinner(player.GameID)
	}

	return nil
}

// decideWinner
func decideWinner(gameID string) {
	p1, p2 := games[gameID].Players[0], games[gameID].Players[1]
	winnerMsg := "It's a tie!"

	rules := map[string]string{
		"rock":     "scissors",
		"paper":    "rock",
		"scissors": "paper",
	}

	if p1.Choice != p2.Choice {
		if rules[p1.Choice] == p2.Choice {
			winnerMsg = fmt.Sprintf("%s wins!", p1.Name)
			p1.Score++
		} else {
			winnerMsg = fmt.Sprintf("%s wins!", p2.Name)
			p2.Score++
		}
	}

	msg := &Message{
		Key:  PUB_GAME_RESULT,
		Data: fmt.Sprintf("%s Scores - %s: %d, %s: %d", winnerMsg, p1.Name, p1.Score, p2.Name, p2.Score),
	}
	for _, player := range games[gameID].Players {
		player.WritMessage(msg)
	}
	resetGame(gameID)
}

// resetGame
func resetGame(gameID string) {
	if !CheckGameID(gameID) {
		return
	}

	games[gameID].Mutex.Lock()
	for _, player := range games[gameID].Players {
		player.Choice = ""
		player.GameID = ""
	}
	games[gameID].Mutex.Unlock()

	gamesMutex.Lock()
	delete(games, gameID)
	gamesMutex.Unlock()
}

func CheckGameID(gameID string) bool {
	return games[gameID] != nil
}

func CheckGameStatus(gameID string, status GameStatus) bool {
	return games[gameID].Status == status
}

func CheckGamePlayersNum(gameID string, num int) bool {
	return len(games[gameID].Players) == num
}
