# SNISID National Cutover Technical Runbook
**Code de Projet :** SNISID-CUT-2026  
**Propriétaire :** Directeur Technique du SNISID / Responsable Infrastructure  
**Version :** 1.4  

---

## 1. Calendrier d'Exécution (T-24h à T+48h)

Ce runbook décrit la séquence exacte d'opérations chirurgicales à exécuter par l'équipe d'infrastructure durant la fenêtre de bascule nationale.

```
       [T-24h] Préparation & Sauvegardes
          |
       [T-2h] Gel de l'ancienne base ONI (Read-Only)
          |
       [T-0] Point d'Inflexion (Bascule DNS & Routage API)
          |
       [T+2h] Test d'Intégrité & Validation ABIS
          |
       [T+24h] Hypercare & Surveillance Active
```

---

## 2. Séquence d'Opérations Étape par Étape

### PHASE A : Préparation Finale (T-24h à T-6h)

- [ ] **A.1 Vérification de l'espace disque du Datacenter (DC-1 & DC-2)**
  *   *Commande :* `df -h /data`
  *   *Résultat attendu :* Espace libre $> 45\%$ (Minimum 1.2 TB disponible pour les indices d'identité).
  *   *Responsable :* Équipe Infra-Ops.

- [ ] **A.2 Exécution de la sauvegarde à froid de l'ancien système ONI**
  *   *Commande :* `pg_dump -h legacy-db.oni.gov.ht -U postgres -d onidb -F c -b -v -f /backups/oni_legacy_cold_T24.backup`
  *   *Résultat attendu :* Code de sortie `0`. Génération du hash SHA-256 de vérification.
  *   *Responsable :* DBA ONI.

- [ ] **A.3 Vérification de la synchronisation de l'heure NTP**
  *   *Commande :* `chronyc sources -v`
  *   *Résultat attendu :* Dérive de temps $< 1$ ms sur l'ensemble des clusters.
  *   *Responsable :* Équipe SecOps.

---

### PHASE B : Le Point de Bascule (T-2h à T-0)

- [ ] **B.1 Passage de l'ancien système ONI en mode Lecture Seule (Read-Only)**
  *   *Commande :* `psql -c "ALTER DATABASE onidb SET default_transaction_read_only = on;"`
  *   *Rôle :* Bloquer définitivement toute nouvelle écriture ou modification de carte d'identité sur l'ancienne base historique durant le cutover.
  *   *Responsable :* DBA ONI.

- [ ] **B.2 Synchronisation finale du lot différentiel (Le Delta)**
  *   *Description :* Lancement du pipeline d'extraction finale pour récupérer les fiches modifiées/créées durant les dernières 24 heures.
  *   *Commande :* `python3 /app/migration_factory/extract_delta.py --since-hours=24 --output=/data/deltas/oni_delta_last24h.json`
  *   *Responsable :* Équipe Migration Factory.

- [ ] **B.3 Redirection des enregistrements DNS de l'API d'authentification**
  *   *Action :* Mettre à jour la cible de l'API Gateway nationale d'identification de l'ancienne IP `10.5.120.40` vers l'IP virtuelle de l'API Gateway SNISID `10.50.10.10`.
  *   *Commande :* `terraform apply -target=cloudflare_record.api_gateway`
  *   *Responsable :* DevSecOps Lead.

---

### PHASE C : Validation et Go-Live (T-0 à T+4h)

- [ ] **C.1 Test d'intégrité de la passerelle API Gateway (gRPC / REST)**
  *   *Action :* Envoyer une requête d'authentification de test simulée pour vérifier que le SNISID répond à la place de l'ancienne base.
  *   *Commande :* `curl -X POST https://api.snisid.gouv.ht/v1/identity/verify -H "Authorization: Bearer TEST_TOKEN" -d '{"cin": "111-222-333-44"}'`
  *   *Résultat attendu :* Code `200 OK` avec le flag `verified: true` (réponse émise par le SNISID).
  *   *Responsable :* QA Team Lead.

- [ ] **C.2 Activation des Edge Nodes départementaux**
  *   *Action :* Envoyer l'ordre d'activation à distance aux 10 Local Edge Nodes du pays pour lancer le service d'enrôlement en mode Production.
  *   *Commande :* `ansible-playbook -i production_nodes.ini playbooks/activate_offline_nodes.yml`
  *   *Responsable :* SRE Lead.

---

## 3. Déclencheurs de Retour Arrière (Rollback Triggers)

L'équipe d'infrastructure doit impérativement déclencher la procédure de rollback si l'une des conditions suivantes est remplie dans les 4 heures post-cutover :

```
                                 ROLLBACK DECISION
                                 
+------------------------------------+       +------------------------------------+
|  Condition 1: API Response Latency |  OR   | Condition 2: Sync Failure Rate     |
|       (Average > 5000ms)           |       |       (Failure rate > 12%)         |
+------------------------------------+       +------------------------------------+
                   \                                          /
                    \                                        /
                     v                                      v
              +----------------------------------------------------+
              |               TRIGGER IMMEDIATE ROLLBACK           |
              |   - Run rollback playbook                          |
              |   - Revert DNS to 10.5.120.40 (Old server)         |
              |   - Put ONI Legacy back to Read-Write              |
              +----------------------------------------------------+
```

1. **Latence insupportable des requêtes d'API :** Temps moyen de réponse de l'API d'identification supérieur à 5000 ms pendant plus de 15 minutes consécutives (paralysie des banques et de la DIE).
2. **Échec de synchronisation des Edge Nodes :** Plus de 12% des transactions hors-ligne rejetées par le central en raison d'erreurs de chiffrement ou de signature matérielle.
3. **Instabilité du Cluster Central :** Redémarrage intempestif (CrashLoopBackOff) des pods Kubernetes hébergeant l'ABIS central.

### 3.1 Exécution de la Procédure de Retour Arrière (Rollback Execution Playbook)
En cas de décision de rollback, exécuter dans l'ordre strict :
1. **Rétablir les DNS vers l'ancien serveur :** `terraform destroy -target=cloudflare_record.api_gateway` (Bascule le trafic sous 2 minutes).
2. **Remettre l'ancienne base ONI en écriture :** `psql -c "ALTER DATABASE onidb SET default_transaction_read_only = off;"`
3. **Informer la cellule de crise :** Publier un message d'incident majeur sur le canal Slack de la War Room et envoyer un SMS d'alerte à l'ensemble du comité national.
