package models

import (
	"time"

	"github.com/droptheplot/rand_chat/env"
	"github.com/droptheplot/rand_chat/telegram"
	"github.com/droptheplot/rand_chat/vk"
)

type Message struct {
	ID        int
	RoomID    int
	Room      Room
	CreatedAt time.Time

	Text    string `gorm:"-"`
	UserID  int64  `gorm:"-"`
	UserApp string `gorm:"-"`
}

func CreateMessage(roomID int) (message Message) {
	message = Message{RoomID: roomID}

	env.DB.Create(&message)

	return message
}

func (message Message) Handle() {
	switch message.Text {
	case "/start":
		JoinRoom(message.UserID, message.UserApp)
	case "/stop":
		StopRoom(message.UserID, message.UserApp)
	default:
		room, targetID, targetApp := FindRoom(message.UserID, message.UserApp)

		go CreateMessage(room.ID)

		switch targetApp {
		case "vk":
			vk.SendMessage(targetID, message.Text)
		case "telegram":
			telegram.SendMessage(targetID, message.Text)
		default:
			panic("unknown app.")
		}
	}
}
