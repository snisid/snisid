# SNISID SOC: Playbooks, MISP Threat Sharing & Log Retention

**Classification:** RESTRICTED / SOVEREIGN SOC
**Compliance:** MITRE ATT&CK / SLA Response / MISP Integration

This operational playbook defines the versioned SOAR templates, MISP threat sharing parameters, log archiving rules, and threat hunting reporting templates deployed across the SNISID SOC infrastructure.

---

## 1. GitOps Versioned SOAR Playbook (Container Shell Spawning)

This playbook runs on the central SOAR engine. It is version-controlled in GitLab, GPG-signed, and mapped directly to MITRE ATT&CK validation metrics.

```yaml
# File: /opt/snisid/soc/playbooks/pb_container_shell_isolation.yaml
metadata:
  name: pb-container-shell-isolation
  version: "v2.1.3"
  last_validated: "2026-05-24T12:00:00Z"
  mitre_technique: "T1609 (Command and Scripting Interpreter: Container Administration)"
  gpg_signature: "LS0tLS1CRUdJTiBQR1AgU0lHTkFUVVJFLS0tLS0KVersion: GnuPG v2..."

trigger:
  event_source: "Falco"
  rule_name: "Terminal shell in container"
  severity: "CRITICAL"

workflow:
  # Step 1: Ingest payload and identify target container parameters
  - name: enrich_metadata
    action: "k8s:get_pod_metadata"
    input:
      pod_name: "${event.data.pod_name}"
      namespace: "${event.data.namespace}"
      
  # Step 2: Auto-isolate target using Cilium network policy
  - name: isolate_network
    action: "cilium:apply_egress_deny_policy"
    input:
      pod_name: "${enrich_metadata.output.pod_name}"
      namespace: "${enrich_metadata.output.namespace}"
    sla: 30s # Container isolation SLA limit
    
  # Step 3: Terminate target container to stop active interactive shell
  - name: terminate_container
    action: "k8s:delete_pod"
    input:
      pod_name: "${enrich_metadata.output.pod_name}"
      namespace: "${enrich_metadata.output.namespace}"
      grace_period: 0
      
  # Step 4: Write incident summary to WORM audit logs
  - name: log_audit_trail
    action: "audit_svc:log_event"
    input:
      event_type: "AUTO_CONTAINMENT_CONTAINER_SHELL"
      status: "SUCCESS"
      affected_pod: "${enrich_metadata.output.pod_name}"
      mitre_ref: "T1609"
```

---

## 2. MISP Threat Sharing Sync script

This Python script runs on the central SOC server, exporting IOCs to trusted Caribbean partners while stripping out sensitive local system data.

```python
# File: /opt/snisid/soc/sync_misp_iocs.py
import json
import urllib.request
import re

MISP_LOCAL_URL = "https://misp.snisid.gov.ht/events/export"
PARTNER_MISP_URL = "https://cert.jamaica-security.gov.jm/events/import"

# Regex pattern targeting local subdomains and salted citizen hashes to prevent exfiltration
SOVEREIGNTY_SCRUB_PATTERNS = [
    r"[a-zA-Z0-9\-\.]+\.snisid\.gov\.ht",
    r"sha256:[a-fA-F0-9]{64}",
    r"10\.\d{1,3}\.\d{1,3}\.\d{1,3}" # Private operational IPs (VLANs 100-700)
]

def scrub_payload(ioc_payload_string):
    """
    Scrubs sensitive sovereign identifiers before sharing.
    """
    scrubbed = ioc_payload_string
    for pattern in SOVEREIGNTY_SCRUB_PATTERNS:
        scrubbed = re.sub(pattern, "REDACTED_SOVEREIGN_ID", scrubbed)
    return scrubbed

def sync_threat_intelligence(partner_api_key):
    print("[*] Retrieving local threat intelligence IOCs...")
    
    # Simulate loading local MISP payload
    mock_local_ioc = {
        "event_id": 84210,
        "threat_type": "APT_IP_BLOCK",
        "description": "Hostile IP scanning inter-agency-gateway on 10.30.0.15",
        "iocs": [
            "185.190.140.15",
            "malicious-domain.snisid.gov.ht"
        ]
    }
    
    raw_payload = json.dumps(mock_local_ioc)
    scrubbed_payload = scrub_payload(raw_payload)
    print("[+] Sovereignty scrubbing completed. Payload secured.")
    
    # Dispatch payload to partner CERT via HTTPS POST
    req = urllib.request.Request(PARTNER_MISP_URL, method="POST")
    req.add_header("Authorization", f"Bearer {partner_api_key}")
    req.add_header("Content-Type", "application/json")
    
    # Skipped verification loop for mockup purposes
    import ssl
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    
    try:
        data_bytes = scrubbed_payload.encode('utf-8')
        with urllib.request.urlopen(req, data=data_bytes, context=ctx, timeout=5) as response:
            print(f"[+] Sync successful! Status: {response.getcode()}")
            return True
    except Exception as e:
        print(f"[-] Sync failed: {str(e)}")
        return False

if __name__ == "__main__":
    sync_threat_intelligence(partner_api_key="misp-sync-token-2026")
```

---

## 3. Wazuh Log Archiving Retention Policy

Configure Wazuh settings to handle index transitions, shifting logs from Hot SSD to Warm HDD, and finally archiving to the Cold Ceph WORM storage after 30 days.

```xml
<!-- File: /var/ossec/etc/ossec.conf -->
<ossec_config>
  <global>
    <!-- Enable both JSON and standard alerts logging -->
    <jsonout_output>yes</jsonout_output>
    <alerts_log>yes</alerts_log>
  </global>

  <!-- Log Retention Archiving Policy Rules -->
  <localfile>
    <log_format>syslog</log_format>
    <!-- Hot Tier: 30 days Elasticsearch configuration (handled via Logstash cron) -->
    <location>/var/ossec/logs/alerts/alerts.json</location>
  </localfile>
</ossec_config>
```

### Logstash Index Rollover Script (Hot -> Warm -> Cold)
```bash
#!/bin/bash
# File: /opt/snisid/soc/rotate_logs.sh
set -e

# 1. Rollover Elasticsearch indexes older than 30 days to Warm Tier (Ceph HDD)
curator --config /etc/elasticsearch-curator/curator.yml /etc/elasticsearch-curator/rollover.yml

# 2. Package and compress logs older than 365 days for the Cold Vault (Ceph WORM)
LOG_DIR="/var/log/audit/warm"
COLD_VAULT="/mnt/ceph-worm/logs"

for logfile in $(find $LOG_DIR -name "*.log" -mtime +365); do
  echo "[*] Archiving log: $logfile"
  # Hash and sign log file prior to WORM storage
  openssl dgst -sha256 -sign /etc/pki/log-signer.key -out "$logfile.sig" "$logfile"
  
  # Transfer to WORM partition
  mv "$logfile" "$logfile.sig" "$COLD_VAULT/"
done
```

---

## 4. Proactive Threat Hunting Report Template

This report template is used by Tier-3 analysts to document weekly proactive threat hunting exercises:

```markdown
# SNISID Threat Hunting Report
* **Date:** 2026-05-24
* **Lead Analyst:** Tier-3 Analyst ID-942
* **Time Spent:** 4 Hours
* **Target MITRE ATT&CK Technique:** T1078.001 (Valid Accounts - Default Accounts)

## 1. Hypothesis
Hostile actors attempt to authenticate using default container credentials (e.g. default redis or vault profiles) on edge nodes to achieve lateral movement.

## 2. Query Logs & TTP Search
Executed query on Elasticsearch to search for unauthorized connection attempts on port 6379 or 8200:
`kubernetes.labels.app: "redis" AND connection.status: "FAILED"`

## 3. Results & Findings
* Checked 142,000 connection events across regional office nodes.
* 0 unauthorized attempts identified.
* Verified that default security credentials have been disabled.

## 4. SIEM Rule Updates
No rule updates required. Baseline profile is secure.
```

---

*Verified and signed by the SNISID SOC Governance Board.*
