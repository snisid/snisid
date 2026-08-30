package analytics

import (
	"time"
)

type MetricWindow struct {
	MetricName string    `json:"metric_name"`
	Timestamp  time.Time `json:"timestamp"`
	Count      int64     `json:"count"`
	Sum        float64   `json:"sum"`
	Min        float64   `json:"min"`
	Max        float64   `json:"max"`
}

type FusedEvent struct {
	EventID       string                 `json:"event_id"`
	CorrelationID string                 `json:"correlation_id"`
	Type          string                 `json:"type"`
	Source        string                 `json:"source"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
}
