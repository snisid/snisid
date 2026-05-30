# Data Protection Impact Assessment (DPIA) — SNISID Biometric Processing

## 1. Description du Traitement
**Finalité**: Vérification de l'identité des citoyens via la reconnaissance faciale (1:1 et 1:N) pour l'accès aux services gouvernementaux critiques (eIDAS LoA High).
**Catégories de données**: Données biométriques (Gabarits faciaux, Minuties, Scores de Liveness). *Catégorie particulière (Article 9 RGPD).*
**Personnes concernées**: Citoyens de la nation.

## 2. Base Légale et Nécessité (Article 6 & 9)
**Base légale**: Mission d'intérêt public et obligation légale (Art. 6(1)(c) et (e)).
**Exception Art. 9**: Motifs d'intérêt public important (Art. 9(2)(g)), appuyés par la loi nationale de sécurité.

## 3. Évaluation des Risques pour les Droits et Libertés

| Risque Identifié | Impact | Probabilité | Mesure de Mitigation (Implémentation SNISID) | Risque Résiduel |
| :--- | :--- | :--- | :--- | :--- |
| Usurpation d'identité via Deepfake | Critique | Moyenne | Moteur d'IA "PAD / Liveness Detection" certifié NIST. Rejet >99.5% des attaques de spoofing. | Faible |
| Vol massif de la base biométrique | Critique | Faible | 1. Les images brutes ne sont jamais stockées en base.<br>2. Les gabarits (templates) sont chiffrés avec HashiCorp Vault (Transit Engine).<br>3. Accès limité par micro-segmentation Cilium (Zero Trust). | Très Faible |
| Faux Rejets (Discrimination) | Élevé | Moyenne | Paramétrage algorithmique garantissant un FRR < 1% sur toutes les ethnies, audité annuellement. Mode de secours (fallback) manuel par opérateur disponible. | Faible |

## 4. Mesures de Sécurité Techniques et Organisationnelles
- **Chiffrement au Repos**: AES-256 via TDE sur les volumes Kubernetes + pgcrypto au niveau de la base.
- **Chiffrement en Transit**: mTLS strict appliqué par Istio Service Mesh.
- **Contrôle d'Accès**: MFA obligatoire pour 100% des opérateurs (FIDO2/WebAuthn), géré par Keycloak.
- **Traçabilité**: Journalisation immuable WORM avec chaîne de hachage (Hash Chain).
- **Rétention**: Conservation limitée à la durée de vie du document d'identité, puis suppression cryptographique (Crypto-shredding des clés).

## 5. Avis du Délégué à la Protection des Données (DPO)
*À compléter et signer par le DPO national avant la mise en production du module.*
