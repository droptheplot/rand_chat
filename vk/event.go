package vk

type Event struct {
	Type    string  `json:"type"`
	Message Message `json:"object"`
}

type Message struct {
	UserID int64  `json:"user_id"`
	Body   string `json:"body"`
	Out    int    `json:"out"`
}
