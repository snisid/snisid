---
# ============================================================
# SNISID-Data — Master Data Management (MDM)
# Référentiels Maîtres (Golden Records)
# Document ID: SNISID-DATA-MDM-001
# Version: 1.0.0
# ============================================================

## 1. LA PROBLÉMATIQUE DES RÉFÉRENTIELS

Si le Ministère de la Santé dit qu'un citoyen s'appelle "Jean", et la Police dit "Jehan", qui a raison ? Le Master Data Management résout ce problème en créant le "Golden Record".

## 2. DOMAINES MAÎTRES (Master Domains)

Le MDM gère 3 entités critiques pour l'État :
1. **Citoyen / Identité :** La source de vérité est le registre SNISID (Phase 2). Tous les autres ministères doivent s'y conformer.
2. **Adresses / Territoire :** La source de vérité est le Ministère de l'Intérieur (Codes DPC/IHSI des Départements, Communes, Sections Communales).
3. **Véhicules :** La source de vérité est la Direction de la Circulation (OAVCT).

## 3. SYNCHRONISATION MDM

Lorsqu'un Golden Record est mis à jour (Ex: Le SNISID corrige l'orthographe du nom d'un citoyen), le système MDM publie un événement Kafka `MasterEntityUpdated`. Les systèmes isolés (Impôts, Hôpitaux) s'abonnent à ce topic et mettent à jour leurs bases de données locales automatiquement (Event-Driven Architecture).

---
*Document ID: SNISID-DATA-MDM-001 | Approuvé par: Comité National Interministériel*
