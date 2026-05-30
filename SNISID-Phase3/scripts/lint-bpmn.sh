#!/usr/bin/env bash
# SNISID — BPMN linter (v1.1)
# Vérifie les garde-fous obligatoires définis dans docs/13-BPMN-Repository.md.
# Distingue : workflows utilitaires, offline-first et workflows métier.
# Usage: ./scripts/lint-bpmn.sh [BPMN_DIR]
set -euo pipefail

ROOT="${1:-BPMN}"
errors=0
warnings=0

emit_err()  { echo "❌ [$1] $2 — $3"; errors=$((errors+1)); }
emit_warn() { echo "⚠️  [$1] $2 — $3"; warnings=$((warnings+1)); }

require() {  # require RULE FILE MSG PATTERN
  grep -q "$4" "$2" || emit_err  "$1" "$2" "$3"
}
require_warn() {
  grep -q "$4" "$2" || emit_warn "$1" "$2" "$3"
}

echo "🔎 Linting BPMN under $ROOT/ ..."
echo

for f in $(find "$ROOT" -name '*.bpmn' | sort); do
  echo "→ $f"

  # === Universal rules ===
  require BPMN-001 "$f" "Missing zeebe:versionTag" "zeebe:versionTag"

  # ServiceTask must declare taskDefinition
  if grep -q '<bpmn:serviceTask' "$f" && ! grep -q 'zeebe:taskDefinition' "$f"; then
    emit_err BPMN-008 "$f" "ServiceTask without taskDefinition.type"
  fi

  # UserTask must declare candidateGroups OR candidateUsers
  if grep -q '<bpmn:userTask' "$f" && ! grep -qE 'candidateGroups|candidateUsers' "$f"; then
    emit_err BPMN-007 "$f" "userTask without candidateGroups/candidateUsers"
  fi

  # === Category-based rules ===
  case "$f" in

    # ---- OFFLINE workflows ----
    # They use local audit (audit.local.*) and write to outbox; Kafka emit happens
    # asynchronously via the sync hub (delayed-sync) — so they are exempt from BPMN-003
    # and use audit.local.encrypted instead of audit.emit.
    *Offline/*)
      grep -qE 'audit\.local\.|audit\.emit' "$f" \
        || emit_err BPMN-002 "$f" "Missing audit (local or central)"
      # Kafka not required (deferred via sync hub)
      ;;

    # ---- AUDIT workflow ----
    # It IS the audit pipeline itself
    *Audit/*)
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      ;;

    # ---- ESCALATION workflows ----
    # Emit kafka but don't escalate themselves
    *Escalation/*)
      require BPMN-002 "$f" "Missing audit emit" "audit.emit"
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      ;;

    # ---- FRAUD detection utility ----
    *Fraud/fraud-detection*)
      require BPMN-002 "$f" "Missing audit emit" "audit.emit"
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      ;;

    # ---- Read-only / Verification flows (no document produced) ----
    *Identity/verification*|*Judicial/court-integration*)
      require BPMN-002 "$f" "Missing audit emit" "audit.emit"
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      # PKI signing not needed (no act issued, only a verification token)
      # SLA escalation handled at API gateway layer
      ;;

    # ---- Long-running JUDICIAL with own SLA mechanics ----
    *Judicial/appeal-management*|*Judicial/fraud-investigation*)
      require BPMN-002 "$f" "Missing audit emit" "audit.emit"
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      require BPMN-004 "$f" "Missing pki.sign.qualified" "pki.sign.qualified"
      # SLA in days → escalation managed via separate cron / case management
      ;;

    # ---- Light flows (Tax registration, Health) — recommend escalation ----
    *Tax/*|*Health/*)
      require BPMN-002 "$f" "Missing audit emit" "audit.emit"
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      require BPMN-004 "$f" "Missing pki.sign.qualified" "pki.sign.qualified"
      require_warn BPMN-005 "$f" "Missing SLA/escalation" "escalation.sla.breach\|boundaryEvent"
      require_warn BPMN-006 "$f" "Missing notification.send" "notification.send"
      ;;

    # ---- CRITICAL business workflows (Civil registry, Identity, Judicial validation/suspension) ----
    *)
      require BPMN-002 "$f" "Missing audit emit" "audit.emit"
      require BPMN-003 "$f" "Missing kafka.emit" "kafka.emit"
      require BPMN-004 "$f" "Missing pki.sign.qualified" "pki.sign.qualified"
      require BPMN-005 "$f" "Missing SLA/escalation" "escalation.sla.breach\|escalation.crisis.national\|boundaryEvent"
      require_warn BPMN-006 "$f" "Missing notification.send" "notification.send"
      ;;
  esac

  # CRITIQUE workflows must include fraud detection
  case "$f" in
    *Civil-Registry/*|*Identity/enrollment*|*Identity/recovery*|*Identity/correction*|*Elections/*|*Immigration/*)
      require_warn BPMN-010 "$f" "Missing fraud.detection.automated" "fraud.detection.automated"
      ;;
  esac
done

echo
echo "================================================"
echo "  Errors:   $errors"
echo "  Warnings: $warnings"
echo "================================================"
[ "$errors" -gt 0 ] && exit 2 || exit 0
