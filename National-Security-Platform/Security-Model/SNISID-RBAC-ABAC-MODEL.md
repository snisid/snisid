---
# ============================================================
# SNISID-Security — National Role-Based Security Model
# Modèle d'habilitations ABAC/RBAC
# Document ID: SNISID-RBAC-ABAC-001
# Version: 1.0.0
# ============================================================

## 1. CONCEPT : ZERO TRUST & STRICT SEGREGATION

Le modèle de sécurité SNISID empêche tout abus de pouvoir par une ségrégation stricte. Un policier ne peut pas voir le dossier médical ou fiscal d'un citoyen, et un juge ne peut pas ordonner une incarcération sans motif légal.

## 2. MATRICE RBAC (Role-Based Access Control)

| Rôle | Application (Phase) | Permissions Spécifiques | Restrictions |
|------|---------------------|--------------------------|--------------|
| **Agent PNH** | Police Ops (Ph3) | Vérifier Statut NIU, Enregistrer Main Courante | Aucun accès aux cas judiciaires |
| **DCPJ Analyst** | Intel Platform (Ph3) | Recherche Graphe Criminel, Créer Alerte | Accès limité aux mandats |
| **Juge Instruction** | Judicial Engine (Ph3) | Émettre Mandat, Lire Dossier Assigné | Ne voit que les dossiers de sa juridiction |
| **Officier État Civil** | Civil Registry (Ph2) | Enregistrer Naissance/Décès | Aucun accès criminel |
| **Border Agent** | Immigration (Ph3) | Scan passeport, Vérification Watchlist | Résultat Binaire (Pass/Block) uniquement |
| **System Admin (AND)** | Platform (Ph1) | Configurer Kubernetes/Réseau | **0 accès aux données (Chiffrées au repos)** |

## 3. ABAC (Attribute-Based Access Control) via OPA

Exemple de politique OPA : Un policier (DDO) tente de lire une main courante (DDO).
```rego
package snisid.justice.access

default allow = false

allow {
    input.user.role == "AGENT_PNH"
    input.resource.type == "INCIDENT_REPORT"
    input.user.jurisdiction == input.resource.jurisdiction
}

# La DCPJ peut lire tous les rapports de la PNH
allow {
    input.user.role == "ANALYST_DCPJ"
    input.resource.type == "INCIDENT_REPORT"
}
```

---
*Document ID: SNISID-RBAC-ABAC-001 | Approuvé par: CISO National*
