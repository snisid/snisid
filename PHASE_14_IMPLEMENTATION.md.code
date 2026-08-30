# PHASE_14_IMPLEMENTATION.md

## Nom de la phase
Phase 14 - Governance & Legal Decrees (Policy-as-Code)

## Objectif
Traduire les règles logicielles et l'architecture sécurisée du système SNISID en un cadre juridique robuste et des décrets présidentiels, et rendre ces décrets informatiquement contraignants via l'Open Policy Agent (OPA).

## Fonctionnalités ajoutées
- **Policy-as-Code (Rego)** : Création de 3 politiques OPA pour appliquer automatiquement :
  - L'interdiction aux ministères de contourner la National API Gateway (Décret Interopérabilité).
  - L'autorisation pour le CISO National d'isoler (Drop Traffic) un ministère compromis (Décret Zero Trust).
  - Le respect des quotas et contrats pour la monétisation des banques privées (Décret KYC).
- **Générateur de Décrets** : Création d'un script Go pour transformer ces règles en fichiers HTML/PDF institutionnels officiels, en s'appuyant sur l'Approval Pack de la Phase 13.

## Fichiers créés
- `governance/opa-policies/interoperability.rego`
- `governance/opa-policies/zerotrust.rego`
- `governance/opa-policies/banking.rego`
- `governance/opa-policies/interoperability_test.rego`
- `National-Executive-Operations/scripts/generate_decrees.go`

## Fichiers modifiés
- `PHASE_14_IMPLEMENTATION.md` (Mise à jour pour inclure l'implémentation OPA).

## Dépendances ajoutées
- Moteur OPA (Open Policy Agent) pour l'évaluation du langage Rego.
- Go (pour le script de génération des documents officiels).

## Variables d’environnement
- `OPA_SERVER_URL` (Utilisée par l'API Gateway pour interroger OPA).

## Changements de base de données
- Aucun, l'évaluation des règles Rego s'effectuant en mémoire via OPA.

## Commandes exécutées
Création de l'arborescence et génération de code source :
```bash
# Pour exécuter les tests Rego (Nécessite le binaire OPA et de l'espace disque) :
opa test governance/opa-policies/

# Pour générer physiquement les documents des décrets :
cd National-Executive-Operations/scripts/
go run generate_decrees.go
```

## Instructions de déploiement
L'API Gateway (Phase 10) devra être reconfigurée pour monter le dossier `governance/opa-policies/` en volume partagé et déléguer ses autorisations au conteneur `opa`.

## Procédure de rollback
Pour annuler l'intégration Policy-as-Code de la Phase 14 :
```bash
rm -rf governance/opa-policies
rm National-Executive-Operations/scripts/generate_decrees.go
```

## Risques connus
- Blocage strict du réseau : Toute erreur de syntaxe ou de logique dans `zerotrust.rego` pourrait isoler accidentellement tout le gouvernement. 

## Points à valider manuellement
- S'assurer que le binaire `opa` s'exécute correctement lorsque l'espace disque du serveur sera libéré.
- Les décrets générés dans `National-Executive-Operations/documents/generated_decrees/` doivent être signés via le Parapheur (Phase 13).
