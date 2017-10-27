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
		room, err = FindFreeRoom(db, message.User)

		if err != nil {
			CreateRoom(db, message.User)
			message.User.SendMessage("Ищем собеседника...")
		} else {
			JoinRoom(db, room, message.User)

			room, _ = FindRoom(db, message.User)

			room.Owner().SendMessage("Собеседник найдет, скажите привет!")
			room.Guest().SendMessage("Собеседник найдет, скажите привет!")
		}
	case "/stop":
		StopRoom(db, message.User)
	default:
		room, err = FindRoom(db, message.User)

		if err != nil {
			message.User.SendMessage("Используйте /start чтобы найти собеседника.")
		} else {
			if room.IsEmpty() {
				message.User.SendMessage("Ищем собеседника...")
			} else {
				go CreateMessage(db, room.ID)

				room.Target(message.User).SendMessage(message.Text)
			}
		}
	}
}
