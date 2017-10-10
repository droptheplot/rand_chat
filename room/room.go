package room

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/droptheplot/rand_chat/telegram"
)

var DB *gorm.DB

type Room struct {
	ID        int
	OwnerID   int64
	GuestID   int64 `gorm:"default:NULL"`
	Active    bool  `gorm:"default:TRUE"`
	CreatedAt time.Time
}

func Find(ID int64) (room Room, targetID int64) {
	DB.Where("(owner_id = ? OR guest_id = ?) AND active = TRUE", ID, ID).First(&room)

	if room.OwnerID == ID {
		targetID = room.GuestID
	} else {
		targetID = room.OwnerID
	}

	return room, targetID
}

func Join(ID int64) (room Room) {
	DB.Where("guest_id IS NULL AND owner_id != ? AND active = TRUE", ID).First(&room)

	if DB.NewRecord(room) {
		room = Create(ID)
	} else {
		DB.Model(&room).Update("guest_id", ID)

		telegram.SendMessage(room.OwnerID, "Someone found, say hello!")
		telegram.SendMessage(room.GuestID, "Someone found, say hello!")
	}

	return room
}

func Create(ownerID int64) (room Room) {
	room = Room{OwnerID: ownerID}

	DB.Create(&room)

	telegram.SendMessage(room.OwnerID, "Waiting for someone.")

	return room
}

func Stop(ID int64) {
	var room Room

	DB.Where("(owner_id = ? OR guest_id = ?) AND active = TRUE", ID, ID).First(&room)

	telegram.SendMessage(room.OwnerID, "Disconnected.")
	telegram.SendMessage(room.GuestID, "Disconnected.")

	DB.Model(&room).Update("active", false)
}
