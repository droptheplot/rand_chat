package models

import (
	"time"

	"github.com/droptheplot/rand_chat/telegram"
	"github.com/droptheplot/rand_chat/vk"
	"github.com/jinzhu/gorm"
)

type Message struct {
	ID        int
	RoomID    int
	Room      Room
	CreatedAt time.Time

	Text string `gorm:"-"`
	User User   `gorm:"-"`
}

func CreateMessage(db *gorm.DB, roomID int) (message Message) {
	message = Message{RoomID: roomID}

	db.Create(&message)

	return message
}

func (message Message) Handle(db *gorm.DB) {
	switch message.Text {
	case "/start":
		JoinRoom(db, message.User)
	case "/stop":
		StopRoom(db, message.User)
	default:
		room, target := FindRoom(db, message.User)

		go CreateMessage(db, room.ID)

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
