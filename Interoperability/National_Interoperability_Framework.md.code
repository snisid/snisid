# National Interoperability Framework
**Version** : 1.0 | **Date** : 2025-03-27 | **Statut** : 🟡 En relecture

## Standards

### APIs
- Format : RESTful (OpenAPI 3.0) ou GraphQL
- Authentification : mTLS + JWT signé par la CA nationale
- Rate limiting : OUI, défini par profil d’agence

### Events (Kafka)
- Format : Protobuf (schema registry mandatory)
- Topics par domaine : `snisid.citizen.enrollment`, `snisid.security.flag`
- Rétention : 7 jours, puis archivage dans le data lake

### Data exchange
- Toutes les données échangées doivent être chiffrées (TLS 1.3 ou mTLS)
- Les données en transit ne doivent jamais être en clair

### Audit
- Chaque mutation est tracée avec timestamp, utilisateur, IP, action
- Obligatoire : endpoint dédié `/audit` pour consultation

### Versioning
- API versionnée (v1, v2) avec support simultané 12 mois
- Schéma Kafka évolutif (backward compatible)

### SLA national
- Disponibilité API : 99,9 %
- Temps de réponse médian < 200 ms
- Temps de synchronisation offline max : 24h
