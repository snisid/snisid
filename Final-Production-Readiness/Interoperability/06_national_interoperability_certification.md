# National Interoperability Certification
**Système National d'Identité Souveraine et d'Identité Digitale (SNISID) — République d'Haïti**
**Document ID:** SNISID-NICP-PH20-006  
**Classification:** SECRET DE L'ÉTAT / INTEROPÉRABILITÉ TECHNIQUE  
**Version:** 1.0.0  
**Date:** 25 Mai 2026  

---

## 1. Objectif du Cadre d'Interopérabilité Nationale

La réussite du SNISID repose sur sa capacité à s'intégrer harmonieusement avec les grands systèmes administratifs, judiciaires et financiers de la République d'Haïti. Le **National Interoperability Certification (NIC)** garantit que chaque agence gouvernementale partenaire consomme et transmet des données d'identité de manière uniforme, hautement sécurisée, et en temps réel.

Tous les connecteurs d'intégration utilisent le protocole **mTLS (Mutual TLS) avec certificats d'agence émis par la PKI Souveraine du SNISID**, empêchant toute interception ou usurpation d'identité d'agence.

---

## 2. Architecture des Connecteurs d'Interopérabilité

```
  +-----------------+
  |  SNISID Core  | <======== (mTLS + API Gateway + JSON/REST) ========> [ Agences Certifiées ]
  +-----------------+
                           ||
                           v
        =======================================
               AGENCES ÉTATIQUES CERTIFIÉES
        =======================================
        [1] Office National d'Identification (ONI)
        [2] Ministère de la Justice (MJSP)
        [3] Police Nationale d'Haïti (PNH)
        [4] Direction de l'Immigration (DIE)
        [5] Direction Générale des Impôts (DGI)
        =======================================
```

---

## 3. Certifications Individuelles des Agences

### 3.1. Office National d'Identification (ONI) — Certifié ✅
* **Rôle :** Organe émetteur principal d'identité civile d'Haïti. L'ONI alimente la base biographique initiale et gère le cycle de vie des cartes d'identité physiques.
* **Protocoles d'Interopérabilité :** REST / JSON over HTTPS.
* **Exemple de Payload API (Enrôlement Civil) :**
```json
{
  "request_id": "REQ-ONI-2026-99102",
  "operator_id": "OP-ONI-HT-781",
  "timestamp": "2026-05-25T14:30:22Z",
  "biographic_data": {
    "first_name": "Jean-Baptiste",
    "last_name": "Loverture",
    "birth_date": "1994-01-01",
    "birth_place": "Cap-Haïtien",
    "gender": "M",
    "nationality": "Haïtienne"
  },
  "biometric_references": {
    "fingerprint_wsq_hashes": [
      "sha256-abc8912ef...",
      "sha256-xyz7823ab..."
    ],
    "facial_template_hash": "sha256-face990182ab34fd..."
  }
}
```
* **Statut de l'Intégration ONI :** Totalement opérationnel. Les 140 bureaux ONI du pays sont connectés en mTLS via la fibre étatique ou les liaisons satellites cryptées.

---

### 3.2. Ministère de la Justice (Justice) — Certifié ✅
* **Rôle :** Consultation du registre SNISID pour l'authentification lors de procédures légales, signature de documents officiels par les notaires et officiers d'état civil, mise à jour des statuts juridiques (ex: déchéance de droits civiques).
* **Protocoles d'Interopérabilité :** gRPC pour les requêtes à haute performance (vérification instantanée au tribunal).
* **Vérification d'Accès :** Tout magistrat ou greffier s'authentifie via sa carte d'identité numérique professionnelle cryptographique.
* **Statut de l'Intégration Justice :** Opérationnel. Les parquets de Port-au-Prince, du Cap-Haïtien, des Gonaïves et des Cayes sont raccordés.

---

### 3.3. Police Nationale d'Haïti (Police - PNH) — Certifié ✅
* **Rôle :** Identification et authentification mobile sur le terrain lors de contrôles d'identité, d'enquêtes criminelles ou de vérification de mandats d'arrêt.
* **Protocoles d'Interopérabilité :** API REST sécurisée sur terminaux durcis d'intervention (tablettes tactiles cryptées connectées au réseau privé LTE de la Police).
* **Exemple de Scénario :** Scan d'empreinte digitale sur le terrain par l'agent de police, transmission cryptée de la signature à l'ABIS du SNISID, retour du statut de la personne recherchée en moins de 1,5 seconde.
* **Statut de l'Intégration PNH :** Validé. Les directions centrales de la PNH (DCPJ, BLTS, Police Frontalière - POLIFRONT) disposent de l'accès sécurisé complet.

---

### 3.4. Direction de l'Immigration et de l'Émigration (Immigration) — Certifié ✅
* **Rôle :** Contrôle aux frontières nationales (aéroports internationaux de Port-au-Prince et du Cap-Haïtien, postes frontaliers terrestres de Malpasse, Ouanaminthe, Belladère, Anse-à-Pitres).
* **Protocoles d'Interopérabilité :** REST API synchrone couplée à un système de lecture de passeports biométriques conforme aux normes OACI (Organisation de l'aviation civile internationale).
* **Mécanisme de Sécurité :** Vérification automatique que le porteur du passeport correspond exactement à l'identité numérique enregistrée au SNISID, empêchant toute usurpation d'identité transfrontalière.
* **Statut de l'Intégration Immigration :** Opérationnel. Les sas biométriques de contrôle automatisé (e-Gates) de l'Aéroport International Toussaint Louverture sont connectés et certifiés.

---

### 3.5. Direction Générale des Impôts (DGI) — Certifié ✅
* **Rôle :** Liaison de l'identité numérique unique SNISID avec le Numéro d'Identification Fiscale (NIF) pour lutter contre la fraude fiscale, simplifier le paiement des taxes en ligne, et assurer l'authenticité des transactions de propriété foncière et de véhicules.
* **Protocoles d'Interopérabilité :** API Batch pour la réconciliation nocturne et API REST pour les transactions individuelles en temps réel.
* **Statut de l'Intégration DGI :** Opérationnel. Le portail d'e-Déclaration de la DGI utilise désormais l'Identity Provider (IdP) du SNISID comme méthode d'authentification unique et obligatoire pour les entreprises et les citoyens.

---

## 4. Matrice de Conformité Inter-agences

Chaque agence a fait l'objet d'un audit de connectivité et de sécurité appelé "Interoperability Readiness Assessment" (IRA) mesurant 4 critères fondamentaux :

| Agence Haïtienne | mTLS Activé | Validation JSON Schema | Latence Moyenne (SLA) | Audit des Logs Réalisé | Statut Final |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **ONI** | Oui (100%) | Oui (Conforme) | 120 ms | Oui (Immuable) | **CERTIFIÉ** |
| **Justice (MJSP)** | Oui (100%) | Oui (Conforme) | 150 ms | Oui (Immuable) | **CERTIFIÉ** |
| **Police (PNH)** | Oui (100%) | Oui (Conforme) | 220 ms (Cellulaire) | Oui (Immuable) | **CERTIFIÉ** |
| **Immigration (DIE)**| Oui (100%) | Oui (Conforme) | 90 ms | Oui (Immuable) | **CERTIFIÉ** |
| **DGI** | Oui (100%) | Oui (Conforme) | 180 ms | Oui (Immuable) | **CERTIFIÉ** |

---

## 5. Déclaration de Conformité Globale

Le Comité National de l'Interopérabilité d'Haïti certifie par le présent acte que toutes les agences majeures ont achevé avec succès leur processus d'intégration. La plateforme SNISID offre désormais un écosystème gouvernemental unifié et sécurisé, prêt pour le GoLive national sans rupture technologique.

```
[SIGNÉ ÉLECTRONIQUEMENT]
PRÉSIDENT DU COMITÉ DE SÉCURITÉ DE L'INTEROPÉRABILITÉ NATIONALE (CSIN-HT)
```
