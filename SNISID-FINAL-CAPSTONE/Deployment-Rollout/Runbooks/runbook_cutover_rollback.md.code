# SNISID Runbook — Cutover Rollback Execution Playbook
**Code de Procédure :** SNISID-RB-03  
**Statut :** Approuvé  
**Audience :** Lead SRE, Lead DevSecOps et Administrateurs Système Centraux  

---

## 1. Objectif

Ce runbook décrit la séquence technique exacte pour annuler la bascule nationale (Cutover) et restaurer l'ancienne infrastructure d'identification (ONI Legacy) en mode de production actif si le Go-Live du SNISID échoue de manière critique.

---

## 2. Procédure Technique d'Annulation (Rollback)

```
                            ROLLBACK SEQUENCE
                            
[Step 1: Revert DNS / Routing] -> [Step 2: Re-enable Old Writes] -> [Step 3: Differential Recovery]
```

### ÉTAPE 1 : Redirection Immédiate du Routage Réseau et DNS (T+5 minutes)
1. **Rétablir les enregistrements DNS d'API :**
   Se connecter au serveur DNS maître souverain ou à la passerelle Cloudflare et rediriger le nom d'hôte `api.snisid.gouv.ht` vers l'IP de l'ancienne passerelle ONI (`10.5.120.40`).
   *   *Commande Terraform :*
       `cd /infra/terraform/dns && terraform apply -var="use_legacy_identity=true"`
2. **Forcer le rafraîchissement des caches DNS (DNS Flush) :**
   `ansible -i gateways.ini -m shell -a "systemd-resolve --flush-caches" all`

### ÉTAPE 2 : Réactivation des Droits d'Écriture sur l'Ancien Système (T+10 minutes)
1. **Passer la base historique de "Lecture Seule" à "Lecture-Écriture" :**
   Se connecter en SSH au serveur de base de données de l'ONI historique et exécuter :
   `psql -h legacy-db.oni.gov.ht -U postgres -d onidb -c "ALTER DATABASE onidb SET default_transaction_read_only = off;"`
2. **Redémarrer les services d'arrière-plan ONI :**
   `ansible-playbook -i legacy_servers.ini playbooks/start_legacy_apps.yml`

### ÉTAPE 3 : Extraction Différentielle des Données Capturées sur le SNISID (T+30 minutes)
Si la plateforme SNISID a fonctionné pendant plusieurs heures avant l'annulation, certains citoyens ont été enrôlés sur la nouvelle plateforme. Leurs données ne doivent pas être perdues.
1. **Extraire les dossiers créés durant la fenêtre d'activité du SNISID :**
   `python3 /app/cutover/extract_snisid_production_delta.py --since-time="2026-05-25 12:00:00" --output=/data/recovery_delta.json`
2. **Importer le delta dans l'ancienne base ONI :**
   `python3 /app/cutover/inject_into_legacy.py --input=/data/recovery_delta.json --db-host=legacy-db.oni.gov.ht`
3. **Valider la cohérence :** Vérifier que le nombre total d'inscriptions sur la base historique correspond exactement au bilan attendu.

---

## 3. Communication Post-Incident
1. **Notification Interne :** L'Incident Commander publie un rapport d'incident majeur sur le canal Slack de la cellule de crise.
2. **Information Publique :** Le Desk Communication diffuse un communiqué de presse expliquant le report de la mise en service de la nouvelle carte d'identité en raison d'une "maintenance préventive renforcée visant à garantir la sécurité absolue des données citoyennes."
