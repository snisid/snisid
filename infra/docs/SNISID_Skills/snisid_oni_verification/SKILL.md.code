# Skill: SNISID ONI Verification

## Description

Ce skill permet à Claude AI de valider l'identité unique des citoyens via la Carte d'Identification Nationale (CIN) auprès de l'Office National d'Identification (ONI). Il vérifie les données biométriques et biographiques, détecte les doublons et les tentatives d'usurpation.

## Capabilities

*   **`verify_national_id(cin_number: str, nnu: str, first_name: str, last_name: str, dob: str)`**
    *   **Description:** Vérifie une Carte d'Identification Nationale (CIN) en utilisant le numéro de carte, le Numéro National Unique (NNU), le prénom, le nom et la date de naissance.
    *   **Parameters:**
        *   `cin_number` (string, required): Le numéro de la Carte d'Identification Nationale.
        *   `nnu` (string, required): Le Numéro National Unique (NNU).
        *   `first_name` (string, required): Le prénom figurant sur la CIN.
        *   `last_name` (string, required): Le nom figurant sur la CIN.
        *   `dob` (string, required): La date de naissance au format YYYY-MM-DD.
    *   **Returns:** Un objet JSON indiquant la validité de la CIN, la correspondance des données biométriques et toute anomalie détectée.

## Usage Example

```python
print(snisid_oni_verification.verify_national_id(cin_number="987654321", nnu="123-456-789", first_name="Marie", last_name="Jeanne", dob="1985-05-20"))
```
```

## Intégration Technique

Ce skill s'intègre avec la base de données biométrique centrale de l'ONI via une API sécurisée. Les requêtes sont authentifiées et chiffrées. Les réponses de l'API sont analysées pour confirmer l'identité de l'individu et détecter toute tentative de fraude ou d'usurpation. Le skill est conçu pour être appelé par Claude AI dans le cadre du processus de vérification croisée du SNISID.
