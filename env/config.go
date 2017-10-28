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
	TLS        tls      `json:"tls"`
	Migrations string
	Templates  string
	Static     string
}

type telegram struct {
	Webhook string `json:"webhook"`
	Token   string `json:"token"`
}

type vk struct {
	Confirmation string `json:"confirmation"`
	Token        string `json:"token"`
}

type tls struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
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
	Config.Templates = path.Join(root, "/templates/*.html")
	Config.Static = path.Join(root, "/static")
}
