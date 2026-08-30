# Skill: SNISID ANH Verification

## Description

Ce skill permet à Claude AI de vérifier l'authenticité des extraits de naissance auprès des Archives Nationales d'Haïti (ANH) via une API sécurisée. Il valide le numéro d'acte, les informations biographiques et l'identifiant unique de l'extrait de naissance.

## Capabilities

*   **`verify_birth_certificate(act_number: str, first_name: str, last_name: str, dob: str)`**
    *   **Description:** Vérifie un extrait de naissance en utilisant le numéro d'acte, le prénom, le nom et la date de naissance.
    *   **Parameters:**
        *   `act_number` (string, required): Le numéro d'acte de l'extrait de naissance.
        *   `first_name` (string, required): Le prénom figurant sur l'extrait de naissance.
        *   `last_name` (string, required): Le nom figurant sur l'extrait de naissance.
        *   `dob` (string, required): La date de naissance au format YYYY-MM-DD.
    *   **Returns:** Un objet JSON indiquant la validité de l'extrait de naissance et toute anomalie détectée.

## Usage Example

```python
print(snisid_anh_verification.verify_birth_certificate(act_number="123456", first_name="Jean", last_name="Pierre", dob="1990-01-15"))
```
```

## Intégration Technique

Ce skill s'intègre avec le système de l'ANH via une API RESTful. Les requêtes sont authentifiées et chiffrées pour garantir la sécurité et la confidentialité des données. Les réponses de l'API sont traitées pour extraire les informations de validité et les éventuelles incohérences. Le skill est conçu pour être appelé par Claude AI dans le cadre du processus de vérification croisée du SNISID.
