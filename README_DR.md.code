# SNISID Disaster Recovery (DR) Playbook

## 📊 Recovery Objectives
- **RPO (Recovery Point Objective)**: < 15 minutes (Data loss limit).
- **RTO (Recovery Time Objective)**: < 60 minutes (Downtime limit).

## 🌋 Scenario 1: Primary Regional Outage
1.  **Detection**: Monitoring alerts show 100% packet loss to the primary region.
2.  **Activation**: The SOC Lead authorizes regional failover.
3.  **Execution**:
    - Run `./scripts/failover_trigger.sh dr-cluster-ctx`.
    - Verify Kafka MirrorMaker has synced latest offsets.
    - Validate identity API connectivity in the DR region.

## 💾 Scenario 2: Data Corruption / Ransomware
1.  **Detection**: GNN fraud engine detects massive anomalous state changes.
2.  **Activation**: Infrastructure Lead initiates point-in-time recovery.
3.  **Execution**:
    - Identify latest clean backup from `backups/`.
    - Run `./scripts/restore_offline.sh snisid_backup_20240501.tar.gz`.
    - Perform data integrity check on Neo4j identity graph.

## 📶 Scenario 3: Total Connectivity Loss (Air-Gapped)
1.  **Condition**: National backbone failure.
2.  **Activation**: Local agencies switch to **Offline SOC Mode**.
3.  **Execution**:
    - Deploy local k3d cluster via `.\install.ps1`.
    - Import latest available `offline_images.tar` and `postgres.sql` dump.
    - Continue local identity enrollment in standalone mode.

## 🛠️ Maintenance & Testing
- **Quarterly DR Drill**: Full failover to DR cluster and back.
- **Monthly Backup Audit**: Automated validation of `postgres.sql` and `neo4j.dump` integrity.
