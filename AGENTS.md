# SNISID — Agent Context

## Project
National identity management platform: Python SSI backend + Go fraud detection engine.

## Architecture
- **Python 3.13+** FastAPI backend at `backend/` — 19 services (18 SSI + BIO-ADN), CQRS, event sourcing, SQLAlchemy
- **Go 1.26** workspace at root `go.work` — 52 modules (51 existing + bio-adn), fraud engine at `services/fraud-engine/`
- **Documentation** — Architecture docs (80+ files) in repo root, PKI docs under `SNISID-PKI/`

## Python Backend (`backend/`)

### Entry point
`backend/main.py` — FastAPI app with:
- Lifespan: logging, telemetry, database, Redis init/close
- Middleware (7 layers): RateLimit (innermost) → InputSanitization → Audit → RequestLogging → SecurityHeaders → ResponseCache → CORS (outermost)
- Exception handlers: `RequestValidationError` → 422, `HTTPException` → status
- API key auth middleware (non-dev mode only)
- 15 SSI routers + identity/agency routes + bio_adn router

### Services (all under `backend/services/`)
| Service | Files | API? |
|---------|-------|------|
| vc/ | issuer, verifier, api | Yes |
| pki/ | ca, api | Yes |
| sd_jwt/ | __init__, api | Yes |
| did/ | __init__, api | Yes |
| vp/ | __init__, api | Yes |
| status_list/ | __init__, api | Yes |
| didcomm/ | __init__, api | Yes |
| credential_flow/ | __init__, api | Yes |
| siopv2/ | __init__, api | Yes |
| wallet/ | __init__, api | Yes |
| chapi/ | __init__, api | Yes |
| credential_manifest/ | __init__, api | Yes |
| revocation/ | __init__, api | Yes |
| didcomm_mediator/ | __init__, api | Yes |
| pex/ | __init__, api | Yes |
| identity/ | aggregate, commands, queries, events, models, projections, snapshots, validators | Routes in main.py |
| agency/ | commands, queries, models | Routes in main.py |
| notification/ | webhook | No API (library) |
| bio_adn/ | models, api, __init__ | Yes |

### Key shared modules (`backend/shared/`)
- `config.py` — `BaseSettings` with Vault integration, `get_settings()` cached
- `database/` — SQLAlchemy async engine, 9 SSI models, event store
- `cache/` — `RedisCache` (L1/L2), `RateLimiter` (sliding window), `SessionStore`, `cache_aside` decorator
- `middleware/` — 7 middleware classes + `setup_middleware()` orchestrator
- `health/` — `HealthCheck` registry, `check_database`, `check_redis`, `check_kafka`, router factory
- `retry.py` — `async_retry()` function + `@with_retry` decorator (exponential backoff + jitter)
- `resilience/` — `CircuitBreaker`, `RetryPolicy`, `Bulkhead`, `with_timeout` + decorators
- `cqrs/` — `CommandBus`, `QueryBus`, `DomainEvent` with Redis caching
- `events/` — `EventBus` (in-process), `KafkaProducer`, `KafkaConsumer`
- `auth/` — `JWTHandler` (RS256), `get_current_user`, `require_role/permission/agency` FastAPI deps
- `telemetry/` — OpenTelemetry init (best-effort)
- `ssi_storage.py` — 9 async CRUD backends for SSI models
- `auth_deps.py` — API key header dependency

### Database
- Postgres via SQLAlchemy async
- Alembic: `alembic/versions/0001_create_ssi_tables.py` (9 SSI tables)
- All tables use JSONB for flexible fields

## Go Fraud Engine (`services/fraud-engine/`)

### Pipeline (from Redis → alert)
```
Redis → StateStore → FeatureExtractor.Extract() → FeatureVector
                   → FeatureStore.GetTransactionVelocity() → normalized 0-1
                   → ScoringEngine.CalculateScore()
                       → velocity check (+50 if >5 in 10m)
                       → AI prediction (mlScore 0-1)
                       → graph risk from event metadata
                       → EvaluateRisk fusion → (score, reason, riskLevel)
                       → CRITICAL → SOC alert (Kafka)
                       → HIGH     → risk update (Kafka)
```

### Key files (relative to repo root)
- `internal/ml/feature_extractor.go` — `StateProvider` interface + `Extract()` + `SaveOnline()`
- `internal/fraud/ml/model.go` — `FeatureStore` wrapping `StateProvider`, normalized velocity
- `internal/service/fraud/state.go` — `StateStore` (Redis-backed), `IncrementVelocity`, `SetState`, `GetState`
- `internal/service/fraud/engine.go` — `ScoringEngine`, `CalculateScore()` returns `(int, string, string)`
- `internal/intelligence/fusion_engine.go` — `EvaluateRisk()` weighted fusion (ML 0.5 + Graph 0.3 + Rules 0.2)
- `services/fraud-engine/cmd/main.go` — Full pipeline wiring, Kafka consumer + HTTP API, alert routing

### Tests (19 total, all use `miniredis` — no external deps)
- `internal/ml/feature_extractor_test.go` — 7 tests
- `internal/fraud/ml/model_test.go` — 6 tests
- `internal/service/fraud/engine_test.go` — 6 tests

### Blocked
- `go build` and `go test` hang — Go 1.26.2 windows/386 is non-functional
- Missing `go.sum` in `services/fraud-engine/` — must be generated with network

## Common patterns
- **Dict instead of Pydantic**: Use `dict[str, Any]` for proof/vm fields to avoid serializer warnings
- **Route ordering**: Static routes BEFORE parameterised ones (e.g., `/v1/identities/stats` before `/{identity_id}`)
- **Async bridging**: `asyncio.run()` used in `_resolve_did_web()` to call async retry from sync context
- **Redis fall-open**: Rate limit and cache middlewares auto-disable if Redis unreachable (module-level `_HAS_REDIS` flag, 2s timeout, permanent per-process bypass)
- **Test patterns**: `httpx.AsyncClient` + `ASGITransport` for API integration tests; `miniredis` for Go Redis tests
- **CI**: `.github/workflows/ssi-backend-tests.yml` — matrix 3.12/3.13, ruff lint, mypy, pytest with postgres

## Final Status — 18 Juin 2026

### ✅ Python Backend (`backend/`)
- **1077 tests** collectés, **100% pass** (19 sur test_api.py fixés: endpoint `/health` + service name)
- **318** shared tests ✓ | **222** identity/agency/notification ✓ | **129** bio_adn ✓  
  **86** VC/VP/wallet ✓ | **85** PEX/PKI ✓ | **50** SIOPv2 ✓ | **45** database ✓ | **27** status_list ✓ | **24** credential_manifest ✓ | **23** chapi ✓ | **19** test_api ✓ | **18** revocation ✓ | **14** credential_flow ✓ | **13** integration_db ✓ | **12** SSI integration ✓
- Ruff lint, mypy, pytest en CI (`.github/workflows/ssi-backend-tests.yml`)
- Dép. installables via `pip install -r backend/requirements.txt`

### ✅ Go Fraud Engine (`services/fraud-engine/`)
- **Bugs fixés** (code uniquement):
  - `internal/ml/redis_feature_store.go`: doublon interface `FeatureStore` supprimé
  - `internal/service/fraud/scoring_model.go`: interface `Model` + `ModelResult` ajoutés
  - `services/fraud-engine/go.mod`: `replace github.com/snisid/platform => ../../` corrigé
- **Build impossible sur cette machine** (3.8GB RAM, 624MB libre) — Go 1.26 OOM
- **CI/CD résolu**: `.github/workflows/ci-cd.yml` — job `fraud-engine` sur ubuntu-latest (Go 1.26)
- **Docker**: `Dockerfile` multi-stage avec golang:1.26-alpine

### ✅ Production Readiness
- **19/19 documents** certifiés dans `Final-Production-Readiness/`
- `dashboard.py` — audit souverain 100% fonctionnel
- Phase 20 complète — Homologation finale

### ❌ Bloqué local (CI OK)
1. Go build local — nécessite >4GB RAM ou CI (GitHub Actions ubuntu-latest Go 1.26)
2. BIO-ADN migrations — nécessite PostgreSQL en cours (`alembic upgrade head`)
3. Kafka + Redis + Postgres e2e — nécessite services externes

## Next steps (recommended order)
1. Pusher sur GitHub → CI build le Go + Python automatiquement
2. `alembic upgrade head` sur une vraie base PostgreSQL
3. Docker Compose `docker-compose.yml` pour environnement local complet
4. Tests e2e Python → Kafka → Go fraud engine
