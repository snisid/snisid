# SNISID Parallel Operations Model
## Modèle de Coexistence Parallèle et Transition Graduelle

---

## 1. Philosophie de Transition Sans Rupture (Zero-Downtime Transition)

Un arrêt soudain des systèmes existants (legacy) d'identification en Haïti pour basculer sur le SNISID représenterait un risque d'interruption démocratique, économique et administrative catastrophique. Le **Parallel Operations Model** définit la méthodologie opérationnelle qui régit la coexistence pacifique et la synchronisation continue des anciens systèmes avec le SNISID pendant la phase de transition.

```
                         PERIODE DE TRANSITION PARALLELE
                         
[Ancien Système : ONI Legacy]   =====================> [Désarmement Progressif]
                                     \
                                      \ (Dual Capture / Sync Pipeline)
                                       v
[Nouveau Système : SNISID]      =====================> [Plein Régime Opérationnel]
```

---

## 2. Principes Fondamentaux de la Coexistence

1. **La Double Saisie Limitée (Dual Entry Principle) :**
   Les bureaux de liaison pilotes effectuent l'enrôlement sur les terminaux SNISID modernes. Les données sont automatiquement sérialisées et réinjectées en temps réel (ou via micro-lots) dans les bases de données ONI historiques par le biais d'un connecteur d'écriture inverse (Reverse Connector), évitant aux agents d'avoir à saisir deux fois les mêmes données.
2. **Le Principe d'Asymétrie de Confiance (Trust Asymmetry) :**
   En cas de divergence entre l'ancienne base et la base SNISID, c'est la donnée SNISID qui fait foi (car elle a subi le nettoyage de la *Migration Factory* et la validation de l'ABIS biométrique de dernière génération).
3. **Le Maintien du Mode Secours (Fallback Availability) :**
   Pendant toute la période de transition (fixée à 6 mois par département), l'ancienne infrastructure réseau et matérielle reste en veille active (Standby) et peut être réactivée sous 2 heures si une faille critique non détectée du SNISID survenait.

---

## 3. Stratégie de Bascule Contrôlée (Controlled Cutover Framework)

La bascule du mode "Coexistence parallèle" au mode "SNISID Exclusif" se fait par paliers d'activités administratives.

```
       [Phase 1 : Enrôlement Dual]
                  | (Durée : 4 semaines. Validation taux d'erreur < 0.05%)
                  v
       [Phase 2 : Consultation Active (Lecture)]
                  | (Durée : 4 semaines. API SNISID ouverte aux banques/administrations)
                  v
       [Phase 3 : Émission Exclusive (Production)]
                  | (L'ancien système ne produit plus de documents physiques)
                  v
       [Phase 4 : Extinction du Legacy]
                  | (Archivage sécurisé à froid et mise hors service des vieux serveurs)
```

### 3.1 Les Quatre Phases de la Bascule Départementale

| Étape de Bascule | Activités Clés | Critères de Succès pour Passer à l'Étape Suivante |
| :--- | :--- | :--- |
| **Phase 1 : Enrôlement Dual** | Saisie sur SNISID, écriture automatique vers ONI Legacy. | 100% des dossiers synchronisés avec succès. Taux de rejet biométrique inférieur à 1%. |
| **Phase 2 : Consultation Active** | Les banques commerciales (BRH, Sogebank, Unibank) et les ministères lisent les identités uniquement via les API sécurisées du SNISID. | Plus de 50 000 requêtes d'API d'identification par jour traitées avec un temps de réponse < 200 ms. |
| **Phase 3 : Émission Exclusive** | Seul le SNISID émet la nouvelle Carte Nationale d'Identification Biométrique Unique (CNIBU). Les anciennes cartes restent valides mais aucune nouvelle n'est produite sur l'ancienne infrastructure. | Impression et remise physique de 98% des cartes demandées dans un délai de 5 jours ouvrés. |
| **Phase 4 : Archivage & Extinction** | Les vieux serveurs ONI régionaux sont arrêtés. Les bases de données historiques sont sauvegardées sur bande magnétique chiffrée et placées sous clé scellée aux Archives Nationales. | Signature du procès-verbal de transfert de responsabilité souveraine. |

---

## 4. Périodes de Validation et Vérification de Cohérence

À la fin de chaque journée opérationnelle en phase parallèle, un script automatisé de rapprochement (*Reconciliation Pipeline*) vérifie que le nombre d'inscriptions comptabilisées dans le système historique correspond exactement aux dossiers validés dans le SNISID.

```
                     RAPPROCHEMENT DE COEXISTENCE
                     
+------------------------+                     +------------------------+
|   SNISID Daily Logs    |                     |    ONI Legacy Logs     |
+------------------------+                     +------------------------+
            \                                              /
             \                                            /
              v                                          v
       +--------------------------------------------------------+
       |             Reconciliation Checker Script              |
       |  (Vérification d'intégrité, de volume et d'identité)  |
       +--------------------------------------------------------+
                                   |
                +------------------+------------------+
                |                                     |
         [Match : Validé]                   [Mismatch : ALERTE]
         - Clôture de la journée            - Blocage de la bascule
         - Synchro au DC principal          - Rapport d'écart généré
```

En cas d'écart (ex: un dossier présent sur le système historique mais rejeté ou manquant sur le SNISID), la cellule de crise suspend immédiatement la transition pour la commune affectée jusqu'à la résolution du conflit d'identité par l'équipe technique sous 24 heures.
