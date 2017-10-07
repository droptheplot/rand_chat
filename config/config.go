package config

import (
	"os"
	"encoding/json"
)

type store struct {
	Webhook string `json:"webhook"`
	Database string `json:"database"`
	Telegram string `json:"telegram"`
}

var Store store

func init() {
	file, err := os.Open("config.json")

	if err != nil {
		panic(err)
	}

	json.NewDecoder(file).Decode(&Store)

	if err != nil {
		panic(err)
	}
}
