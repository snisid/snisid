---
# ============================================================
# SNISID-Security — National Digital Evidence Platform
# Gestion inaltérable des preuves numériques
# Document ID: SNISID-DIGITAL-EVIDENCE-001
# Version: 1.0.0
# ============================================================

## 1. CONCEPT : CHAIN OF CUSTODY INALTÉRABLE

Le système de preuves numériques gère le stockage, le hachage et la traçabilité de tout élément numérique utilisé en justice (vidéos de surveillance, extractions téléphoniques Cellebrite, documents scannés).

### 1.1 Exigences Légales (Forensic Readiness)
- Toute preuve numérique doit avoir un hash cryptographique (SHA-256 ou supérieur) calculé au moment exact de la saisie.
- Le hash est inséré dans le Criminal Event Store (Kafka/CockroachDB). Il ne peut plus jamais être modifié.
- Le fichier original est stocké sur un système WORM (Write Once, Read Many) empêchant toute altération, même par un administrateur système (Protection contre l'Insider Threat).

## 2. ARCHITECTURE DE STOCKAGE (MinIO WORM)

Le stockage s'appuie sur MinIO avec la fonctionnalité "Object Lock" (Compliance Mode).

```yaml
# Configuration conceptuelle du Bucket "snisid-evidence"
Bucket: snisid-evidence
ObjectLocking: Enabled
Mode: COMPLIANCE # Personne, pas même le root, ne peut supprimer l'objet avant expiration
RetentionPeriod: 10 YEARS # Durée de prescription criminelle standard
Versioning: Enabled
Encryption: SSE-KMS (Vault Transit Engine)
```

## 3. WORKFLOW D'INTÉGRITÉ

1. L'enquêteur upload le fichier `.mp4` via le portail sécurisé.
2. Le backend calcule le `SHA-256`.
3. Le fichier est poussé sur MinIO WORM.
4. L'événement `EvidenceUploaded(CaseID, Hash, MinioURI)` est signé par la carte PKI de l'enquêteur et publié sur Kafka.
5. **Au Tribunal :** Lors de la projection de la vidéo, le logiciel du tribunal recalcule le Hash de la vidéo tirée de MinIO et le compare au Hash enregistré dans le ledger Kafka. Une correspondance certifie l'intégrité absolue.

---
*Document ID: SNISID-DIGITAL-EVIDENCE-001 | Approuvé par: Laboratoire de Police Scientifique*
