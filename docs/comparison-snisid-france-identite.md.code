# Comparaison SNISID vs France Identité

> Analyse comparative : **France Identité** (ANTS/DINUM, certifié eIDAS niveau élevé, 2M+ utilisateurs, mai 2025)
> vs **SNISID** (Système National d'Identification Sécurisée et Intégrée d'Haïti, architecture en cours de finalisation)
> Date : Juin 2026

---

## Résumé exécutif

| Dimension | France Identité | SNISID |
|-----------|----------------|--------|
| Score global | **87 / 100** | **54 / 100** |
| Lacunes critiques | 0 | **7** |
| Lacunes importantes | 0 | **9** |
| Points à améliorer | 0 | **8** |
| Avantages SNISID | — | **5 domaines** |

**Conclusion :** France Identité est un système opérationnel mature, certifié et interopérable à l'échelle européenne. SNISID a une architecture technique plus ambitieuse (fraude IA, biométrie NPU, fédération multi-pays, gouvernance formelle) mais manque des couches standards de confiance, conformité légale et interopérabilité.

---

## Priorités

| Symbole | Priorité | Description |
|---------|----------|-------------|
| `[CRITIQUE]` | Bloquant | Empêche la reconnaissance légale ou l'interopérabilité internationale |
| `[IMPORTANT]` | Élevé | Limite sérieusement les usages ou la sécurité |
| `[MOYEN]` | Moyen | À intégrer avant mise en production complète |

---

## Domaine 1 — Confiance & Certification légale

**France Identité : 95/100 | SNISID : 20/100**

### LAC-01 `[CRITIQUE]` Aucune certification de sécurité

- France Identité : CSPN ANSSI, eIDAS niveau élevé, FranceConnect+ qualifié ANSSI
- SNISID : Aucun processus de certification formelle (ni ISO 27001, ni SOC 2, ni équivalent)

### LAC-02 `[CRITIQUE]` Pas de qualification autorité protection des données

- France Identité : Supervisé CNIL, RGPD by design, PIA publié, DPO désigné
- SNISID : Aucun DPO, pas de PIA biométrique, `legal-oversight/validator.go` est un stub

### LAC-03 `[CRITIQUE]` Aucun décret ou base légale documentée

- France Identité : Décret n°2022-676 SGIN, règlement eIDAS (UE) n°910/2014
- SNISID : Aucune référence à une loi haïtienne autorisant le traitement biométrique national

---

## Domaine 2 — Cryptographie & Standards ouverts

**France Identité : 92/100 | SNISID : 45/100**

### LAC-04 `[CRITIQUE]` Pas de support Verifiable Credentials W3C (VC 2.0)

- France Identité : VCs W3C v2.0, OID4VP, OID4VCI, ARF v2.0 EUDI Wallet
- SNISID : Aucune couche VC W3C

### LAC-05 `[CRITIQUE]` Divulgation sélective absente

- France Identité : SD-JWT, preuves ZKP (prouver âge sans date exacte)
- SNISID : Toutes les données en clair dans le JWT

### LAC-06 `[IMPORTANT]` Pas d'infrastructure PKI nationale

- France Identité : CA racine d'État ANSSI, certificats X.509 qualifiés
- SNISID : JWT RS256 auto-signé uniquement, pas de CRL/OCSP

### LAC-07 `[IMPORTANT]` Pas de signature électronique qualifiée (QES)

- France Identité : QES eIDAS reconnue dans toute l'UE
- SNISID : Aucun service de signature qualifiée

---

## Domaine 3 — Hardware & NFC

**France Identité : 90/100 | SNISID : 30/100**

### LAC-08 `[CRITIQUE]` Pas de support NFC / lecture de puce physique

- France Identité : Lecture NFC CNIe, ICAO 9303, BAC/PACE/EAC
- SNISID : Aucune couche NFC

### LAC-09 `[MOYEN]` Pas de stockage clés dans l'élément sécurisé

- France Identité : Secure Enclave iOS, StrongBox Android
- SNISID : Clés en mémoire / SecureStore basique

---

## Domaine 4 — Authentification multi-niveaux

**France Identité : 88/100 | SNISID : 60/100**

### LAC-10 `[IMPORTANT]` Pas d'authentification graduée eIDAS

- France Identité : 3 niveaux eIDAS (faible/substantiel/élevé), flux différents
- SNISID : JWT RS256 uniforme pour tous les cas

### LAC-11 `[IMPORTANT]` Pas de support FIDO2 / WebAuthn

- France Identité : TOTP, POP, carte agent PIN + X.509
- SNISID : Pas de FIDO2, pas de clé physique pour admins

---

## Domaine 5 — Fédération d'identité

**France Identité : 95/100 | SNISID : 25/100**

### LAC-12 `[CRITIQUE]` Pas de fournisseur d'identité OpenID Connect

- France Identité : FranceConnect, 1400+ services, 7 fournisseurs d'identité
- SNISID : Pas d'endpoints OIDC (`.well-known`, `/authorize`, `/token`, `/userinfo`)

---

## Domaine 6 — Portefeuille numérique

**France Identité : 85/100 | SNISID : 30/100**

### LAC-13 `[IMPORTANT]` Pas de wallet multi-documents

- France Identité : CNIe + permis + carte vitale, 2M utilisateurs, EUDI Wallet
- SNISID : App mobile ne gère que l'identité nationale

### LAC-14 `[IMPORTANT]` Pas de justificatif d'identité à usage unique

- France Identité : QR code éphémère, one-time use, attributs minimaux
- SNISID : QR code avec données brutes

---

## Domaine 7 — Interopérabilité internationale

**France Identité : 90/100 | SNISID : 40/100**

### LAC-15 `[IMPORTANT]` Pas de protocole ICAO / EES

- France Identité : EES Schengen, ICAO 9303 NFC, Interpol
- SNISID : `federation-gateway` ne parle pas ICAO

### LAC-16 `[MOYEN]` Pas de support ISO 18013-5 (mDL)

- France Identité : mDL permis de conduire mobile, consortium APTITUDE
- SNISID : Pas de couche mdoc

---

## Domaine 8 — Gouvernance & Audit

**France Identité : 92/100 | SNISID : 50/100**

### LAC-17 `[IMPORTANT]` Hash chain d'audit non certifiée

- France Identité : Horodatage qualifié RFC 3161, journaux opposables
- SNISID : Hash chain correcte mais non qualifiée

### LAC-18 `[MOYEN]` Pas de comité de pilotage inter-agences

- France Identité : DINUM + ANTS + Ministère Intérieur
- SNISID : Pas de structure documentée

### LAC-19 `[MOYEN]` Pas de processus de révocation officiel

- France Identité : France Titres, révocation instantanée propagée
- SNISID : États revoked/suspended non connectés à un registre

---

## Domaine 9 — UX & Accessibilité

**France Identité : 80/100 | SNISID : 45/100**

### LAC-20 `[MOYEN]` App mobile sans mode dégradé

- France Identité : Mode offline partiel, zones faible couverture
- SNISID : Pas de cache local identité

### LAC-21 `[MOYEN]` Pas de gestion des mineurs

- France Identité : Cadre 18 ans+, représentants légaux
- SNISID : Pas de modèle mineur / tuteur

---

## Domaine 10 — Monitoring & Partenaires

**France Identité : 85/100 | SNISID : 40/100**

### LAC-22 `[MOYEN]` Pas de SLA public

- France Identité : Status page publique, SLA 99.9%
- SNISID : Aucune visibilité partenaires

### LAC-23 `[MOYEN]` Pas de sandbox développeurs

- France Identité : Playground, documentation partenaires
- SNISID : Pas d'environnement de test isolé

### LAC-24 `[MOYEN]` Pas de processus d'onboarding partenaires

- France Identité : Habilitation structurée, bac à sable → production
- SNISID : Pas de processus documenté

---

## Avantages SNISID

| # | Domaine | Détail |
|---|---------|--------|
| 1 | Détection de fraude IA | Causal inference (Do-calculus), behavioral profiling, Neo4j, RL agent — France Identité n'a que du liveness detection basique |
| 2 | Biométrie offline | Runtime NPU (ONNX), matching 1:N FAISS offline, galeries chiffrées — adapté aux zones sans connectivité |
| 3 | Architecture multi-pays native | Federation gateway avec mTLS, échange cross-country — plus avancé que France Identité (national) |
| 4 | Gouvernance formelle | DSL compiler → Rego/OPA, TLA+ formal verification, digital twin économique — aucun équivalent France Identité |
| 5 | Résilience opérationnelle | Auto-healing K8s, snapshot/restore, critical-runtime monitor — infrastructure plus robuste |

---

## Plan de priorité

### Phase 1 — Fondations légales (3-6 mois)

| Priorité | Lacune | Action |
|----------|--------|--------|
| 1 | LAC-03 | Faire adopter un décret d'application haïtien |
| 2 | LAC-02 | Désigner un DPO, réaliser les PIA |
| 3 | LAC-06 | Créer la PKI nationale (CA racine) |
| 4 | LAC-01 | Lancer certification ISO 27001 |
| 5 | LAC-12 | Implémenter IdP OpenID Connect minimal |

### Phase 2 — Standards cryptographiques (3-4 mois)

| Priorité | Lacune | Action |
|----------|--------|--------|
| 6 | LAC-04 | Verifiable Credentials W3C + OID4VCI/OID4VP |
| 7 | LAC-05 | SD-JWT + divulgation sélective |
| 8 | LAC-08 | Support NFC ICAO 9303 |
| 9 | LAC-07 | Signature électronique qualifiée |
| 10 | LAC-10 | Niveaux d'assurance eIDAS |

### Phase 3 — Fonctionnalités utilisateur (2-3 mois)

| Priorité | Lacune | Action |
|----------|--------|--------|
| 11 | LAC-13 | Wallet multi-documents |
| 12 | LAC-14 | QR éphémère à usage unique |
| 13 | LAC-09 | Secure Enclave / Android Keystore |
| 14 | LAC-11 | FIDO2 / WebAuthn |
| 15 | LAC-20 | Mode offline app mobile |

### Phase 4 — Interopérabilité (2-3 mois)

| Priorité | Lacune | Action |
|----------|--------|--------|
| 16 | LAC-15 | Protocole ICAO / contrôle frontières |
| 17 | LAC-17 | Horodatage qualifié RFC 3161 |
| 18 | LAC-19 | Registre de révocation officiel |
| 19 | LAC-22 | Status page + SLA public |
| 20 | LAC-23 | Sandbox et portail développeurs |
| 21 | LAC-16 | ISO 18013-5 (mDL) |
| 22 | LAC-18 | Comité de pilotage inter-agences |
| 23 | LAC-21 | Gestion mineurs / tuteurs |
| 24 | LAC-24 | Processus onboarding partenaires |

---

## Références

- France Identité : https://france-identite.fr
- FranceConnect : https://franceconnect.gouv.fr
- eIDAS 2 : Règlement (UE) 2024/1183
- EUDI Wallet ARF v2.0 : https://github.com/eu-digital-identity-wallet/architecture-and-reference-framework
- ICAO 9303 : https://www.icao.int/publications/pages/publication.aspx?docnum=9303
- ISO 18013-5 : https://www.iso.org/standard/82772.html
