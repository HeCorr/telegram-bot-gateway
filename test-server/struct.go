package main

type TelegramUpdate struct {
	Message Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Username string `json:"username"`
}
