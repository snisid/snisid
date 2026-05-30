# SNISID Sovereign Digital Doctrine
**Version** : 1.0 | **Date** : 2025-03-27 | **Statut** : ✅ Approuvé

## Règles fondamentales

### 1. Données
- **Hébergement** : exclusivement dans des datacenters certifiés situés sur le territoire national.
- **Chiffrement** : AES-256 au repos, TLS 1.3 en transit.
- **Export** : toute copie hors du territoire est interdite sans autorisation expresse du Conseil National de Souveraineté Numérique.

### 2. Identité et Accès (IAM)
- **Contrôle** : système d’identité centralisé (SNISID IAM) géré par l’État.
- **Fédération** : OK avec systèmes tiers, mais authentification finale toujours validée par SNISID.

### 3. PKI
- **Root CA** : unique, nationale, gérée par l’Autorité de Certification Nationale (ACN).
- **Certificats** : tous les services SNISID utilisent des certificats émis par cette CA.

### 4. APIs
- **Standard** : API RESTful (OpenAPI 3.0) + GraphQL pour requêtes complexes.
- **Transports** : mTLS obligatoire.
- **Événements** : Kafka (protobuf) pour les notifications temps réel.

### 5. Cloud
- **Public** : interdit pour les données critiques.
- **Privé** : Kubernetes on-premise ou cloud souverain (ex : cloud national haïtien).

### 6. Offline-first
- **Obligation** : toute fonctionnalité critique (enrôlement, vérification, délivrance) doit avoir un mode offline documenté et testé.
- **Sync** : asynchrone, avec file d’attente et résolution de conflits.

### 7. Cybersécurité
- **Zero Trust** : jamais de confiance implicite, micro-segmentation, authentification continue.
- **Audit** : immutable, horodaté (Kafka + blockchain optionnelle).

### 8. Interopérabilité
- **Objet** : toutes les agences doivent adopter le SNISID Interoperability Framework.
- **Non-respect** : exclusion du réseau inter-agences.

## Interdictions absolues
- Dépendance cloud étranger critique (AWS, Azure, GCP) pour les données d’identité.
- Données biométriques stockées hors du territoire.
- Accès root non journalisé.
- API sans authentification mTLS.
- Certification « home-made » (self-signed) hors environnement de test.
