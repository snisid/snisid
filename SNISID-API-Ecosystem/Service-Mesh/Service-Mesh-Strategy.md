# SNISID Government Service Mesh

## Objective
To secure, manage, and monitor internal communications (East-West traffic) between government microservices.

## Technology: Istio
Istio provides a uniform way to secure, connect, and monitor microservices.

### 1. Security (mTLS)
All communications within the mesh are encrypted by default.
```yaml
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: default
  namespace: snisid-mesh
spec:
  mtls:
    mode: STRICT
```

### 2. Traffic Management
Circuit breaking to prevent cascading failures.
```yaml
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: oni-circuit-breaker
spec:
  host: oni-service.snisid.svc.cluster.local
  trafficPolicy:
    outlierDetection:
      consecutive5xxErrors: 5
      interval: 10s
      baseEjectionTime: 30s
      maxEjectionPercent: 100
```

### 3. Observability
Integration with OpenTelemetry for distributed tracing.

## Features
- **Mutual TLS (mTLS)**: Automatic encryption and identity.
- **Retry Policies**: Automatic retries for transient failures.
- **Traffic Splitting**: Canaries for new API versions.
- **Fault Injection**: Testing resilience under stress.
