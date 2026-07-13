# SNISID: Go Backend Core Implementation Master Prompts

This document serves as the absolute implementation blueprint for the SNISID Go Backend. It translates the high-level architecture into concrete Go packages, interfaces, and implementation steps. These "Master Prompts" are designed to be fed directly to engineering teams (or AI coders) to scaffold the system cleanly.

## 🧱 1. FOUNDATION: Monorepo & Clean Architecture

**Target Stack:** Go 1.22+, Uber Fx (Dependency Injection), Zerolog (Logging), Viper (Config).

**Monorepo Structure:**
```text
snisid-backend/
├── cmd/
│   ├── identity-svc/      # Main entry point (Wire/Fx setup)
│   ├── api-gateway/       # Main entry point for Gateway
│   └── fraud-svc/         # Main entry point for Fraud
├── internal/
│   ├── identity/          # Domain: Identity
│   │   ├── domain/        # Entities, Value Objects, Repository Interfaces
│   │   ├── usecase/       # Business Logic (CRUD, Validation)
│   │   ├── delivery/      # HTTP/gRPC Handlers
│   │   └── repository/    # Postgres / Neo4j Implementations
│   └── shared/            # Shared internal logic
├── pkg/
│   ├── kafka/             # Kafka producer/consumer wrappers
│   ├── logger/            # Zerolog JSON logging setup
│   ├── opa/               # Open Policy Agent Go SDK wrapper
│   └── telemetry/         # OpenTelemetry tracing setup
├── api/                   # OpenAPI / Protobuf definitions
└── deployments/           # Helm charts, Dockerfiles
```

**Implementation Prompt:**
> "Scaffold a Go 1.22 monorepo using standard layout. Implement Uber Fx for dependency injection in `cmd/identity-svc`. Create a `pkg/logger` package using Zerolog configured for JSON output, enforcing a `correlation_id` field in all logs. Create a `pkg/config` package using Viper to load `.env` variables."

---

## 🌐 2. API GATEWAY

**Target Stack:** Go `net/http` (1.22 multiplexer) or Gin, OpenTelemetry, Redis (Rate Limiting).

**Implementation Prompt:**
> "Implement the API Gateway service in `cmd/api-gateway`. Create an HTTP router utilizing strict middleware chains: 
> 1. `RequestIDMiddleware` (injects UUID v7 into context).
> 2. `TracingMiddleware` (OpenTelemetry span creation).
> 3. `AuthMiddleware` (strips Bearer token, validates JWT signature using JWKS, checks Redis blocklist).
> 4. `RateLimitMiddleware` (Token bucket algorithm via Redis).
> Map standard REST routes (e.g., `/api/v1/identities`) to proxy requests to internal gRPC/HTTP services."

---

## 🪪 3. IDENTITY SERVICE

**Implementation Prompt:**
> "Implement the Identity Service using Clean Architecture in `internal/identity`. 
> 1. **Domain:** Define the `Citizen` struct matching the SNISID JSON schema. Define a `CitizenRepository` interface.
> 2. **Repository:** Implement the interface using `sqlc` to connect to PostgreSQL.
> 3. **Usecase:** Implement the `EnrollCitizen` logic. It must persist the data to Postgres (Status: `PENDING_VERIFICATION`) and use the `pkg/kafka` producer to publish an `identity.citizen.enrolled` event.
> 4. **Delivery:** Create REST/gRPC endpoints to expose this usecase."

---

## 🔐 4. AUTH SERVICE

**Implementation Prompt:**
> "Implement the Authentication Service. 
> 1. Create a `GenerateAccessToken` function that creates a 10-minute JWT (RS256) containing `sub`, `aud`, and `roles`.
> 2. Create an `IssueRefreshToken` function that generates a cryptographically secure opaque string, hashes it, and stores it in Postgres/Redis with a 7-day TTL.
> 3. Implement the `Refresh` endpoint: it must validate the opaque string, instantly delete it (rotation), issue a new Access Token, and issue a *new* Refresh Token."

---

## 🛡️ 5. AUTHORIZATION SERVICE

**Implementation Prompt:**
> "Implement the Authorization middleware in `pkg/opa`. 
> 1. Integrate the `github.com/open-policy-agent/opa/rego` Go SDK.
> 2. Write a function `EvaluateAccess(ctx, jwtClaims, requestResource, clientIP) bool`.
> 3. The function must execute a local Rego policy evaluating RBAC (JWT roles) and ABAC (IP, Time) rules.
> 4. If the request contains the `X-Emergency-Override-Reason` header, bypass ABAC but instantly publish a `soc.alert.critical` Kafka event via `pkg/kafka`."

---

## 📡 6. EVENT INFRASTRUCTURE

**Target Stack:** `confluent-kafka-go`, Protobuf Schema Registry.

**Implementation Prompt:**
> "Implement the `pkg/kafka` module.
> 1. **Producer:** Create a thread-safe `EventPublisher` utilizing `confluent-kafka-go`. It must wrap all outgoing messages in the standard SNISID Envelope (including `event_id`, `correlation_id` from Context). Enable `idempotence=true`.
> 2. **Consumer:** Create an `EventSubscriber` framework. It must execute message handlers asynchronously.
> 3. **Retry/DLQ:** If a handler returns an error, the framework must retry 3 times with exponential backoff (using Go channels or `time.Sleep`). On the 4th failure, publish the raw message to `<topic>.dlq` and commit the offset on the main topic."

---

## 🧠 7. ORCHESTRATOR

**Implementation Prompt:**
> "Implement the Saga Orchestrator in `internal/orchestrator`.
> 1. Listen for the `identity.citizen.enrolled` Kafka event.
> 2. Trigger asynchronous gRPC calls to the `AI Inference Service` (for biometrics) and the `Graph Intelligence Service`.
> 3. Wait for the `fraud.case.scored` event.
> 4. Based on the score, send an `ActivateIdentityCommand` or `SuspendIdentityCommand` back to the Identity Service. Do not use 2PC (Two-Phase Commit); rely strictly on compensating events."

---

## 📜 8. AUDIT SYSTEM

**Implementation Prompt:**
> "Implement the Audit Interceptor. 
> 1. Create a gRPC Interceptor and an HTTP Middleware in `pkg/audit`.
> 2. For every mutating request (`POST`, `PUT`, `DELETE`), extract the `ActorID` (from JWT), `Action`, `ResourceID`, and `CorrelationID`.
> 3. Asynchronously publish an `audit.record.logged` event to Kafka using `pkg/kafka`. Ensure this publish operation is non-blocking to the main HTTP response thread, but guaranteed delivery via a local outbox pattern if Kafka is temporarily down."

---

## ⚙️ 9. DATABASE LAYER

**Target Stack:** `sqlc` (PostgreSQL), `neo4j-go-driver`, `go-redis/redis`.

**Implementation Prompt:**
> "Implement the infrastructure adapters in `internal/shared/db`.
> 1. **Postgres:** Use `sqlc` to generate type-safe Go queries from raw SQL schemas. Implement connection pooling using `pgxpool`.
> 2. **Neo4j:** Initialize the `neo4j-go-driver` with routing protocols (`neo4j://`). Implement a `CreateCitizenNode` function that opens a write transaction and executes Cypher queries.
> 3. **Redis:** Use `go-redis` to implement a `TokenBlocklistRepository` checking for JWT `jti` existence in `< 1ms`."
