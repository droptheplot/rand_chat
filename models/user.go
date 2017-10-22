package models

import (
	"github.com/droptheplot/rand_chat/telegram"
	"github.com/droptheplot/rand_chat/vk"
)

type User struct {
	ID  int64
	App string
}

func (user User) SendMessage(text string) {
	switch user.App {
	case "vk":
		vk.SendMessage(user.ID, text)
	case "telegram":
		telegram.SendMessage(user.ID, text)
	default:
		panic("unknown app.")
	}
}
