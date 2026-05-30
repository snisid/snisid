---
# ============================================================
# SNISID-Edge — Edge Observability & Audit
# Télémétrie Régionale et Logs Immuables Offline
# Document ID: SNISID-EDGE-AUDIT-001
# Version: 1.0.0
# ============================================================

## 1. OBSERVABILITÉ DÉCONNECTÉE

Gérer un serveur que l'on ne peut pas pinger est complexe. La plateforme d'observabilité Edge (K3s) résout ce problème en inversant le paradigme : c'est le noeud qui pousse (Push) ses métriques quand il le peut.

## 2. DISCONNECTED TRACEABILITY (Audit Immuable Local)

Si un agent de police effectue une recherche illégale dans la base locale alors que le commissariat est hors-ligne :
1. L'action est enregistrée dans le journal d'audit local.
2. Pour empêcher l'agent de supprimer ses traces (Tampering), le journal utilise une structure en **chaîne de blocs légère (Hash Chain)**. Chaque log inclut le hash du log précédent.
3. Si un log est supprimé, la chaîne est brisée.
4. Au retour de la connexion, le SOC central valide la chaîne cryptographique complète. Si elle est brisée, le commissariat entier est isolé (Incident Response).

## 3. REGIONAL DASHBOARDS

Le Préfet ou le Chef de la Police Départementale possède un accès (Grafana) limité à son département, lui permettant de voir l'état des noeuds Edge de ses commissariats sans avoir besoin d'interroger la capitale.

---
*Document ID: SNISID-EDGE-AUDIT-001 | Approuvé par: Directeur de l'Audit*
