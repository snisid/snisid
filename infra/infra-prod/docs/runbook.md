# Runbook: Sovereign Infrastructure (SNISID)

## 1. Incident Response (SOC)
### High-Risk Fraud Alert
1. **Detection**: Kafka event published to `fraud.detected`.
2. **Analysis**: SOC engine identifies target entity and severity.
3. **Containment**: Autonomous isolation of the entity identity node via Istio mesh policies.
4. **Escalation**: Alert pushed to Command Center UI for human validator review.

## 2. Backup & Continuity
### Neo4j Graph Database
- **Frequency**: Every 4 hours.
- **Command**: `neo4j-admin backup --to=/backup/neo4j/$(date +%Y%m%d)`
### MinIO Object Storage
- **Command**: `mc mirror minio/production-bucket backup-server/vault/`

## 3. Disaster Recovery (DR)
### Complete Cluster Failure
1. **Infrastructure Rebuild**: Run `terraform apply` to provision new node pools.
2. **Storage Restore**: Re-mount S3/MinIO volumes and restore Neo4j from last verified snapshot.
3. **Application Sync**: Trigger ArgoCD manual sync: `argocd app sync snisid`.
4. **Traffic Re-routing**: Update global load balancer to point to new cluster ingress.
