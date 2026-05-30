# 🛠 SNISID Workflow Engine

> Implémentation **Temporal.io + Camunda 8 Zeebe** des workflows nationaux.
> Phase 3 — Étape 1 et orchestrateur des Étapes 3 à 12.

---

## Stack

| Composant | Techno | Rôle |
|-----------|--------|------|
| BPMN engine | **Camunda 8 / Zeebe** | Exécution des `.bpmn` |
| Long-running orchestration | **Temporal.io** | Sagas, retries, signaux |
| Workers | **TypeScript** (NodeJS 20) + **Java 21** | Activities |
| Job workers Zeebe | **TS via `@camunda8/sdk`** | Service tasks |
| Schema | **Avro** + Confluent Schema Registry | Events Kafka |
| Kafka client | **kafkajs** | Producer / Consumer |
| Logs/Traces | **OpenTelemetry** | Tracing E2E |
| Config | **GitOps** (ArgoCD) | Déploiement |

---

## Arborescence

```
workflow-engine/
├── README.md
├── package.json
├── tsconfig.json
├── src/
│   ├── activities/                # Implémentations (audit, pki, kafka, biometric, ...)
│   ├── workers/                   # Workers Temporal + Zeebe job workers
│   ├── workflows/                 # Workflows code-first Temporal (long-running)
│   ├── kafka/                     # Producer/consumer helpers
│   ├── pki/                       # Signature + TSA
│   ├── observability/             # OpenTelemetry setup
│   └── governance/                # Hooks de versioning, signature BPMN
├── deploy/
│   ├── temporal/                  # k8s manifests
│   ├── zeebe/                     # k8s manifests
│   └── argocd/                    # Apps
└── tests/
    ├── unit/
    ├── integration/
    └── chaos/
```

## Démarrage rapide (dev)

```bash
docker compose -f deploy/dev/docker-compose.yml up -d
npm install
npm run worker:zeebe
npm run worker:temporal
```
