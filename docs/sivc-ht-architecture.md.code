# SIVC-HT Architecture

## Système d'Intelligence Véhiculaire et Criminalité d'Haïti

### Vue d'ensemble

SIVC-HT est le module MP-18 du système SNISID, conçu pour gérer les véhicules criminels, les plaques volées, et les alertes véhiculaires en temps réel pour la Police Nationale d'Haïti (PNH).

### Microservices

| Service | Port | Rôle |
|---------|------|------|
| vehicle-criminal-svc | 8090 | Service principal - alertes, vérification plaques, CRUD |
| vehicle-alert-svc | 8092 | Dispatch alertes (radio, SMS, push) |
| interpol-sync-svc | 8091 | Synchronisation INTERPOL SMV/SAD |

### Stack technique

- **Backend**: Go 1.22
- **Base de données**: PostgreSQL 16
- **Cache**: Redis 7 (hotlist temps réel)
- **Graphe**: Neo4j (analyse criminelle)
- **Analytics**: ClickHouse
- **Messagerie**: Apache Kafka
- **API**: REST (chi) + gRPC
- **Sécurité**: mTLS + SPIFFE
- **Container**: Docker + Kubernetes

### Flux critique: Vérification plaque LAPI

```
Caméra LAPI → OCR → vehicle-criminal-svc/check/plate/:plate
                          ↓
                    Redis Hotlist (< 1ms)
                          ↓ (miss)
                    PostgreSQL (< 50ms)
                          ↓ (hit)
                    Alerte dispatchée:
                      - Radio PNH
                      - SMS agents
                      - Push notification
```

### PNH Unités

- BLVV: Brigade de Lutte Contre le Vol de Véhicules
- BLTS: Bureau de Lutte contre le Trafic de Stupéfiants
- BAC: Bureau des Affaires Criminelles
- DCPJ: Direction Centrale de la Police Judiciaire
- CAE: Cellule Anti-Enlèvement
- BRI: Brigade de Recherche et d'Intervention
- GIPNH: Groupe d'Intervention de la Police Nationale
