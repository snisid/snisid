# Data Access Governance Model

## Objectif
Contrôler l'accès aux données de manière contextuelle, justifiée, limitée par finalité et auditée.

## Capacités

| Fonction | Support |
|---|---:|
| RBAC | Oui |
| ABAC | Oui |
| Purpose limitation | Oui |
| Consent enforcement | Oui |
| Data masking | Oui |
| Just-in-time access | Oui |
| Break-glass contrôlé | Oui |

## Modèle de décision

```text
Sujet + Rôle + Attributs + Ressource + Classification + Finalité + Consentement + Risque -> ALLOW/DENY/MASK
```

## Attributs ABAC

| Catégorie | Attributs exemples |
|---|---|
| Sujet | agence, rôle, habilitation, MFA, localisation |
| Ressource | classification, domaine, owner, sensibilité |
| Contexte | heure, réseau, device, niveau risque |
| Finalité | service_delivery, fraud_control, audit, legal_case |
| Consentement | statut, portée, expiration, base légale |

## Masquage

| Donnée | Masquage par défaut |
|---|---|
| National ID | XXX-XXX-1234 |
| Téléphone | +509 **** 1234 |
| Adresse | Commune uniquement |
| Biométrie | Interdiction export, token uniquement |
| Date naissance | Année uniquement pour analytics |

## Break-glass

Autorisé uniquement si :
- urgence validée,
- MFA renforcé,
- justification obligatoire,
- durée limitée,
- notification owner/security,
- revue post-incident.
