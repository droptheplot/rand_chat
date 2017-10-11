package models

import (
	"time"

	"github.com/droptheplot/rand_chat/env"
)

type Message struct {
	ID        int
	RoomID    int
	Room      Room
	CreatedAt time.Time
}

func CreateMessage(roomID int) (message Message) {
	message = Message{RoomID: roomID}

	env.DB.Create(&message)

	return message
}
