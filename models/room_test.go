package models

import (
	"testing"

	"github.com/droptheplot/rand_chat/env"
	"github.com/stretchr/testify/assert"
)

func TestStopRoom(t *testing.T) {
	env.Reset()

	env.DB.Create(&Room{OwnerID: 1, OwnerApp: "vk"})
	env.DB.Create(&Room{OwnerID: 2, OwnerApp: "telegram", GuestID: 1, GuestApp: "vk"})
	env.DB.Create(&Room{OwnerID: 3, OwnerApp: "telegram"})
	env.DB.Create(&Room{OwnerID: 4, OwnerApp: "telegram"})

	StopRoom(User{ID: 1, App: "vk"})

	var room Room
	env.DB.Where(&Room{OwnerID: 1, OwnerApp: "vk"}).First(&room)
	assert.False(t, room.Active)

	var count int
	env.DB.Model(&Room{}).Where(&Room{Active: true}).Count(&count)
	assert.Equal(t, 2, count)
}

func TestJoinRoom(t *testing.T) {
	env.Reset()

	env.DB.Create(&Room{OwnerID: 1, OwnerApp: "vk"})
	env.DB.Create(&Room{OwnerID: 2, OwnerApp: "telegram", GuestID: 1, GuestApp: "vk"})
	env.DB.Create(&Room{OwnerID: 3, OwnerApp: "vk", Active: false})

	user := User{ID: 3, App: "vk"}

	JoinRoom(user)

	var room Room
	env.DB.Where(&Room{OwnerID: 1, OwnerApp: "vk"}).First(&room)
	assert.Equal(t, user, room.Guest())

	var count int
	env.DB.Model(&Room{}).Where(&Room{GuestID: 1, GuestApp: "vk", Active: true}).Count(&count)
	assert.Equal(t, 1, count)
}

func TestFindRoom(t *testing.T) {
	env.Reset()

	env.DB.Create(&Room{OwnerID: 1, OwnerApp: "vk", GuestID: 2, GuestApp: "vk"})
	env.DB.Create(&Room{OwnerID: 1, OwnerApp: "telegram", GuestID: 2, GuestApp: "vk"})

	user := User{ID: 1, App: "vk"}
	room, target := FindRoom(user)

	assert.Equal(t, user, room.Owner())
	assert.Equal(t, User{ID: 2, App: "vk"}, target)
}
