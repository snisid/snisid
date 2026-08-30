# SNISID Runbook — Connectivity Failure Offline Fallback Procedure
**Code de Procédure :** SNISID-RB-04  
**Statut :** Approuvé  
**Audience :** Opérateurs de Bureau de Liaison, Équipe Réseau SRE  

---

## 1. Objectif

Ce runbook régit le comportement opérationnel des Bureaux de Liaison Communaux (BLC) et des équipes réseau en cas de perte de connectivité complète (coupure fibre optique, panne d'antenne satellite Starlink ou absence de signal 4G) sur un site de déploiement.

---

## 2. Transition en Mode Hors-Ligne (Offline Fallback)

Le basculement en mode d'enrôlement et vérification hors-ligne est **automatique** au niveau logiciel, mais nécessite une procédure d'assurance qualité humaine sur le terrain.

```
                      CONNECTIVITY FAILURE WORKFLOW
                      
  [Loss of Connectivity Detected by LEN]
                    |
                    v (Automatic local routing)
  [Local Edge Node starts operating in Offline-First mode]
  - Alerts operator on workstation screen
  - Blocks central ABIS 1-to-N real-time check
  - Fallbacks to local DB check and generates temporary IUI-T
                    |
                    +-------------------+-------------------+
                                        |
                                        v (Manual action)
  [Operator proceeds with Offline Enrollment & Local Signing]
  - Store encrypted transactions on local SSD disk LUKS
  - Safe-guard physical keys
```

---

## 3. Étapes de Remédiation pour l'Opérateur Terrain

1. **Vérifier l'Alerte Écran :** S'assurer que le message *"Mode Hors-ligne Activé"* s'affiche bien en haut de l'interface d'enrôlement SNISID.
2. **Continuer les Enrôlements Civils :** Ne pas renvoyer les citoyens chez eux. Expliquer calmement que le système fonctionne en mode d'enregistrement autonome sécurisé.
3. **Capture Biométrique Renforcée :** En mode hors-ligne, le moteur de déduplication central (ABIS 1-to-N) n'est pas consultable en temps réel. L'opérateur doit être d'autant plus rigoureux lors de la capture des empreintes digitales et du visage (respect strict des critères de qualité OACI) pour minimiser les rejets ultérieurs lors de la synchronisation centrale.
4. **Vérification d'Alimentation Électrique :** S'assurer que le Local Edge Node est correctement branché sur le circuit d'alimentation solaire/batterie pour garantir son fonctionnement ininterrompu.

---

## 4. Rétablissement et Synchronisation Différée (Delayed Sync)

Une fois la connectivité rétablie (le voyant réseau du routeur redevient vert stable) :
1. **Déclenchement Automatique :** Le Local Edge Node détecte le retour du tunnel gRPC TLS chiffré vers le Datacenter Central.
2. **Exécution du Protocole de Synchronisation :** Le LEN lance la vidange sécurisée de sa file d'attente locale transactionnelle vers le NIRE central.
3. **Contrôle du Rapport de Synchronisation :** L'opérateur doit vérifier la notification de synthèse envoyée sur son poste, attestant du succès d'intégration et du taux de dossiers acceptés (ex: `Promoted: 120, Quarantined: 1` - voir simulation `offline_sync_protocol.py`).
4. **Traitement des Quarantaines :** Si un dossier est placé en quarantaine centrale (suspicion de doublon détectée par l'ABIS central), l'opérateur doit planifier l'ajournement du dossier selon les directives de la cellule de crise.
