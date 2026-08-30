# SNISID National Cutover Strategy
## Stratégie de Bascule Nationale et de Minimisation du Downtime

---

## 1. Cadre Général du National Cutover

Le **Cutover** (bascule définitive) est l'étape technique la plus sensible de la phase 15. Elle consiste à transférer l'autorité décisionnelle et la production active d'identité de l'ancien système de l'ONI vers la nouvelle infrastructure souveraine du SNISID. Ce document définit la stratégie de bascule, les portes de validation, et les mécanismes de retour arrière (*Rollback*).

```
                      NATIONAL CUTOVER TIMELINE (72 HOURS)
                      
T-24h               T-0 (Cutover Point)                  T+24h               T+48h
  |----------------------|---------------------------------|-------------------|
Readiness Checks    Read-Only Block                   Integration Sync    Go-Live Active
- Backups Secured   - Stop old writes                 - Run NIRE          - Activate APIs
- Network Tunneling - DNS switch                      - Verify hashes     - Hypercare Start
```

---

## 2. Principes pour Minimiser l'Indisponibilité (Downtime Minimization)

La continuité des services administratifs et de sécurité (comme la vérification aux frontières ou l'authentification bancaire) exige un plan de bascule limitant le downtime à **moins de 4 heures** pour les services en lecture seule.

### 2.1 Les Stratégies Clés de Continuité
1. **Bascule Graduelle DNS (Blue-Green Deployment) :**
   Les API d'identification ne subissent pas d'arrêt brutal. L'ancien serveur ("Blue") et le SNISID ("Green") fonctionnent en parallèle. Le routage réseau (via proxy inverse Nginx ou serveurs DNS Cloudflare souverains) migre progressivement la charge (10% -> 50% -> 100%).
2. **Synchronisation d'Échauffement (Warm Synchronization) :**
   Avant la bascule, 99.9% de l'historique ONI a déjà été migré de manière asynchrone par la *Migration Factory*. Seul le différentiel des fiches créées durant les dernières 72 heures (le Delta) est injecté durant la fenêtre de bascule, réduisant le temps de traitement de plusieurs jours à quelques minutes.

---

## 3. Portes de Validation de Bascule (Validation Checkpoints)

Avant d'initier la bascule (Go/No-Go Decision), le Comité de Gouvernance Nationale doit valider la liste de contrôle d'aptitude suivante :

```
                  GO/NO-GO COMMITMENT GATE
                  
[Checkpoint 1: Backups Verified (100% integrity)] --------+
                                                          |
[Checkpoint 2: Network Latency (DC-1 to DC-2 < 12ms)] ----+---> [GO DECISION]
                                                          |
[Checkpoint 3: Regional Rollout Teams Certified (10/10)] -+
```

### 3.1 Tableau des Checkpoints de Secours

| Checkpoint | Validation Attendue | Rôle Sécuritaire | Seuil d'Échec (No-Go) |
| :--- | :--- | :--- | :--- |
| **Integrity-Check** | Sauvegarde à froid de la base ONI complétée et hash SHA-256 signé cryptographiquement. | Protection absolue des données historiques de l'état civil. | Écart de hash ou sauvegarde incomplète. |
| **Network-Check** | Tunnels IPSec redondants stables entre le DC Central et tous les Edge Nodes départementaux. | Prévention des déconnexions massives durant la bascule. | Plus de 2 Edge Nodes déconnectés. |
| **ABIS-Check** | Moteur biométrique d'identification 1-to-N en ligne avec un temps de réponse moyen $< 1.2$ sec. | Évitement de la paralysie des centres d'enrôlement dès la réouverture. | Latence moyenne $> 3.5$ sec. |
| **Legal-Check** | Décret d'autorisation de transfert de responsabilité souveraine signé par le Président du Conseil d'Administration du SNISID. | Couverture juridique de la bascule d'autorité. | Absence de décret officiel. |

---

## 4. Stratégie de Retour Arrière (Rollback Capability)

Si une anomalie bloquante non détectée survient après la bascule (ex : corruption de la base transactionnelle principale, défaillance réseau totale du DC central), la procédure de retour arrière est exécutée.

### 4.1 Niveaux de Rollback (N-L)

```
[Incident < T+4h]  ===> Niveau 1 : Restauration DNS (Bascule instantanée sur l'ancien serveur)
[Incident > T+24h] ===> Niveau 2 : Reconstructive Rollback (Extraction différentielle des fiches SNISID et réinjection dans l'ancien système)
```

1. **Niveau 1 : Restauration DNS (Instant DNS Restore) :**
   *   *Déclenchement :* Incident détecté dans les 4 premières heures post-cutover.
   *   *Mécanisme :* Redirection instantanée des serveurs DNS d'API d'authentification vers l'ancienne infrastructure (toujours en veille active). Le temps d'exécution est inférieur à 5 minutes sans aucune perte de données, car aucune nouvelle transaction d'envergure n'a été validée sur le SNISID.
2. **Niveau 2 : Rétablissement Reconstructif (Reconstructive Rollback) :**
   *   *Déclenchement :* Incident critique détecté après 24 heures de production sur le SNISID.
   *   *Mécanisme :* L'ancien système est réactivé. Les fiches d'enrôlement qui ont été capturées sur le SNISID durant les 24 heures d'activité sont extraites sous forme de lots XML différentiels, formatées pour s'adapter à l'ancienne base, puis importées manuellement. Cette procédure garantit **zéro perte de données d'identité**, même en cas de repli tardif.
