# PHASE_05_IMPLEMENTATION.md

## Nom de la phase
Phase 5 - Public Key Infrastructure (PKI)

## Objectif
Déployer l'infrastructure à clé publique (PKI) souveraine garantissant l'identité numérique, le chiffrement mTLS, la signature des documents nationaux et l'intégration sécurisée aux modules cryptographiques matériels (HSM).

## Fonctionnalités ajoutées
- Infrastructure Kubernetes pour la gestion des certificats via `cert-manager`.
- Stratégies d'émission et de révocation des clés (OCSP/CRL).
- Playbooks et Runbooks pour les cérémonies de clés hors-ligne (Root CA) et compromissions (HSM failure).
- Tableaux de bord de surveillance de l'infrastructure de chiffrement.

## Fichiers créés / intégrés
L'ensemble de l'architecture `pki/` a été intégré au projet sous le nom `SNISID-PKI/` :
- `SNISID-PKI/cert-manager/` (Manifests Kubernetes)
- `SNISID-PKI/certificates/` (Modèles de certificats)
- `SNISID-PKI/hsm/` (Spécifications Hardware)
- `SNISID-PKI/root-ca/` (Procédure de Cérémonie)
- `SNISID-PKI/runbooks/` (Gestion de crise PKI)

## Fichiers modifiés
Aucun. L'arborescence vient s'ajouter aux socles définis lors des Phases 1 et 4.

## Dépendances ajoutées
L'exécution et la mise en production nécessiteront :
- Un cluster Kubernetes (provisionné en Phase 4).
- HashiCorp Vault (pour le `vault-issuer`).
- Modules de Sécurité Matériels (HSM) physiques.

## Variables d’environnement
- `VAULT_TOKEN` et `VAULT_ADDR` pour que `cert-manager` s'authentifie auprès de Vault.
- Les URIs d'accès OCSP/CRL pour les vérifications de révocation.

## Migrations ou changements de base de données
- Aucun changement applicatif direct.

## Commandes de test / build / déploiement
Le déploiement des manifests `cert-manager` sur le cluster :
```bash
kubectl apply -f "SNISID-PKI\cert-manager\cluster-issuers.yaml"
```

## Procédure de rollback
Pour supprimer l'infrastructure PKI :
```bash
Remove-Item -Path "c:\Users\sopil\Desktop\snisid system\SNISID-PKI" -Recurse -Force
```
*(Attention : en production, la suppression des CAs entraînera la révocation et le dysfonctionnement de tous les services reposant sur mTLS).*

## Risques connus
- La compromission de la clé privée de la Root CA (gérée dans les Runbooks) entraînerait l'effondrement de toute la chaîne de confiance nationale. L'isolement hors-ligne est obligatoire.

## Points à valider manuellement
- Exécuter physiquement la *Root CA Ceremony* selon le manuel fourni.
- S'assurer que le portail de révocation OCSP est accessible par tous les services internes et externes.
