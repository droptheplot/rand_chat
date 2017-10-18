package telegram

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/droptheplot/rand_chat/env"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.ParseForm()

		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.EscapedPath(), "/bottoken/sendMessage")
		assert.Equal(t, r.Form.Get("chat_id"), "1")
		assert.Equal(t, r.Form.Get("text"), "hello")
	}))
	defer s.Close()

	base = s.URL

	SendMessage(1, "hello")
}

func TestSetWebhook(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.ParseForm()

		assert.Equal(t, r.Method, "POST")
		assert.Equal(t, r.URL.EscapedPath(), "/bottoken/setWebhook")
		assert.Equal(t, r.Form.Get("url"), env.Config.Telegram.Webhook)
	}))
	defer s.Close()

	base = s.URL

	SetWebhook()
}
