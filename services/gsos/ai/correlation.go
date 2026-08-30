package ai

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/snisid/platform/internal/platform/logger"
	"go.uber.org/zap"
)

type SecurityEvent struct {
	ID          string  `json:"id"`
	Country     string  `json:"country"`
	EventType   string  `json:"event_type"`
	Severity    float64 `json:"severity"`
	Timestamp   int64   `json:"timestamp"`
	SourceIP    string  `json:"source_ip,omitempty"`
	TargetID    string  `json:"target_id,omitempty"`
	Description string  `json:"description"`
}

type ThreatCorrelation struct {
	ThreatID        string   `json:"threat_id"`
	Score           float64  `json:"score"`
	Countries       []string `json:"countries"`
	EventTypes      []string `json:"event_types"`
	EventCount      int      `json:"event_count"`
	TimeWindowStart int64    `json:"time_window_start"`
	TimeWindowEnd   int64    `json:"time_window_end"`
	IsCoordinated   bool     `json:"is_coordinated"`
	Technique       string   `json:"technique"`
}

type CorrelationLayer struct {
	ModelVersion string
	mu           sync.Mutex
	eventHistory []SecurityEvent
	windowSize   time.Duration
}

func NewCorrelationLayer(modelVersion string) *CorrelationLayer {
	return &CorrelationLayer{
		ModelVersion: modelVersion,
		eventHistory: []SecurityEvent{},
		windowSize:   15 * time.Minute,
	}
}

func (c *CorrelationLayer) IngestEvent(event SecurityEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.eventHistory = append(c.eventHistory, event)

	cutoff := time.Now().Add(-1 * time.Hour).Unix()
	var recent []SecurityEvent
	for _, e := range c.eventHistory {
		if e.Timestamp >= cutoff {
			recent = append(recent, e)
		}
	}
	c.eventHistory = recent
}

func (c *CorrelationLayer) AnalyzeGlobalThreats(events []SecurityEvent) []ThreatCorrelation {
	logger.Info(context.Background(), "GSOS-AI: analyzing global threat correlation",
		zap.Int("events", len(events)),
		zap.String("model", c.ModelVersion),
	)

	c.mu.Lock()
	c.eventHistory = append(c.eventHistory, events...)
	allEvents := make([]SecurityEvent, len(c.eventHistory))
	copy(allEvents, c.eventHistory)
	c.mu.Unlock()

	if len(allEvents) < 2 {
		return nil
	}

	return c.correlateEvents(allEvents)
}

func (c *CorrelationLayer) correlateEvents(events []SecurityEvent) []ThreatCorrelation {
	countryGroups := make(map[string][]SecurityEvent)
	for _, e := range events {
		countryGroups[e.Country] = append(countryGroups[e.Country], e)
	}

	typeGroups := make(map[string][]SecurityEvent)
	for _, e := range events {
		typeGroups[e.EventType] = append(typeGroups[e.EventType], e)
	}

	var correlations []ThreatCorrelation

	for eventType, group := range typeGroups {
		if len(group) < 3 {
			continue
		}

		countries := make(map[string]int)
		var minTime, maxTime int64 = 1<<63 - 1, 0
		for _, e := range group {
			countries[e.Country]++
			if e.Timestamp < minTime {
				minTime = e.Timestamp
			}
			if e.Timestamp > maxTime {
				maxTime = e.Timestamp
			}
		}

		spread := maxTime - minTime
		countryCount := len(countries)
		avgSeverity := 0.0
		for _, e := range group {
			avgSeverity += e.Severity
		}
		avgSeverity /= float64(len(group))

		score := 0.0
		score += float64(countryCount) * 0.15
		score += avgSeverity * 0.3
		score += float64(len(group)) * 0.02

		if spread < 300 && countryCount > 1 {
			score += 0.25
		}

		countryList := make([]string, 0, len(countries))
		for c := range countries {
			countryList = append(countryList, c)
		}
		sort.Strings(countryList)

		score = math.Min(1.0, score)

		coordinated := spread < 300 && countryCount > 2 && len(group) > 5
		technique := c.classifyTechnique(eventType, group)

		correlation := ThreatCorrelation{
			ThreatID:        "GSOS-" + eventType + "-" + formatTime(minTime),
			Score:           math.Round(score*100) / 100,
			Countries:       countryList,
			EventTypes:      []string{eventType},
			EventCount:      len(group),
			TimeWindowStart: minTime,
			TimeWindowEnd:   maxTime,
			IsCoordinated:   coordinated,
			Technique:       technique,
		}
		correlations = append(correlations, correlation)

		if correlation.Score > 0.7 {
			logger.Warn(context.Background(), "GSOS-AI: high-threat correlation detected",
				zap.String("threat_id", correlation.ThreatID),
				zap.Float64("score", correlation.Score),
				zap.Any("countries", countryList),
				zap.String("technique", technique),
			)
		}
	}

	sort.Slice(correlations, func(i, j int) bool {
		return correlations[i].Score > correlations[j].Score
	})

	return correlations
}

func (c *CorrelationLayer) DetectCoordinatedAttack() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	cutoff := time.Now().Add(-30 * time.Minute).Unix()
	var recent []SecurityEvent
	for _, e := range c.eventHistory {
		if e.Timestamp >= cutoff {
			recent = append(recent, e)
		}
	}

	if len(recent) < 10 {
		return false
	}

	countryCount := make(map[string]int)
	typeCount := make(map[string]int)
	for _, e := range recent {
		countryCount[e.Country]++
		typeCount[e.EventType]++
	}

	if len(countryCount) >= 3 && len(typeCount) >= 2 {
		totalEvents := len(recent)
		maxCountry := 0
		for _, c := range countryCount {
			if c > maxCountry {
				maxCountry = c
			}
		}
		if float64(maxCountry)/float64(totalEvents) < 0.5 {
			logger.Warn(context.Background(), "GSOS-AI: coordinated attack detected across multiple countries")
			return true
		}
	}

	return false
}

func (c *CorrelationLayer) classifyTechnique(eventType string, events []SecurityEvent) string {
	typePatterns := map[string]string{
		"LOGIN_FAILURE":      "CREDENTIAL_BRUTE_FORCE",
		"SUSPICIOUS_ACCESS":  "PRIVILEGE_ESCALATION",
		"DATA_EXFILTRATION":  "DATA_EXFIL",
		"BIOMETRIC_MISMATCH": "IDENTITY_SPOOFING",
		"FRAUD_TRANSACTION":  "SYNTHETIC_FRAUD",
		"LAPI_HIT":           "VEHICLE_THEFT",
		"FPR_HIT":            "WANTED_PERSON",
	}

	if technique, ok := typePatterns[eventType]; ok {
		return technique
	}

	highSeverityCount := 0
	for _, e := range events {
		if e.Severity > 0.8 {
			highSeverityCount++
		}
	}
	if highSeverityCount > len(events)/2 {
		return "TARGETED_ATTACK"
	}

	return "UNCLASSIFIED"
}

func (c *CorrelationLayer) GetStats() map[string]interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	countries := make(map[string]int)
	for _, e := range c.eventHistory {
		countries[e.Country]++
	}

	return map[string]interface{}{
		"total_events":        len(c.eventHistory),
		"countries_observed":  len(countries),
		"model_version":       c.ModelVersion,
		"window_size_minutes": c.windowSize.Minutes(),
	}
}

func formatTime(ts int64) string {
	return time.Unix(ts, 0).Format("20060102-150405")
}
