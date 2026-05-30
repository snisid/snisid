# SNISID RUNBOOK — Root CA Compromise / Destruction / Reconstruction
**Classification:** TOP-SECRET / EMERGENCE NATIONALE  
**Version:** 5.0.0  
**Trigger:** Détection compromission Root CA clé privée, destruction HSM Core, catastrophe nationale, ou attaque étatique sur PKI  
**RTO cible:** < 72 heures pour reconstruction partielle, < 7 jours pour pleine restauration

---

## 0. Déclencheurs d'urgence nationale

| Code | Scénario | Déclencheur | Niveau |
|------|----------|-------------|--------|
| PKI-E1 | **Compromission Root CA** | Clé privée extraite (théoriquement impossible avec HSM, mais attaque zero-day HSM, insider, ou extraction Shamir) | **CATASTROPHE** |
| PKI-E2 | **Destruction HSM Core** | Incendie / explosion / EMP / inondation salle blanche | **CATASTROPHE** |
| PKI-E3 | **Compromission massive Intermediate CA** | Clé Intermediate extraite + certificats frauduleux massifs émis | **CATASTROPHE** |
| PKI-E4 | **Shamir parts compromise** | 4+ parts récupérées par adversaire | **CATASTROPHE** |
| PKI-E5 | **Attaque étatique PKI** | APT nation-state avec accès réseau + physique PKI | **CATASTROPHE** |
| PKI-E6 | **Erreur cérémonie Root** | Erreur opérationnelle irréparable lors cérémonie | **CRITIQUE** |

---

## 1. Phase 1 — IMMÉDIAT (H+0 à H+1h) : Confinement

### 1.1 SOC National — Activation
```bash
# 1.1.1 Déclencher Cellule de Crise PKI Nationale
# Notification chaîne de commandement : SOC → IGC → Présidence → Justice → Défense

# 1.1.2 Isolation totale PKI Core DC
kubectl taint nodes -l snisid.gov/tier=tier-0 snisid.gov.emergency=root-compromise:NoSchedule
kubectl cordon -l snisid.gov/tier=tier-0
kubectl drain -l snisid.gov/tier=tier-0 --ignore-daemonsets --delete-emptydir-data --force

# 1.1.3 Isolation réseau Cilium — deny all PKI traffic
cat << 'EOF' | kubectl apply -f -
apiVersion: cilium.io/v2
kind: CiliumClusterwideNetworkPolicy
metadata:
  name: emergency-pki-lockdown
spec:
  endpointSelector:
    matchLabels:
      snisid.gov/tier: tier-0
  ingressDeny:
    - {}
  egressDeny:
    - {}
EOF

# 1.1.4 Sceller physiquement salle blanche (garde armée, aucun accès sans Présidence + IGC)
# 1.1.5 Confiscation / gel des Shamir parts (tous les sites — alerte diplomatique si nécessaire)
```

### 1.2 IGC — Évaluation initiale
- [ ] Vérification tamper-evident HSM Core (photos, logs, vidéos)
- [ ] Vérification logs Vault (dernières 24h) — anomalie issuance ?
- [ ] Vérification CRL actuelle — certificats frauduleux émis ?
- [ ] Vérification OCSP logs — requêtes anormales ?
- [ ] Vérification signatures récentes — documents frauduleux ?
- [ ] Premiers forensics HSM (si accessible) — read-only audit log HSM

### 1.3 Notification nationale
- [ ] Présidence (TOP-SECRET verbal)
- [ ] Conseil National Sécurité
- [ ] Justice (procédure judiciaire d'urgence si nécessaire)
- [ ] Ambassades (protection Shamir parts à l'étranger)
- [ ] Forces armées (protection sites critiques)
- [ ] Alliés stratégiques (notification PKI compromise si cross-trust)

---

## 2. Phase 2 — CONFINEMENT ACTIF (H+1h à H+6h) : Mitigation

### 2.1 Révocation massive (même si Root CA compromise — minimiser fenêtre)
```bash
# 2.1.1 Émettre CRL Emergency pour TOUS les Intermediate CAs
# Cette CRL signée par Root CA (si HSM encore sain) ou par procédure offline

# 2.1.2 Si HSM Core détruit : utiliser HSM DR (Root CA partition backup)
# HSM-DR-01 : activer partition SNISID_ROOT_CA_BACKUP (4-of-6 custodians DR)
# Signer CRL Emergency avec clé Root CA backup

# 2.1.3 Publication CRL Emergency
s3cmd put --acl-public --mime-type="application/pkix-crl" \
  /tmp/emergency-root-revocation.crl \
  s3://snisid-crl-distribution/emergency/root-ca-revoked.crl

# 2.1.4 Broadcast Kafka topic national
curl -X POST https://kafka.national.snisid.gouv.local/topics/pki.emergency.revocation \
  -H "Authorization: Bearer ${EMERGENCY_KAFKA_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"event":"root_ca_compromise","timestamp":"'$(date -u +%Y-%m-%dT%H:%M:%SZ)'","crl_url":"https://crl.snisid.gouv.local/emergency/root-ca-revoked.crl","action":"REVOKE_ALL_INTERMEDIATES"}'

# 2.1.5 Edge nodes sync emergency CRL via USB burst / satellite
# Déployer équipes terrain avec USB chiffrés contenant CRL + nouveau trust anchor temporaire
```

### 2.2 Arrêt des services dépendants
- [ ] **Arrêt immédiat Vault PKI Core** (kubectl delete deployment vault — namespace snisid-identity)
- [ ] **Arrêt cert-manager** (empêcher nouvelles émissions)
- [ ] **Arrêt Digital Signature Service** (pas de signature tant que confiance non restaurée)
- [ ] **Arrêt OCSP Core** (remplacé par CRL offline + notification manuelle)
- [ ] **Activation OCSP DR** si HSM DR sain
- [ ] **Notification citoyens** : "Service identité numérique temporairement suspendu — présentez pièce d'identité physique"

### 2.3 Forensics & investigation
- [ ] Équipe forensics IGC sur site (HSM Core, réseau, logs)
- [ ] Analyse SIEM : timeline exacte compromise
- [ ] Analyse vidéo surveillance salle blanche (dernières 72h)
- [ ] Interrogatoire custodians (protocole judiciaire)
- [ ] Analyse matérielle HSM (si possible — sinon destruction contrôlée + envoi constructeur)
- [ ] Analyse réseau : exfiltration données ?

---

## 3. Phase 3 — RECONSTRUCTION (H+6h à H+72h) : Nouvelle Root CA

### 3.1 Décision politico-juridique
- [ ] **Décret présidentiel** : autorisation reconstruction PKI nationale
- [ ] **Décision Justice** : validité juridique des signatures émises avant compromission
- [ ] **Décision IGC** : plan technique reconstruction
- [ ] **Notification internationale** : nouveau trust anchor SNISID (diplomatique)

### 3.2 Cérémonie de reconstruction Root CA (accélérée)
```
# 3.2.1 Rassembler 4+ custodians (si Shamir parts intactes) OU
#        Nommer 6 nouveaux custodians (si anciens compromise ou suspicion)
#        Custodians doivent être de confiance absolue, background check renforcé

# 3.2.2 Si HSM Core détruit :
#        - Activer HSM DR (Root CA partition) si sain
#        - OU générer NOUVELLE paire clé Root CA dans HSM DR (nouveau HSM commandé express)
#        - Nouveau HSM : vérification firmware, attestation constructeur, inspection IGC

# 3.2.3 Si compromission clé (extraction) :
#        - Générer NOUVELLE paire clé Root CA (ancienne clé irrévocablement compromise)
#        - Nouvelle Root CA : SNISID Root CA Nationale v2
#        - Nouvelle chaîne : Root v2 → nouveaux Intermediates v2
#        - Ancienne chaîne : publiquement révoquée, archivée forensics

# 3.2.4 Cérémonie accélérée (24h instead of 1 journée)
#        Phase 1 (H+6h) : Vérification salle DR / HSM DR / participants
#        Phase 2 (H+12h) : Génération nouvelle Root CA dans HSM DR
#        Phase 3 (H+18h) : Vérifications croisées + publication empreinte
#        Phase 4 (H+24h) : Shamir parts nouvelle génération → distribution 6 sites
```

### 3.3 Reconstruction chaîne PKI
```bash
# 3.3.1 Nouveaux Intermediate CAs (tous révoqués, tous reconstruits)
# Government CA v2 → certifié par Root CA v2
# Citizen CA v2 → certifié par Root CA v2
# Infrastructure CA v2 → certifié par Root CA v2
# Judicial CA v2 → certifié par Root CA v2
# Device CA v2 → certifié par Root CA v2

# 3.3.2 Génération certificats Intermediates dans HSM DR
# Commandes Vault PKI (nouveau mount : snisid-pki-infra-v2, etc.)
vault write snisid-pki-infra-v2/intermediate/generate/internal \
  common_name="SNISID Infrastructure CA v2" \
  ttl=43800h \
  key_type=ec \
  key_bits=384 \
  signature_algorithm=ecdsa-with-SHA384

# Signer par Root CA v2 (via HSM DR)
vault write snisid-root-ca-v2/root/sign-intermediate \
  csr=@/tmp/infra-v2.csr \
  common_name="SNISID Infrastructure CA v2" \
  ttl=43800h \
  use_csr_values=true

# 3.3.3 Distribution bundle nouveau trust anchor
# Edge nodes, citizen wallets, enrollment stations, government smartcards, judicial devices
# Méthodes : satellite burst + USB/SD + diplomatic courier (alliés) + broadcast radio national
```

### 3.4 Restauration services PKI
```bash
# 3.4.1 Redémarrer Vault PKI avec nouveaux mounts
# 3.4.2 Mettre à jour cert-manager ClusterIssuers (Root v2 certs)
# 3.4.3 Force-renewal TOUTES les certificats infrastructure (90j batch)
# 3.4.4 Redémarrer Istio mesh (SDS push nouvelle Root CA)
# 3.4.5 Redémarrer OCSP responders (nouvelle chaîne)
# 3.4.6 Redémarrer Digital Signature Service
```

---

## 4. Phase 4 — RESTAURATION CITOYENNE (H+72h à H+7j)

### 4.1 Renouvellement citoyen
- [ ] **Notification massive** : SMS, radio, TV nationale, agents terrain
- [ ] **Re-enrollment obligatoire** : tous les citoyens doivent obtenir nouveau certificat Citizen CA v2
- [ ] **Grace period** : 90 jours pour re-enrollment (anciens certificats Citizen CA v1 révoqués mais valables pour services basiques pendant grace period)
- [ ] **Campagne terrain** : mobile enrollment units déployées partout
- [ ] **Portail en ligne** : re-enrollment avec biométrie + ID (si infrastructure en ligne restaurée)

### 4.2 Renouvellement gouvernement
- [ ] **Nouveaux smartcards** pour tous les agents gouvernementaux
- [ ] **Nouvelles signatures officielles** : tous les actes post-compromission invalidés doivent être re-signés
- [ ] **Archives** : signatures pré-compromission marquées "V1 — validité historique à vérifier"

### 4.3 Renouvellement judiciaire
- [ ] **Nouvelles signatures judiciaires** : magistrats re-certifiés avec Judicial CA v2
- [ ] **Validation rétroactive** : jugements V1 valides si signés AVANT timestamp de compromission (TSA log)
- [ ] **Invalidation** : jugements/signatures post-compromission sans nouvelle chaîne = INVALIDES

### 4.4 Renouvellement international
- [ ] **Nouveaux accords de reconnaissance** avec pays alliés (Root CA v2 trust)
- [ ] **Notification INTERPOL / eIDAS équivalent** : ancien trust anchor SNISID révoqué
- [ ] **Mise à jour systèmes frontaliers** : nouveaux certificats device / enrollment

---

## 5. Phase 5 — POST-MORTEM & RENFORCEMENT

- [ ] **Rapport IGC** : causes, timeline, responsabilités, mesures (30 jours)
- [ ] **Rapport Justice** : procédures judiciaires si nécessaire (sabotage, trahison, négligence)
- [ ] **Audit externe souverain** : validité technique reconstruction
- [ ] **Renforcement HSM** : passage Level 4, geo-distribution +3 sites, air-gap renforcé
- [ ] **Nouvelle procédure Shamir** : 5-of-9 au lieu de 4-of-6, parts geo-distribuées +10 sites
- [ ] **Formation custodians** : nouveau protocole, rotation custodians
- [ ] **Test cérémonie** : drill reconstruction annuel (simulation, pas de nouvelle clé)
- [ ] **Mise à jour standards** : algorithmes renforcés (P-521 ou post-quantique si mature)

---

*Ce runbook est testé annuellement via simulation tabletop. Jamais de divulgation hors Cellule Crise PKI.*
