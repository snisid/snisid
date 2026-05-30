package wsgateway

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/snisid/platform/backend/internal/platform/logger"
	"net/http"
	"sync"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.http.Request) bool { return true },
	}
	clients = make(map[*websocket.Conn]bool)
	mu      sync.Mutex
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WS-GATEWAY: Upgrade failed", err)
		return
	}
	defer ws.Close()

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, ws)
			mu.Unlock()
			break
		}
	}
}

func StartKafkaConsumer(brokers []string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-group",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.Error("WS-GATEWAY: Kafka read failed", err)
			continue
		}

		Broadcast(msg.Value)
	}
}

func Broadcast(msg []byte) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}
