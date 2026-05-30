# 🔧 RUNBOOK — Corrupted Datasets Recovery

**ID** : RB-ANL-003
**Sévérité max** : CRITICAL
**Propriétaire** : Data Engineering + Data Governance
**SLA résolution** : 4h

---

## 1. SYMPTÔMES

- DQ score chute brutale (< 60)
- Anomalies massives détectées (Great Expectations FAIL)
- Plaintes utilisateurs (chiffres incohérents)
- Modèles IA produisent prédictions aberrantes

## 2. IMPACT

- Décisions erronées potentielles
- Doit être traité en **PRIORITÉ ABSOLUE** si dataset Gold/Platinum

## 3. DIAGNOSTIC

```bash
# Vérifier dernier checkpoint Delta
spark-sql -e "DESCRIBE HISTORY delta.\`s3a://snisid-lake/gold/<table>\` LIMIT 20;"

# Identifier la version corrompue
spark-sql -e "SELECT version, timestamp, operation, userName \
              FROM (DESCRIBE HISTORY delta.\`s3a://...\`) \
              ORDER BY version DESC;"

# Vérifier suite DQ
great_expectations checkpoint run gold_<table>_checkpoint

# Lineage
marquez query --dataset gold.<table> --upstream
```

## 4. PROCÉDURE DE REMÉDIATION

### 4.1 Isolement immédiat
```sql
-- Marquer le dataset comme STALE dans le catalogue
UPDATE openmetadata.tables
SET status = 'QUARANTINED', message = 'Corruption détectée — investigation'
WHERE fqn = 'lakehouse.gold.<table>';
```
- Bannière rouge sur dashboards consommateurs
- Désactivation jobs ML downstream

### 4.2 Time-travel Delta (rollback dataset)
```sql
-- Restaurer version saine
RESTORE TABLE delta.`s3a://snisid-lake/gold/<table>`
TO VERSION AS OF <last_good_version>;

-- Vérifier
SELECT COUNT(*) FROM delta.`s3a://snisid-lake/gold/<table>`;
```

### 4.3 Iceberg snapshot rollback
```sql
CALL system.rollback_to_snapshot('gold.<table>', <snapshot_id>);
```

### 4.4 Restauration depuis Bronze/Silver (si médaillon intact)
```bash
airflow dags trigger rebuild_gold_<table> \
  --conf '{"from_layer": "silver", "start": "<ts>", "end": "<ts>"}'
```

### 4.5 Restauration depuis backup MinIO/Ceph
```bash
mc cp --recursive snisid-backup/lake/gold/<table>/<date>/ \
                  snisid-minio/lake/gold/<table>/
```

### 4.6 Cas extrême : restauration source PostgreSQL
1. Restaurer dump WAL au point T-1
2. Re-jouer CDC depuis ce point
3. Reconstruire médaillon complet

## 5. VÉRIFICATION

- [ ] DQ score ≥ 90
- [ ] Suite Great Expectations PASS
- [ ] Row count cohérent avec historique (±5 %)
- [ ] Échantillon métier validé par data steward
- [ ] Lever quarantaine dans catalogue
- [ ] Bannière retirée
- [ ] Modèles downstream relancés

## 6. POST-MORTEM

**Obligatoire** quelle que soit la sévérité.
Analyser :
- Cause racine (bug code ? source corrompue ? action humaine ?)
- Pourquoi la DQ n'a pas bloqué en amont ?
- Renforcer tests Great Expectations
- Améliorer alerting fraîcheur / volume
