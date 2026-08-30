# SNISID Offline API Resilience Model

## Strategy
Ensure that critical government services remain functional even during intermittent network connectivity.

## 1. Retry Queues
- **Local Outbox Pattern**: Services write to a local database before attempting to call an external API.
- **Backoff Strategy**: Exponential backoff (1s, 2s, 4s, 8s...) for retries.

## 2. Event Replay
- Kafka persistence allows services to "catch up" on missed events once connectivity is restored.

## 3. Delayed Synchronization
- Non-critical data updates are queued and synced in batches during off-peak hours or when bandwidth is available.

## 4. Offline Edge APIs
- Deployment of local API caches at agency field offices.
- **Identity Caching**: Local copy of frequently accessed identity records with limited TTL.

## 5. Circuit Breakers
- Implementation of circuit breakers (via Service Mesh) to prevent hanging requests from depleting local resources during an outage.
