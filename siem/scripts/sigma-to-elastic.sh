#!/bin/bash
# SNISID Script: Compile generic Sigma rules into Elastic Detection Engine format
# Requires: sigma-cli (pip install sigmacli)

set -e

SIGMA_RULES_DIR="../../soc/rules/sigma"
OUTPUT_DIR="./compiled_elastic_rules"
KIBANA_URL="https://kibana.snisid-internal.svc.cluster.local:5601"
# Ensure API_KEY is passed securely via env vars

mkdir -p "$OUTPUT_DIR"

echo "[*] Converting Sigma rules to Elastic Query Language (EQL)..."

# Use the official sigmacli with the elasticsearch backend and ECS mapping
sigma convert \
  --target elasticsearch \
  --pipeline ecs_windows \
  --format detection_rule \
  --output "$OUTPUT_DIR/elastic_rules.ndjson" \
  "$SIGMA_RULES_DIR"/*.yml

echo "[*] Conversion complete."

if [ -n "$KIBANA_API_KEY" ]; then
    echo "[*] Uploading rules to Kibana Detection Engine..."
    curl -X POST "$KIBANA_URL/api/detection_engine/rules/_import" \
      -H "kbn-xsrf: true" \
      -H "Authorization: ApiKey $KIBANA_API_KEY" \
      -F file=@"$OUTPUT_DIR/elastic_rules.ndjson"
    echo -e "\n[*] Upload complete."
else
    echo "[!] KIBANA_API_KEY not set. Skipping auto-upload."
    echo "You can manually import $OUTPUT_DIR/elastic_rules.ndjson into Kibana."
fi
