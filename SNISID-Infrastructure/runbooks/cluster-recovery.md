# SNISID RUNBOOK — Cluster Recovery (Kubernetes)
**Classification:** RESTREINT DEFENSE  
**Version:** 4.0.0  
**RTO cible:** < 45 minutes  
**Fréquence test:** Trimestriel

---

## 1. Prérequis

- Accès bastion national (MFA + certificat smartcard)
- `kubectl` version 1.28.x
- Vault token valide (`vault login -method=cert`)
- Journal d'opérations enregistré sur système SIEM

## 2. Scénarios de déclenchement

| Code | Scénario | Déclencheur |
|------|----------|-------------|
| K8S-001 | Perte master unique (2/3) | Alertmanager : `K8s_ControlPlaneDown` |
| K8S-002 | Perte quorum etcd | Vault / API indisponible, etcd logs `leader election failed` |
| K8S-003 | Corruption datastore etcd | `etcdctl endpoint health` retourne `false` |
| K8S-004 | Cluster-wide failure (incendie, blackout) | DR automatique ou manuel selon runbook DR |

## 3. Procédure K8S-001 : Remplacement master

### 3.1 Identifier le master défaillant
```bash
kubectl get nodes -l node-role.kubernetes.io/control-plane
# Master-02 NotReady
```

### 3.2 Retirer le master du cluster
```bash
# Sur bastion avec accès master survivants
MASTER="core-prod-t1-master-02"
kubectl drain ${MASTER} --ignore-daemonsets --delete-emptydir-data
kubectl delete node ${MASTER}
```

### 3.3 Reprovisionner via Terraform
```bash
cd /infrastructure/terraform/environments/core
terraform taint 'module.core_masters["02"]'
terraform apply -target='module.core_masters["02"]'
```

### 3.4 Join le nouveau master
```bash
# Récupérer join command depuis master survivant
ssh core-prod-t1-master-01 "sudo kubeadm token create --print-join-command --certificate-key \$(sudo kubeadm init phase upload-certs --upload-certs | tail -1)"
# Exécuter la commande sur le nouveau nœud provisionné
```

### 3.5 Validation
```bash
kubectl get nodes
kubectl get pods -n kube-system
etcdctl endpoint health --endpoints=https://10.1.10.11:2379,https://10.1.10.12:2379,https://10.1.10.13:2379 --cacert=/etc/kubernetes/pki/etcd/ca.crt
```

## 4. Procédure K8S-002 : Recovery etcd quorum

### 4.1 État du quorum
```bash
ETCDCTL_API=3 etcdctl member list \
  --endpoints=https://10.1.10.11:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key
```

### 4.2 Si 2 membres perdus (quorum rompu)
> **ATTENTION :** Cette procédure est destructive. Toute opération est loggée SIEM.

```bash
# 1. Choisir le membre avec le data le plus récent (raft index le plus élevé)
# 2. Forcer un nouveau cluster depuis ce nœud
sudo systemctl stop etcd

# Backup préalable (obligatoire)
etcdctl snapshot save /var/backups/etcd/emergency-$(date +%s).db \
  --endpoints=https://10.1.10.11:2379 \
  --cacert=/etc/kubernetes/pki/etcd/ca.crt \
  --cert=/etc/kubernetes/pki/etcd/server.crt \
  --key=/etc/kubernetes/pki/etcd/server.key

# 3. Reset membre sain comme nouveau cluster
etcdctl snapshot restore /var/backups/etcd/emergency-*.db \
  --name core-prod-t1-master-01 \
  --initial-cluster core-prod-t1-master-01=https://10.1.10.11:2380 \
  --initial-cluster-token snisid-etcd-recovery \
  --initial-advertise-peer-urls https://10.1.10.11:2380 \
  --data-dir /var/lib/etcd

# 4. Redémarrer etcd
sudo systemctl start etcd

# 5. Ajouter les nouveaux membres (reprovisionnés)
etcdctl member add core-prod-t1-master-02 --peer-urls=https://10.1.10.12:2380
# Puis kubeadm reset + join sur les nœuds remplacés
```

## 5. Procédure K8S-003 : Corruption etcd

```bash
# 1. Arrêter tous les etcd
ansible -i inventory/core.ini etcd -m systemd -a "name=etcd state=stopped"

# 2. Identifier le membre avec le MVCC le plus récent
for node in 11 12 13; do
  ETCDCTL_API=3 etcdctl --endpoints=https://10.1.10.${node}:2379 endpoint status --write-out=json
done

# 3. Restaurer depuis snapshot validé
# Snapshot le plus récent non corrompu (vérif hash)
SNAP=$(ls -t /var/backups/etcd/*.db | head -1)
sha256sum -c ${SNAP}.sha256 || exit 1

etcdctl snapshot restore ${SNAP} --name core-prod-t1-master-01 ... (cf K8S-002)

# 4. Reconstruire le cluster étape par étape
# 5. Valider avec : etcdctl endpoint status
```

## 6. Post-recovery

- [ ] Vérifier Vault unseal status
- [ ] Vérifier ArgoCD sync status (aucune dérive)
- [ ] Vérifier Ceph HEALTH_OK
- [ ] Lancer `kubectl get pods -A` — tous les pods doivent être Running/Completed
- [ ] Exécuter le script `validate-national-health.sh`
- [ ] Remplir le post-mortem dans le SIEM national
- [ ] Notifier le SOC National et l'IGC

---

*Ce runbook est testé trimestriellement. Toute modification nécessite validation IGC.*
