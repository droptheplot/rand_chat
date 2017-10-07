package main

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"github.com/droptheplot/rand_chat/config"
	"github.com/droptheplot/rand_chat/room"
	"github.com/droptheplot/rand_chat/telegram"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("postgres", config.Store.Database)

	if err != nil {
		panic(err)
	}

	DB.LogMode(true)

	room.DB = DB
}

func main() {
	driver, _ := postgres.WithInstance(DB.DB(), &postgres.Config{})
	migrations, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	migrations.Up()

	telegram.SetWebhook()

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api" || r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	var update telegram.Update
	json.NewDecoder(r.Body).Decode(&update)

	switch update.Message.Text {
	case "/start":
		room.Join(update.Message.Chat.ID)
	case "/stop":
		room.Destroy(update.Message.Chat.ID)
	default:
		_, targetID := room.Find(update.Message.Chat.ID)

		telegram.SendMessage(targetID, update.Message.Text)
	}

	w.WriteHeader(http.StatusOK)
}
