---
# ============================================================
# SNISID-Security — National Immigration & Border Security
# Contrôle des Frontières (Aéroports, Ports, Terrestre)
# Document ID: SNISID-BORDER-001
# Version: 1.0.0
# ============================================================

## 1. CONTRÔLE DES FRONTIÈRES (DGIE & POLIFRONT)

Le système de contrôle aux frontières est la barrière de sécurité de l'État. Il s'assure qu'aucun individu recherché ou interdit de sortie ne puisse quitter le territoire (ex: Aéroport Toussaint Louverture, frontière Malpasse).

## 2. WORKFLOW DE VÉRIFICATION FRONTALIER

Chaque passage (Entrée ou Sortie) déclenche le workflow suivant en quelques millisecondes :

1. **Capture d'Identité :**
   - Haïtien : Scan passeport ou Carte SNISID + Biométrie faciale 1:1.
   - Étranger : Scan passeport + Prise d'empreintes/Photo (génère un NIU temporaire).
2. **Consultation Identité (SNISID-Core) :**
   - Le statut de l'identité est-il `ACTIF` ? (Bloque si `DECEDE` ou `SUSPENDU`).
3. **Vérification Watchlist (Correlation Engine) :**
   - Interdiction de quitter le territoire (IST) ?
   - Mandat d'arrêt actif ?
   - Signalement Interpol ?
4. **Décision et Logging :**
   - Si Vert : Événement `BorderCrossing` enregistré dans Kafka.
   - Si Rouge : Alerte discrète générée pour POLIFRONT / DCPJ, événement `BorderRefusal` enregistré.

## 3. ARCHITECTURE DES KIOSQUES (OFFLINE-READY)

Comme les frontières terrestres (Ouanaminthe, Belladère) peuvent souffrir de pannes réseau, les kiosques frontaliers utilisent l'architecture Offline-First (K3s + NATS).

- **Cache de Sécurité :** Le K3s frontalier maintient un cache chiffré de l'ensemble de la *No-Fly List* et des *Mandats d'arrêt actifs*.
- **Opération Offline :** En cas de coupure, la frontière n'est pas ouverte à l'aveugle. Les voyageurs sont vérifiés contre le cache local.
- **Sync :** Dès le rétablissement, les `BorderCrossingEvents` sont envoyés au serveur central.

## 4. API D'IMMIGRATION

```yaml
openapi: 3.1.0
paths:
  /border/crossings:
    post:
      summary: Enregistrer un passage frontalier
      requestBody:
        content:
          application/json:
            schema:
              properties:
                niu: { type: string }
                port_id: { type: string }
                direction: { enum: [ENTRY, EXIT] }
                biometric_match_score: { type: number }
  /border/watchlists/check:
    post:
      summary: Vérifier si l'individu est sur une liste de restriction
```

---
*Document ID: SNISID-BORDER-001 | Approuvé par: DGIE*
