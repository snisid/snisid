# Charte Graphique et Guide de Design UI/UX SNISID (v2.0)

**Classification :** RESTREINT / DESIGN SYSTEM  
**Recommandation de Référence :** SNISID v2.0 — MP-012

Ce document définit les tokens de conception visuelle, la palette de couleurs officielle, et les règles d'expérience utilisateur (UX) obligatoires pour toutes les interfaces Web et Mobile du **Système National d'Identification Sécurisé et d'Interopérabilité Digitale (SNISID)**.

---

## 📂 Fichiers de Référence dans le Workspace

* **Styles & Intégration** :
  * [index.css](file:///c:/Users/sopil/Desktop/snisid%20system/frontend/src/index.css) — Intégration des variables CSS globales de couleurs.
  * [verify_ux_rules.py](file:///c:/Users/sopil/Desktop/snisid%20system/pki/scripts/verify_ux_rules.py) — Script d'audit de la palette CSS et du dictionnaire i18n.

---

## 1. Palette de Couleurs Officielle SNISID

Toutes les interfaces doivent se conformer aux variables HEX définies ci-dessous pour garantir l'identité de marque et le confort visuel (mode sombre par défaut).

| Élément | Code HEX | Classe CSS Variable | Usage UI |
|---|---|---|---|
| **Background Principal** | `#0D1B2A` | `--snisid-bg` | Fond d'écran global des dashboards |
| **Accent Primaire** | `#1565C0` | `--snisid-primary` | Boutons d'action principaux (CTAs), liens actifs |
| **Accent Secondaire** | `#00BCD4` | `--snisid-secondary` | Graphiques réseau, indicateurs AIOps |
| **Statut OK** | `#2E7D32` | `--snisid-ok` | Succès d'authentification, services opérationnels |
| **Alerte** | `#E65100` | `--snisid-warning` | Alertes SLA, déconnexions mineures, avertissements |
| **Critique** | `#C62828` | `--snisid-critical` | Incidents de sécurité P1, intrusions, alertes SOC |
| **Texte Principal** | `#ECEFF1` | `--snisid-text` | Texte brut, labels de formulaires |
| **Surface Cards** | `#1E3A5F` | `--snisid-surface` | Conteneurs de cartes, bordures de fenêtres |

---

## 2. Règles UX Obligatoires (Garanties Système)

### 2.1. Indicateur Réseau Permanent (Online / Offline)
* Un badge d'état de connectivité réseau doit être visible à tout moment en haut à droite de l'écran. 
* L'indicateur passe à l'état `Offline` en rouge (`#C62828`) dès que la connexion NATS/HTTP est perdue, rappelant à l'agent de terrain que les enrôlements sont en attente de synchronisation locale.

### 2.2. Feedback Action < 300ms
* Pour toute action de l'utilisateur (clic sur enrôlement, validation biométrique), le système doit fournir un retour visuel en moins de 300 ms.
* Si le temps de traitement de l'API dépasse 300 ms, un spinner animé avec un message d'attente s'affiche pour éviter les clics répétés.

### 2.3. Traduction Systématique (Français / Créole Haïtien)
Toutes les chaînes de caractères de l'interface utilisateur, notamment les messages d'erreur critiques, doivent être disponibles dans les deux langues officielles de la République d'Haïti.

#### Structure type de dictionnaire de traduction :
```json
{
  "errors": {
    "auth_failed": {
      "fr": "Échec de l'authentification : empreinte non reconnue.",
      "ht": "Echèk otantifikasyon : anprent pa rekonèt."
    },
    "network_offline": {
      "fr": "Terminal hors-ligne. Enregistrement en attente de synchronisation.",
      "ht": "Kip la deploge nan rezo. Enskripsyon an ap tann pou senkronize."
    }
  }
}
```

### 2.4. Protection des Données Biométriques
* **Zéro exposition brute :** Aucune image d'empreinte digitale, scan d'iris ou photographie ICAO brute ne doit être affichée à l'écran après la capture.
* Seuls les indicateurs de qualité (ex : "Qualité Empreinte : 92%") ou des masques vectorisés floutés de contrôle d'alignement sont autorisés.

### 2.5. Déconnexion Automatique de Sécurité
* Pour les accès de niveau substantiel ou critique (AAL2+ et AAL3), le frontend déclenche un décompte d'inactivité.
* Après **5 minutes sans interaction**, la session est révoquée, les clés de déchiffrement en RAM sont purgées et l'utilisateur est redirigé vers l'écran de login.

### 2.6. Double Confirmation des Actions Irréversibles
* Toute action destructive (effacement de cache, modification de NNI, suspension de certificat) requiert une double confirmation :
  1. Un clic initial sur l'action.
  2. Une fenêtre modale demandant la saisie manuelle d'un mot de code (ex : écrire "SUSPENDRE" ou "EFFACER") ou la confirmation biométrique.

### 2.7. Accessibilité et Contraste Élevé
* Un interrupteur d'accessibilité permet de basculer en mode **Contraste Élevé** (`.high-contrast`), appliquant des contrastes de type noir pur et jaune vif pour les personnes malvoyantes.
* Tous les rapports émis par le système intègrent un filigrane diagonal en opacité réduite marqué `"CONFIDENTIEL — SNISID"`.

---

## 3. Recommandations Spécifiques (v2.0 — MP-012)

### 3.1. Protocole de Tests Terrain en Haïti
* **Scénarios de Stress Physiques :** Organiser des sessions d'utilisabilité dans les communes rurales (ex : Artibonite, Grand'Anse) en conditions réelles : écran sous lumière directe du soleil (test anti-reflet), pluie fine tropicale, niveau de batterie inférieur à 15%.
* **Profils Utilisateurs :** Inclure des citoyens n'ayant jamais interagi avec des terminaux numériques et des agents ONI fatigués en fin de journée pour simplifier les flux de navigation.

### 3.2. Accessibilité WCAG 2.1 AA et Lecteurs d'Écran
* Avant le passage à la Phase 3, l'interface du Portail Citoyen doit être entièrement navigable au clavier seul (sans souris) à l'aide des touches `Tab` et `Enter`.
* Tous les éléments interactifs doivent inclure des attributs d'accessibilité ARIA (ex: `aria-label`, `aria-live="polite"`) compatibles avec les lecteurs d'écran NVDA et JAWS.

---

*Ce guide de style UX est approuvé pour le déploiement des interfaces de production SNISID.*
