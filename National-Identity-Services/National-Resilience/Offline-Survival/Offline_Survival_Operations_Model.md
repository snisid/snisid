# Offline Survival Operations Model

## 1. Objectif
Garantir la survie SNISID sans Internet national via vérification identité offline, enrôlement offline, synchronisation edge et survivabilité locale des données.

## 2. Capacités
| Fonction | Support | Description |
|---|---:|---|
| Offline identity verification | Oui | caches signés/chiffrés |
| Offline enrollment | Oui | dossiers locaux signés, réconciliation future |
| Edge synchronization | Oui | sync différée dès retour connectivité |
| Local data survivability | Oui | stockage local chiffré, sauvegarde locale |

## 3. Architecture
```text
National Core ⇄ Regional Edge Node ⇄ Mobile/Field Kits
               ├ Identity Verification Cache signé
               ├ Offline Enrollment Module
               ├ Local Encrypted Store
               ├ Audit Log append-only
               └ Sync Agent différé
```

## 4. Principes de sécurité
Caches minimisés et expirables, données locales chiffrées, opérations signées, privilèges temporaires, journalisation append-only, réconciliation obligatoire.

## 5. Kits offline
| Kit | Contenu |
|---|---|
| Regional Edge Kit | mini-serveur, stockage chiffré, UPS, satellite optionnel, docs papier |
| Mobile Enrollment Kit | tablette durcie, lecteur biométrique, imprimante sécurisée, batteries |
| Verification Kit | app offline, cache signé, lecteur QR/NFC/biométrique selon politique |
| Crisis Admin Kit | comptes break-glass, runbooks papier, clés d'activation scellées |

## 6. Processus vérification offline
Authentifier opérateur → charger cache signé → vérifier identifiant/statut/expiration → produire décision locale → journaliser → synchroniser.

## 7. Processus enrôlement offline
Ouvrir session urgence → collecter données minimales → capturer preuves → émettre preuve temporaire/dossier pending → signer/chiffrer → synchroniser → résoudre conflits.

## 8. Gestion des conflits
| Conflit | Résolution |
|---|---|
| doublon identité | revue centrale + matching biométrique/biographique |
| statut contradictoire | priorité registre central sauf preuve crise validée |
| opération non signée | rejet ou enquête |
| cache expiré | mode restreint avec approbation renforcée |

## 9. Exercices
Vérification offline mensuelle, enrôlement offline trimestriel, fonctionnement batteries trimestriel, sync post-crise trimestrielle, activation kit semestrielle.
