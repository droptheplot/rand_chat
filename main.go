package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/droptheplot/rand_chat/models"
	"github.com/droptheplot/rand_chat/telegram"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	telegram.SetWebhook()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/api/update", UpdateHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"Title": "qwe"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var update telegram.Update
	json.NewDecoder(r.Body).Decode(&update)

	switch update.Message.Text {
	case "/start":
		models.JoinRoom(update.Message.Chat.ID)
	case "/stop":
		models.StopRoom(update.Message.Chat.ID)
	default:
		room, targetID := models.FindRoom(update.Message.Chat.ID)

		models.CreateMessage(room.ID)

		telegram.SendMessage(targetID, update.Message.Text)
	}

	w.WriteHeader(http.StatusOK)
}
