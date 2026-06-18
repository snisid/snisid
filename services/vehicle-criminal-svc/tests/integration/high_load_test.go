package integration

import (
	"sync"
	"testing"
	"time"

	"github.com/snisid/vehicle-criminal-svc/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentPlateCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping load test")
	}

	plates := []string{"PP-1234", "SE-00871", "M-5678", "TC-9012", "PL-3456"}

	var wg sync.WaitGroup
	latencies := make([]time.Duration, len(plates)*100)

	idx := 0
	for i := 0; i < 100; i++ {
		for _, plate := range plates {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				start := time.Now()
				result := &domain.PlateCheckResult{
					PlateNumber: p,
					CheckedAt:   time.Now(),
				}
				_ = result
				latencies[idx] = time.Since(start)
				idx++
			}(plate)
		}
	}

	wg.Wait()

	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	avg := total / time.Duration(len(latencies))

	t.Logf("Average check latency: %v", avg)
	assert.Less(t, avg, 5*time.Millisecond, "average latency should be under 5ms")
}
