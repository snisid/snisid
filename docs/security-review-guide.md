# Guide de Revue de Code Sécurisé — Security Owner

## Principes Fondamentaux

1. **Zero Trust dans le Code** : Ne jamais faire confiance aux entrées, toujours valider.
2. **Defense in Depth** : Multiples couches de sécurité (validation, chiffrement, audit).
3. **Least Privilege** : Chaque service a les permissions minimales strictement nécessaires.
4. **Secure by Default** : Les configurations par défaut doivent être sécurisées.

## Checklist de Revue

### 1. Secrets et Credentials
- [ ] Aucun mot de passe, clé API, token JWT en dur dans le code
- [ ] Les secrets utilisent Vault (via `vault:` annotations ou External Secrets Operator)
- [ ] Les `.env` et `secrets.*` sont dans `.gitignore`
- [ ] Les certificats et clés privées ne sont jamais commités
- [ ] Les mots de passe par défaut sont interdits

**Patterns à rejeter:**
```go
// REJETÉ
password := "mon-mot-de-passe"
apiKey := "sk-1234567890"
jwtSecret := "change-me-in-production"
```

**Patterns acceptés:**
```go
// ACCEPTÉ
password := os.Getenv("DB_PASSWORD")
// ou via Vault:
client, _ := vault.NewClient(vault.WithAddress(os.Getenv("VAULT_ADDR")))
secret, _ := client.Logical().Read("kv-v2/data/database")
password := secret.Data["password"]
```

### 2. Injection SQL
- [ ] Aucune utilisation de `fmt.Sprintf` ou `+` pour construire des requêtes SQL
- [ ] Toutes les requêtes utilisent des paramètres positionnels (`$1`, `$2`, ...)
- [ ] Pas de `ORDER BY` ou `LIMIT` dynamique sans validation stricte
- [ ] En Python: utiliser SQLAlchemy ORM ou `cursor.execute("...", params)`

**À rejeter:**
```go
db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userID))
db.Exec(fmt.Sprintf("UPDATE users SET %s = TRUE WHERE id = $1", column))
```

**Accepté:**
```go
db.Query("SELECT * FROM users WHERE id = $1", userID)
// Pour colonnes dynamiques: utiliser une whitelist
validColumns := map[string]bool{"active": true, "verified": true}
if !validColumns[column] { return error }
db.Exec(fmt.Sprintf("UPDATE users SET %s = TRUE WHERE id = $1", column), userID)
```

### 3. Communications Sécurisées
- [ ] Toutes les connexions gRPC utilisent mTLS (pas de `insecure.NewCredentials()`)
- [ ] Les connexions HTTP utilisent TLS 1.3 minimum
- [ ] Les connexions Redis, PostgreSQL, Kafka sont chiffrées
- [ ] Pas de protocoles non chiffrés (HTTP -> 301 HTTPS)

**À rejeter:**
```go
conn, _ := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
```

**Accepté:**
```go
creds := credentials.NewTLS(&tls.Config{
  MinVersion: tls.VersionTLS13,
  ServerName: "biometrics.snisid.svc.cluster.local",
})
conn, _ := grpc.Dial(address, grpc.WithTransportCredentials(creds))
```

### 4. Gestion d'Erreurs
- [ ] Aucune erreur ignorée avec `_ =` ou `_, _ =`
- [ ] Toutes les erreurs sont loggées (avec `zap` ou `logur`)
- [ ] Les `panic()` sont interdites dans les chemins de production
- [ ] Les erreurs sensibles ne sont jamais exposées aux clients

**À rejeter:**
```go
_ = db.QueryRow(...)  // erreur ignorée
_, _ = w.Write([]byte("ok"))  // write error ignoré
panic("failed to init: " + err.Error())  // panic en production
```

**Accepté:**
```go
if err := db.QueryRow(...).Scan(&result); err != nil {
  logger.Error("database query failed", zap.Error(err))
  return fmt.Errorf("query user: %w", err)
}
```

### 5. Authentification et Authorisation
- [ ] Les tokens JWT sont signés avec une clé forte (RS256 ou ES256)
- [ ] Les sessions ont une durée de vie limitée (TTL court)
- [ ] Les mots de passe utilisent bcrypt/argon2 (pas SHA256)
- [ ] L'autorisation est vérifiée à chaque requête (pas seulement à l'auth)
- [ ] Les endpoints API ont des rate limits

### 6. Validation des Entrées
- [ ] Toutes les entrées utilisateur sont validées (type, format, limites)
- [ ] Les IDs sont des UUID, jamais des IDs séquentiels
- [ ] Les téléchargements de fichiers sont limités en taille et type
- [ ] Les payloads JSON ont une taille maximale

### 7. Audit et Logging
- [ ] Toutes les opérations sensibles sont auditées (création, modification, suppression)
- [ ] Les logs ne contiennent jamais de secrets (passwords, tokens, PII)
- [ ] Les logs sont centralisés (Loki, Elasticsearch)
- [ ] Les événements d'authentification sont toujours loggés

### 8. Dépendances
- [ ] Aucune dépendance avec des vulnérabilités connues (CVE)
- [ ] Les dépendances sont épinglées à des versions exactes
- [ ] Les dépendances Go utilisent `go.sum` (vérification des checksums)
- [ ] Les images Docker sont scannées pour les vulnérabilités

## Processus de Revue

1. **Automated Gate** : Les scans (gitleaks, trivy, golangci-lint) bloquent le merge
2. **Peer Review** : Chaque PR nécessite au moins 1 approbation
3. **Security Review** : Les changements sensibles nécessitent l'approbation du Security Owner
4. **Sign-off Tags** : Les releases sont signées avec GPG

## Bloqueurs Critiques (DO NOT MERGE)

Toute PR contenant:
- Secrets hardcodés
- Requêtes SQL non paramétrées
- Communications non chiffrées
- Mocks en chemins de production
- panic() dans des handlers de requêtes
