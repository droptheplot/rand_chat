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
	var room Room
	var err error

	switch message.Text {
	case "/start":
		tx := db.Begin()

		room, err = FindFreeRoom(tx, message.User)

		if err != nil {
			room = CreateRoom(tx, message.User)
		}

		JoinRoom(tx, room, message.User)

		tx.Commit()
	case "/stop":
		StopRoom(db, message.User)
	default:
		room, err = FindRoom(db, message.User)

		if err == nil {
			go CreateMessage(db, room.ID)

			room.Target(message.User).SendMessage(message.Text)
		}
	}
}
