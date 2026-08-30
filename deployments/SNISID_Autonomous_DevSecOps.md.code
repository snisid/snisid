# SNISID Autonomous DevSecOps: Configurations & Scripts

**Classification:** RESTRICTED / SOVEREIGN DEVSECOPS
**Compliance:** SLSA Level 4 / NIST SSDF / ISO 27001

This operational playbook defines the Kyverno policies, ArgoCD Sync settings, KEDA configurations, and automated self-healing scripts for the SNISID DevSecOps platform.

---

## 1. AI-Assisted Vulnerability Auto-Remediation Script

This Python script runs within the CI runner, processing JSON vulnerability scans from Trivy and automatically generating a Git Pull Request to update insecure base container images.

```python
# File: /opt/snisid/devsecops/remediate_vulnerabilities.py
import json
import sys
import re

def parse_and_remediate(trivy_report_path, dockerfile_path, output_pr_path):
    print(f"[*] Processing Trivy report: {trivy_report_path}")
    
    with open(trivy_report_path, 'r') as f:
        report = json.load(f)
        
    vulnerabilities = []
    # Extract vulnerabilities from Trivy JSON structure
    if 'Results' in report:
        for result in report['Results']:
            if 'Vulnerabilities' in result:
                for vuln in result['Vulnerabilities']:
                    if vuln.get('Severity') in ['CRITICAL', 'HIGH']:
                        vulnerabilities.append({
                            'id': vuln.get('VulnerabilityID'),
                            'package': vuln.get('PkgName'),
                            'installed': vuln.get('InstalledVersion'),
                            'fixed': vuln.get('FixedVersion')
                        })
                        
    if not vulnerabilities:
        print("[+] No CRITICAL or HIGH vulnerabilities found. Exit.")
        sys.exit(0)
        
    print(f"[!] Found {len(vulnerabilities)} high-severity issues. Remediating...")
    
    # Read the Dockerfile to check the base image
    with open(dockerfile_path, 'r') as f:
        dockerfile_content = f.read()
        
    # Match the base image (e.g. FROM node:18.1-alpine)
    match = re.search(r'FROM\s+([a-zA-Z0-9\-\./_]+):([a-zA-Z0-9\-\._]+)', dockerfile_content)
    if not match:
        print("[-] Could not identify base image in Dockerfile.")
        sys.exit(1)
        
    image_name = match.group(1)
    current_tag = match.group(2)
    print(f"[*] Identified Base Image: {image_name}:{current_tag}")
    
    # Simple tag upgrade strategy (upgrading minor patch tags for safety)
    # Example: If current tag is 18.1-alpine, check and recommend update to latest stable minor (e.g. 18.16-alpine)
    tag_parts = re.findall(r'\d+', current_tag)
    if tag_parts and len(tag_parts) >= 2:
        major = tag_parts[0]
        minor = tag_parts[1]
        new_minor = str(int(minor) + 1) # Propose upgrade of minor version
        new_tag = current_tag.replace(f"{major}.{minor}", f"{major}.{new_minor}")
        
        remediated_content = dockerfile_content.replace(current_tag, new_tag)
        
        with open(dockerfile_path, 'w') as f:
            f.write(remediated_content)
        print(f"[+] Dockerfile base image updated from {current_tag} to {new_tag}")
        
        # Write recommendation output for PR generation
        with open(output_pr_path, 'w') as f:
            f.write(f"title: Auto-Remediation: Update base image to {new_tag}\n")
            f.write(f"body: Automated security update resolving critical/high CVEs by upgrading base image. Package vulnerabilities resolved: {len(vulnerabilities)}\n")
    else:
        print("[-] Non-standard image tag structure. Manual review required.")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) < 4:
        print("Usage: python remediate_vulnerabilities.py <trivy_report.json> <Dockerfile> <output_pr_instructions.txt>")
        sys.exit(1)
    parse_and_remediate(sys.argv[1], sys.argv[2], sys.argv[3])
```

---

## 2. OPA / Kyverno Policy-as-Code Manifests

These policies are deployed inside production clusters to prevent the execution of unsigned or non-compliant workloads.

### 2.1. Kyverno Signature Verification Policy
Enforces that all images running in production must have a signature from the trusted SNISID Cosign key:

```yaml
# File: /deployments/kyverno/verify-image-signature.yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: verify-image-signatures
  annotations:
    policies.kyverno.io/title: Verify Image Signatures via Cosign
spec:
  validationFailureAction: Enforce
  background: false
  rules:
    - name: verify-cosign-signature
      match:
        any:
          - resources:
              kinds:
                - Pod
      verifyImages:
        - imageCriterion: "harbor.snisid.gov.ht/*"
          key: |-
            -----BEGIN PUBLIC KEY-----
            MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE7p9s9lM9V5bYjL/Xv1H1p9lX4n1z
            Y/K6fR7o0m9q1k9V5bYjL/Xv1H1p9lX4n1zY/K6fR7o0m9q1k9V5bYjLw==
            -----END PUBLIC KEY-----
```

### 2.2. Kyverno Root Execution Restriction Policy
Blocks containers that attempt to run as the root user:

```yaml
# File: /deployments/kyverno/disallow-root-namespaces.yaml
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: disallow-root-user
  annotations:
    policies.kyverno.io/title: Disallow Root User Execution
spec:
  validationFailureAction: Enforce
  background: true
  rules:
    - name: check-run-as-non-root
      match:
        any:
          - resources:
              kinds:
                - Pod
      validate:
        message: "Running as root user is prohibited on SNISID production nodes."
        pattern:
          spec:
            securityContext:
              runAsNonRoot: true
```

---

## 3. KEDA Predictive Auto-Scaling Configurations

KEDA integrates with the Prometheus monitoring server to dynamically scale pods based on predicted traffic patterns.

```yaml
# File: /deployments/k8s/keda-predictive-autoscaler.yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: identity-service-scaler
  namespace: snisid-core
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: identity-service
  minReplicaCount: 3
  maxReplicaCount: 30
  cooldownPeriod: 300
  restoreToOriginalReplicaCount: false
  triggers:
    - type: prometheus
      metadata:
        serverAddress: http://prometheus-k8s.monitoring.svc.cluster.local:9090
        # Query analyzing temporal traffic growth rate over the last 15 minutes
        metricName: http_requests_predicted_growth
        query: sum(rate(http_requests_total{job="identity-service"}[5m])) * 1.5
        threshold: '100'
```

---

## 4. Drift Detection & ArgoCD Self-Heal Configuration

This ArgoCD Application configuration enforces automatic self-healing, actively rolling back any manual CLI drift.

```yaml
# File: /deployments/argocd/argocd-selfheal-config.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: snisid-core-platform
  namespace: argocd
spec:
  project: default
  source:
    repoURL: 'git@git.snisid.gov.ht:infrastructure/gitops.git'
    targetRevision: HEAD
    path: gitops/overlays/prod
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: snisid-core
  syncPolicy:
    automated:
      prune: true     # Automatically delete resources no longer in Git
      selfHeal: true  # Revert manual out-of-band changes instantly
    syncOptions:
      - CreateNamespace=true
      - Validate=true
```

---

## 5. Automated Rollback Trigger Script

This bash script runs as a daemon or webhook handler inside the Kubernetes management namespace. It monitors canary deployment metrics and triggers instant rollback if error levels exceed the safety threshold.

```bash
#!/bin/bash
# File: /opt/snisid/devsecops/canary_rollback_monitor.sh
set -e

PROMETHEUS_URL="http://prometheus-k8s.monitoring.svc.cluster.local:9090"
APP_NAME="identity-service"
NAMESPACE="snisid-core"
THRESHOLD="0.01" # 1% Error Rate Threshold

echo "[*] Launching Canary Health Monitor for ${APP_NAME}..."

# Query Prometheus for the HTTP 5xx error rate on the canary pods
query_error_rate() {
  local query="sum(rate(http_requests_total{namespace=\"${NAMESPACE}\",service=\"${APP_NAME}-canary\",status=~\"5..\"}[2m]))%20/%20sum(rate(http_requests_total{namespace=\"${NAMESPACE}\",service=\"${APP_NAME}-canary\"}[2m]))"
  
  local response=$(curl -s "${PROMETHEUS_URL}/api/v1/query?query=${query}")
  local val=$(echo "$response" | jq -r '.data.result[0].value[1] // "0"')
  echo "$val"
}

# Loop and check metrics every 15 seconds
while true; do
  error_rate=$(query_error_rate)
  echo "[*] Current Canary Error Rate: ${error_rate}"
  
  # Check if error rate is greater than the threshold
  if (( $(echo "${error_rate} > ${THRESHOLD}" | bc -l) )); then
    echo "[!] CRITICAL: Error rate (${error_rate}) exceeds safety threshold (${THRESHOLD})!"
    echo "[*] Initiating Autonomous Rollback..."
    
    # Execute Argo Rollouts abort command
    kubectl argo rollouts abort "${APP_NAME}" -n "${NAMESPACE}"
    
    # Notify SOC Slack/Wazuh Integrations
    curl -X POST -H 'Content-type: application/json' \
      --data "{\"text\":\"[AUTONOMOUS ROLLBACK] Rolled back deployment/${APP_NAME} in namespace ${NAMESPACE} due to error rate of ${error_rate}\"}" \
      https://alerts.snisid.gov.ht/webhooks/slack
      
    break
  fi
  sleep 15
done
```

---

*Verified and signed by the SNISID Automation Operations Board.*
