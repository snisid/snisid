# SNISID Runbook — Failed Migration Recovery Procedure
**Code de Procédure :** SNISID-RB-02  
**Statut :** Approuvé  
**Audience :** Ingénieurs Migration Factory et DBA Centraux  

---

## 1. Objectif

Ce runbook décrit la procédure d'intervention immédiate en cas d'échec de chargement d'un lot d'importation (*Migration Batch*), de détection d'une corruption de données massive, ou de défaillance majeure du moteur de réconciliation démographique durant le traitement de l'historique ONI.

---

## 2. Diagramme de Décision Opérationnelle

```
                     MIGRATION BATCH FAILURE DETECTED
                                    |
            +-----------------------+-----------------------+
            | (Anomalie < 5% du lot)                        | (Anomalie > 5% / Corruption)
            v                                               v
[Mettre en Quarantaine (Q1/Q2)]                   [STOP & TRIGGER ROLLBACK]
- Isoler les fiches corrompues                    - Suspendre le pipeline ETL
- Traitement manuel asynchrone                    - Exécuter le script de retour arrière
- Continuer le traitement du reste du lot         - Nettoyer la base de staging
```

---

## 3. Étapes de Résolution d'Incident

### ÉTAPE 1 : Identification de l'Origine et Isolement (Triage)
1. **Analyser les Logs de la Migration Factory :** Lire les logs d'erreurs du lot affecté via Grafana Loki ou en interrogeant directement la table d'audit de migration :
   `SELECT batch_id, error_code, count(*) FROM snisid_audit.migration_log WHERE status = 'FAILED' GROUP BY batch_id, error_code;`
2. **Qualifier la Gravité :**
   *   *Si Erreurs de format mineures (ex: format de date mal interprété pour quelques individus) :* Les fiches sont isolées en table de quarantaine. Le lot est validé pour le reste des dossiers sains.
   *   *Si Corruption de schéma (ex: décalage de colonnes, erreur d'encodage de caractères accentués sur tout le fichier) :* Bloquer le pipeline.

### ÉTAPE 2 : Procédure de Nettoyage et Restauration (Rollback du Lot)
Si le lot a été partiellement inséré en base de production avant d'échouer, exécuter l'annulation transactionnelle :
1. **Identifier le point de restauration du lot :** Récupérer le `batch_id` fautif.
2. **Exécuter la procédure d'annulation d'insertion :**
   `python3 /app/migration_factory/rollback_batch.py --batch-id=MIG-BATCH-20260525-04`
   *Ce script va supprimer de manière atomique toutes les fiches rattachées à ce lot spécifique et réinitialiser les index de séquence d'IUI.*
3. **Valider le nettoyage :** S'assurer qu'aucun résidu orphelin ne subsiste :
   `SELECT count(*) FROM snisid_prod.citizens WHERE migration_batch_id = 'MIG-BATCH-20260525-04';`
   *Le résultat de la requête doit être strictement égal à `0`.*

### ÉTAPE 3 : Correction de l'ETL et Réexécution
1. **Corriger le parseur :** Si l'erreur provient d'un encodage de caractères, mettre à jour le script d'assainissement de la *Migration Factory* (ex: appliquer le convertisseur unicode `DataCleansingProgram.clean_name`).
2. **Tester sur échantillon :** Exécuter le pipeline sur un sous-ensemble de 100 fiches de test.
3. **Relancer la migration globale :** Après validation par le responsable technique, planifier la réexécution du lot corrigé.
