# SNISID Enterprise Platform Bootstrap (Windows)

Write-Host "🚀 Initializing SNISID Enterprise DevOps Stack..." -ForegroundColor Cyan

# 1. Cluster Creation
$clusterExists = k3d cluster list | Select-String "snisid"
if ($clusterExists) {
    Write-Host "⚠️ Cluster already exists. Skipping creation."
} else {
    Write-Host "📦 Creating k3d cluster with LoadBalancer ports..."
    k3d cluster create snisid --port "80:80@loadbalancer" --port "443:443@loadbalancer"
}

# 2. Namespace setup
kubectl create namespace snisid --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -

# 3. Install ArgoCD
Write-Host "⚓ Installing ArgoCD..."
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 4. Install Nginx Ingress
Write-Host "🌐 Installing Nginx Ingress Controller..."
helm upgrade --install ingress-nginx ingress-nginx --repo https://kubernetes.github.io/ingress-nginx --namespace ingress-nginx --create-namespace

# 5. Deploy Monitoring (LGTM Stack)
Write-Host "📊 Deploying Observability Stack (Prometheus/Grafana/Loki)..."
helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack --namespace monitoring

# 6. Deploy SNISID via ArgoCD (GitOps)
Write-Host "🚀 Bootstrapping SNISID Application via ArgoCD..."
kubectl apply -f ./deploy/argocd/application.yaml

Write-Host "✅ SNISID DevOps Stack is ready!" -ForegroundColor Green
Write-Host "📍 ArgoCD UI: http://localhost:80/argocd (Forward port 8080 if needed)"
Write-Host "📍 Grafana: http://localhost:80/grafana"
Write-Host "📍 SNISID UI: http://snisid.local (Map to 127.0.0.1 in hosts)"
