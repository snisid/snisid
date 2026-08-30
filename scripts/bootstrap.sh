#!/bin/bash

# SNISID Production Cluster Bootstrapper
# Sets up K3d, namespaces, and core services.

set -e

echo "----------------------------------------------------"
echo "🌐 BOOTSTRAPPING SNISID PRODUCTION PLATFORM"
echo "----------------------------------------------------"

# 1. Create Cluster
if ! k3d cluster get snisid > /dev/null 2>&1; then
    echo "🏗️ Creating K3d cluster 'snisid'..."
    k3d cluster create snisid --servers 3 --agents 2
fi

# 2. Create Namespaces
echo "📦 Creating Namespaces (SOC, SIM, DATA)..."
kubectl create ns snisid-soc || true
kubectl create ns snisid-sim || true
kubectl create ns snisid-data || true

# 3. Install ArgoCD
echo "🔄 Installing ArgoCD Control Plane..."
kubectl create namespace argocd || true
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 4. Apply GitOps Root
echo "🚀 Applying GitOps App-of-Apps..."
kubectl apply -f deployments/gitops/app-of-apps.yaml

echo "----------------------------------------------------"
echo "✅ SYSTEM BOOTSTRAPPED: ACCESS ARGOCD AT http://localhost:8080"
echo "----------------------------------------------------"
