# 📡 SNISID — Offline-First National Strategy

**Document N° :** SNISID-OFF-006
**Étape Phase 0 :** 6/16
**Principe :** *Haïti nécessite une résilience terrain native.*

---

## 1. Contexte Haïtien

- Couverture Internet partielle (zones reculées : Grand'Anse, Nord-Ouest, Sud-Est, Bas-Plateau)
- Coupures électriques fréquentes (réseau EDH instable)
- Catastrophes naturelles récurrentes (cyclones, séismes)
- Insécurité épisodique limitant la mobilité des agents

> **Conclusion :** L'offline n'est pas une option, c'est la **fondation**. Online-first = échec garanti.

---

## 2. Doctrine "Offline-First"

| Principe | Application |
|----------|-------------|
| Le terrain fonctionne **toujours**, même 30+ jours sans connectivité | Tous workflows critiques dispo localement |
| La synchronisation est **éventuelle**, jamais bloquante | Store-and-forward natif |
| Les conflits sont **détectés et résolus** automatiquement ou manuellement | CRDT, vector clocks, règles métier |
| L'intégrité légale est **préservée** offline | Signatures locales, horodatage, chaîne de confiance |

---

## 3. Architecture Offline

```
┌─────────────────────────────────────────────────┐
│  CENTRAL CORE (Datacenter PaP + DR Cap-Haïtien) │
└──────────────────────┬──────────────────────────┘
                       │ Sync différée (4G/Satellite/VSAT)
                       ▼
        ┌──────────────────────────────┐
        │  EDGE NODES (10 départements)│
        │  - Cache lecture             │
        │  - Acceptation écritures     │
        │  - Buffer sync               │
        └──────────────┬───────────────┘
                       │ Sync locale (WiFi, BLE, USB)
                       ▼
        ┌──────────────────────────────┐
        │  OFFLINE KITS terrain        │
        │  - Tablette durcie + Mini PC │
        │  - Capteurs biométriques     │
        │  - Imprimante portable       │
        │  - Batterie + panneau solaire│
        │  - Stockage chiffré 512 Go   │
        └──────────────────────────────┘
```

---

## 4. Fonctions Offline Supportées

| Fonction | Offline | Sync |
|----------|---------|------|
| **Enrôlement biométrique** | ✅ Capture complète locale | Upload différé |
| **Vérification 1:1 d'identité** | ✅ Si carte présentée + template local | — |
| **Vérification 1:N** | ⚠️ Sur sous-ensemble local (commune) | Match complet lors sync |
| **Délivrance acte naissance** | ✅ Avec QR de vérification | Validation centrale ultérieure |
| **Consultation registres** | ✅ Lecture cache | Mise à jour pull |
| **Signature électronique** | ✅ HSM portable / smartcard | — |
| **Notifications citoyens** | ⚠️ Buffer SMS | Envoi à reconnexion |

---

## 5. Composants Offline Kit (par équipe mobile)

| Composant | Spec |
|-----------|------|
| Mini PC durci | Intel NUC i7, 32 Go RAM, SSD 1 To NVMe chiffré |
| Tablette opérateur | Android durci (IP67), 10", lecteur empreintes intégré |
| Scanner empreintes | 4-4-2 (FBI Appendix F certifié) |
| Caméra photo ID | ICAO 9303 compliant |
| Lecteur iris (optionnel) | Selon protocole national |
| Imprimante portable | A4, alimentation 12V, papier sécurisé |
| Batterie | 200 Wh autonomie 8h |
| Panneau solaire pliable | 100 W |
| Routeur 4G/satellite | Carte SIM multi-opérateur + slot Starlink/Eutelsat |
| HSM USB | YubiHSM ou équivalent FIPS 140-2 L3 |
| Coffre transport | Pelican IP67 |

---

## 6. Synchronisation — Mécanismes

### 6.1 Store-and-Forward
- Toute écriture créée en local = stockée en **outbox** chiffrée
- Tentative d'envoi périodique (cron + back-off exponentiel)
- Acquittement par le central avant suppression locale

### 6.2 Sync Bi-directionnelle
- **Push** : événements créés localement → central
- **Pull** : référentiels mis à jour (ex. liste personnes recherchées, listes communales) → kit

### 6.3 Résolution de Conflits
| Cas | Règle |
|-----|-------|
| Même enrôlement créé 2 fois (lieux différents) | Match biométrique → fusion, alerte data steward |
| Conflit champ donné | Last-write-wins horodaté + audit |
| Conflit légal (ex. 2 actes de naissance) | Escalade manuelle Officier État Civil |

### 6.4 Garanties d'Intégrité
- Chaque enregistrement signé par opérateur (clé privée HSM kit)
- Horodatage local + horodatage central à la sync (double timestamp)
- Hash chaîné (style Merkle/blockchain légère) → détection altération

---

## 7. Protocoles Réseau de Sync

Par ordre de préférence :
1. Fibre / 4G LTE (quand dispo)
2. 3G / Edge
3. Satellite (Starlink, Eutelsat) pour zones blanches
4. Sneakernet (transport physique disque chiffré + courrier sécurisé) — solution de dernier recours

---

## 8. Sécurité Offline

- Disques **LUKS** + clé dérivée passphrase opérateur + token HSM
- **Effacement à distance** possible si kit signalé volé
- **Auto-destruction logique** après 30 j sans authentification
- **Tamper detection** physique (sceau électronique)
- **Quarantaine** à la sync — anti-malware scan avant ingestion

---

## 9. Formation Terrain

- Module obligatoire 5 jours pour tout opérateur
- Certification opérateur SNISID renouvelable annuellement
- Manuel terrain en français + créole
- Hotline support 24/7 (USSD + voix)

---

## 10. KPI Offline

| KPI | Cible |
|-----|-------|
| Autonomie sans connectivité | ≥ 30 jours |
| Taux de sync réussie | ≥ 98 % |
| Latence sync (du terrain → central) | ≤ 24h en moyenne |
| Taux de conflit nécessitant arbitrage humain | ≤ 0,5 % |
| Disponibilité enrôlement offline | ≥ 99 % du temps opérationnel |

---
*Fin du document — Étape 6/16*
