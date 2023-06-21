package service

type Message struct {
	Command int    `json:"command"`
	Data    string `json:"data"`
}
