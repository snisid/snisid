# National Cyber Resilience Model

## 1. Objectif
Transformer la cybersécurité SNISID en résilience active contre ransomwares, DDoS, attaques internes et compromissions d'identité.

## 2. Capacités
| Fonction | Support | Mécanisme |
|---|---:|---|
| Ransomware resilience | Oui | immutabilité, segmentation, clean restore |
| DDoS resilience | Oui | scrubbing, rate limit, mode dégradé |
| Insider attack resilience | Oui | PAM, journaux immuables, double contrôle |
| Identity compromise containment | Oui | révocation sessions, rotation clés, isolation IAM |

## 3. Doctrine
Assume breach → isoler vite → restaurer proprement → maintenir P0 → prouver l'intégrité.

## 4. Ransomware resilience
Backups immuables/offline, segmentation, moindre privilège, détection chiffrement anormal, blocage comptes suspects, restauration isolée, scan malware avant retour.

## 5. DDoS resilience
Protection upstream, limitation débit API, files d'attente, mode lecture seule, priorisation agences P0/P1, endpoints crise et capacité offline régionale.

## 6. Insider resilience
PAM, just-in-time access, séparation devoirs, dual control, immutable logs, monitoring comportemental administrateur.

## 7. Containment identité compromise
Isoler IAM affecté → révoquer sessions/tokens → activer break-glass propre → suspendre workflows à risque → vérifier registre → restaurer point sain si nécessaire → rotation secrets/clés → audit.

## 8. Clean room recovery
Infrastructure neuve par IaC, images signées, secrets régénérés, backups scannés, journalisation renforcée, ouverture progressive.

## 9. Scénarios à tester
Ransomware DC primaire, compte admin compromis, DDoS APIs identité, corruption registre, fuite token inter-agence.
