package audit

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/snisid/platform/backend/internal/platform/logger"
)

type AccessAuditLog struct {
	Timestamp     int64                  `json:"timestamp"`
	User          string                 `json:"user"`
	Resource      string                 `json:"resource"`
	Action        string                 `json:"action"`
	Decision      string                 `json:"decision"`
	Justification string                 `json:"justification"`
	Context       map[string]interface{} `json:"context"`
}

type ImmutableLogger struct {
	Sink string // e.g. "S3-WORM" or "KAFKA-AUDIT"
}

func (l *ImmutableLogger) LogAccess(log AccessAuditLog) {
	log.Timestamp = time.Now().Unix()
	data, _ := json.Marshal(log)

	// In a real system, this would push to an immutable WORM (Write Once Read Many) storage
	fmt.Printf("📜 NEXUS-AUDIT: Committing immutable access log: %s\n", string(data))
	logger.Info(fmt.Sprintf("ACCESS-AUDIT: %s", string(data)))
}
