package app

type Message struct {
	Key  MessageKey  `json:"key"`
	Data interface{} `json:"data,omitempty"`
}

type MessageKey string

const (
	WELCOME     MessageKey = "webcome"
	NEW_GAME    MessageKey = "new_game"
	ADD_GAME    MessageKey = "add_game"
	GAME_CHOICE MessageKey = "game_choice"

	GET_SCORE MessageKey = "get_score"

	PUB_GAME_STAUTS MessageKey = "pub_game_status"
	PUB_GAME_RESULT MessageKey = "pub_game_result"
)
