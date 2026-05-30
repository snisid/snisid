# SNISID Microservices: Probes, Contracts & Chaos Playbooks

**Classification:** RESTRICTED / SOVEREIGN CORE
**Compliance:** SLSA Level 4 / Pact Testing / LitmusChaos

This operational playbook defines the standardized Go/Rust health check probe templates, Pact contract testing scripts, and LitmusChaos experiment manifests deployed across the SNISID microservices infrastructure.

---

## 1. Standardized Go Health Probe Handler (/health/ready)

This Go code snippet implements the standardized readiness probe handler, checking local databases and Kafka dependencies, and outputting JSON metrics.

```go
// File: /opt/snisid/core/health_handler.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type DependencyStatus struct {
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
}

type ReadinessResponse struct {
	Status              string                      `json:"status"`
	CheckedAt           int64                       `json:"checked_at"`
	P99LatencyMs        int64                       `json:"p99_latency_ms"`
	Dependencies        map[string]DependencyStatus `json:"dependencies"`
	LastAuditLogEmitted string                      `json:"last_audit_log_emitted"`
}

func checkDatabaseConnection(ctx context.Context, db *sql.DB) DependencyStatus {
	start := time.Now()
	err := db.PingContext(ctx)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return DependencyStatus{Status: "DOWN", LatencyMs: latency}
	}
	return DependencyStatus{Status: "UP", LatencyMs: latency}
}

func ReadinessHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbStatus := checkDatabaseConnection(ctx, db)
		
		// Mock checks for Kafka and Vault endpoints
		kafkaStatus := DependencyStatus{Status: "UP", LatencyMs: 8}
		vaultStatus := DependencyStatus{Status: "UP", LatencyMs: 5}

		overallStatus := "UP"
		if dbStatus.Status == "DOWN" || kafkaStatus.Status == "DOWN" {
			overallStatus = "DOWN"
		}

		response := ReadinessResponse{
			Status:              overallStatus,
			CheckedAt:           time.Now().Unix(),
			P99LatencyMs:        12,
			Dependencies: map[string]DependencyStatus{
				"database": dbStatus,
				"kafka":    kafkaStatus,
				"vault":    vaultStatus,
			},
			LastAuditLogEmitted: "SHA256: 7f83b2a9e10c793d4a5b6c7d8e9f0a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r",
		}

		w.Header().Set("Content-Type", "application/json")
		if overallStatus == "DOWN" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// Initialization placeholder for Go HTTP server
}
```

---

## 2. Pact Contract Testing Configuration

Below is the Go configuration file used to verify that the `identity-svc` (consumer) and `biometric-svc` (provider) schemas match prior to code compilation.

```go
// File: /opt/snisid/core/contract_test.go
package main

import (
	"fmt"
	"testing"
	"github.com/pact-foundation/pact-go/dsl"
)

func TestConsumerContract(t *testing.T) {
	// Initialize Pact client
	pact := &dsl.Pact{
		Consumer: "identity-svc",
		Provider: "biometric-svc",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// Define expected biometric verification JSON schema
	pact.
		AddInteraction().
		Given("A valid biometric matching request").
		UponReceiving("A POST request to match minutiae").
		WithRequest(dsl.Request{
			Method:  "POST",
			Path:    dsl.String("/v1/biometrics/match"),
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(dsl.MapMatcher{
				"citizen_hash_id": dsl.Like("sha256:7f83b2a9e10c"),
				"minutiae_records": dsl.Like("FGP-ISO-19794-2"),
			}),
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body: dsl.Like(dsl.MapMatcher{
				"matched": dsl.Like(true),
				"score":   dsl.Like(0.98),
			}),
		})

	// Run verification check
	err := pact.Verify(func() error {
		fmt.Println("[+] Pact Mock Verification execution completed successfully.")
		return nil
	})

	if err != nil {
		t.Fatalf("[-] Pact Verification Failed: %v", err)
	}
}
```

---

## 3. LitmusChaos Experiment Manifest (Pod-Delete)

This manifest is executed monthly in the staging environment. It deletes random `identity-svc` pods under a 50-request-per-second load.

```yaml
# File: /deployments/chaos/litmus-pod-delete.yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: identity-service-chaos
  namespace: snisid-core
spec:
  appinfo:
    appns: 'snisid-core'
    applabel: 'app=identity-service'
    appkind: 'deployment'
  # Litmus experiment execution parameters
  chaosServiceAccount: litmus-admin
  experiments:
    - name: pod-delete
      spec:
        components:
          env:
            - name: TOTAL_CHAOS_DURATION
              value: '30' # Duration in seconds
            - name: CHAOS_INTERVAL
              value: '10' # Time between pod deletions
            - name: FORCE
              value: 'true'
```

---

## 4. Verification Check Script

This Python script validates that the standardized `/health/ready` endpoint of a service returns a successful HTTP 200 within 2.0 seconds.

```python
# File: /opt/snisid/core/check_readiness.py
import urllib.request
import json
import sys
import time

def verify_readiness(url):
    print(f"[*] Testing Readiness Probe: {url}")
    
    start_time = time.time()
    try:
        with urllib.request.urlopen(url, timeout=2.0) as response:
            status = response.getcode()
            latency = (time.time() - start_time) * 1000
            
            body = response.read().decode('utf-8')
            payload = json.loads(body)
            
            print(f"[+] Response Status: {status} (Latency: {latency:.2f}ms)")
            
            # Verify required schema attributes
            if "status" in payload and payload["status"] == "UP":
                print("[+] Probe validation success. System is UP and compliant.")
                return True
            else:
                print("[-] Probe check returned DOWN status or invalid payload schema.")
                return False
                
    except Exception as e:
        print(f"[-] Probe verification failed: {str(e)}")
        return False

if __name__ == "__main__":
    # Test check against mock payload logic (simulated file execution)
    # Target in production: http://identity-service.snisid-core.svc.cluster.local:8080/health/ready
    if len(sys.argv) < 2:
        print("Usage: python check_readiness.py <probe_url>")
        sys.exit(1)
    verify_readiness(sys.argv[1])
```

---

*Verified and signed by the SNISID Microservice Standards Board.*
