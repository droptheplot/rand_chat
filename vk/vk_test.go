package vk

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
		assert.Equal(t, r.URL.EscapedPath(), "/method/messages.send")
		assert.Equal(t, r.Form.Get("peer_id"), "1")
		assert.Equal(t, r.Form.Get("message"), "hello")
		assert.Equal(t, r.Form.Get("access_token"), env.Config.VK.Token)
	}))
	defer s.Close()

	base = s.URL

	SendMessage(1, "hello")
}
