# Skill: SNISID Core Verification

## Description

Ce skill central orchestre le processus de vérification croisée des documents d'identité et fiscaux en Haïti, en intégrant les informations provenant des Archives Nationales d'Haïti (ANH), de l'Office National d'Identification (ONI), de la Direction Générale des Impôts (DGI) et de la Direction Centrale de la Police Judiciaire (DCPJ/BRJ). Il est responsable de la détection des incohérences, de la génération d'alertes en temps réel et de la journalisation immuable de toutes les actions.

## Capabilities

*   **`perform_cross_verification(birth_certificate_data: dict, national_id_data: dict, tax_id_data: dict, police_certificate_data: dict)`**
    *   **Description:** Effectue une vérification croisée complète de tous les documents fournis et génère un rapport consolidé avec les alertes.
    *   **Parameters:**
        *   `birth_certificate_data` (dict, required): Données de l'extrait de naissance, incluant `act_number`, `first_name`, `last_name`, `dob`.
        *   `national_id_data` (dict, required): Données de la Carte d'Identification Nationale, incluant `cin_number`, `nnu`, `first_name`, `last_name`, `dob`.
        *   `tax_id_data` (dict, required): Données fiscales, incluant `nif`, `receipt_number`, `amount`.
        *   `police_certificate_data` (dict, required): Données du Certificat de Police, incluant `certificate_number`, `request_id`, `first_name`, `last_name`, `dob`.
    *   **Returns:** Un objet JSON contenant un rapport de vérification consolidé, la liste des alertes détectées et le statut global de la demande.

## Usage Example

```python
birth_cert = {"act_number": "123456", "first_name": "Jean", "last_name": "Pierre", "dob": "1990-01-15"}
national_id = {"cin_number": "987654321", "nnu": "123-456-789", "first_name": "Jean", "last_name": "Pierre", "dob": "1990-01-15"}
tax_id = {"nif": "000-111-222-3", "receipt_number": "REC-2026-001", "amount": 1500.00}
police_cert = {"certificate_number": "CP-2026-0001", "request_id": "REQ-BRJ-54321", "first_name": "Jean", "last_name": "Pierre", "dob": "1990-01-15"}

print(snisid_core_verification.perform_cross_verification(
    birth_certificate_data=birth_cert,
    national_id_data=national_id,
    tax_id_data=tax_id,
    police_certificate_data=police_cert
))
```

## Intégration Technique

Ce skill agit comme un orchestrateur. Il appelle séquentiellement ou en parallèle les skills spécifiques (ANH, ONI, DGI, DCPJ/BRJ) pour collecter les résultats de vérification. Il compare ensuite les données et les statuts retournés par chaque skill pour identifier les incohérences. Toutes les actions et les résultats sont enregistrés dans un système de logs immuables pour assurer la traçabilité et la conformité avec le principe d'"Anti-Corruption by Design". Les alertes sont générées en temps réel et peuvent être configurées pour bloquer automatiquement les processus de demande en cas de fraude avérée ou de divergences critiques.
