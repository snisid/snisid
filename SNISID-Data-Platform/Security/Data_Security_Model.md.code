# Data Security Model

## Objectif
Protéger le patrimoine numérique national avec chiffrement, segmentation, gestion des clés, stockage immutable et contrôles d'accès.

## Capacités

| Fonction | Support |
|---|---:|
| Encryption at rest | Oui |
| Encryption in transit | Oui |
| Key management | Oui |
| Data segmentation | Oui |
| Immutable storage | Oui |
| Secrets management | Oui |
| DLP | Oui |

## Contrôles

| Contrôle | Exigence |
|---|---|
| Chiffrement repos | AES-256 ou équivalent, clés État |
| Chiffrement transit | TLS 1.3/mTLS pour services critiques |
| KMS/HSM | Rotation, séparation des rôles, audit |
| Segmentation | Zones Public/Internal/Restricted/Secret |
| Immutable storage | WORM pour audit, archives critiques |
| Tokenisation | Pour identifiants sensibles |
| DLP | Détection export PII/biométrie |
| Backup | Chiffré, testé, isolé |

## Séparation des environnements

- Production, préproduction et développement séparés.
- Données de production interdites en dev sans anonymisation.
- Accès administrateur via bastion, MFA et session recording.
- Secrets jamais stockés dans code ou notebooks.

## Gestion des clés

- Clés racines en HSM ou équivalent souverain.
- Rotation périodique et rotation d'urgence.
- Dual control pour opérations critiques.
- Journalisation de toute utilisation de clé.
- Plan de révocation et récupération.
