#!/bin/bash
# SNISID Unified Backup Utility
# Version: 1.0.0

set -e
TIMESTAMP=$(date +%Y%m%d%H%M%S)
BACKUP_DIR="./backups/$TIMESTAMP"
mkdir -p "$BACKUP_DIR"

echo "🚀 Starting SNISID Full Backup [$TIMESTAMP]..."

# 1. PostgreSQL Backup
echo "🐘 Backing up PostgreSQL..."
kubectl exec -n snisid deploy/snisid-postgresql -- pg_dump -U snisid snisid > "$BACKUP_DIR/postgres.sql"

# 2. Neo4j Backup
echo "🌿 Backing up Neo4j Graph..."
kubectl exec -n snisid deploy/snisid-neo4j -- bin/neo4j-admin database dump neo4j --to-path=/var/lib/neo4j/backups
kubectl cp snisid/snisid-neo4j:/var/lib/neo4j/backups/neo4j.dump "$BACKUP_DIR/neo4j.dump"

# 3. Kafka Metadata (Optional, for topic reconstruction)
echo "🚇 Backing up Kafka Topic Metadata..."
kubectl exec -n snisid deploy/snisid-kafka -- kafka-topics.sh --bootstrap-server localhost:9092 --list > "$BACKUP_DIR/kafka_topics.txt"

# 4. Kubernetes Manifests
echo "☸️ Backing up K8s Resources..."
kubectl get all -n snisid -o yaml > "$BACKUP_DIR/k8s_resources.yaml"
kubectl get secret -n snisid -o yaml > "$BACKUP_DIR/k8s_secrets.yaml"

# 5. Compress
echo "📦 Finalizing Backup Archive..."
tar -czf "snisid_backup_$TIMESTAMP.tar.gz" -C "$BACKUP_DIR" .
rm -rf "$BACKUP_DIR"

echo "✅ Backup Complete: snisid_backup_$TIMESTAMP.tar.gz"
