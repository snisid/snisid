---
# ============================================================
# SNISID-Cyber — Zero Trust Architecture (ZTA)
# Microsegmentation et Vérification Continue
# Document ID: SNISID-ZTA-001
# Version: 1.0.0
# ============================================================

## 1. PARADIGME "NEVER TRUST, ALWAYS VERIFY"

Dans les architectures classiques, une fois le pare-feu franchi (VPN), l'utilisateur a accès à tout le réseau interne. Dans le programme SNISID, l'adresse IP ne confère **aucune confiance**.

## 2. PILIERS DU ZERO TRUST SNISID

1. **Identity is the Perimeter :** L'authentification est requise pour *chaque* requête, évaluée par des politiques contextuelles (Heure, Lieu, Device, Anomalie comportementale).
2. **Device Trust :** Même avec le bon mot de passe, un accès au back-office est refusé si le laptop n'est pas un matériel "Government Issued" (certificat client vérifié) avec un EDR à jour.
3. **Mutual TLS (mTLS) :** Toute communication entre deux serveurs (ex: API Gateway vers CockroachDB) est chiffrée. Un hacker écoutant le trafic interne ne verra qu'une bouillie cryptographique.
4. **Microsegmentation (Least Privilege) :** Les serveurs/pods ne peuvent parler qu'à ce dont ils ont strictement besoin. Si un pod "Frontend-Web" est piraté, il ne pourra jamais ouvrir une connexion SSH vers un autre serveur, car la politique réseau eBPF le bloque au niveau du noyau (Kernel).

## 3. CONTEXTUAL AUTHORIZATION (OPA)

Le moteur Open Policy Agent (OPA) vérifie chaque action.
```rego
package snisid.access
default allow = false

# Un médecin ne peut accéder à un dossier médical que depuis le réseau de son hôpital,
# pendant ses heures de travail.
allow {
    input.user.role == "medecin"
    input.request.time > "08:00"
    input.request.time < "18:00"
    input.network.location == "hopital_general_pap"
}
```

---
*Document ID: SNISID-ZTA-001 | Approuvé par: Architecte Zero Trust*
