package models

import (
	"time"

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

		if !db.NewRecord(room) {
			go CreateMessage(db, room.ID)

			target.SendMessage(message.Text)
		}
	}
}
