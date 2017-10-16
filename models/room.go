package models

import (
	"time"

	"github.com/droptheplot/rand_chat/env"
	"github.com/droptheplot/rand_chat/telegram"
)

type Room struct {
	ID        int
	OwnerID   int64
	OwnerApp  string
	GuestID   int64  `gorm:"default:NULL"`
	GuestApp  string `gorm:"default:NULL"`
	Active    bool   `gorm:"default:TRUE"`
	CreatedAt time.Time
	Messages  []Message
}

func FindRoom(ID int64, app string) (room Room, targetID int64, targetApp string) {
	env.DB.Where(`((owner_id = ? AND owner_app = ?) OR (guest_id = ? AND guest_app = ?))
									AND active = TRUE`, ID, app, ID, app).First(&room)

	if room.OwnerID == ID {
		targetID = room.GuestID
		targetApp = room.GuestApp
	} else {
		targetID = room.OwnerID
		targetApp = room.OwnerApp
	}

	return room, targetID, targetApp
}

func JoinRoom(ID int64, app string) (room Room) {
	env.DB.Where("guest_id IS NULL AND owner_id != ? AND active = TRUE", ID).First(&room)

	if env.DB.NewRecord(room) {
		room = CreateRoom(ID, app)
	} else {
		env.DB.Model(&room).Updates(Room{GuestID: ID, GuestApp: app})

		telegram.SendMessage(room.OwnerID, "Someone found, say hello!")
		telegram.SendMessage(room.GuestID, "Someone found, say hello!")
	}

	return room
}

func CreateRoom(ownerID int64, app string) (room Room) {
	room = Room{OwnerID: ownerID, OwnerApp: app}

	env.DB.Create(&room)

	telegram.SendMessage(room.OwnerID, "Waiting for someone.")

	return room
}

func StopRoom(ID int64, app string) {
	var room Room

	env.DB.Where(`((owner_id = ? AND owner_app = ?) OR (guest_id = ? AND guest_app = ?))
									AND active = TRUE`, ID, ID).First(&room)

	if env.DB.NewRecord(room) {
		return
	}

	telegram.SendMessage(room.OwnerID, "Disconnected.")
	telegram.SendMessage(room.GuestID, "Disconnected.")

	env.DB.Model(&room).Update("active", false)
}
