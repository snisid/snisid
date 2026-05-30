package handler

import (
	"encoding/json"
	"net/http"
	"nexus-snisid/pkg/eventbus"
	"nexus-snisid/pkg/kafka"
	"time"
)

var producer = kafka.NewProducer("kafka:9092", "events.risk")

func RiskHandler(w http.ResponseWriter, r *http.Request) {

	event := eventbus.Event{
		Type:   "RISK_CALCULATED",
		Source: "risk-engine",
		Payload: map[string]interface{}{
			"score": 0.82,
		},
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(event)
	producer.Publish("risk", data)

	w.Write([]byte("risk calculated"))
}
