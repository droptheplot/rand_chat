package models

import (
	"time"

	"github.com/droptheplot/rand_chat/env"
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

func (room Room) Owner() User {
	return User{ID: room.OwnerID, App: room.OwnerApp}
}

func (room Room) Guest() User {
	return User{ID: room.GuestID, App: room.GuestApp}
}

// FindRoom returns Room and User to send message.
func FindRoom(user User) (Room, User) {
	var room Room
	var target User

	env.DB.Where(`((owner_id = ? AND owner_app = ?) OR (guest_id = ? AND guest_app = ?))
									AND active = TRUE`, user.ID, user.App, user.ID, user.App).First(&room)

	if room.Owner() == user {
		target = room.Guest()
	} else {
		target = room.Owner()
	}

	return room, target
}

func JoinRoom(user User) (room Room) {
	env.DB.Where("guest_id IS NULL AND owner_id != ? AND active = TRUE", user.ID).First(&room)

	if env.DB.NewRecord(room) {
		room = CreateRoom(user)
	} else {
		env.DB.Model(&room).Updates(Room{GuestID: user.ID, GuestApp: user.App})
	}

	return room
}

func CreateRoom(user User) (room Room) {
	room = Room{OwnerID: user.ID, OwnerApp: user.App}

	env.DB.Create(&room)

	return room
}

func StopRoom(user User) {
	env.DB.Model(&Room{}).Where(
		`((owner_id = ? AND owner_app = ?) OR (guest_id = ? AND guest_app = ?)) AND active = TRUE`,
		user.ID, user.App, user.ID, user.App,
	).Update("active", false)
}
