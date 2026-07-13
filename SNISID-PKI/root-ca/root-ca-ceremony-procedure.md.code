# SNISID — Root CA Key Ceremony Procedure
**Classification:** TOP-SECRET  
**Version:** 5.0.0  
**Statut:** Procédure obligatoire — Toute déviation = incident national

---

## Contexte

La Root CA SNISID est la racine absolue de confiance de l'État. Sa compromission = compromission de l'identité nationale entière. Cette procédure garantit que la clé privée Root CA :

1. **N'a jamais existé en dehors du HSM**
2. **N'a été générée que sous supervision multi-parties**
3. **Est révocable uniquement via procédure nationale d'urgence**

---

## Participants obligatoires (6 Key Custodians)

| # | Rôle | Nom (completer jour J) | Fonction | Justification |
|---|------|----------------------|----------|---------------|
| 1 | **Custodian A** | _________________ | Présidence / Haute autorité | Détenteur souverain |
| 2 | **Custodian B** | _________________ | IGC — Directeur Cyber | Maître PKI |
| 3 | **Custodian C** | _________________ | IGC — Chef Audit PKI | Audit continu |
| 4 | **Custodian D** | _________________ | Justice — Représentant légal | Validité juridique |
| 5 | **Custodian E** | _________________ | Technique — Architecte national | Expert technique |
| 6 | **Custodian F** | _________________ | Technique — Ops HSM national | Opérations HSM |

**Exigences :** Biométrie enregistrée HSM. Smartcard personnelle. Aucun conflit d'intérêts familial/professionnel.

---

## Phase 1 — Pré-cérémonie (J-7 jours)

### 1.1 Vérification salle blanche
- [ ] Accès biométrique validé pour les 6 custodians
- [ ] Caméras 24/7 fonctionnelles (test enregistrement + conservation)
- [ ] Sceaux tamper-evident sur rack HSM (photos signées)
- [ ] Générateur diesel testé (charge 100% pendant 4h)
- [ ] UPS testé (switchover < 10ms)
- [ ] Aucun équipement électronique non-autorisé (scan RF, bagages)
- [ ] Aéroportation positive vérifiée (pas de poussière critique)

### 1.2 Vérification HSM
- [ ] Firmware Thales Luna 7 version attestée (hash SHA-256 publié)
- [ ] Attestation HSM Thales validée (certificat d'authenticité)
- [ ] HSM zéroisé (zeroization) + testé
- [ ] Clé partition admin HSM créée avec M-of-N = 4-of-6
- [ ] Test : créer clé test, signer, vérifier, détruire
- [ ] Backup HSM DR identifié (sérial numéro : __________)

### 1.3 Matériel cérémonie
- [ ] 6 laptops air-gapped (OS durci, pas de WiFi/BT/Ethernet)
- [ ] 6 smartcard readers USB isolés (pas de hub)
- [ ] 6 clés USB optique (write-once vérifié)
- [ ] Journal papier continu numéroté (tampon IGC)
- [ ] 2 imprimantes thermiques (air-gapped)
- [ ] Kit destruction d'urgence (thermite / broyeur NSA-approved)

### 1.4 Shamir configuration
- [ ] Threshold : **4** (minimum pour reconstruire)
- [ ] Parts : **6** (1 par custodian)
- [ ] Génération vérifiée mathématiquement (libre N-of-M vérifié)
- [ ] Aucun custodian ne connaît d'autre custodian la génération a eu lieu

---

## Phase 2 — Cérémonie jour J (H-2h)

### 2.1 Entrée séquentielle (06:00)
06:00 — Custodian A entre, biométrie, scan complet
06:15 — Custodian B entre, idem
06:30 — Custodian C entre, idem
06:45 — Custodian D entre, idem
07:00 — Custodian E entre, idem
07:15 — Custodian F entre, idem

**Si un custodian est retard > 30min :** Annulation, report J+7.
**Si un custodian est indisponible :** Report J+7, remplacement selon liste alternates (3 noms pré-approuvés).

### 2.2 Vérification mutuelle (07:30)
- [ ] Chaque custodian présente smartcard + biométrie à 2 autres custodians
- [ ] Vérification d'identité papier + digitale (cross-check)
- [ ] Tous signent le registre d'entrée physique

### 2.3 Inspection HSM (08:00)
- [ ] Photos comparées sceaux rack (J-7 vs J)
- [ ] Vérification tamper-evident HSM (led status, étiquettes)
- [ ] Si tamper-evident altéré : **CÉRÉMONIE ANNULÉE** — incident national, investigation IGC
- [ ] Custodian E vérifie firmware hash (lecture HSM console)
- [ ] Custodian F vérifie partition admin avec M-of-N = présents

### 2.4 Génération clé Root CA (09:00)
```bash
# Commandes exécutées sur laptop air-gapped H-1 (préparé par Custodian E)
# JAMAIS connecté au réseau. JAMAIS de clé exportée.

# 1. Authentification M-of-N sur HSM partition admin
lunash:> hsm login -partition SNISID_ROOT -role admin
# > Custodian A : smartcard + PIN
# > Custodian B : smartcard + PIN
# > Custodian C : smartcard + PIN
# > Custodian D : smartcard + PIN
# (Threshold 4 atteint)

# 2. Création partition Root CA (dédiée, jamais utilisée autrement)
lunash:> partition create -partition SNISID_ROOT_CA_2026
lunash:> partition policy -partition SNISID_ROOT_CA_2026 -min_pwd_len 16 -max_login_fail 3 -lockout_duration 86400

# 3. Génération clé ECDSA P-384 (génération DANS le HSM, jamais export)
lunash:> key generate -partition SNISID_ROOT_CA_2026 -algorithm EC -curve secp384r1 -label "snisid-root-ca-2026-key" -extractable no -modifiable no -size 384

# 4. Vérification : la clé n'est PAS extractable
lunash:> key show -partition SNISID_ROOT_CA_2026 -label "snisid-root-ca-2026-key"
# Attendu : Extractable = NO, Modifiable = NO, Class = PRIVATE_KEY

# 5. Test signature (data "SNISID ROOT CA 2026 CEREMONY TEST")
lunash:> sign -partition SNISID_ROOT_CA_2026 -label "snisid-root-ca-2026-key" -data "SNISID ROOT CA 2026 CEREMONY TEST" -mechanism ECDSA_SHA384 -out /tmp/test.sig
# Vérification publique avec pubkey exportable
lunash:> key export -partition SNISID_ROOT_CA_2026 -label "snisid-root-ca-2026-key" -pubkey -out /tmp/root-pub.pem
openssl dgst -sha384 -verify /tmp/root-pub.pem -signature /tmp/test.sig <(echo "SNISID ROOT CA 2026 CEREMONY TEST")
# DOIT retourner : Verified OK
```

### 2.5 Création certificat auto-signé Root CA (10:30)
```bash
# CSR généré dans HSM (clé privée jamais sortie)
lunash:> csr generate -partition SNISID_ROOT_CA_2026 -label "snisid-root-ca-2026-key" -subject "CN=SNISID Root CA Nationale,O=État,C=HT" -out /tmp/root-ca.csr

# Signature du certificat par la clé Root CA elle-même (auto-sign)
lunash:> cert sign -partition SNISID_ROOT_CA_2026 -label "snisid-root-ca-2026-key" -csr /tmp/root-ca.csr -out /tmp/root-ca.crt -serial 2026052500000001 -days 3650 -sha384

# Vérification openssl
openssl x509 -in /tmp/root-ca.crt -text -noout
# Attendu : Subject == Issuer, KeyUsage = keyCertSign, cRLSign, BasicConstraints CA:TRUE pathlen:3
```

### 2.6 Vérifications croisées (11:30)
- [ ] Custodian A vérifie empreinte SHA-256 certificat
- [ ] Custodian B vérifie validité ASN.1 / OID / policies
- [ ] Custodian C vérifie empreinte SHA-256 clé publique
- [ ] Custodian D vérifie conformité cadre légal (format, O, C)
- [ ] Custodian E vérifie robustesse cryptographique (P-384, SHA-384, pas d'extensions dangereuses)
- [ ] Custodian F vérifie HSM state (key non-extractable, log génération)

### 2.7 Publication empreinte nationale (12:00)
- [ ] Hash SHA-256 du certificat Root CA publié dans registre national officiel (multi-sites)
- [ ] Hash SHA-256 publié dans journal officiel papier (date, numéro)
- [ ] Hash SHA-256 broadcast via SMS d'urgence présidentiel (redondance)
- [ ] Hash SHA-256 inscrit dans blockchain nationale d'intégrité (si existante)

---

## Phase 3 — Shamir Secret Sharing (14:00)

### 3.1 Principe
La clé Root CA **ne sort jamais du HSM**. Néanmoins, en cas de destruction physique du HSM (catastrophe, incendie, attaque), une reconstruction est nécessaire.

Le HSM génère une **key backup blob** chiffrée qui est splitée via Shamir Secret Sharing (SSS) en 6 parts, threshold 4.

### 3.2 Procédure
```bash
# 1. Génération backup blob chiffré (clé de backup wrapping key, jamais clé privée en clair)
lunash:> backup -partition SNISID_ROOT_CA_2026 -label "snisid-root-ca-2026-key" -out /tmp/root-ca-backup.blob
# Ce blob est chiffré par une clé de backup stockée dans HSM admin partition

# 2. Export blob sur 6 supports chiffrés individuellement
# Chaque custodian génère une passphrase personnelle (20+ caractères, diceware FR)
# Chaque support est chiffré avec la passphrase du custodian + clé publique IGC

# 3. Shamir SSS sur la passphrase de déchiffrement du backup blob
# Utilisation libre : ssss-split -t 4 -n 6
# Les 6 parts sont écrites sur papier, scellées dans enveloppes inviolables
```

### 3.3 Distribution parts
| Partie | Support | Localisation | Accès |
|--------|---------|--------------|-------|
| Part 1 | Papier scellé + USB optique | Coffre présidence (site A) | Custodian A seul |
| Part 2 | Papier scellé + USB optique | Coffre IGC siège (site B) | Custodian B seul |
| Part 3 | Papier scellé + USB optique | Banque centrale / institution monétaire (site C) | Custodian C seul |
| Part 4 | Papier scellé + USB optique | Ambassade A (étranger, pays allié) | Custodian D seul |
| Part 5 | Papier scellé + USB optique | Ambassade B (étranger, pays allié) | Custodian E seul |
| Part 6 | Papier scellé + USB optique | Base militaire sécurisée (site D intérieur) | Custodian F seul |

---

## Phase 4 — Clôture (16:00)

### 4.1 Sécurisation
- [ ] HSM partition Root CA verrouillée (auto-lock 30s d'inactivité)
- [ ] Clés admin HSM détruites (zeroized) — recréation nécessite nouvelle cérémonie
- [ ] Laptops air-gapped zéroisés (shred -n 35 -z disques SSD)
- [ ] Impressions thermiques scellées dans sac sécurisé IGC
- [ ] Vidéos 24/7 extraites sur 3 supports, scellés, dispersés

### 4.2 Signatures registre
- [ ] Chaque custodian signe registre papier (signature manuscrite + empreinte digitale)
- [ ] Registre scanné, hash SHA-256 publié
- [ ] Registre original scellé IGC (conservation 25 ans)

### 4.3 Notifications
- [ ] Rapport présidence (TOP-SECRET, 48h)
- [ ] Rapport IGC (archivé SIEM national)
- [ ] Rapport audit externe (si mandaté)

---

## Phase 5 — Post-cérémonie (J+30j)

- [ ] Audit externe complet (cabinet souverain ou alliance technique)
- [ ] Publication trust anchor dans tous les systèmes SNISID
- [ ] Distribution bundle Root CA (certificat uniquement, pas clé) aux edge nodes
- [ ] Tests cross-validation : signature Intermediate CA → Root CA OK
- [ ] Procédure révocation Root CA rédigée, approuvée, archivée TOP-SECRET

---

*Procédure validée par IGC. Toute violation = incident de sécurité nationale.*
