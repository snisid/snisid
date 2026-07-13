# SNISID-BIO-ADN — NDIS-HT Index National
**Document ID :** SNISID-NDIS-NAT-001 | **Version :** 1.0.0

---

## 1. RÔLE DU NDIS-HT

Le **NDIS-HT** (National Data Index System — Haïti) est le niveau le plus élevé
de la hiérarchie SNISID-BIO-ADN. Il héberge tous les profils ADN remontés
par les 10 SDIS et permet le **matching cross-départemental**.

Il est l'équivalent du NDIS américain géré par le FBI.

**Opéré par :** SNISID Central (Direction Nationale) + DCPJ Nationale  
**Hébergement :** Datacenter SNISID Port-au-Prince (redondance Cap-Haïtien)

---

## 2. MATCHING CROSS-DÉPARTEMENTAL

```
Nouveau profil BIO-FSC (scène de crime — Saint-Marc, Artibonite)
        │
        ▼
Matching SDIS-ARTIBONITE (pas de hit local)
        │
        ▼
Upload NDIS-HT (hebdomadaire)
        │
        ▼
Matching NDIS cross-départemental :
    ├── vs BIO-CON de tous les 10 SDIS
    ├── vs BIO-ARR de tous les 10 SDIS
    └── vs BIO-FSC d'autres affaires
        │
        ├── HIT BIO-CON Ouest (condamné Port-au-Prince)
        │       └── Alerte CRITIQUE → DCPJ Nationale → DCPJ Artibonite
        │
        └── Pas de hit → Rapport hebdomadaire négatif
```

---

## 3. CONNEXION INTERPOL DNA GATEWAY

Le NDIS-HT se connecte à l'**INTERPOL DNA Gateway** via le
Bureau Central National (BCN) haïtien de la DCPJ.

```
NDIS-HT
    │
    │  (mTLS + certificat INTERPOL I-24/7)
    ▼
BCN Haïti DCPJ
    │
    │  (I-24/7 réseau sécurisé INTERPOL)
    ▼
INTERPOL DNA Gateway (Lyon)
    │
    ├── 87 pays membres
    └── Matching profils internationaux non identifiés
```

**Cas d'usage Haïti :**
- Identification de victimes de catastrophes (séismes 2010, 2021)
- Fugitifs haïtiens identifiés à l'étranger
- Traite des personnes — victimes retrouvées dans d'autres pays

---

## 4. RAPPORTS NDIS HEBDOMADAIRES

Chaque semaine, le NDIS-HT génère automatiquement :

| Rapport | Destinataire | Format |
|---------|-------------|--------|
| Statistiques profils uploadés | Directeur SNISID | PDF chiffré |
| Hits de la semaine | DCPJ Nationale | PDF chiffré + Kafka |
| Profils sans correspondance | Directeurs SDIS | Dashboard SNISID |
| Alertes qualité labs | Superviseurs LDIS | Email |
| Rapport INTERPOL | BCN DCPJ | Format I-24/7 |

---

## 5. STATISTIQUES CIBLES (Année 1)

| Métrique | Cible |
|----------|-------|
| Profils BIO-CON enregistrés | 5 000 |
| Profils BIO-FSC soumis | 2 000 |
| Taux de correspondance BIO-FSC → BIO-CON | ≥ 15% |
| Délai hit → notification agence | < 24h |
| Disponibilité système matching | ≥ 99.5% |
| Délai LAPI response P99 | < 200ms |
