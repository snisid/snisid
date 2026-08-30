#!/bin/bash
# SNISID v1.0 — Production Build System

echo "🚀 Starting SNISID Production Build..."

# 1. Clean previous artifacts
rm -rf ./dist

# 2. Build Go Microservices
services=("iam" "siem" "gateway" "federation")
for s in "${services[@]}"; do
    echo "🏗️ Building service: $s"
    go build -o "./dist/$s" "./services/$s"
done

# 3. Build UI
echo "🎨 Building SOC Dashboard UI..."
cd frontend/soc-dashboard && npm run build && cd ../..

# 4. Finalize
echo "✅ SNISID v1.0 Build Complete. Artifacts located in /dist"
