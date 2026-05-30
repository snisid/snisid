#!/bin/bash
# MISP Bootstrap Configuration for SNISID SOC

MISP_URL="https://misp.snisid-security.svc.cluster.local"
AUTH_KEY=$1

if [ -z "$AUTH_KEY" ]; then
  echo "Usage: ./misp-config.sh <API_KEY>"
  exit 1
fi

# 1. Enable default feeds (e.g., CIRCL OSINT)
curl -s -k -X POST "$MISP_URL/feeds/enable/1" \
  -H "Authorization: $AUTH_KEY" \
  -H "Accept: application/json"

# 2. Setup Elasticsearch integration for IOC export
curl -s -k -X POST "$MISP_URL/servers/edit" \
  -H "Authorization: $AUTH_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "Server": {
      "name": "SNISID_Elastic_SIEM",
      "url": "https://elasticsearch-master:9200",
      "push": true,
      "pull": false
    }
  }'

echo "MISP Configuration applied successfully."
