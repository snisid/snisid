package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("WebSocket Gateway starting on :8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
