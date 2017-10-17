package vk

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/droptheplot/rand_chat/env"
)

var base = "https://api.vk.com"

func SendMessage(peerID int64, message string) *http.Response {
	response, _ := http.PostForm(
		build("messages.send"),
		url.Values{"peer_id": {strconv.FormatInt(peerID, 10)}, "message": {message},
			"access_token": {env.Config.VK.Token},
			"v":            {"5.67"}},
	)

	return response
}

func build(suffix string) string {
	return fmt.Sprintf("%s/method/%s", base, suffix)
}
