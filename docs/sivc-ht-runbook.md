# SIVC-HT Runbook

## Opérations et Maintenance

### Démarrage du service

```bash
# Variables d'environnement requises
export SIVC_DB_HOST=localhost
export SIVC_DB_PORT=5432
export SIVC_DB_NAME=snisid_sivc
export SIVC_DB_USER=sivc_svc
export SIVC_DB_PASSWORD=<secret>
export SIVC_REDIS_ADDR=redis-master:6379
export SIVC_KAFKA_BROKERS=kafka:9092
export SIVC_SERVICE_PORT=8090

# Lancement
go run ./cmd/server/main.go
```

### Vérification de santé

```bash
# Health check
curl http://localhost:8090/api/v1/sivc/check/plate/PP-1234

# Vérification plaque rapide
curl http://localhost:8090/api/v1/sivc/check/plate/SE-00871
```

### Monitoring

- Métriques Prometheus sur `/metrics`
- Logs structurés JSON (zap)
- Traces distribuées via OpenTelemetry

### Procédures d'urgence

#### Rafraîchissement hotlist Redis

Si la hotlist est corrompue ou obsolète:

```bash
# Force reload depuis PostgreSQL
curl -X POST http://localhost:8090/api/v1/sivc/hotlist/refresh
```

#### Alerte INTERPOL en échec

Vérifier les syncs en attente:

```bash
curl http://localhost:8090/api/v1/sivc/interpol/sync-status
```

#### Plaque SE non enregistrée détectée

Le système génère automatiquement une alerte CRITICAL pour les plaques SE non enregistrées dans FOVeS. Vérifier:

1. Le statut de l'alerte
2. Le véhicule associé dans FOVeS/SIV
3. L'historique des sightings

### Escalade

| Niveau | Action | Unité |
|--------|--------|-------|
| INFO | Surveillance passive | BAC |
| CAUTION | Vérification recommandée | BLVV |
| WANTED | Interception avec précaution | BLVV, BRI |
| CRITICAL | Opération coordonnée | GIPNH, DCPJ |

### Sauvegarde

```bash
# Backup PostgreSQL SIVC
pg_dump -h localhost -U sivc_svc snisid_sivc > sivc_backup_$(date +%Y%m%d).sql
```
