# SNISID National API Security Framework

## 1. Authentication
- **OIDC/OAuth2**: Centralized identity provider for service-to-service and user-to-service auth.
- **mTLS**: Mutual TLS certificates issued by the National Root CA for all mesh communication.

## 2. Authorization
- **Role-Based Access Control (RBAC)**: Specific roles for specific agency actions.
- **Attribute-Based Access Control (ABAC)**: Finer grain control (e.g., access only to records from a specific region).
- **Open Policy Agent (OPA)**: Centralized policy engine.

## 3. Data Protection
- **Encryption in Transit**: TLS 1.3 only.
- **Encryption at Rest**: AES-256 for all API-related storage.
- **API Signing**: JWS (JSON Web Signature) for sensitive payloads.

## 4. Threat Protection
- **WAF (Web Application Firewall)**: Protects the National Gateway against OWASP Top 10.
- **DDoS Mitigation**: Integrated into the gateway layer.
- **Bot Management**: Detection of non-human traffic.

## 5. Zero Trust Principles
- **Never Trust, Always Verify**: Every request must be authenticated and authorized.
- **Least Privilege**: Services only have access to the APIs they strictly need.
- **Micro-segmentation**: Using Service Mesh to isolate workloads.
