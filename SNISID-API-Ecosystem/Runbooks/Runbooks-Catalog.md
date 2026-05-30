# SNISID API Runbooks

## Runbook: API Outage Recovery
- **Symptoms**: 503 Service Unavailable, Connection Refused.
- **Steps**:
  1. Check Gateway health status.
  2. Verify Service Mesh connectivity.
  3. Inspect backend service logs for crashes.
  4. Restart pods if necessary.
  5. Scale out if resources are exhausted.

## Runbook: Gateway Overload
- **Symptoms**: High latency, 429 Too Many Requests.
- **Steps**:
  1. Identify top consumers (IP/Agency).
  2. Adjust rate-limiting policies dynamically.
  3. Enable caching for read-only endpoints.
  4. Increase Gateway replicas.

## Runbook: Certificate Failure
- **Symptoms**: TLS Handshake errors, expired certificates.
- **Steps**:
  1. Verify certificate expiration date.
  2. Force renewal via Cert-Manager.
  3. Verify CA chain integrity.

## Runbook: Kafka Disruption
- **Symptoms**: Producer timeouts, consumer lag increasing.
- **Steps**:
  1. Check Kafka broker health.
  2. Verify disk space on brokers.
  3. Rebalance partitions.
  4. Perform replay of events once stable.
