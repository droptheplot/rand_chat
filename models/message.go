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

	Text string `gorm:"-"`
	User User   `gorm:"-"`
}

func CreateMessage(roomID int) (message Message) {
	message = Message{RoomID: roomID}

	env.DB.Create(&message)

	return message
}

func (message Message) Handle() {
	switch message.Text {
	case "/start":
		JoinRoom(message.User)
	case "/stop":
		StopRoom(message.User)
	default:
		room, target := FindRoom(message.User)

		go CreateMessage(room.ID)

		switch target.App {
		case "vk":
			vk.SendMessage(target.ID, message.Text)
		case "telegram":
			telegram.SendMessage(target.ID, message.Text)
		default:
			panic("unknown app.")
		}
	}
}
