package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/droptheplot/rand_chat/app"
	"github.com/droptheplot/rand_chat/env"
	"github.com/droptheplot/rand_chat/telegram"
	"github.com/droptheplot/rand_chat/vk"
)

var templates = template.Must(template.ParseGlob(env.Config.Templates))
var db, logger = env.Init()

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/api/chart", ChartHandler).Methods("GET")
	r.HandleFunc("/api/telegram", TelegramHandler).Methods("POST")
	r.HandleFunc("/api/vk", VKHandler).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(env.Config.Static))))

	srv := &http.Server{
		Handler:      handlers.RecoveryHandler()(r),
		Addr:         ":443",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	srv.ListenAndServeTLS(env.Config.TLS.Cert, env.Config.TLS.Key)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"Title": "Рандомный чат"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TelegramHandler(w http.ResponseWriter, r *http.Request) {
	var update telegram.Update
	json.NewDecoder(r.Body).Decode(&update)

	app.Message{
		Text: update.Message.Text,
		User: app.User{
			ID:  update.Message.User.ID,
			App: "telegram",
		},
	}.Handle(db, logger)

	w.WriteHeader(http.StatusOK)
}

func VKHandler(w http.ResponseWriter, r *http.Request) {
	var event vk.Event
	json.NewDecoder(r.Body).Decode(&event)

	if event.Type == "confirmation" {
		w.Write([]byte(env.Config.VK.Confirmation))
		return
	}

	if event.Message.Out == 1 {
		w.Write([]byte("ok"))
		return
	}

	app.Message{
		Text: event.Message.Body,
		User: app.User{
			ID:  event.Message.UserID,
			App: "vk",
		},
	}.Handle(db, logger)

	w.Write([]byte("ok"))
}

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	var charts = make(map[string][]interface{})

	for _, chart := range app.GetCharts(db) {
		charts["dates"] = append(charts["dates"], chart.Date.Format("2 Jan"))
		charts["counts"] = append(charts["counts"], chart.Count)
	}

	result, _ := json.Marshal(charts)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}
