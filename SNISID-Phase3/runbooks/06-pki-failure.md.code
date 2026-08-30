# 🔐 Runbook 06 — Panne PKI / TSA

**Severity :** Sev1
**Owner :** Astreinte PKI + WGO

## 1. Symptômes
- Signature `pki.sign.qualified` timeout > 1 s
- Erreurs `PKI_UNAVAILABLE` dans logs workflow engine
- TSA RFC 3161 inaccessible
- Tous les workflows critiques bloqués sur étape signature

## 2. Remédiation

### A. Bascule HA PKI
```bash
# Forcer la promotion du secondaire
kubectl -n pki patch deploy pki-signer -p '{"spec":{"template":{"metadata":{"labels":{"role":"primary","dc":"dc2"}}}}}'
```

### B. Buffer signature
- Activer le mode "deferred signature" (signature en différé) :
  ```bash
  kubectl set env deploy/snisid-workflow-engine PKI_DEFERRED=true
  ```
- Les workflows continuent ; signatures rattrapées dès retour PKI.
- ⚠️ N'utiliser que **24 h max** sous décision WGO.

### C. TSA secondaire
- Switcher vers TSA secondaire (`tsa-dc2.snisid.ht`) :
  ```bash
  kubectl set env deploy/snisid-workflow-engine PKI_TSA_ENDPOINT=https://tsa-dc2.snisid.ht/rfc3161
  ```

## 3. Vérification
- Signature p99 < 250 ms
- Workflows progressent

## 4. Suivi
- Audit complet des signatures émises pendant l'incident
- Re-tampon RFC 3161 si nécessaire pour preuve juridique
