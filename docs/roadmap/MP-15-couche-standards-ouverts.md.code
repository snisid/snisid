# MP-15 — Couche Standards Ouverts

> **Objectif :** Implémenter les standards ouverts critiques manquants (VC W3C, OID4VP, SD-JWT, PKI, IdP OIDC)
> **Lacunes couvertes :** LAC-04, LAC-05, LAC-06, LAC-12
> **Priorité :** Phase 2 — Standards cryptographiques
> **Effort estimé :** 3-4 mois

---

## Modules à développer

### 1. Identity Provider OpenID Connect (`services/identity-provider/`)

Fournisseur d'identité complet avec endpoints OIDC standards.

```
services/identity-provider/
├── cmd/main.go
├── internal/
│   ├── oidc/discovery.go         # .well-known/openid-configuration
│   ├── oidc/authorize.go         # Endpoint autorisation OAuth2
│   ├── oidc/token.go             # Endpoint token (code, refresh, client_creds)
│   ├── oidc/userinfo.go          # Endpoint userinfo avec scopes SNISID
│   ├── oidc/introspect.go        # Introspection de token
│   ├── client/client_manager.go  # Gestion clients OAuth2 (agences)
│   ├── consent/consent_engine.go # Gestion consentement utilisateur
│   └── session/sso_session.go    # Session SSO partagée
├── go.mod
└── Dockerfile
```

### 2. Verifiable Credentials Issuer (`services/vc-issuer/`)

Émission de credentials vérifiables W3C v2.0 signés.

```
services/vc-issuer/
├── cmd/main.go
├── internal/
│   ├── issuer/vc_issuer.go       # Émettre VCs signées (Ed25519/ES256)
│   ├── issuer/vc_schema.go       # Schémas VC (identité, biométrie, document)
│   └── handler/oidc4vci.go       # Endpoint OID4VCI
├── go.mod
└── Dockerfile
```

### 3. VC Verifier avec OID4VP (`services/vc-verifier/`)

Vérification cryptographique des VCs avec protocole OID4VP.

```
services/vc-verifier/
├── cmd/main.go
├── internal/
│   ├── verifier/vc_verifier.go   # Vérification cryptographique
│   └── handler/oid4vp.go         # Endpoint OID4VP
├── go.mod
└── Dockerfile
```

### 4. Selective Disclosure SD-JWT (`services/selective-disclosure/`)

Divulgation sélective d'attributs avec SD-JWT et preuves ZKP basiques.

```
services/selective-disclosure/
├── internal/
│   ├── sdjwt/issuer.go           # Émettre SD-JWT avec claims sélectifs
│   ├── sdjwt/verifier.go         # Vérifier présentation SD-JWT
│   ├── zkp/age_proof.go          # Preuve ZKP âge >= X sans date exacte
│   └── zkp/nationality_proof.go  # Preuve nationalité sans données brutes
├── go.mod
└── Dockerfile
```

### 5. PKI Nationale (`pki/`)

Infrastructure à clé publique nationale avec CA racine haïtienne.

```
pki/
├── scripts/create-root-ca.sh         # Génération CA racine SNISID
├── scripts/create-intermediate-ca.sh # CA intermédiaires par usage
├── scripts/issue-service-cert.sh     # Certificats pour chaque service
├── config/openssl-root.cnf           # Configuration OpenSSL CA racine
├── config/openssl-intermediate.cnf
├── crl/crl-update.sh                 # Mise à jour CRL automatique
├── ocsp/ocsp-responder.go            # Service OCSP pour vérification en ligne
└── docs/pki/pki-architecture.md      # Architecture PKI documentée
```

---

## Dépendances

| Module | Dépend de | Fournit |
|--------|-----------|---------|
| pki | — | CA racine, certificats services, CRL, OCSP |
| identity-provider | pki | IdP OIDC complet (authorize, token, userinfo) |
| vc-issuer | pki | Émission VCs W3C avec OID4VCI |
| vc-verifier | vc-issuer | Vérification VCs avec OID4VP |
| selective-disclosure | vc-issuer | SD-JWT + preuves ZKP |

---

## Standards cibles

- **OIDC** : OpenID Connect Core 1.0, Discovery 1.0
- **OAuth2** : RFC 6749, RFC 7009 (introspection), RFC 7662 (révocation)
- **VC** : W3C Verifiable Credentials Data Model v2.0
- **OID4VCI** : OpenID for Verifiable Credential Issuance
- **OID4VP** : OpenID for Verifiable Presentations
- **SD-JWT** : Selective Disclosure JWT (draft-ietf-oauth-sd-jwt)
- **X.509** : RFC 5280 (PKI), RFC 6960 (OCSP), RFC 5280 (CRL)
