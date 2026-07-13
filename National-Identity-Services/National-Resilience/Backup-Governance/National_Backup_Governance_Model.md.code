# National Backup Governance Model

## 1. Objectif
Industrialiser les sauvegardes nationales SNISID afin que les données critiques restent récupérables, intègres, chiffrées, auditables et testées.

> Une sauvegarde non testée n'est pas une sauvegarde.

## 2. Capacités
| Domaine | Support | Exigence |
|---|---:|---|
| Automated backups | Oui | planification centralisée et supervision |
| Air-gapped backups | Oui | copies déconnectées contre ransomware/catastrophe |
| Encrypted backups | Oui | chiffrement au repos et en transit |
| Backup audits | Oui | journaux, rapports, conformité |
| Recovery validation | Oui | restauration testée et preuve d'intégrité |

## 3. Politique 3-2-1-1-0
3 copies, 2 supports, 1 copie hors site, 1 copie offline/immutable, 0 erreur de restauration vérifiée.

## 4. Classification
| Classe | Données | Fréquence | Rétention | Test |
|---|---|---:|---:|---:|
| B0 Vital | IAM, registre identité, clés recovery, configs P0 | continu + quotidien | 7 ans ou loi | hebdo |
| B1 Critique | transactions, enrôlement, audits | horaire/quotidien | 3-7 ans | mensuel |
| B2 Essentiel | portails, reporting, documents | quotidien | 1-3 ans | trimestriel |
| B3 Standard | données non critiques | hebdo | politique | semestriel |

## 5. Gouvernance
| Rôle | Responsabilité |
|---|---|
| Backup Owner National | politique, exceptions, budget |
| Data Owner | classification et exigences légales |
| Platform Team | jobs backup/restore |
| Security Team | chiffrement, immutabilité, accès |
| Resilience Office | exercices, KPI, conformité |
| Internal Audit | vérification indépendante |

## 6. Validation obligatoire
Chaque backup critique produit : statut job, checksum, contrôle complétude, restauration échantillon, test applicatif, RPO réel et ticket de preuve.

## 7. Air gap et immutabilité
Copies chiffrées, versionnées, WORM/Object Lock ou déconnectées physiquement/logiquement. Connexion du vault uniquement pendant fenêtres contrôlées, avec double contrôle et journal de chaîne de garde.
