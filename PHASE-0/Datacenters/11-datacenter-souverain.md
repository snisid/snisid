# 🏢 SNISID — Sovereign Datacenter Strategy

**Document N° :** SNISID-DC-011
**Étape Phase 0 :** 11/16
**Principe :** *Les données critiques doivent rester sous contrôle national.*

---

## 1. Doctrine

Toute donnée à valeur d'**identité, état civil, biométrie, justice, sécurité nationale** doit :
- Être **stockée physiquement** sur le territoire haïtien
- Être **opérée** par du personnel haïtien habilité
- Être **possédée juridiquement** par l'État haïtien
- Ne **jamais** transiter par un hyperscaler étranger sans chiffrement souverain (clés HSM Haïti)

Les usages non sensibles (sites web publics, archives ouvertes) peuvent recourir à du cloud externe.

---

## 2. Topologie Cible

```
                      ┌────────────────────────┐
                      │   PRIMARY DATACENTER   │
                      │   Port-au-Prince       │
                      │   Tier III, ~500 m²    │
                      └───────────┬────────────┘
                                  │ Liaison fibre dédiée
                                  │ + Microwave backup
                                  ▼
                      ┌────────────────────────┐
                      │   DR DATACENTER        │
                      │   Cap-Haïtien          │
                      │   Tier III, hot standby│
                      └───────────┬────────────┘
                                  │
        ┌─────────────────────────┼────────────────────────┐
        │                         │                        │
   ┌────▼─────┐              ┌────▼─────┐            ┌────▼─────┐
   │ EDGE     │     ...      │ EDGE     │            │ EDGE     │
   │ Cayes    │              │ Jacmel   │            │ Hinche   │
   └────┬─────┘              └────┬─────┘            └────┬─────┘
        │                         │                        │
        └─────── 10 nodes départementaux Tier I/II ────────┘
                                  │
                            Offline kits terrain (50+)
```

---

## 3. Datacenter Primaire (PaP)

| Spec | Valeur |
|------|--------|
| Tier (Uptime Institute) | **III** minimum (visée IV) |
| Superficie utile | 500 m² |
| Puissance IT | 500 kW (évolutif 1 MW) |
| PUE cible | < 1,6 |
| Refroidissement | Free cooling + chillers redondants N+1 |
| Énergie | Réseau + onduleurs N+1 + groupes électrogènes (72h autonomie) + solaire 200 kWc |
| Connectivité | 3 opérateurs distincts, fibres entrant par 2 chemins physiques |
| Sécurité physique | Périmètre clôturé, mantrap, biométrie, vidéosurveillance, garde 24/7 |
| Conformité | ISO 27001, ISO 22301, conformité TIA-942 |
| Résistance | Sismique zone 4 (Haïti), cyclonique cat 5 |

**Localisation candidate :** zone élevée hors plaine inondable, proche aéroport mais hors trajectoire approche.

---

## 4. Datacenter de Reprise (Cap-Haïtien)

- **Hot standby** avec réplication synchrone des bases critiques (RPO ~0)
- Capacité ≥ 70 % du primaire pour assurer continuité dégradée
- Distance > 250 km du primaire (résilience séisme)
- Bascule automatique testée semestriellement
- Sert aussi de site de calcul pour traitements lourds (AFIS batch)

---

## 5. Edge Nodes (10 départements)

Un node par chef-lieu de département :

| Spec | Valeur |
|------|--------|
| Tier | I/II (selon chef-lieu) |
| Superficie | ~30 m² (rack room dans bâtiment public) |
| Puissance | 10-30 kW |
| Compute | 3-5 serveurs (clusters K3s/RKE2) |
| Stockage | 50-100 To NVMe chiffré |
| Réseau | Fibre + 4G + satellite backup |
| Fonctions | Cache lecture, ingestion enrôlements, AFIS local subset, sync hub |

---

## 6. Offline Nodes Terrain

Voir document **Offline-First Strategy (06)**. Synthèse :
- 50+ kits déployables (montée en charge)
- Chaque kit = mini-DC mobile autonome 30+ jours
- Sync vers edge node de rattachement

---

## 7. Modèle Cloud Souverain

Stack open-source recommandée pour le datacenter :
- **OpenStack** ou **Proxmox VE** pour la virtualisation
- **Ceph** pour le stockage distribué
- **Kubernetes** par-dessus pour les charges cloud-native
- **MinIO** pour stockage objet S3-compatible
- **PostgreSQL Patroni** pour HA SGBD

> Évite la dépendance VMware/AWS/Azure/GCP pour le cœur de SNISID.

---

## 8. Stratégie Sauvegardes

| Niveau | Fréquence | Rétention | Localisation |
|--------|-----------|-----------|--------------|
| Snapshots applicatifs | Toutes les heures | 7 jours | PaP |
| Backups incrémentaux | Quotidien | 30 jours | PaP + Cap-Haïtien |
| Backups complets | Hebdomadaire | 1 an | PaP + Cap-Haïtien |
| Archives long terme | Mensuel | 10 ans | Cap-Haïtien + offline WORM |
| Disaster offline copy | Trimestriel | 3 ans | Coffre-fort bancaire chiffré |

**Règle 3-2-1-1-0** : 3 copies, 2 supports, 1 hors-site, 1 immuable, 0 erreur restoration testée.

---

## 9. Énergie & Résilience Climatique

- Mix énergétique : réseau + solaire + groupes diesel/HVO
- Stockage batteries Li-ion 4h + 72h diesel
- Suivi consommation et bilan carbone publié annuellement
- Cible : 30 % renouvelable d'ici 2028, 50 % d'ici 2030

---

## 10. Modèle de Gouvernance Datacenter

- Exploitation par une **entité publique dédiée** (filiale de l'AND ou contrat de gestion)
- Personnel haïtien habilité (clearance national)
- Audits semestriels
- Transparence : rapport public annuel (sauf éléments classifiés)

---

## 11. Phasage de Construction

| Phase | Période | Action |
|-------|---------|--------|
| 1 | 2026 Q3-Q4 | Choix site PaP + études + appel d'offres |
| 2 | 2027 Q1-Q3 | Construction PaP (génie civil + IT) |
| 3 | 2027 Q4 | Mise en service PaP + chargement données |
| 4 | 2028 Q1-Q2 | Construction Cap-Haïtien (DR) |
| 5 | 2028 Q3-Q4 | Edge nodes (10 départements) |
| 6 | 2029-2030 | Densification + capacité extension |

---

## 12. Budget Indicatif (USD)

| Poste | Estimation |
|-------|-----------|
| Datacenter PaP (génie civil + IT) | 15-25 M$ |
| Datacenter DR Cap-Haïtien | 10-15 M$ |
| 10 Edge nodes | 5-8 M$ |
| 50 Offline kits | 3-5 M$ |
| OPEX annuel exploitation | 4-6 M$/an |

> Financement mixte : budget national + BID/BM/UE.

---
*Fin du document — Étape 11/16*
