# SNISID — Migrations de base de données

## Structure

| Migration | Description | Statut |
|-----------|-------------|--------|
| 001 | Schéma identité core (citoyens, events, snapshots) | À appliquer |
| 002 | Chaîne d'audit immutable (Merkle chain) | À appliquer |
| 003 | Références biométriques et documents | À appliquer |
| 004 | Tables SSI (DID, VC, status lists) | À appliquer |
| 005 | Schéma criminel (cases, warrants, evidence) | À appliquer |
| 006 | Cache ML features + model registry | À appliquer |
| 007 | Index de performance et optimisations | À appliquer |

## Application des migrations

### Développement (Docker Compose)
```bash
docker-compose up -d postgres
export PGPASSWORD=dev_password
psql -h localhost -p 5432 -U snisid -d snisid_db -f database/migrations/001_identity_core.up.sql
psql -h localhost -p 5432 -U snisid -d snisid_db -f database/migrations/002_audit_chain.up.sql
```

### Avec golang-migrate (recommandé)
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path database/migrations -database "postgresql://snisid:${DB_PASSWORD}@localhost:5432/snisid_db?sslmode=disable" up
```

### Rollback
```bash
migrate -path database/migrations -database "${DATABASE_URL}" down 1
```

## Convention de nommage
- `NNN_description.up.sql` — application de la migration
- `NNN_description.down.sql` — rollback de la migration
- NNN commence à 001, incrémenté de 1

## Règles
- Jamais de `DROP` dans un `.up.sql` sans migration inverse dans `.down.sql`
- Les tables d'audit et d'events sont **append-only** — ne pas ajouter de DELETE
- Toute nouvelle table doit avoir `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- Utiliser `uuid_generate_v4()` pour les PK UUID
