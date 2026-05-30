# 🎛️ NATIONAL DECISION DASHBOARD SYSTEM

> **Objectif** : Cockpit décisionnel national — vue consolidée pour tous les décideurs.

---

## 1. CAPACITÉS

| Domaine | Support |
|---------|:-------:|
| Presidential dashboards | ✅ |
| Ministry dashboards | ✅ |
| Crisis dashboards | ✅ |
| Regional dashboards | ✅ |

---

## 2. HIÉRARCHIE DES COCKPITS

```
🏛️ PRÉSIDENCE
   └── Cockpit Présidentiel (vue souveraine consolidée)
        ├── 🏛️ Ministères (10+)
        │     ├── Intérieur
        │     ├── Justice
        │     ├── Santé Publique (MSPP)
        │     ├── Éducation
        │     ├── Finances
        │     ├── Affaires Étrangères
        │     ├── Communication
        │     ├── Économie & Commerce
        │     ├── Agriculture
        │     └── Défense
        ├── 🗺️ Régions (10 départements)
        └── 🚨 Crises (actif sur déclenchement)
```

---

## 3. COCKPIT PRÉSIDENTIEL

```
┌────────────────────────────────────────────────────────────────┐
│   🇭🇹  COCKPIT PRÉSIDENTIEL — SNISID                            │
│   Mise à jour : 25 mai 2026 14:32 (auto-refresh 60s)            │
├────────────────────────────────────────────────────────────────┤
│  POPULATION IDENTIFIÉE        COUVERTURE NATIONALE              │
│  9 482 117 / 11 800 000        80.4 %                           │
│  +12 487 cette semaine         ████████████████░░░░             │
├────────────────────────────────────────────────────────────────┤
│  SERVICES RÉGALIENS (SLO 99.9%)                                 │
│   🟢 Identification    99.97 %                                  │
│   🟢 État civil        99.94 %                                  │
│   🟡 Registre électoral 99.71 %                                 │
│   🟢 Vérification CIN  99.99 %                                  │
├────────────────────────────────────────────────────────────────┤
│  RISQUES NATIONAUX (NRIC)                                       │
│   Cyber       🟢 0.18      Fraude       🟡 0.42                 │
│   Opérations  🟢 0.21      Infrastructure 🟢 0.15               │
│   Menaces nat. 🟡 0.38                                          │
├────────────────────────────────────────────────────────────────┤
│  ALERTES CRITIQUES (24h)                                        │
│   🔴 0    🟠 3    🟡 12                                         │
├────────────────────────────────────────────────────────────────┤
│  CARTE NATIONALE                  PRÉVISIONS 30j                │
│   [Heatmap GEOINT]                Demande CIN: ↑ +18 %          │
│                                   Bureaux saturés: 4 régions    │
└────────────────────────────────────────────────────────────────┘
```

---

## 4. COCKPIT MINISTÉRIEL (template)

Panneaux standard :
- KPI ministère (volumes, SLA, satisfaction)
- Workflows en cours / backlog
- Performances par bureau / agent
- Budget et consommables
- Alertes spécifiques au domaine
- Comparaison période précédente
- Prévisions sectorielles

---

## 5. COCKPIT CRISE (activable)

Mode crise déclenche un cockpit étendu :
- Carte impact temps réel
- Capacités vs besoins
- Continuité services régaliens
- Journal décisions
- Communications gouvernementales
- Statut équipes terrain

---

## 6. COCKPIT RÉGIONAL (départemental)

Pour chacun des 10 départements :
- Démographie & couverture identification
- Bureaux & agents
- Demandes en cours
- Performance régionale vs nationale
- Risques locaux (fraude, infra, crise)

---

## 7. PERSONNALISATION & ACCÈS

| Profil | Vue par défaut | Permissions |
|--------|----------------|-------------|
| Président | Cockpit présidentiel | Lecture totale, drilldown |
| Ministre | Cockpit ministère + national synthétique | Lecture domaine + public |
| Directeur SNISID | Tous cockpits | Lecture + actions |
| Chef région | Cockpit régional | Lecture région + public |
| Analyste NRIC | Cockpits risque/crise | Lecture + annotations |

Toutes consultations loguées (audit).

---

## 8. SUPPORTS

| Support | Détail |
|---------|--------|
| Web responsive | Superset + Grafana embed |
| Mobile sécurisé | App SNISID Decision (iOS/Android MDM) |
| Mur d'écrans war-room | Mode kiosk |
| Briefing PDF auto | Quotidien matin pour décideurs |
| Voix / assistant IA | Demande "Lis-moi les alertes rouges" |

---

## 9. FIABILITÉ

| Indicateur | Cible |
|------------|-------|
| Disponibilité cockpit présidentiel | > 99.99 % |
| Latence chargement | < 2 s |
| Fraîcheur données critiques | < 60 s |
| Backup cockpit | mode dégradé HTML statique cache |

---

## 10. PRINCIPE DESIGN

- Information hiérarchisée (3 niveaux max par écran)
- Code couleur unifié 🟢 🟡 🟠 🔴
- Toujours afficher : timestamp, source, niveau confiance
- Drilldown 1 clic vers détail
- Pas de surcharge cognitive
