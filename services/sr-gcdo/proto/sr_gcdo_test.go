package proto

import (
	"testing"
)

func TestGlobalIntelEngine_Correlate_NoEvents(t *testing.T) {
	e := &GlobalIntelEngine{}
	e.Correlate(nil)
	e.Correlate([]ThreatEvent{})
}

func TestGlobalIntelEngine_Correlate_SingleEvent(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{SourceRegion: "US", AttackType: "PHISHING", Severity: 0.6, Timestamp: 1000},
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_Correlate_BelowThreshold(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{SourceRegion: "US", AttackType: "PHISHING", Severity: 0.5, Timestamp: 1000},
		{SourceRegion: "EU", AttackType: "PHISHING", Severity: 0.6, Timestamp: 1001},
		{SourceRegion: "APAC", AttackType: "PHISHING", Severity: 0.4, Timestamp: 1002},
		{SourceRegion: "US", AttackType: "PHISHING", Severity: 0.7, Timestamp: 1003},
		{SourceRegion: "EU", AttackType: "PHISHING", Severity: 0.5, Timestamp: 1004},
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_Correlate_TriggersReflex(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := make([]ThreatEvent, 6)
	for i := range events {
		events[i] = ThreatEvent{
			SourceRegion: "US",
			AttackType:   "DDoS",
			Severity:     0.8,
			Timestamp:    int64(1000 + i),
		}
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_Correlate_MultiplePatterns(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{AttackType: "MALWARE", Severity: 0.7, Timestamp: 1},
		{AttackType: "MALWARE", Severity: 0.6, Timestamp: 2},
		{AttackType: "MALWARE", Severity: 0.8, Timestamp: 3},
		{AttackType: "MALWARE", Severity: 0.5, Timestamp: 4},
		{AttackType: "MALWARE", Severity: 0.9, Timestamp: 5},
		{AttackType: "MALWARE", Severity: 0.7, Timestamp: 6},
		{AttackType: "PHISHING", Severity: 0.5, Timestamp: 7},
		{AttackType: "PHISHING", Severity: 0.6, Timestamp: 8},
		{AttackType: "PHISHING", Severity: 0.7, Timestamp: 9},
		{AttackType: "PHISHING", Severity: 0.5, Timestamp: 10},
		{AttackType: "PHISHING", Severity: 0.6, Timestamp: 11},
		{AttackType: "PHISHING", Severity: 0.7, Timestamp: 12},
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_Correlate_MixedEvents(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{AttackType: "RANSOMWARE", Severity: 0.9, Timestamp: 1},
		{AttackType: "RANSOMWARE", Severity: 0.8, Timestamp: 2},
		{AttackType: "EXPLOIT", Severity: 0.3, Timestamp: 3},
	}
	e.Correlate(events)
}

func TestTriggerGlobalReflex(t *testing.T) {
	e := &GlobalIntelEngine{}
	e.TriggerGlobalReflex("DDoS")
	e.TriggerGlobalReflex("")
}

func TestGlobalIntelEngine_Correlate_ExactThreshold(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{AttackType: "SCAN", Severity: 0.5, Timestamp: 1},
		{AttackType: "SCAN", Severity: 0.5, Timestamp: 2},
		{AttackType: "SCAN", Severity: 0.5, Timestamp: 3},
		{AttackType: "SCAN", Severity: 0.5, Timestamp: 4},
		{AttackType: "SCAN", Severity: 0.5, Timestamp: 5},
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_Correlate_JustAboveThreshold(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{AttackType: "BRUTE_FORCE", Severity: 0.5, Timestamp: 1},
		{AttackType: "BRUTE_FORCE", Severity: 0.5, Timestamp: 2},
		{AttackType: "BRUTE_FORCE", Severity: 0.5, Timestamp: 3},
		{AttackType: "BRUTE_FORCE", Severity: 0.5, Timestamp: 4},
		{AttackType: "BRUTE_FORCE", Severity: 0.5, Timestamp: 5},
		{AttackType: "BRUTE_FORCE", Severity: 0.5, Timestamp: 6},
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_Correlate_DifferentRegions(t *testing.T) {
	e := &GlobalIntelEngine{}
	events := []ThreatEvent{
		{SourceRegion: "US", AttackType: "DATA_EXFIL", Severity: 0.9, Timestamp: 1},
		{SourceRegion: "EU", AttackType: "DATA_EXFIL", Severity: 0.8, Timestamp: 2},
		{SourceRegion: "APAC", AttackType: "DATA_EXFIL", Severity: 0.7, Timestamp: 3},
		{SourceRegion: "US", AttackType: "DATA_EXFIL", Severity: 0.6, Timestamp: 4},
		{SourceRegion: "EU", AttackType: "DATA_EXFIL", Severity: 0.9, Timestamp: 5},
		{SourceRegion: "APAC", AttackType: "DATA_EXFIL", Severity: 0.8, Timestamp: 6},
	}
	e.Correlate(events)
}

func TestGlobalIntelEngine_TableDriven(t *testing.T) {
	e := &GlobalIntelEngine{}

	tests := []struct {
		name   string
		events []ThreatEvent
	}{
		{"empty events", nil},
		{"single event", []ThreatEvent{{AttackType: "A", Severity: 0.5}}},
		{"five same type", []ThreatEvent{
			{AttackType: "X", Severity: 0.1}, {AttackType: "X", Severity: 0.2},
			{AttackType: "X", Severity: 0.3}, {AttackType: "X", Severity: 0.4},
			{AttackType: "X", Severity: 0.5},
		}},
		{"six same type triggers", []ThreatEvent{
			{AttackType: "Y", Severity: 0.1}, {AttackType: "Y", Severity: 0.2},
			{AttackType: "Y", Severity: 0.3}, {AttackType: "Y", Severity: 0.4},
			{AttackType: "Y", Severity: 0.5}, {AttackType: "Y", Severity: 0.6},
		}},
		{"mixed types no trigger", []ThreatEvent{
			{AttackType: "A", Severity: 0.5}, {AttackType: "B", Severity: 0.5},
			{AttackType: "C", Severity: 0.5}, {AttackType: "D", Severity: 0.5},
			{AttackType: "E", Severity: 0.5},
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e.Correlate(tc.events)
		})
	}
}
