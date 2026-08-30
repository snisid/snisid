---
# ============================================================
# SNISID-Infra — Primary Sovereign Datacenter
# Topologie Physique et Sécurité (Tier III+)
# Document ID: SNISID-DC-PRI-001
# Version: 1.0.0
# ============================================================

## 1. ARCHITECTURE TIER III+

Le Datacenter Primaire (Port-au-Prince) héberge le `SNISID-Core`, l'Identity Registry, et le système central PNH. Il doit garantir une disponibilité de 99.982% (soit un maximum de 1.6 heure d'arrêt par an).

### 1.1 Redondance Physique (N+1 minimum)
- **Refroidissement (Cooling)** : Systèmes HVAC N+1.
- **Énergie** : Double alimentation vers chaque baie de serveurs (A/B Feeds) avec onduleurs séparés.
- **Réseau** : Deux entrées fibres optiques distinctes (Point of Presence Natcom / Digicel / FAI alternatif) routées physiquement sur des chemins séparés.

## 2. AIR-GAPPED SECURITY ZONES (Ségrégation Physique)

Le Datacenter est divisé en plusieurs zones de sécurité croissantes (Modèle en Oignon).

```mermaid
graph TD
    subgraph "Zone 1: DMZ (Public Facing)"
        FW1[Firewalls Périmétriques]
        WAF[Web Application Firewalls]
        KONG[API Gateway / Ingress]
    end

    subgraph "Zone 2: Government Internal"
        K8S_APPS[Kubernetes Worker Nodes\n(Services PNH, Justice, etc.)]
        ISTIO[Istio Control Plane]
    end

    subgraph "Zone 3: Secure Enclave (Air-Gapped)"
        VAULT[HashiCorp Vault / HSM]
        COCKROACH[CockroachDB Identity Nodes]
        CA[Root Certificate Authority]
    end

    Zone1 -- "mTLS Traffic Only" --> Zone2
    Zone2 -- "Strict RBAC/OPA Queries" --> Zone3
```

**Règle "Air-Gap" (Zone 3) :**
Les serveurs hébergeant la base de données d'identité (CockroachDB) n'ont aucune route réseau vers l'Internet public. Ils ne peuvent être contactés QUE par les microservices de la Zone 2, via des ports spécifiques et chiffrés.
Les serveurs HSM (Hardware Security Module) contenant la racine PKI (Root CA) sont littéralement déconnectés du réseau et nécessitent une intervention physique en binôme (Two-man rule) pour la signature de certificats intermédiaires.

## 3. RACK ARCHITECTURE (Kubernetes Bare-Metal)

Les baies (Racks) utilisent le concept de "Failure Domain" (Domaine de panne).
- Un cluster Kubernetes (RKE2) de 12 serveurs est réparti sur 3 racks différents (4 serveurs par rack).
- Si un Rack entier prend feu (Perte du switch Top-of-Rack ou de la PDU), le cluster Kubernetes perd 1/3 de sa capacité, mais reste en ligne (High Availability).

---
*Document ID: SNISID-DC-PRI-001 | Approuvé par: Directeur Infrastructure (AND)*
