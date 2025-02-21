package app

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Player struct {
	Conn   *websocket.Conn
	Name   string
	Choice string
	Score  int
	GameID string
}

// NewPlayer new player
func NewPlayer(conn *websocket.Conn, name string) *Player {
	return &Player{Conn: conn, Name: name}
}

// logf log format
func (p *Player) logf(format string, v ...any) {
	prefixStr := fmt.Sprintf("[%s] ", p.Name)
	log.Printf(prefixStr+format, v...)
}

// Disconnect
func (p *Player) Disconnect() {
	p.logf("disconnect !!\n")
	p.Conn.Close()
}

// Handler
func (p *Player) Handler() {
	p.logf("join!")
	p.WritMessage(&Message{
		Key:  WELCOME,
		Data: fmt.Sprintf("welcom game %s", p.Name),
	})

	defer p.Disconnect()
	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			p.logf("read: %s\n", err)
			break
		}

		p.handleMessage(message)
	}
}

// handleMessage
func (p *Player) handleMessage(message []byte) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		p.logf("%s\n", err)
		return
	}
	p.logf("receive: %+v\n", msg)

	switch msg.Key {
	case NEW_GAME:
		err := NewGame(p)
		key := NEW_GAME + "_ack"
		if err != nil {
			p.WritMessage(&Message{
				Key:  key,
				Data: err.Error(),
			})
			return
		}
		p.WritMessage(&Message{
			Key:  key,
			Data: p.GameID,
		})
	case ADD_GAME:
		key := ADD_GAME + "_ack"
		gameID, ok := msg.Data.(string)
		if !ok {
			p.WritMessage(&Message{
				Key:  key,
				Data: "Data error",
			})
		}
		err := AddGame(gameID, p)
		if err != nil {
			p.WritMessage(&Message{
				Key:  key,
				Data: err.Error(),
			})
			return
		}
		p.WritMessage(&Message{
			Key:  key,
			Data: p.GameID,
		})
	case GAME_CHOICE:
		key := GAME_CHOICE + "_ack"
		choice, ok := msg.Data.(string)
		if !ok {
			p.WritMessage(&Message{
				Key:  key,
				Data: "Data error",
			})
		}
		err := HandleGame(p, choice)
		if err != nil {
			p.WritMessage(&Message{
				Key:  key,
				Data: err.Error(),
			})
			return
		}
		p.WritMessage(&Message{
			Key:  key,
			Data: fmt.Sprintf("you choice: %s", choice),
		})
	case GET_SCORE:
		key := GET_SCORE + "_ack"
		p.WritMessage(&Message{
			Key:  key,
			Data: fmt.Sprintf("your score is: %d", p.Score),
		})
	}
}

// WritMessage
func (p *Player) WritMessage(message *Message) {
	msg, err := json.Marshal(message)
	if err != nil {
		p.logf("%s\n", err)
		return
	}

	err = p.Conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		p.logf("write: %s\n", err)
		return
	}
}
