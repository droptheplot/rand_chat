package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/droptheplot/rand_chat/models"
	"github.com/droptheplot/rand_chat/telegram"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	telegram.SetWebhook()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/api/chart", ChartHandler).Methods("GET")
	r.HandleFunc("/api/update", UpdateHandler).Methods("POST")
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"Title": "RandChat"})

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

		go models.CreateMessage(room.ID)

		telegram.SendMessage(targetID, update.Message.Text)
	}

	w.WriteHeader(http.StatusOK)
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
