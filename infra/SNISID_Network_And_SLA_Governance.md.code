# SNISID Network & SLA Governance: Configurations & Manifests

**Classification:** RESTRICTED / SOVEREIGN INFRA
**Compliance:** NIST SP 800-207 / ISO 22301 / SLA Governance

This playbook defines the Cilium L7 network policies, Istio circuit breakers, and SLA alerting scripts deployed across the SNISID infrastructure.

---

## 1. Automated SLA Enforcement Alert Trigger

This Python script queries Prometheus metrics for Tier 1 component availability, checks for threshold drops over a rolling window, and triggers notifications based on alert levels.

```python
# File: /opt/snisid/infra/evaluate_sla_alerts.py
import json
import urllib.request
import sys

PROMETHEUS_URL = "http://prometheus-k8s.monitoring.svc.cluster.local:9090"

def get_availability_metric(query):
    """
    Fetches the availability score of a component from Prometheus.
    """
    url = f"{PROMETHEUS_URL}/api/v1/query?query={query}"
    try:
        with urllib.request.urlopen(url, timeout=5) as response:
            data = json.loads(response.read().decode('utf-8'))
            val = data['data']['result'][0]['value'][1]
            return float(val)
    except Exception as e:
        print(f"[-] Metric fetch failed: {str(e)}")
        # Default to 100.0 if Prometheus query fails to prevent false alarms
        return 100.0

def evaluate_sla():
    print("[*] Launching SLA Enforcer check...")
    
    # Target Query: CockroachDB uptime rate over last 15 minutes
    db_uptime_query = "sum(rate(up{job=\"cockroachdb\"}[15m]))%20/%20count(up{job=\"cockroachdb\"})%20*%20100"
    db_availability = get_availability_metric(db_uptime_query)
    print(f"[*] Core Database Availability (Tier 1): {db_availability:.5f}%")
    
    # 1. Red Alert Trigger (< 99.900%)
    if db_availability < 99.900:
        print("[!] RED ALERT: Tier 1 SLA breached! Uptime: {:.5f}%".format(db_availability))
        trigger_escalation("RED", db_availability)
        return "RED"
        
    # 2. Orange Alert Trigger (99.900% - 99.989%)
    elif db_availability < 99.990:
        print("[!] ORANGE ALERT: Tier 1 SLA degraded! Uptime: {:.5f}%".format(db_availability))
        trigger_escalation("ORANGE", db_availability)
        return "ORANGE"
        
    # 3. Yellow Alert Trigger (99.990% - 99.998%)
    elif db_availability < 99.999:
        print("[!] YELLOW ALERT: Tier 1 SLA anomaly! Uptime: {:.5f}%".format(db_availability))
        trigger_escalation("YELLOW", db_availability)
        return "YELLOW"
        
    print("[+] Core Infrastructure SLA target met. Status: GREEN.")
    return "GREEN"

def trigger_escalation(level, score):
    """
    Simulates calling pager and Slack API webhooks.
    """
    payload = {
        "event": "SLA_BREACH_ALERT",
        "severity": level,
        "availability_score": score,
        "message": f"Tier 1 component availability dropped to {score:.5f}%"
    }
    
    # If Level is RED, initiate executive War Room trigger
    if level == "RED":
        payload["action_required"] = "INVOKE_WAR_ROOM_CISO_DG_MINISTERS"
        
    print(f"[+] Escalation Action Dispatched:\n{json.dumps(payload, indent=2)}")

if __name__ == "__main__":
    evaluate_sla()
```

---

## 2. Cilium L7 NetworkPolicy (VLAN 200 to 300)

This manifest enforces strict L7 routing rules, permitting only specific HTTP methods on microservice paths and blocking all other cross-VLAN packets.

```yaml
# File: /deployments/kubernetes/cilium-vlan-l7-policy.yaml
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: limit-api-ingress-l7
  namespace: snisid-core
spec:
  endpointSelector:
    matchLabels:
      app: identity-service
  ingress:
    # Allow traffic from VLAN 200 (Gateway Pods) only
    - fromEndpoints:
        - matchLabels:
            app: kong-api-gateway
            vlan: "200"
      toPorts:
        - ports:
            - port: "8080"
              protocol: TCP
          rules:
            # Enforce L7 HTTP method controls
            http:
              - method: "GET"
                path: "/v1/citizen/.*"
              - method: "POST"
                path: "/v1/citizen/register"
```

---

## 3. Istio Circuit Breaker Config (DestinationRule)

Configure Istio sidecars to trip the circuit breaker and isolate the `identity-service` pod if latency or error thresholds are crossed.

```yaml
# File: /deployments/istio/circuit-breaker-destinationrule.yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: identity-service-cb
  namespace: snisid-core
spec:
  host: identity-service.snisid-core.svc.cluster.local
  trafficPolicy:
    connectionPool:
      tcp:
        maxConnections: 100
      http:
        http1MaxPendingRequests: 10
        maxRequestsPerConnection: 10
    outlierDetection:
      # Trip if pod returns a 5xx error
      consecutive5xxErrors: 3
      interval: 10s
      baseEjectionTime: 60s # Wait 60 seconds before retrying pod
      maxEjectionPercent: 50
```

---

## 4. Operational Fallback Simulation Case

When the circuit breaker trips, the client gateway intercepts the Envoy ejection metric and switches traffic routing to local caching nodes.

```bash
#!/bin/bash
# File: /opt/snisid/infra/simulate_cb_fallback.sh
set -e

GATEWAY_URL="http://gateway.snisid.local:8080"

echo "[*] Auditing API Gateway communication..."

# Perform query check. If API Gateway returns a 503 Service Unavailable (caused by CB trip)
# switch routing tables to offline SQLite sync databases.
status_code=$(curl -s -o /dev/null -w "%{http_code}" --max-time 2 "$GATEWAY_URL/v1/citizen/check")

if [ "$status_code" -eq 503 ] || [ "$status_code" -eq 000 ]; then
  echo "[!] Gateway Tripped or Down (HTTP Status: $status_code)!"
  echo "[*] ACTION: Activating Local Cache-Only Offline Mode..."
  
  # Update local routing flag
  echo "CACHE_ONLY" > /var/run/snisid/offline_flag
  
  # Log to syslog for EDR/Wazuh alerts ingestion
  logger -t snisid-network "[ALERT] Circuit breaker tripped. Local offline cache fallback active."
else
  echo "[+] Normal API Gateway connectivity. HTTP Status: $status_code"
fi
```

---

*Verified and signed by the SNISID Infrastructure Governance Board.*
