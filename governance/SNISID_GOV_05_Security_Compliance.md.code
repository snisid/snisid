# VOLUME 5 : Sécurité de la Gouvernance et Conformité (Security & Compliance)
## Architecture de Gouvernance — SNISID

La confiance étatique ne se déclare pas, elle se prouve. La gouvernance sécuritaire du SNISID est conçue pour prouver la conformité aux standards internationaux à chaque instant, sans intervention humaine (Continuous Compliance).

---

## 🔒 CHAPITRE 1 : GOUVERNANCE ZERO TRUST ET CONTRÔLE D'ACCÈS

Le système opère sous le principe de "Moindre Privilège Absolu". L'identité d'un agent de l'État ne suffit pas à lui accorder l'accès ; le contexte de sa requête est évalué dynamiquement.

### 1.1 RBAC (Role-Based Access Control)
*   Les rôles sont définis légalement, pas techniquement. Un "Officier d'État Civil" possède le droit de valider une naissance.
*   **Séparation des Pouvoirs (Segregation of Duties) :** L'ingénieur informatique (Sysadmin) qui maintient la base de données n'a **jamais** accès aux données en clair. Les données sont chiffrées au repos (TDE) et en transit.

### 1.2 ABAC (Attribute-Based Access Control)
*   **Gouvernance Contextuelle :** Un officier de Jacmel (Sud-Est) ne peut consulter le dossier d'un citoyen du Cap-Haïtien (Nord) que s'il peut justifier légalement (via numéro de dossier Workflow) que ce citoyen se présente physiquement dans son bureau.
*   **Heures Ouvrables :** L'accès aux interfaces de modification est verrouillé en dehors des heures d'ouverture de l'État, sauf pour les rôles d'urgence (Hôpitaux, Police).

---

## 🛡️ CHAPITRE 2 : GOUVERNANCE DES MENACES INTERNES (INSIDER THREAT)

La plus grande menace pour une infrastructure identitaire est la corruption interne.

### 2.1 Approbation Multiple (Two-Man Rule)
*   Pour toute action critique (ex: modification d'une filiation, révocation d'une nationalité), le système exige la signature cryptographique (FIDO2) de **deux officiers assermentés différents**.
*   **Audit WORM :** Si un agent de l'Immigration effectue une recherche massive (ex: 50 requêtes en 1 heure sur des citoyens spécifiques), le SOC est alerté (UEBA) et l'accès de l'agent est suspendu préventivement.

---

## 📋 CHAPITRE 3 : CONFORMITÉ NATIONALE ET INTERNATIONALE

Le SNISID est audité face aux normes de cybersécurité mondiales pour faciliter l'interopérabilité des passeports haïtiens à l'international (OACI).

### 3.1 Cadre Normatif
*   **ISO/IEC 27001 :** Système de Management de la Sécurité de l'Information (SMSI).
*   **ISO 22301 :** Gestion de la Continuité d'Activité (Garantie par la conception Multi-Cluster).
*   **NIST Cybersecurity Framework (CSF) :** Utilisé par le SOC pour structurer ses opérations (Identify, Protect, Detect, Respond, Recover).

### 3.2 Conformité Automatisée (Policy-as-Code)
La gouvernance ne repose pas sur des manuels papier audités une fois par an.
*   Le moteur **Open Policy Agent (OPA)** vérifie en temps réel que chaque configuration déployée sur le cluster Kubernetes est conforme à l'ISO 27001.
*   Si un développeur tente de désactiver le chiffrement TLS d'une base de données, la plateforme rejette techniquement la demande. L'auditabilité est intégrée au code.
