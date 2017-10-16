package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/droptheplot/rand_chat/env"
	"github.com/droptheplot/rand_chat/models"
	"github.com/droptheplot/rand_chat/telegram"
	"github.com/droptheplot/rand_chat/vk"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	telegram.SetWebhook()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/api/chart", ChartHandler).Methods("GET")
	r.HandleFunc("/api/telegram", TelegramHandler).Methods("POST")
	r.HandleFunc("/api/vk", VKHandler).Methods("POST")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"Title": "Рандомный чат"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TelegramHandler(w http.ResponseWriter, r *http.Request) {
	var update telegram.Update
	json.NewDecoder(r.Body).Decode(&update)

	message := Message{
		Text:    update.Message.Text,
		UserID:  update.Message.User.ID,
		UserApp: "telegram",
	}

	MessageHandler(message)

	w.WriteHeader(http.StatusOK)
}

func VKHandler(w http.ResponseWriter, r *http.Request) {
	var event vk.Event
	json.NewDecoder(r.Body).Decode(&event)

	if event.Type == "confirmation" {
		w.Write([]byte(env.Config.VK))
		return
	}

	if event.Message.Out == 1 {
		w.Write([]byte("ok"))
		return
	}

	message := Message{
		Text:    event.Message.Body,
		UserID:  event.Message.UserID,
		UserApp: "vk",
	}

	MessageHandler(message)

	w.Write([]byte("ok"))
}

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	var charts = make(map[string][]interface{})

	for _, chart := range models.GetCharts() {
		charts["dates"] = append(charts["dates"], chart.Date.Format("2 Jan"))
		charts["counts"] = append(charts["counts"], chart.Count)
	}

	result, _ := json.Marshal(charts)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

type Message struct {
	Text    string
	UserID  int64
	UserApp string
}

func MessageHandler(message Message) {
	switch message.Text {
	case "/start":
		models.JoinRoom(message.UserID, message.UserApp)
	case "/stop":
		models.StopRoom(message.UserID, message.UserApp)
	default:
		room, targetID, targetApp := models.FindRoom(message.UserID, message.UserApp)

		go models.CreateMessage(room.ID)

		switch targetApp {
		case "vk":
			vk.SendMessage(targetID, message.Text)
		case "telegram":
			telegram.SendMessage(targetID, message.Text)
		default:
			panic("unknown app.")
		}
	}
}
