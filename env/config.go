package env

import (
	"encoding/json"
	"flag"
	"os"
)

type config struct {
	Database string   `json:"database"`
	Telegram telegram `json:"telegram"`
	VK       vk       `json:"vk"`
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
	var path = os.Getenv("RAND_CHAT_ROOT")

	if flag.Lookup("test.v") == nil {
		path += "/config.json"
	} else {
		path += "/config.test.json"
	}

	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	json.NewDecoder(file).Decode(&Config)

	if err != nil {
		panic(err)
	}
}
