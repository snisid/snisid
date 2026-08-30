#!/usr/bin/env bash
# SNISID — Idempotent topic creation
# Reads kafka/topics.yaml and applies on the prod cluster.
# Usage: BOOTSTRAP=broker1:9093 ./apply-topics.sh
set -euo pipefail

: "${BOOTSTRAP:?BOOTSTRAP=host:port required}"
: "${KAFKA_CLI:=kafka-topics.sh}"
: "${CMDCFG:=client.properties}"   # mTLS config (security.protocol=SSL, ssl.keystore.location=...)
YAML="${YAML:-kafka/topics.yaml}"

if ! command -v yq >/dev/null; then
  echo "Please install yq (mikefarah)."; exit 2
fi

count=$(yq '.spec.topics | length' "$YAML")
echo "📦 Applying $count topics from $YAML against $BOOTSTRAP"

for i in $(seq 0 $((count - 1))); do
  name=$(yq ".spec.topics[$i].name"             "$YAML")
  parts=$(yq ".spec.topics[$i].partitions // 6" "$YAML")
  rf=$(yq ".spec.topics[$i].replicationFactor // .spec.defaults.replicationFactor" "$YAML")
  cleanup=$(yq ".spec.topics[$i].cleanupPolicy // .spec.defaults.cleanupPolicy" "$YAML")
  retention=$(yq ".spec.topics[$i].retentionMs // .spec.defaults.retentionMs" "$YAML")
  minISR=$(yq ".spec.topics[$i].minInSyncReplicas // .spec.defaults.minInSyncReplicas" "$YAML")

  echo "  → $name (parts=$parts, rf=$rf, cleanup=$cleanup, retention=$retention, min.isr=$minISR)"

  $KAFKA_CLI --bootstrap-server "$BOOTSTRAP" --command-config "$CMDCFG" \
    --create --if-not-exists \
    --topic "$name" \
    --partitions "$parts" \
    --replication-factor "$rf" \
    --config "cleanup.policy=$cleanup" \
    --config "retention.ms=$retention" \
    --config "min.insync.replicas=$minISR" \
    --config "compression.type=zstd" \
    || true

  # ensure configs even if topic exists
  $KAFKA_CLI --bootstrap-server "$BOOTSTRAP" --command-config "$CMDCFG" \
    --alter --topic "$name" \
    --config "retention.ms=$retention" \
    --config "min.insync.replicas=$minISR" 2>/dev/null || true
done

echo "✅ Done."
