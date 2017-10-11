package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/droptheplot/rand_chat/env"
)

var base = "https://api.telegram.org"

func SendMessage(chatID int64, text string) *http.Response {
	response, _ := http.PostForm(
		build("sendMessage"),
		url.Values{"chat_id": {strconv.FormatInt(chatID, 10)}, "text": {text}},
	)

	return response
}

func SetWebhook() *http.Response {
	response, _ := http.PostForm(build("setWebhook"), url.Values{"url": {env.Config.Webhook}})

	return response
}

func build(suffix string) string {
	return fmt.Sprintf("%s/bot%s/%s", base, env.Config.Telegram, suffix)
}
