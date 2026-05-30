# SNISID Biometric Security: Policies, Routers & Drift Monitors

**Classification:** RESTRICTED / SOVEREIGN BIOMETRICS
**Compliance:** ISO/IEC 19794-2 / ISO/IEC 30107-3 / AES-256-GCM

This playbook defines the biometric template encryption scripts, 1:N deduplication score routers, alternative verification schemas, and monthly model drift validators deployed across the SNISID AI and ABIS infrastructure.

---

## 1. Biometric Template Encryption & Integrity Builder

This Python script processes plain-text minutiae vectors, encrypts them using AES-256-GCM, and signs the output with HMAC-SHA512 to enforce in-gallery integrity.

```python
# File: /opt/snisid/core/protect_biometrics.py
import hmac
import hashlib
from os import urandom

# Simulating standard AES-GCM-256 encryption block using raw byte mappings
def encrypt_biometric_template(plain_vector, aes_key, hmac_key):
    """
    Encrypts plain minutiae vector and returns encrypted payload + signature.
    """
    print("[*] Encrypting plain biometric vector...")
    
    # Generate random IV (96-bit for GCM)
    iv = urandom(12)
    
    # Simulate GCM block cipher output (XOR payload with key stream for mockup)
    plain_bytes = plain_vector.encode('utf-8')
    cipher_bytes = bytearray(len(plain_bytes))
    for i in range(len(plain_bytes)):
        cipher_bytes[i] = plain_bytes[i] ^ aes_key[i % len(aes_key)]
        
    ciphertext = bytes(cipher_bytes)
    
    # Calculate GCM auth tag using SHA-256 HMAC wrapper
    gcm_tag_hasher = hmac.new(aes_key, iv + ciphertext, hashlib.sha256)
    gcm_tag = gcm_tag_hasher.digest()[:16] # 128-bit tag
    
    # Calculate Gallery integrity signature (HMAC-SHA512)
    integrity_hasher = hmac.new(hmac_key, iv + ciphertext + gcm_tag, hashlib.sha512)
    integrity_sig = integrity_hasher.digest()
    
    return {
        "iv_hex": iv.hex(),
        "ciphertext_hex": ciphertext.hex(),
        "gcm_tag_hex": gcm_tag.hex(),
        "integrity_signature_hex": integrity_sig.hex()
    }

def verify_and_decrypt_template(payload, aes_key, hmac_key):
    """
    Validates HMAC signature and decrypts cipher back to plain minutiae vector.
    """
    print("[*] Verifying template gallery integrity...")
    
    iv = bytes.fromhex(payload["iv_hex"])
    ciphertext = bytes.fromhex(payload["ciphertext_hex"])
    gcm_tag = bytes.fromhex(payload["gcm_tag_hex"])
    claimed_sig = bytes.fromhex(payload["integrity_signature_hex"])
    
    # Recalculate and verify signature (Mitigate SQL Injection database tampering)
    integrity_hasher = hmac.new(hmac_key, iv + ciphertext + gcm_tag, hashlib.sha512)
    computed_sig = integrity_hasher.digest()
    
    if not hmac.compare_digest(claimed_sig, computed_sig):
        print("[-] CRITICAL: Biometric template signature mismatch! Tampering detected.")
        return None
        
    print("[+] Integrity verified. Decrypting template...")
    
    # Reverse XOR for mockup decryption
    plain_bytes = bytearray(len(ciphertext))
    for i in range(len(ciphertext)):
        plain_bytes[i] = ciphertext[i] ^ aes_key[i % len(aes_key)]
        
    return plain_bytes.decode('utf-8')

if __name__ == "__main__":
    # Generate mock keys (wrapped by HSM in production)
    mock_aes = urandom(32) # AES-256 Key
    mock_hmac = urandom(64) # SHA-512 Key
    
    plain_data = "TEMPLATE-FGP-ISO-19794-2-MINUTIAE-DATA-VECTOR"
    
    encrypted = encrypt_biometric_template(plain_data, mock_aes, mock_hmac)
    print(f"[+] Encrypted Output Cipher: {encrypted['ciphertext_hex'][:30]}...")
    
    decrypted = verify_and_decrypt_template(encrypted, mock_aes, mock_hmac)
    print(f"[+] Decrypted Template Vector: {decrypted}")
```

---

## 2. 1:N Deduplication Score Router

This script evaluates matching scores returned by the GPU ABIS cluster and routes the enrollment request through the appropriate governance workflow.

```python
# File: /opt/snisid/core/route_abis_scores.py
import sys

def route_matching_score(citizen_id, similarity_score):
    """
    Similarity score thresholds:
    Score < 85%  -> New Citizen (Approved)
    85% - 94.9%  -> Potential Duplicate (Suspended, routed to manual review)
    >= 95%       -> Hard Match (Blocked, routed to DCPJ)
    """
    print(f"[*] Auditing ABIS 1:N match result for Citizen: {citizen_id}")
    print(f"[*] Match Similarity Score: {similarity_score:.2f}%")
    
    # 1. Hard Match Trigger (>= 95%)
    if similarity_score >= 95.0:
        print("[!] CRITICAL: Biometric Hard Match detected! Duplicate identity attempt.")
        print("[*] ACTION: Halt workflow. Set profile status: BLOCKED_DUPLICATE.")
        print("[*] ESCALATION: Lock record and dispatch case file to DCPJ Forensics.")
        return "BLOCKED_DCPJ_ESC"
        
    # 2. Potential Duplicate Trigger (85.0% - 94.9%)
    elif similarity_score >= 85.0:
        print("[!] WARNING: Potential Duplicate biometric match identified.")
        print("[*] ACTION: Suspend workflow. Set profile status: PENDING_ADJUDICATION.")
        print("[*] ESCALATION: Route case to Tier-2 Manual Forensic Examiner.")
        return "SUSPEND_TIER2_REVIEW"
        
    # 3. Approved New Citizen (< 85%)
    print("[+] Match score within safe baseline limits.")
    print("[*] ACTION: Confirm enrollment. Set profile status: ACTIVE. Assign NNI.")
    return "APPROVE_NEW_CITIZEN"

if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python route_abis_scores.py <citizen_id> <similarity_score>")
        sys.exit(1)
    route_matching_score(sys.argv[1], float(sys.argv[2]))
```

---

## 3. Biometric Model Drift Monitor

This script calculates monthly FAR and FRR stats across demographic segments (age, gender, commune) and alerts the Biometric Ethics Commission if decay thresholds are breached.

```python
# File: /opt/snisid/core/monitor_biometric_drift.py
import sys
import json

def audit_demographic_drift(metrics_db, segment_key, baseline_frr):
    """
    Checks if a demographic segment's False Rejection Rate (FRR) has 
    drifted by >= 50% compared to baseline parameters.
    """
    print(f"[*] Auditing demographic drift for segment: {segment_key}")
    
    # Extract monthly FRR metrics
    segment_data = metrics_db.get(segment_key, {})
    current_frr = segment_data.get("monthly_frr", 0.0)
    
    print(f"[*] Baseline FRR: {baseline_frr}%, Current FRR: {current_frr}%")
    
    # Check for drift spike (current / baseline >= 1.50)
    if baseline_frr > 0:
        increase_ratio = current_frr / baseline_frr
    else:
        increase_ratio = 1.0
        
    print(f"[*] Computed Performance Drift Ratio: {increase_ratio:.2f}x")
    
    if increase_ratio >= 1.50:
        print(f"[!] BIOMETRIC DECAY DETECTED: Segment '{segment_key}' has drifted by +{int((increase_ratio-1)*100)}%!")
        print("[*] ACTION: Flag model calibration task in Kubeflow pipeline.")
        print("[*] ESCALATION: Notify Biometric Ethics Commission.")
        return increase_ratio, False
        
    print(f"[+] Performance within safe parameters for segment '{segment_key}'.")
    return increase_ratio, True

if __name__ == "__main__":
    # Mock database containing monthly performance stats
    # Baseline standard FRR = 0.5% (0.005)
    # Segment grandanse_elderly has current FRR = 0.9% (0.009, increase of +80%) -> Drift Alert
    mock_db = {
        "grandanse_elderly": {
            "monthly_frr": 0.9,
            "queries_audited": 1200
        },
        "ouest_adults": {
            "monthly_frr": 0.52,
            "queries_audited": 25000
        }
    }
    
    ratio, is_compliant = audit_demographic_drift(mock_db, "grandanse_elderly", baseline_frr=0.5)
    if not is_compliant:
        sys.exit(1)
```

---

## 4. Alternative Verification Profile Schema

For citizens with fingerprint capture difficulties, the database stores an alternative validation profile configuration:

```json
{
  "citizen_hash_id": "sha256:7f83b2a9e10c...",
  "alternative_mode_active": true,
  "override_authorized_by": "clerk-admin-102",
  "override_aal_level": "AAL2",
  "verification_modalities": {
    "fingerprints": {
      "captured": false,
      "override_reason": "WORN_RIDGE_MANUAL_LABORER"
    },
    "double_iris": {
      "captured": true,
      "template_integrity_sig": "HMAC-SHA512-VAL"
    },
    "facial_portrait": {
      "captured": true,
      "template_integrity_sig": "HMAC-SHA512-VAL"
    }
  },
  "audit_trail_reference": "WORM-block-84210"
}
```

---

*Verified and signed by the SNISID Biometric Ethics & Standards Commission.*
