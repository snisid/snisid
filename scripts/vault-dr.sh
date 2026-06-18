#!/bin/bash
# SNISID Vault Disaster Recovery Script
set -euo pipefail

VAULT_ADDR_PRIMARY="${VAULT_ADDR_PRIMARY:-https://vault-primary:8200}"
VAULT_ADDR_DR="${VAULT_ADDR_DR:-https://vault-dr:8200}"
VAULT_TOKEN="${VAULT_TOKEN}"

echo "=== SNISID Vault DR ==="

case "${1:-status}" in
  status)
    echo "Primary cluster status:"
    VAULT_ADDR=$VAULT_ADDR_PRIMARY vault status
    echo "DR cluster status:"
    VAULT_ADDR=$VAULT_ADDR_DR vault status
    ;;
  promote)
    echo "[PROMOTE] Promoting DR cluster to primary..."
    VAULT_ADDR=$VAULT_ADDR_DR vault operator raft promote
    echo "[COMPLETE] DR cluster promoted to primary"
    ;;
  demote)
    echo "[DEMOTE] Demoting current primary..."
    VAULT_ADDR=$VAULT_ADDR_PRIMARY vault operator raft demote
    echo "[COMPLETE] Primary demoted"
    ;;
  snapshot)
    SNAPSHOT_FILE="vault-snapshot-$(date +%Y%m%d-%H%M%S).snap"
    echo "[SNAPSHOT] Taking Vault snapshot..."
    VAULT_ADDR=$VAULT_ADDR_PRIMARY vault operator raft snapshot save "$SNAPSHOT_FILE"
    echo "[SAVE] Snapshot saved to $SNAPSHOT_FILE"
    ;;
  restore)
    SNAPSHOT_FILE="${2:-vault-snapshot-latest.snap}"
    echo "[RESTORE] Restoring Vault from $SNAPSHOT_FILE..."
    VAULT_ADDR=$VAULT_ADDR_DR vault operator raft snapshot restore "$SNAPSHOT_FILE"
    echo "[COMPLETE] Vault restored from snapshot"
    ;;
  rekey)
    echo "[REKEY] Rekeying Vault..."
    VAULT_ADDR=$VAULT_ADDR_PRIMARY vault operator rekey -init -key-shares=5 -key-threshold=3
    echo "[COMPLETE] Vault rekey initiated. Follow the rekey process."
    ;;
  *)
    echo "Usage: $0 {status|promote|demote|snapshot|restore|rekey}"
    exit 1
    ;;
esac
