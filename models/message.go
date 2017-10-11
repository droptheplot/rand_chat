package models

import (
	"time"
)

type Message struct {
	ID        int
	RoomID    int
	Room      Room
	CreatedAt time.Time
}

func CreateMessage(roomID int) (message Message) {
	message = Message{RoomID: roomID}

	DB.Create(&message)

	return message
}
