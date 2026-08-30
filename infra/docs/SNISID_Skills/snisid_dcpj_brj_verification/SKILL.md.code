# Skill: SNISID DCPJ/BRJ Verification

## Description

Ce skill permet à Claude AI de vérifier l'authenticité et la validité du Certificat de Police auprès de la Direction Centrale de la Police Judiciaire (DCPJ) via le Bureau de Renseignements Judiciaires (BRJ). Il valide le numéro de demande/identité unique et les informations du citoyen pour s'assurer de l'absence de casier judiciaire et de la correspondance avec les autres documents d'identité.

## Capabilities

*   **`verify_police_certificate(certificate_number: str, request_id: str, first_name: str, last_name: str, dob: str)`**
    *   **Description:** Vérifie un Certificat de Police en utilisant son numéro, l'identifiant unique de la demande, le prénom, le nom et la date de naissance.
    *   **Parameters:**
        *   `certificate_number` (string, required): Le numéro du Certificat de Police.
        *   `request_id` (string, required): L'identifiant unique de la demande de Certificat de Police.
        *   `first_name` (string, required): Le prénom figurant sur le Certificat de Police.
        *   `last_name` (string, required): Le nom figurant sur le Certificat de Police.
        *   `dob` (string, required): La date de naissance au format YYYY-MM-DD.
    *   **Returns:** Un objet JSON indiquant la validité du certificat, l'état du casier judiciaire et toute anomalie détectée.

## Usage Example

```python
print(snisid_dcpj_brj_verification.verify_police_certificate(certificate_number="CP-2026-0001", request_id="REQ-BRJ-54321", first_name="Jean", last_name="Baptiste", dob="1978-11-03"))
```

## Intégration Technique

Ce skill s'intègre avec le Fichier Central de la DCPJ/BRJ via une API sécurisée. Les requêtes sont authentifiées et chiffrées. Les réponses de l'API sont traitées pour valider l'authenticité du certificat et les antécédents judiciaires du citoyen. Le skill est conçu pour être appelé par Claude AI dans le cadre du processus de vérification croisée du SNISID, permettant une traçabilité et une vérification instantanée.
