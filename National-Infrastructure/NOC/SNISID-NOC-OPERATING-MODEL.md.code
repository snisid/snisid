---
# ============================================================
# SNISID-Infra — National Operations Center (NOC)
# Modèle Opérationnel 24/7 et Escalade
# Document ID: SNISID-NOC-001
# Version: 1.0.0
# ============================================================

## 1. RÔLE DU NOC (National Operations Center)

Le NOC est la salle de contrôle (Control Room) physique du gouvernement numérique. Situé de manière sécurisée (Zone 2) dans le Datacenter Primaire, il abrite les équipes SRE (Site Reliability Engineers) et Opérateurs Réseau (NetOps) travaillant en rotation 3x8 (24 heures sur 24, 7 jours sur 7).

## 2. MODÈLE OPÉRATIONNEL (Tiering)

- **Tier 1 (Frontline) :** Opérateurs "Eyes on Glass". Surveillent les écrans Grafana. Acquittent les alertes de bas niveau (ex: "Batterie Edge Jérémie à 20%").
- **Tier 2 (Incident Response) :** Administrateurs Systèmes/Réseaux. Interviennent sur les pannes complexes (ex: "BGP Route flapping sur le lien Digicel").
- **Tier 3 (Subject Matter Experts) :** Architectes SNISID, Ingénieurs K8s/DBA. Contactés (PagerDuty) uniquement pour les crises majeures (ex: "CockroachDB Split Brain").

## 3. SLA MANAGEMENT ET ESCALADE D'INCIDENT (Incident Command System)

### Scénario : Perte totale de connectivité avec le département du Sud.
1. **T0 :** Prometheus détecte une perte de ping/BGP sur le Edge Node "Cayes" et déclenche l'alerte `EdgeOffline_CRITICAL`.
2. **T0+5m :** Tier 1 valide que ce n'est pas un faux positif et crée un ticket JIRA.
3. **T0+15m :** Tier 2 appelle le FAI (Natcom) et constate une coupure physique de la fibre optique nationale.
4. **T0+20m :** L'incident est escaladé au Chef de Quart (Incident Commander). Le "Crisis Operations Model" est activé.
5. **T0+30m :** Basculement ordonné du noeud "Cayes" sur le lien VSAT Satellite secondaire. Le service est rétabli en mode dégradé (Bande passante réduite).

---
*Document ID: SNISID-NOC-001 | Approuvé par: Directeur des Opérations (NOC)*
