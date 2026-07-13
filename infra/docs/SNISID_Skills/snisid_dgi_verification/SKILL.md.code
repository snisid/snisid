# Skill: SNISID DGI Verification

## Description

Ce skill permet à Claude AI de vérifier la cohérence fiscale d'un individu et la validité des reçus de paiement officiels auprès de la Direction Générale des Impôts (DGI). Il utilise le Numéro d'Identification Fiscale (NIF) pour ces vérifications.

## Capabilities

*   **`verify_tax_id_and_payment(nif: str, receipt_number: str, amount: float)`**
    *   **Description:** Vérifie le Numéro d'Identification Fiscale (NIF) et la validité d'un reçu de paiement.
    *   **Parameters:**
        *   `nif` (string, required): Le Numéro d'Identification Fiscale de l'individu.
        *   `receipt_number` (string, required): Le numéro du reçu de paiement officiel.
        *   `amount` (float, required): Le montant du paiement.
    *   **Returns:** Un objet JSON indiquant la validité du NIF, l'authenticité du reçu et toute incohérence fiscale.

## Usage Example

```python
print(snisid_dgi_verification.verify_tax_id_and_payment(nif="000-111-222-3", receipt_number="REC-2026-001", amount=1500.00))
```

## Intégration Technique

Ce skill s'intègre avec le système de la DGI via une API sécurisée. Les requêtes sont authentifiées et chiffrées. Les réponses de l'API sont traitées pour valider le NIF et le reçu de paiement, assurant ainsi la conformité fiscale. Le skill est conçu pour être appelé par Claude AI dans le cadre du processus de vérification croisée du SNISID.
