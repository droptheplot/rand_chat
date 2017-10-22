package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStopRoom(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	tx.Create(&Room{OwnerID: 1, OwnerApp: "vk"})
	tx.Create(&Room{OwnerID: 2, OwnerApp: "telegram", GuestID: 1, GuestApp: "vk"})
	tx.Create(&Room{OwnerID: 3, OwnerApp: "telegram"})
	tx.Create(&Room{OwnerID: 4, OwnerApp: "telegram"})

	StopRoom(tx, User{ID: 1, App: "vk"})

	var room Room
	tx.Where(&Room{OwnerID: 1, OwnerApp: "vk"}).First(&room)
	assert.False(t, room.Active)

	var count int
	tx.Model(&Room{}).Where(&Room{Active: true}).Count(&count)
	assert.Equal(t, 2, count)
}

func TestJoinRoom(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	tx.Create(&Room{OwnerID: 1, OwnerApp: "vk"})
	tx.Create(&Room{OwnerID: 2, OwnerApp: "telegram", GuestID: 1, GuestApp: "vk"})
	tx.Create(&Room{OwnerID: 3, OwnerApp: "vk", Active: false})

	user := User{ID: 3, App: "vk"}

	JoinRoom(tx, user)

	var room Room
	tx.Where(&Room{OwnerID: 1, OwnerApp: "vk"}).First(&room)
	assert.Equal(t, user, room.Guest())

	var count int
	tx.Model(&Room{}).Where(&Room{GuestID: 1, GuestApp: "vk", Active: true}).Count(&count)
	assert.Equal(t, 1, count)
}

func TestFindRoom(t *testing.T) {
	tx := db.Begin()
	defer tx.Rollback()

	tx.Create(&Room{OwnerID: 1, OwnerApp: "vk", GuestID: 2, GuestApp: "vk"})
	tx.Create(&Room{OwnerID: 1, OwnerApp: "telegram", GuestID: 2, GuestApp: "vk"})

	user := User{ID: 1, App: "vk"}
	room, target := FindRoom(tx, user)

	assert.Equal(t, user, room.Owner())
	assert.Equal(t, User{ID: 2, App: "vk"}, target)
}
