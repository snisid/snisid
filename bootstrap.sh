#!/bin/bash
# SNISID Platform Master Bootstrap Script
# Version: 1.0.0

set -e

echo "🚀 Initializing SNISID National Identity & Intelligence Platform..."

# 1. Prerequisite Check
command -v docker >/dev/null 2>&1 || { echo >&2 "❌ Error: Docker is required. Aborting."; exit 1; }
command -v k3d >/dev/null 2>&1 || { echo >&2 "❌ Error: k3d is required. Aborting."; exit 1; }
command -v helm >/dev/null 2>&1 || { echo >&2 "❌ Error: Helm is required. Aborting."; exit 1; }

# 2. Create k3d Cluster
echo "🏗️ Creating local k3d cluster 'snisid'..."
k3d cluster create snisid --api-port 6443 -p "80:80@loadbalancer" -p "443:443@loadbalancer" --agents 2

# 3. Create Namespaces
echo "☸️ Creating namespaces..."
kubectl create namespace snisid
kubectl create namespace monitoring
kubectl create namespace vault

# 4. Deploy Infrastructure (Local Mode)
echo "🚇 Deploying core infrastructure (Postgres, Kafka, Neo4j, Redis)..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install snisid-db bitnami/postgresql -n snisid --set auth.database=snisid --set auth.password=snisidpass
helm install snisid-kafka bitnami/kafka -n snisid
helm install snisid-redis bitnami/redis -n snisid

# 5. Build and Import Images
echo "🐳 Building SNISID microservices..."
make docker-build
k3d image import snisid/api-gateway:latest snisid/identity-api:latest snisid/ai-face:latest -c snisid

# 6. Deploy SNISID Platform
echo "🚀 Deploying SNISID services..."
make k8s-deploy

# 7. Final Status
echo "✅ SNISID Platform is bootstrapping!"
echo "📍 Access the Dashboard at: http://localhost"
echo "📊 Monitoring available at: http://grafana.localhost (local-only)"
echo "--------------------------------------------------------"
kubectl get pods -n snisid
