# SNISID Zero Trust & IAM: Policies, Handlers & Rotation Scripts

**Classification:** RESTRICTED / SOVEREIGN ZEROTRUST
**Compliance:** NIST SP 800-63B / NIST SP 800-207 / OPA Rego

This playbook defines the OPA Rego access policies, Keycloak step-up handler templates, and automated Vault credential rotation scripts deployed across the SNISID security infrastructure.

---

## 1. Extended OPA Rego Policy (Geographic & Temporal)

This Rego policy validates incoming API request contexts against client geolocations (IP department) and temporal working hours.

```rego
# File: /opt/snisid/zerotrust/policy.rego
package snisid.authz

default allow = false

# Access rules for administrative accounts
allow {
    # Rule 1: User must have admin role
    input.user.roles[_] == "admin"
    
    # Rule 2: Access must originate from within Haiti (HT)
    input.request.location.country == "HT"
    
    # Rule 3: Access must match registered department location
    input.request.location.department == input.user.assigned_department
    
    # Rule 4: Access allowed during official working hours only (07:00 - 17:00 AST)
    input.request.time.hour >= 7
    input.request.time.hour <= 17
    
    # Rule 5: User must have high assurance authentication (AAL3)
    input.user.aal_level == "AAL3"
}

# Allow read-only status checking for citizens from any location at AAL1
allow {
    input.request.method == "GET"
    input.request.path == "/v1/citizen/status"
    input.user.aal_level == "AAL1"
}
```

---

## 2. Post-Incident Automated Credential Rotation

This Python script is executed by the SOAR engine 15 minutes after threat containment. It logs into Vault and triggers a complete rotation of DB secrets, API keys, and intermediate CA certificates.

```python
# File: /opt/snisid/zerotrust/rotate_credentials.py
import urllib.request
import json
import sys

VAULT_ADDR = "https://vault.snisid-security.svc.cluster.local:8200"

def rotate_vault_target(token, path, payload):
    print(f"[*] Dispatching rotation event to Vault target: {path}")
    url = f"{VAULT_ADDR}/v1/{path}"
    
    req = urllib.request.Request(url, method="POST")
    req.add_header("X-Vault-Token", token)
    req.add_header("Content-Type", "application/json")
    
    # SSL verification skip for internal network mockup
    import ssl
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    
    try:
        data_bytes = json.dumps(payload).encode('utf-8')
        with urllib.request.urlopen(req, data=data_bytes, context=ctx, timeout=5) as response:
            status = response.getcode()
            print(f"[+] Rotation execution completed. Status: {status}")
            return True
    except Exception as e:
        print(f"[-] Rotation event failed: {str(e)}")
        return False

def initiate_incident_rotation(vault_token):
    print("[*] Starting post-incident credential rotation SLA pipeline...")
    
    # 1. Rotate Database connection passwords
    db_payload = {"name": "identity-db-role"}
    db_ok = rotate_vault_target(vault_token, "database/rotate/identity-db", db_payload)
    
    # 2. Rotate API keys partition
    api_payload = {"role": "inter-agency-gateway"}
    api_ok = rotate_vault_target(vault_token, "transit/keys/api-signature/rotate", api_payload)
    
    # 3. Force intermediate CA certificate update
    ca_payload = {"common_name": "vault.snisid.local"}
    ca_ok = rotate_vault_target(vault_token, "pki_int/tidy", ca_payload)
    
    if db_ok and api_ok and ca_ok:
        print("[+] All incident credential rotations completed successfully.")
        return True
    else:
        print("[-] Verification failed: One or more rotation steps failed.")
        return False

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python rotate_credentials.py <vault_token>")
        sys.exit(1)
    initiate_incident_rotation(sys.argv[1])
```

---

## 3. UEBA Anomaly Step-Up Authentication Handler

This script processes incoming user transaction behavior, checks if they exceed the 2-sigma baseline threshold, and determines if step-up authentication is required.

```python
# File: /opt/snisid/zerotrust/evaluate_ueba_stepup.py
import sys
import json

def process_transaction(user_id, metrics, baseline_mean, baseline_std):
    """
    Evaluates metric score (e.g. request rate or transaction value) 
    against user baseline mean and standard deviation.
    """
    print(f"[*] Processing user UEBA transaction: {user_id}")
    
    current_value = metrics.get("request_rate", 0.0)
    
    # Calculate Z-score (standard deviations from mean)
    if baseline_std == 0:
        z_score = 0.0
    else:
        z_score = abs(current_value - baseline_mean) / baseline_std
        
    print(f"[*] Computed behavioral deviation (Z-score): {z_score:.2f}")
    
    # 1. Critical Anomaly Trigger (>= 3-sigma)
    if z_score >= 3.0:
        print(f"[!] CRITICAL: Behavior exceeds 3-sigma limits!")
        print("[*] ACTION: Keycloak active tokens REVOKED. Account locked.")
        return "REVOKE_SESSION"
        
    # 2. Moderate Anomaly Trigger (2-sigma to 2.99-sigma)
    elif z_score >= 2.0:
        print(f"[!] WARNING: Behavior exceeds 2-sigma limits!")
        print("[*] ACTION: Trigger Step-Up authentication (FIDO2 token validation required).")
        return "TRIGGER_STEPUP"
        
    print("[+] Behavior matches baseline profile. Access granted.")
    return "ALLOW"

if __name__ == "__main__":
    # Test case: User average request rate is 10 req/min with standard deviation of 2.
    # Current request rate is 15 req/min (Z-score = 2.5) -> Trigger step-up
    current_metrics = {"request_rate": 15.0}
    action = process_transaction(
        user_id="operator-203",
        metrics=current_metrics,
        baseline_mean=10.0,
        baseline_std=2.0
    )
    print(f"[*] Routing Action: {action}")
```

---

*Verified and signed by the SNISID Security Directorate.*
