# Déploiement Kubernetes / DevSecOps

## Pipeline

1. SAST/secret scanning.
2. `npm ci && npm run typecheck && npm run test`.
3. SBOM + signature image.
4. Deploy staging avec policy-as-code.
5. Promotion prod avec approbation sécurité.

## Kubernetes

- Namespace isolé `snisid-mcp`.
- NetworkPolicies deny-all puis allow Gateway/SIEM.
- Secrets via External Secrets + KMS/HSM.
- Pod Security Standards restricted.
- mTLS service mesh.
- HPA basé CPU/RPS/latence.
