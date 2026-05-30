# SNISID Playbook: Notification de Violation de Données (< 72h)

Conformément à l'Article 33 du RGPD, toute violation de données à caractère personnel doit être notifiée à l'autorité de contrôle compétente dans les **72 heures** suivant la découverte de la violation.

## 1. Détection et Triage (H+0)
**Trigger**: Alerte CRITIQUE du SIEM Elastic (ex: exfiltration de base de données, compromission d'un opérateur avec accès aux PII).
**Action SOC**:
- Isoler immédiatement le nœud ou le compte compromis (via TheHive / SOAR).
- Créer un ticket Incident Majeur (P1).
- Notifier immédiatement le CISO et le DPO (Data Protection Officer) via PagerDuty.

## 2. Analyse et Qualification (H+1 à H+12)
**Action Forensic Team**:
- Analyser l'Audit Trail immuable pour déterminer l'étendue exacte des données compromises.
- Les données étaient-elles chiffrées (ex: Gabarits biométriques) ?
    - *Si oui, et que les clés Vault sont intactes, le risque pour les droits des personnes est minime.*
    - *Si non (ex: export CSV en clair), le risque est critique.*

## 3. Prise de Décision DPO (H+12 à H+24)
**Action DPO**:
- Évaluer le risque pour les droits et libertés des citoyens concernés.
- Si le risque est avéré : Préparer le formulaire de notification légale.
- Le formulaire doit contenir :
  - La nature de la violation (catégories et nombre de citoyens concernés).
  - Les conséquences probables.
  - Les mesures correctives prises.

## 4. Notification Officielle (Avant H+72)
- **Autorité de Contrôle**: Dépôt du dossier officiel sur le portail de l'autorité (ex: CNIL / Autorité nationale).
- **Communication Citoyens** (Article 34) : Si le risque est *élevé*, préparer la campagne de communication publique via le portail gouvernemental (SNISID.gov) et les emails d'alerte.

## 5. Post-Mortem et Remédiation
- Mise à jour de la matrice de risques ISO 27001.
- Audit externe forcé des contrôles de sécurité défaillants.
