package env

import (
	"encoding/json"
	"flag"
	"os"
	"path"
)

type config struct {
	Database   string   `json:"database"`
	Telegram   telegram `json:"telegram"`
	VK         vk       `json:"vk"`
	Migrations string
}

type telegram struct {
	Webhook string `json:"webhook"`
	Token   string `json:"token"`
}

type vk struct {
	Confirmation string `json:"confirmation"`
	Token        string `json:"token"`
}

// Config returns current configuration.
var Config config

func init() {
	var configPath string

	root := os.Getenv("RAND_CHAT_ROOT")

	if flag.Lookup("test.v") == nil {
		configPath = path.Join(root, "/config.json")
	} else {
		configPath = path.Join(root, "/config.test.json")
	}

	configFile, err := os.Open(configPath)

	if err != nil {
		panic(err)
	}

	json.NewDecoder(configFile).Decode(&Config)

	Config.Migrations = "file://" + path.Join(root, "/migrations")
}
