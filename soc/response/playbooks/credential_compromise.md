# Playbook: Compromission d'un Opérateur (Credential Compromise)

**ID**: PB-SNISID-001
**SLA MTTD**: < 5 minutes
**SLA MTTR**: < 30 minutes
**Severity**: CRITICAL

## 1. Identification
- **Trigger**: Elastic SIEM alerts via `admin_access_anomalies.yml` or `identity_theft_attempts.yml`.
- **Validation**:
  - Verify if the IP address belongs to a known VPN or authorized proxy.
  - Check with the operator (via Out-of-Band communication) if they initiated the action.

## 2. Containment (Automated via Shuffle)
- **Step 2.1**: Suspend the operator's account in Keycloak.
  - `POST /auth/admin/realms/snisid/users/{id}/logout`
  - Disable account flag.
- **Step 2.2**: Terminate active K8s/Vault sessions if the operator had admin privileges.
- **Step 2.3**: Block the offending source IP at the Ingress/WAF level via automated API call.

## 3. Eradication
- Force credential rotation (password reset).
- Revoke and reissue MFA (TOTP/WebAuthn) tokens.
- Review and revert any modifications made by the compromised account during the timeframe of the attack (check `snisid.audit.events`).

## 4. Recovery
- Notify the operator via official channels once the account is secured.
- Resume normal monitoring.

## 5. Lessons Learned
- Attach all findings to TheHive case.
- Review if MFA bypass was utilized and patch the vulnerability.
