# VOLUME 4 : Opérations Nationales et SRE (National Operations)
## Infrastructure de Production Souveraine — SNISID

L'architecture la plus résiliente ne sert à rien si elle n'est pas exploitée par des processus de niveau militaire. Les équipes d'ingénierie opèrent sous la doctrine du "Site Reliability Engineering" (SRE) adaptée aux contraintes étatiques d'Haïti.

---

## 📊 CHAPITRE 1 : SRE ET ERROR BUDGETS

Le SNISID ne tolère aucune interruption imprévue de l'état civil.

### 1.1 SLI (Service Level Indicators) et SLO (Service Level Objectives)
*   **SLO Disponibilité API d'Enrôlement :** 99.99% (Uptime requis: Seulement 52 minutes d'arrêt tolérées par an).
*   **SLO Latence Base de Données :** 95% des requêtes de lecture doivent s'exécuter en moins de 10ms (P95 < 10ms).
*   **Error Budget :** Si l'équipe de développement consomme son "budget d'erreur" (ex: trop d'échecs de déploiement en production ayant causé des indisponibilités), les déploiements de nouvelles fonctionnalités sont bloqués. L'équipe est contrainte de faire uniquement du "Bug Fixing" et de la stabilisation jusqu'au mois suivant.

---

## 📖 CHAPITRE 2 : RUNBOOKS ET PLAYBOOKS

Pour garantir des réactions prévisibles H24 par les équipes d'astreinte, toute alerte technique critique déclenche un "Runbook" (Procédure étape par étape).

### 2.1 Escalation Chain (Chaîne d'Escalade Automatisée)
Utilisation d'outils type PagerDuty / Grafana OnCall.
*   **Temps T = 0 :** Alerte Prometheus (ex: Utilisation disque K8s > 90%). Notification Slack/Teams au SRE L1 de garde.
*   **Temps T + 15 min :** Pas d'acquittement. Appel vocal automatisé au SRE L2.
*   **Temps T + 30 min :** Pas de résolution. Escalade à l'Incident Commander (IC) et appel vocal au Directeur Technique.

### 2.2 Automatisation de la Réponse (Auto-Remediation)
Le SOC et le NOC automatisent la résolution des problèmes courants :
*   Si le CPU d'un pod de la passerelle API dépasse 85%, le contrôleur HPA (Horizontal Pod Autoscaler) lance automatiquement de nouveaux conteneurs sans intervention humaine.

---

## 📜 CHAPITRE 3 : GOUVERNANCE OPÉRATIONNELLE ET CONFORMITÉ

### 3.1 Policy as Code (Kyverno / OPA Gatekeeper)
L'erreur humaine est la première cause de panne. L'infrastructure bloque activement les mauvaises pratiques.
*   Un développeur tente de déployer un conteneur en mode `root` ? Le cluster K8s rejette le fichier YAML automatiquement.
*   Un conteneur provient d'un registre public (ex: DockerHub) au lieu du registre souverain chiffré ? Déploiement bloqué.

### 3.2 Chaos Engineering (Test de Résilience Continu)
*   Des scripts "Chaos Monkey" tuent aléatoirement des pods Kubernetes de l'état civil pendant la journée pour s'assurer que les développeurs ont bien implémenté l'auto-healing et la redondance.
*   L'équipe simule physiquement la coupure de l'alimentation électrique du Datacenter de Port-au-Prince deux fois par an pour valider le basculement Active-Active sur Cap-Haïtien.
