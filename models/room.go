package models

import (
	"time"

	"github.com/droptheplot/rand_chat/telegram"
)

type Room struct {
	ID        int
	OwnerID   int64
	GuestID   int64 `gorm:"default:NULL"`
	Active    bool  `gorm:"default:TRUE"`
	CreatedAt time.Time
	Messages  []Message
}

func FindRoom(ID int64) (room Room, targetID int64) {
	DB.Where("(owner_id = ? OR guest_id = ?) AND active = TRUE", ID, ID).First(&room)

	if room.OwnerID == ID {
		targetID = room.GuestID
	} else {
		targetID = room.OwnerID
	}

	return room, targetID
}

func JoinRoom(ID int64) (room Room) {
	DB.Where("guest_id IS NULL AND owner_id != ? AND active = TRUE", ID).First(&room)

	if DB.NewRecord(room) {
		room = CreateRoom(ID)
	} else {
		DB.Model(&room).Update("guest_id", ID)

		telegram.SendMessage(room.OwnerID, "Someone found, say hello!")
		telegram.SendMessage(room.GuestID, "Someone found, say hello!")
	}

	return room
}

func CreateRoom(ownerID int64) (room Room) {
	room = Room{OwnerID: ownerID}

	DB.Create(&room)

	telegram.SendMessage(room.OwnerID, "Waiting for someone.")

	return room
}

func StopRoom(ID int64) {
	var room Room

	DB.Where("(owner_id = ? OR guest_id = ?) AND active = TRUE", ID, ID).First(&room)

	telegram.SendMessage(room.OwnerID, "Disconnected.")
	telegram.SendMessage(room.GuestID, "Disconnected.")

	DB.Model(&room).Update("active", false)
}
