# SNISID National API Gateway Configuration

## Configuration Principles
- **Centralization**: All cross-agency traffic must pass through the gateway.
- **Security**: Mandatory JWT validation and OPA (Open Policy Agent) checks.
- **Resilience**: Rate limiting per agency to prevent DDoS and resource exhaustion.

## Sample Kong Configuration (Declarative)
```yaml
_format_version: "3.0"
services:
  - name: oni-service
    url: http://oni-internal.gov.ht/v1
    routes:
      - name: oni-identity-route
        paths:
          - /api/v1/identity
        methods: [GET, POST]
    plugins:
      - name: jwt
      - name: rate-limiting
        config:
          minute: 100
          policy: local
      - name: key-auth
      - name: acl
        config:
          allow: [ "trusted-agencies" ]

  - name: dgi-service
    url: http://dgi-internal.gov.ht/v1
    routes:
      - name: dgi-tax-route
        paths:
          - /api/v1/tax
    plugins:
      - name: jwt
      - name: request-transformer
```

## Security Enforcement
- **mTLS**: Required between Gateway and Backend Services.
- **Audit**: All requests are logged to the Central Audit Fabric (Loki).
- **Versioning**: Mandatory `/v1/`, `/v2/` in path.
