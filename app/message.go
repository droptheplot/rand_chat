package app

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
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

func (message Message) Handle(db *gorm.DB, logger zerolog.Logger) {
	var room Room
	var err error

	log := logger.
		Info().
		Int64("user_id", message.User.ID).
		Str("user_app", message.User.App)

	switch message.Text {
	case "/start":
		log.Str("type", "start")

		room, err = FindFreeRoom(db, message.User)

		if err != nil {
			room := CreateRoom(db, message.User)
			message.User.SendMessage("Ищем собеседника...")

			log.Str("action", "room_created")
			log.Int("room_id", room.ID)
		} else {
			JoinRoom(db, room, message.User)

			room, _ = FindRoom(db, message.User)

			room.Owner().SendMessage("Собеседник найден, скажите привет!")
			room.Guest().SendMessage("Собеседник найден, скажите привет!")

			log.Str("action", "room_joined")
			log.Int("room_id", room.ID)
		}
	case "/stop":
		log.Str("type", "stop")

		StopRoom(db, message.User)
	default:
		log.Str("type", "text")

		room, err = FindRoom(db, message.User)

		if err != nil {
			message.User.SendMessage("Отправьте /start чтобы найти собеседника.")

			log.Str("error", "no_start")
		} else {
			log.Int("room_id", room.ID)

			if room.IsEmpty() {
				message.User.SendMessage("Ищем собеседника...")

				log.Str("error", "room_empty")
			} else {
				go CreateMessage(db, room.ID)

				target := room.Target(message.User)
				target.SendMessage(message.Text)

				log.Str("action", "message_sent")
				log.Str("target_app", target.App)
				log.Int64("target_id", target.ID)
			}
		}
	}

	log.Msg("Message")
}
