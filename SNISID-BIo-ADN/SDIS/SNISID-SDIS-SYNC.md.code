# SNISID-BIO-ADN — SDIS-HT Synchronisation Départementale
**Document ID :** SNISID-SDIS-SYN-001 | **Version :** 1.0.0

---

## 1. LES 10 SDIS DÉPARTEMENTAUX

| Code SDIS | Département | DC hébergement | Labs LDIS sous tutelle |
|-----------|-------------|----------------|------------------------|
| `SDIS-OUEST` | Ouest | SNISID PAP (DC principal) | LDIS-PAP-001 |
| `SDIS-NORD` | Nord | SNISID CAP (DC secondaire) | LDIS-CAP-001 |
| `SDIS-ARTIBONITE` | Artibonite | SDIS nœud Gonaïves | LDIS-GON-001 |
| `SDIS-SUD` | Sud | SDIS nœud Les Cayes | LDIS-LES-001 |
| `SDIS-SUDEST` | Sud-Est | SDIS nœud Jacmel | LDIS-JAC-001 |
| `SDIS-CENTRE` | Centre | SDIS nœud Hinche | LDIS-HIN-001 |
| `SDIS-NORDEST` | Nord-Est | SDIS nœud Fort-Liberté | À créer |
| `SDIS-NORDOUEST` | Nord-Ouest | SDIS nœud Port-de-Paix | À créer |
| `SDIS-GRANDANSE` | Grand-Anse | SDIS nœud Jérémie | À créer |
| `SDIS-NIPPES` | Nippes | SDIS nœud Miragoâne | À créer |

---

## 2. RÔLE DU SDIS

Le SDIS-HT agrège les profils de tous les labs LDIS de son département
et permet le **matching intra-départemental** avant transmission au NDIS.

### Avantages du matching SDIS
- 80% des crimes sont géographiquement localisés dans un département
- Réduction de la charge NDIS (filtre efficace)
- Résultats plus rapides pour les investigations locales

---

## 3. SYNCHRONISATION SDIS → NDIS

```
SDIS reçoit profils LDIS (quotidien)
        │
        ▼
Validation signatures labs
        │
        ▼
Déduplication (loci_hash unique)
        │
        ▼
Matching intra-département
        │
        ├── HIT → Alerte locale (agences du département)
        │
        ▼
Upload vers NDIS-HT (hebdomadaire — dimanche 03h00)
        │
        ▼
Confirmation NDIS reçue → flag uploaded_ndis = TRUE
```

---

## 4. GESTION DES ERREURS DE SYNCHRONISATION

| Erreur | Action automatique | Notification |
|--------|-------------------|--------------|
| Réseau indisponible | Retry x3 (backoff exponentiel) | Email directeur SDIS |
| Signature invalide | Rejet + alerte sécurité | Kafka `bio.security.alert` |
| Quality score dégradé | File d'attente révision | Superviseur lab LDIS |
| Doublon NDIS | Rejet silencieux + log | Aucune (normal) |
| NDIS indisponible | Queue persistante Kafka | Alerte NDIS ops |
