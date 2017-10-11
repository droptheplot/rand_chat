package env

import (
	"encoding/json"
	"os"
)

type config struct {
	Webhook  string `json:"webhook"`
	Database string `json:"database"`
	Telegram string `json:"telegram"`
}

var Config config

func init() {
	file, err := os.Open("config.json")

	if err != nil {
		panic(err)
	}

	json.NewDecoder(file).Decode(&Config)

	if err != nil {
		panic(err)
	}
}
