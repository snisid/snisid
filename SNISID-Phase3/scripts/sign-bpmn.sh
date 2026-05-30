#!/usr/bin/env bash
# SNISID — BPMN signing (WGO)
# Calcule SHA-384 de chaque .bpmn et le signe via la PKI nationale.
# Sortie : .bpmn-signatures/<chemin>.bpmn.sig
# Usage: WGO_TOKEN=... ./scripts/sign-bpmn.sh
set -euo pipefail

ROOT="${1:-BPMN}"
SIG_DIR="${SIG_DIR:-.bpmn-signatures}"
PKI_ENDPOINT="${PKI_ENDPOINT:-https://pki.snisid.ht/sign}"
: "${WGO_TOKEN:?WGO_TOKEN required (OIDC bearer for WGO signing role)}"

mkdir -p "$SIG_DIR"

for f in $(find "$ROOT" -name '*.bpmn' | sort); do
  hash=$(sha384sum "$f" | awk '{print $1}')
  rel="${f#$ROOT/}"
  out="$SIG_DIR/$rel.sig"
  mkdir -p "$(dirname "$out")"

  echo "🔏 Signing $f  (sha384:$hash)"
  curl --fail -sS -X POST "$PKI_ENDPOINT" \
    -H "Authorization: Bearer $WGO_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"algo\":\"rsa-pss-sha384\",\"hash\":\"$hash\",\"purpose\":\"bpmn-deployment\",\"file\":\"$rel\"}" \
    | jq -r '.signature' > "$out"
  echo "   → $out"
done

echo "✅ All BPMN signed."
