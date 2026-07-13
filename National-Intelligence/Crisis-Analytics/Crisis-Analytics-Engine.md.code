# 🚨 CRISIS ANALYTICS ENGINE

> **Objectif** : Piloter analytiquement les crises nationales en temps réel.

---

## 1. CAPACITÉS

| Fonction | Support | Description |
|----------|:-------:|-------------|
| Disaster impact analysis | ✅ | Évaluation dégâts (humains, infra, services) |
| Infrastructure degradation | ✅ | Suivi temps réel santé infra critique |
| Emergency capacity analysis | ✅ | Adéquation moyens / besoins |
| National continuity monitoring | ✅ | Continuité services régaliens |

---

## 2. SCÉNARIOS COUVERTS

- Cyclone / ouragan
- Séisme
- Inondation
- Crise sanitaire (épidémie)
- Crise sécuritaire
- Panne énergétique nationale
- Cyber-attaque majeure

---

## 3. ARCHITECTURE

```
┌──────────────────────────────────────────────────┐
│       WAR-ROOM CRISIS COCKPIT (Grafana XL)       │
└────────────────────┬─────────────────────────────┘
                     │
   ┌─────────────────┼──────────────────┐
   ▼                 ▼                  ▼
┌──────────┐  ┌──────────────┐  ┌────────────────┐
│ Impact   │  │ Infra Health │  │ Capacity       │
│ Analysis │  │ Monitor      │  │ Analyzer       │
└────┬─────┘  └──────┬───────┘  └────┬───────────┘
     │               │                │
     └───────┬───────┴────────────────┘
             ▼
   ┌────────────────────────┐
   │ Streaming (Kafka/Flink)│  + GEOINT + IoT terrain
   └────────────────────────┘
             │
             ▼
   ┌────────────────────────┐
   │ Lakehouse Gold Crisis  │
   └────────────────────────┘
```

---

## 4. DISASTER IMPACT MODEL

```python
# disaster_impact.py
def impact_score(zone):
    return (
        0.4 * normalize(zone.population_affected) +
        0.2 * normalize(zone.critical_infra_down) +
        0.2 * normalize(zone.service_outage_hours) +
        0.1 * normalize(zone.medical_load) +
        0.1 * normalize(zone.economic_loss_est)
    )

# Output : carte chaleur impact par commune
```

Visualisé en GEOINT temps réel.

---

## 5. INFRASTRUCTURE DEGRADATION

Sources :
- Capteurs IoT (énergie, télécom, eau)
- Probes synthétiques services SNISID
- Remontées terrain (agents mobiles)
- Partenaires (EDH, ANATEL équiv., DINEPA)

Indicateurs :
- `infra.power.grid_health_pct`
- `infra.telecom.coverage_pct`
- `infra.snisid.datacenter_status`
- `infra.transport.road_passability`

---

## 6. EMERGENCY CAPACITY ANALYSIS

| Ressource | Suivi | Source |
|-----------|-------|--------|
| Lits hôpitaux | Temps réel | MSPP |
| Stocks médicaments | Quotidien | PROMESS |
| Équipes terrain SNISID | Live | Apps mobiles |
| Carburant | 4h | Distributeurs partenaires |
| Eau potable | 6h | DINEPA |
| Connectivité satellite | Temps réel | Partenaires télécom |

```sql
-- Adéquation besoin vs capacité par département
SELECT d.departement,
       SUM(needs.affected_population) / NULLIF(SUM(cap.beds + cap.shelters), 0)
       AS pressure_ratio
FROM gold.disaster_needs needs
JOIN gold.emergency_capacity cap ON cap.departement = needs.departement
GROUP BY d.departement
ORDER BY pressure_ratio DESC;
```

---

## 7. NATIONAL CONTINUITY MONITORING

Services régaliens suivis 24/7 en mode crise :
- Identification (CIN, NIN lookup)
- État civil (naissances/décès)
- Élections (registre électoral)
- Sécurité (interpol checks)
- Santé (carnet vaccinal)

Statuts : 🟢 nominal, 🟡 dégradé, 🟠 partiel, 🔴 indisponible.

---

## 8. WAR-ROOM COCKPIT (panneaux)

1. Carte nationale impact (GEOINT)
2. Timeline événements
3. Capacités vs besoins par région
4. Continuité services régaliens
5. Alertes critiques en cours
6. Décisions prises (journal)
7. Communications gouvernementales
8. Prévisions évolution crise

---

## 9. PLAYBOOK ACTIVATION

Déclenchement automatique si :
- Risque NRIC = ROUGE
- Alerte officielle (DPC, MSPP, météo)
- Plus de 3 services régaliens dégradés > 30 min
- Demande Présidence / Premier Ministère

Actions auto :
- Ouverture canal war-room
- Snapshot lakehouse
- Activation runbooks crise
- Notification chaîne de commandement

---

## 10. POST-CRISE

- Rapport d'impact final (PDF signé)
- Lessons learned dans repository
- Mise à jour modèles prédictifs
- Révision capacités et runbooks
