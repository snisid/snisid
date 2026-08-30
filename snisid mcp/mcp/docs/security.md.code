# Sécurité SNISID

## Contrôles implémentés

- JWT avec issuer/audience/expiration.
- RBAC et permissions fines.
- MFA pour permissions sensibles.
- Device trust et isolation session.
- Chiffrement AES-256-GCM utilitaire.
- API key fingerprint/rotation helpers.
- Audit immuable hash-chain.
- Validation Zod stricte et anti injection.
- Middleware HTTP : Helmet, CORS fermé, rate limit, threat detection.

## Menaces adressées

| Menace | Contrôle |
|---|---|
| Prompt injection | prompts système, sanitation, threat middleware |
| Tool poisoning | registry contrôlé, schemas stricts, audit |
| Privilege escalation | deny-by-default RBAC/MFA |
| RCE | aucune exécution shell dans tools, sanitation |
| API abuse | rate limit, API Gateway, timeouts, retries limités |
| Secret leakage | redaction, no hardcoded secret, .env vault-ready |
