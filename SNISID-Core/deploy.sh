#!/usr/bin/env bash
# ============================================================
# SNISID-Core — Platform Makefile Equivalent (Bash)
# Bootstrap + Deploy + Test + Health Check
# Document ID: SNISID-DEPLOY-001
# Version: 1.0.0
# ============================================================

set -euo pipefail

SNISID_ENV="${ENV:-staging}"
KUBECONFIG_PATH="${KUBECONFIG:-$HOME/.kube/snisid-${SNISID_ENV}.yaml}"
REGISTRY="harbor.snisid.gov.ht"
ARGOCD_URL="https://argocd.snisid.gov.ht"
VAULT_URL="https://vault.snisid.gov.ht"

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() { echo -e "${BLUE}[SNISID]${NC} $1"; }
success() { echo -e "${GREEN}[✅ OK]${NC} $1"; }
warn() { echo -e "${YELLOW}[⚠️  WARN]${NC} $1"; }
error() { echo -e "${RED}[❌ ERROR]${NC} $1"; exit 1; }

# ============================================================
# SECTION 1: PRÉ-REQUIS
# ============================================================
check_prereqs() {
    log "=== Vérification des prérequis SNISID-Core ==="
    
    local tools=("kubectl" "helm" "argocd" "vault" "k3sup" "kubeseal" "cosign" "openssl")
    local missing=()
    
    for tool in "${tools[@]}"; do
        if command -v "$tool" &>/dev/null; then
            success "$tool: $(${tool} version 2>/dev/null | head -1 || echo 'OK')"
        else
            missing+=("$tool")
            warn "$tool: NON TROUVÉ"
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        error "Outils manquants: ${missing[*]}. Installer avant de continuer."
    fi
    
    # Vérifier kubectl connectivity
    if ! kubectl --kubeconfig="$KUBECONFIG_PATH" cluster-info &>/dev/null; then
        error "Impossible de se connecter au cluster K8s ($KUBECONFIG_PATH)"
    fi
    
    success "Tous les prérequis satisfaits!"
}

# ============================================================
# SECTION 2: BOOTSTRAP CLUSTER RKE2
# ============================================================
bootstrap_cluster() {
    log "=== Bootstrap RKE2 Cluster (env: $SNISID_ENV) ==="
    
    if [ -z "${RKE2_MASTER_IP:-}" ]; then
        error "Variable RKE2_MASTER_IP requise"
    fi
    
    log "Installation RKE2 sur $RKE2_MASTER_IP..."
    
    # Installer RKE2 master node
    ssh -o StrictHostKeyChecking=no root@"$RKE2_MASTER_IP" << 'ENDSSH'
        # Configurer le kernel
        cat >> /etc/sysctl.d/99-snisid.conf << 'EOF'
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
fs.inotify.max_user_watches = 524288
fs.inotify.max_user_instances = 512
vm.swappiness = 0
EOF
        sysctl --system

        # Désactiver swap
        swapoff -a
        sed -i '/ swap / s/^/#/' /etc/fstab

        # Installer RKE2
        curl -sfL https://get.rke2.io | INSTALL_RKE2_TYPE=server sh -

        # Configurer RKE2
        mkdir -p /etc/rancher/rke2
        cat > /etc/rancher/rke2/config.yaml << 'CONFIG'
cni: cilium
disable-kube-proxy: true
etcd-expose-metrics: true
protect-kernel-defaults: true
kube-apiserver-arg:
  - audit-log-path=/var/lib/rancher/rke2/server/logs/audit.log
  - audit-log-maxsize=100
  - audit-log-maxbackup=10
  - audit-log-maxage=30
  - audit-policy-file=/etc/rancher/rke2/audit-policy.yaml
  - encryption-provider-config=/etc/rancher/rke2/encryption.yaml
  - enable-admission-plugins=NodeRestriction,PodSecurity
kubelet-arg:
  - protect-kernel-defaults=true
  - event-qps=0
  - anonymous-auth=false
  - authorization-mode=Webhook
  - client-ca-file=/var/lib/rancher/rke2/server/tls/client-ca.crt
CONFIG

        # Politique d'audit Kubernetes
        cat > /etc/rancher/rke2/audit-policy.yaml << 'AUDIT'
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
  - level: Metadata
    namespaces: ["snisid-identity", "snisid-biometrics", "snisid-civil-registry"]
  - level: Request
    resources:
    - group: ""
      resources: ["secrets", "configmaps"]
  - level: None
    users: ["system:kube-proxy"]
    verbs: ["watch"]
  - level: Metadata
    omitStages: ["RequestReceived"]
AUDIT

        systemctl enable rke2-server
        systemctl start rke2-server
        
        echo "✅ RKE2 installé et démarré"
        sleep 30
        /var/lib/rancher/rke2/bin/kubectl --kubeconfig /etc/rancher/rke2/rke2.yaml get nodes
ENDSSH

    success "RKE2 bootstrap terminé sur $RKE2_MASTER_IP"
    
    # Récupérer le kubeconfig
    scp root@"$RKE2_MASTER_IP":/etc/rancher/rke2/rke2.yaml "$KUBECONFIG_PATH"
    sed -i "s/127.0.0.1/$RKE2_MASTER_IP/g" "$KUBECONFIG_PATH"
    chmod 600 "$KUBECONFIG_PATH"
    
    success "Kubeconfig sauvegardé dans $KUBECONFIG_PATH"
}

# ============================================================
# SECTION 3: DÉPLOIEMENT CORE
# ============================================================
deploy_core() {
    log "=== Déploiement SNISID-Core (env: $SNISID_ENV) ==="
    
    # Étape 1: Namespaces & Base K8s
    log "[1/10] Namespaces et base K8s..."
    kubectl --kubeconfig="$KUBECONFIG_PATH" apply -f Kubernetes/base/
    success "Namespaces créés"
    
    # Étape 2: Cert-Manager
    log "[2/10] Cert-Manager..."
    helm repo add jetstack https://charts.jetstack.io --force-update
    helm upgrade --install cert-manager jetstack/cert-manager \
        --namespace cert-manager \
        --create-namespace \
        --set installCRDs=true \
        --wait --timeout 5m
    success "Cert-Manager installé"
    
    # Étape 3: ArgoCD
    log "[3/10] ArgoCD..."
    helm repo add argo https://argoproj.github.io/argo-helm --force-update
    helm upgrade --install argocd argo/argo-cd \
        --namespace argocd \
        --create-namespace \
        --set configs.params."server\.insecure"=false \
        --set server.extraArgs="{--insecure}" \
        --wait --timeout 5m
    
    # Attendre ArgoCD
    kubectl --kubeconfig="$KUBECONFIG_PATH" wait \
        --for=condition=Ready pod \
        -l app.kubernetes.io/name=argocd-server \
        -n argocd \
        --timeout=300s
    success "ArgoCD installé"
    
    # Étape 4: Appliquer App-of-Apps
    log "[4/10] App-of-Apps GitOps..."
    kubectl --kubeconfig="$KUBECONFIG_PATH" apply -f GitOps/argocd/app-of-apps.yaml
    success "App-of-Apps créé — ArgoCD prend le relais"
    
    log "=== Phase de déploiement ArgoCD en cours ==="
    log "⏱️  Attente synchronisation (peut prendre 10-20 minutes)..."
    
    # Surveiller le déploiement
    monitor_deployment
}

# ============================================================
# SECTION 4: HEALTH CHECK
# ============================================================
health_check() {
    log "=== Health Check SNISID-Core (env: $SNISID_ENV) ==="
    
    local failures=0
    
    # Vérifier les namespaces
    for ns in snisid-identity snisid-biometrics snisid-civil-registry snisid-api-gateway snisid-security snisid-event-bus snisid-databases snisid-observability; do
        if kubectl --kubeconfig="$KUBECONFIG_PATH" get namespace "$ns" &>/dev/null; then
            success "Namespace $ns: OK"
        else
            warn "Namespace $ns: MANQUANT"
            ((failures++))
        fi
    done
    
    # Vérifier les pods critiques
    log "Vérification pods..."
    local critical_pods=(
        "snisid-identity:identity-service"
        "snisid-api-gateway:kong"
        "snisid-security:vault"
        "snisid-security:keycloak"
        "snisid-event-bus:snisid-kafka-kafka"
        "snisid-observability:prometheus"
        "argocd:argocd-server"
    )
    
    for pod_spec in "${critical_pods[@]}"; do
        ns="${pod_spec%%:*}"
        label="${pod_spec##*:}"
        
        ready=$(kubectl --kubeconfig="$KUBECONFIG_PATH" \
            get pods -n "$ns" -l "app=$label" \
            -o jsonpath='{.items[*].status.containerStatuses[0].ready}' 2>/dev/null | tr ' ' '\n' | grep -c "true" || echo "0")
        
        total=$(kubectl --kubeconfig="$KUBECONFIG_PATH" \
            get pods -n "$ns" -l "app=$label" \
            --no-headers 2>/dev/null | wc -l || echo "0")
        
        if [ "$ready" -gt 0 ]; then
            success "$ns/$label: $ready/$total pods ready"
        else
            warn "$ns/$label: $ready/$total pods ready"
            ((failures++))
        fi
    done
    
    # Tester l'API Gateway
    log "Test API Gateway..."
    if curl -sf -o /dev/null -w "%{http_code}" https://api.snisid.gov.ht/health 2>/dev/null | grep -q "200"; then
        success "API Gateway: accessible (HTTP 200)"
    else
        warn "API Gateway: non accessible depuis l'extérieur"
    fi
    
    # Résumé
    echo ""
    if [ "$failures" -eq 0 ]; then
        success "=== Health Check PASSÉ — SNISID-Core est opérationnel ==="
    else
        warn "=== Health Check: $failures problèmes détectés ==="
    fi
}

# ============================================================
# SECTION 5: LOAD TEST
# ============================================================
load_test() {
    log "=== Load Test SNISID Identity API ==="
    log "Outil: k6 — Target: https://api.snisid.gov.ht"
    
    cat > /tmp/snisid-k6-test.js << 'K6EOF'
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '2m', target: 100 },    // Montée progressive
    { duration: '5m', target: 500 },    // Charge nominale
    { duration: '2m', target: 1000 },   // Pic de charge
    { duration: '2m', target: 500 },    // Descente
    { duration: '1m', target: 0 },      // Arrêt
  ],
  thresholds: {
    http_req_duration: ['p(99)<2000'],  // P99 < 2s
    http_req_failed: ['rate<0.01'],     // Erreurs < 1%
    errors: ['rate<0.01'],
  },
};

const BASE_URL = 'https://api.snisid.gov.ht/v1';

export default function () {
  // Test: Vérification d'identité (cas le plus fréquent)
  const verifyRes = http.post(
    `${BASE_URL}/verify/identity`,
    JSON.stringify({
      niu: '7392851046',
      requested_level: 1,
      purpose: 'load_test',
    }),
    {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${__ENV.TEST_TOKEN}`,
      },
    }
  );
  
  const verifyOK = check(verifyRes, {
    'verify: status 200': (r) => r.status === 200,
    'verify: latency < 500ms': (r) => r.timings.duration < 500,
    'verify: has decision field': (r) => r.json('decision') !== undefined,
  });
  
  errorRate.add(!verifyOK);
  
  sleep(0.1);  // 100ms entre requêtes
}
K6EOF

    if command -v k6 &>/dev/null; then
        k6 run /tmp/snisid-k6-test.js \
            --env TEST_TOKEN="${TEST_TOKEN:-placeholder}" \
            --out json=/tmp/k6-results.json
        success "Load test terminé. Résultats: /tmp/k6-results.json"
    else
        warn "k6 non installé. Installer depuis: https://k6.io/docs/get-started/installation/"
        warn "Fichier de test sauvegardé: /tmp/snisid-k6-test.js"
    fi
}

monitor_deployment() {
    log "Surveillance du déploiement ArgoCD..."
    for i in $(seq 1 60); do
        sync_status=$(argocd app get snisid-platform \
            --server "$ARGOCD_URL" \
            --auth-token "${ARGOCD_TOKEN:-}" \
            --output json 2>/dev/null | jq -r '.status.sync.status' 2>/dev/null || echo "UNKNOWN")
        
        health_status=$(argocd app get snisid-platform \
            --server "$ARGOCD_URL" \
            --output json 2>/dev/null | jq -r '.status.health.status' 2>/dev/null || echo "UNKNOWN")
        
        echo "  [Tentative $i/60] Sync: $sync_status | Health: $health_status"
        
        if [ "$sync_status" = "Synced" ] && [ "$health_status" = "Healthy" ]; then
            success "ArgoCD: Synced & Healthy!"
            return 0
        fi
        sleep 30
    done
    warn "Timeout — vérifier manuellement: argocd app get snisid-platform"
}

# ============================================================
# SECTION 6: COMMANDES PRINCIPALES
# ============================================================
usage() {
    echo "SNISID-Core Platform Deployment Tool"
    echo ""
    echo "Usage: $0 <command> [ENV=staging|prod]"
    echo ""
    echo "Commandes:"
    echo "  check-prereqs    Vérifier les outils requis"
    echo "  bootstrap        Bootstrap cluster RKE2"
    echo "  deploy-core      Déployer l'infrastructure core"
    echo "  health-check     Vérifier l'état de la plateforme"
    echo "  load-test        Exécuter le test de charge k6"
    echo ""
    echo "Exemples:"
    echo "  ENV=staging ./deploy.sh check-prereqs"
    echo "  ENV=prod RKE2_MASTER_IP=10.0.1.10 ./deploy.sh bootstrap"
    echo "  ENV=prod ./deploy.sh deploy-core"
    echo "  ENV=prod ./deploy.sh health-check"
}

# Main
case "${1:-}" in
    check-prereqs) check_prereqs ;;
    bootstrap) bootstrap_cluster ;;
    deploy-core) deploy_core ;;
    health-check) health_check ;;
    load-test) load_test ;;
    *) usage ;;
esac
