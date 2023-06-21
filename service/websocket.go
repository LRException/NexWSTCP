package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ws *websocket.Conn

func OpenWebsocket() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server started on :9898")
	log.Fatal(http.ListenAndServe(":9898", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ws = conn
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var message Message
		jsonUnmarshalError := json.Unmarshal(p, &message)
		if jsonUnmarshalError != nil {
			log.Println("UnmarshalJsonError:", jsonUnmarshalError)
			return
		}
		if message.Command == 0x0001 {
			StartSocket(string(message.Data))
		} else {
			WS2TCP(message)
		}
	}
}

func WS2TCP(message Message) {
	go SendSocketMessage(message.Data, message.Command)

}

func TCP2WS(message Message) {
	log.Println(message.Data)
	ws.WriteJSON(message)
}
