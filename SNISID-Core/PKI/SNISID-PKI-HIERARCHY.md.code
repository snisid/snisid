---
# ============================================================
# SNISID-Core — PKI Nationale Souveraine
# Hiérarchie PKI + EJBCA + cert-manager + OCSP
# Document ID: SNISID-PKI-001
# Version: 1.0.0
# ============================================================

# ══════════════════════════════════════════════════════════
# HIÉRARCHIE PKI NATIONALE SNISID
# ══════════════════════════════════════════════════════════
#
#   SNISID Root CA (Offline — HSM FIPS 140-2 L3 — Coffre-fort État)
#   ├── SNISID Issuing CA v1 (Online — HSM)
#   │   ├── Certificats citoyens (identité + signature)
#   │   ├── Certificats OEC (officiers état civil)
#   │   ├── Certificats agents ONI
#   │   └── Certificats PKI services SNISID
#   │
#   ├── SNISID TLS CA v1 (Online — Pour HTTPS)
#   │   ├── api.snisid.gov.ht
#   │   ├── *.snisid.gov.ht
#   │   └── Certificats inter-services (mTLS)
#   │
#   └── SNISID Device CA v1 (Online — Pour kits MEK)
#       ├── Certificats tablettes MEK
#       ├── Certificats scanners biométriques
#       └── Certificats edge nodes
#
# ══════════════════════════════════════════════════════════

## 1. ROOT CA — Procédure Cérémonie des Clés

La cérémonie des clés Root CA est le **moment le plus critique** de la PKI nationale.
Elle est réalisée HORS LIGNE dans un environnement sécurisé (coffre-fort ANH).

### 1.1 Matériel Requis

| Élément | Quantité | Description |
|---------|---------|-------------|
| HSM Luna Network HSM A750 | 2 | Primaire + backup |
| Cartes à puce administrateur | 5 (Shamir 3/5) | Contrôle accès HSM |
| Laptop air-gapped (Ubuntu 22.04 LTS) | 1 | Génération certificat |
| Clés USB chiffrées | 5 | Backup fragments Shamir |
| Onduleur 2000 VA | 1 | Protection coupure |
| Caméra vidéo + deux notaires | 1 + 2 | Témoins légaux |
| Représentants: AND, CISO, NDPA | 3 | Autorités |

### 1.2 Commandes de Génération Root CA (OpenSSL)

```bash
#!/bin/bash
# ATTENTION: Exécuter uniquement dans l'environnement air-gapped certifié

# ─── Root CA ───────────────────────────────────────────────
# Générer la clé privée Root CA (RSA-4096 dans HSM)
pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so \
    --slot 0 --pin $HSM_PIN \
    --keypairgen --key-type RSA:4096 \
    --id 01 --label "SNISID-RootCA-Key-v1"

# Créer le certificat Root CA (X.509 v3)
openssl req -x509 -new -nodes \
    -engine pkcs11 -keyform engine \
    -key "pkcs11:token=SNISID-HSM;object=SNISID-RootCA-Key-v1;type=private" \
    -sha512 \
    -days 7300 \
    -subj "/C=HT/O=République d'Haïti/OU=Autorité Nationale Numérique/CN=SNISID Root CA v1" \
    -extensions v3_ca \
    -extfile <(cat <<EOF
[v3_ca]
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
basicConstraints = critical, CA:true, pathlen:2
keyUsage = critical, digitalSignature, keyCertSign, cRLSign
certificatePolicies = 1.3.6.1.4.1.99999.1.1
nameConstraints = critical,permitted;subtree:DNS:.gov.ht,permitted;subtree:DNS:.snisid.gov.ht
EOF
    ) \
    -out snisid-root-ca.crt

# Vérifier
openssl x509 -in snisid-root-ca.crt -text -noout | head -50
openssl verify -CAfile snisid-root-ca.crt snisid-root-ca.crt

echo "✅ Root CA généré: snisid-root-ca.crt"
echo "   SHA-256 fingerprint:"
openssl x509 -in snisid-root-ca.crt -fingerprint -sha256 -noout

# ─── Issuing CA ────────────────────────────────────────────
# Générer la clé Issuing CA (EC P-384 dans HSM)
pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so \
    --slot 0 --pin $HSM_PIN \
    --keypairgen --key-type EC:P-384 \
    --id 02 --label "SNISID-IssuingCA-Key-v1"

# CSR Issuing CA
openssl req -new -engine pkcs11 -keyform engine \
    -key "pkcs11:token=SNISID-HSM;object=SNISID-IssuingCA-Key-v1;type=private" \
    -sha384 \
    -subj "/C=HT/O=République d'Haïti/OU=SNISID PKI/CN=SNISID Issuing CA v1" \
    -out snisid-issuing-ca.csr

# Signer avec Root CA
openssl ca -engine pkcs11 -keyform engine \
    -keyfile "pkcs11:token=SNISID-HSM;object=SNISID-RootCA-Key-v1;type=private" \
    -cert snisid-root-ca.crt \
    -in snisid-issuing-ca.csr \
    -days 3650 \
    -extfile <(cat <<EOF
basicConstraints = critical, CA:true, pathlen:0
keyUsage = critical, digitalSignature, keyCertSign, cRLSign
authorityInfoAccess = OCSP;URI:http://ocsp.snisid.gov.ht, caIssuers;URI:http://crt.snisid.gov.ht/snisid-root-ca.crt
cRLDistributionPoints = URI:http://crl.snisid.gov.ht/snisid-root-ca.crl
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always
EOF
    ) \
    -out snisid-issuing-ca.crt
```

### 1.3 Distribution Shamir 3/5

```python
# key_ceremony/shamir_split.py
# EXÉCUTER UNIQUEMENT LORS DE LA CÉRÉMONIE

from secretsharing import SecretSharer
import qrcode
import base64

def split_hsm_pin(pin: str, total: int = 5, threshold: int = 3):
    """Diviser le PIN HSM en 5 parts — 3 suffisent pour reconstiture"""
    shares = SecretSharer.split_secret(pin, threshold, total)
    print(f"✅ PIN divisé en {total} parts (seuil: {threshold})")
    
    holders = [
        "AND — Directeur Général",
        "CISO National SNISID",
        "NDPA — Directeur",
        "Ministère Justice — Secrétaire d'État",
        "CNN — Représentant Société Civile"
    ]
    
    for i, (share, holder) in enumerate(zip(shares, holders)):
        # Générer QR code pour chaque part
        qr = qrcode.make(share)
        qr.save(f"share_{i+1}_{holder.replace(' ', '_')}.png")
        print(f"  Part {i+1} → {holder}: IMPRIMÉ + USB chiffré")
    
    print("\n⚠️  DÉTRUIRE toute trace numérique de ce script après exécution")
    print("📜  Procès-verbal notarial signé obligatoire")

split_hsm_pin(pin=input("PIN HSM Master: "))
```

---

## 2. CERT-MANAGER — Intégration EJBCA

```yaml
# ejbca-issuer.yaml
apiVersion: ejbca.keyfactor.com/v1alpha1
kind: Issuer
metadata:
  name: snisid-ejbca-issuer
  namespace: snisid-security
spec:
  ejbcaSecretName: ejbca-credentials
  hostname: https://ejbca.snisid.gov.ht
  ejbcaRestCaName: SNISID-IssuingCA-v1
  certificateProfileName: SNISIDCitizenProfile
  endEntityProfileName: SNISIDCitizenEEProfile

---
# CertificateProfile: SNISIDCitizenProfile
# (Configuré dans l'interface EJBCA)
# Paramètres:
#   - Type: End Entity
#   - Algorithme: EC P-384
#   - Validité: 5 ans
#   - Extensions:
#     * Key Usage: digitalSignature, nonRepudiation
#     * Extended Key Usage: id-kp-clientAuth, id-kp-emailProtection
#     * Subject Alternative Name: URI:ht:gov:snisid:citizen:{NIU}
#     * CDP: http://crl.snisid.gov.ht/snisid-issuing-ca.crl
#     * AIA OCSP: http://ocsp.snisid.gov.ht
#   - Publication: No (manual)
#   - Revocation: OCSP + CRL delta (15 min)
```

---

## 3. OCSP RESPONDER

```yaml
# ocsp-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ocsp-responder
  namespace: snisid-security
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ocsp-responder
  template:
    metadata:
      labels:
        app: ocsp-responder
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
      containers:
      - name: ocsp-responder
        image: harbor.snisid.gov.ht/ejbca/ejbca-ce:8.0.0
        args:
        - ocsp-responder
        env:
        - name: OCSP_SIGNING_CERT
          valueFrom:
            secretKeyRef:
              name: ocsp-signing-cert
              key: tls.crt
        - name: OCSP_CACHE_SIZE
          value: "100000"     # 100K certs en cache
        - name: OCSP_CACHE_TTL
          value: "600"        # 10 minutes cache validité
        - name: EJBCA_ENDPOINT
          value: https://ejbca.snisid-security.svc.cluster.local:8443
        ports:
        - name: http
          containerPort: 8080
        resources:
          requests:
            memory: 256Mi
            cpu: 250m
          limits:
            memory: 1Gi
            cpu: 500m
        livenessProbe:
          httpGet:
            path: /healthz
            port: http
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /readyz
            port: http
          initialDelaySeconds: 10
          periodSeconds: 5
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop: [ALL]

---
apiVersion: v1
kind: Service
metadata:
  name: ocsp-responder
  namespace: snisid-security
  annotations:
    # OCSP doit être accessible publiquement via http://ocsp.snisid.gov.ht
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
spec:
  selector:
    app: ocsp-responder
  ports:
  - port: 80
    targetPort: 8080
    name: http
  type: LoadBalancer

---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ocsp-responder-hpa
  namespace: snisid-security
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ocsp-responder
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 60
  - type: External
    external:
      metric:
        name: ocsp_requests_per_second
      target:
        type: AverageValue
        averageValue: "1000"

---
# CRL Distribution — NGINX serving static CRL files
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crl-distributor
  namespace: snisid-security
spec:
  replicas: 2
  selector:
    matchLabels:
      app: crl-distributor
  template:
    metadata:
      labels:
        app: crl-distributor
    spec:
      containers:
      - name: nginx
        image: harbor.snisid.gov.ht/nginx/nginx:1.25-alpine
        ports:
        - containerPort: 80
        volumeMounts:
        - name: crl-files
          mountPath: /usr/share/nginx/html
          readOnly: true
        resources:
          requests:
            memory: 64Mi
            cpu: 50m
          limits:
            memory: 256Mi
            cpu: 200m
      - name: crl-updater
        image: harbor.snisid.gov.ht/snisid/crl-sync:1.0.0
        env:
        - name: EJBCA_ENDPOINT
          value: https://ejbca.snisid-security.svc.cluster.local:8443
        - name: CRL_UPDATE_INTERVAL
          value: "900"    # Sync CRL toutes les 15 minutes
        - name: OUTPUT_DIR
          value: /crl-files
        volumeMounts:
        - name: crl-files
          mountPath: /crl-files
      volumes:
      - name: crl-files
        emptyDir: {}
```

---

## 4. PROFILS DE CERTIFICATS CITOYENS

| Profil | Usage | Algorithme | Validité | Extensions |
|--------|-------|-----------|---------|-----------|
| **SNISIDCitizenAuth** | Authentification NID | EC P-384 | 5 ans | id-kp-clientAuth |
| **SNISIDCitizenSign** | Signature documents | EC P-384 | 5 ans | id-kp-emailProtection, nonRepudiation |
| **SNISIDAgentAuth** | Agents ONI/OEC | EC P-384 | 2 ans | id-kp-clientAuth |
| **SNISIDServiceTLS** | Services SNISID | EC P-384 | 1 an | id-kp-serverAuth |
| **SNISIDDeviceMEK** | Kits terrain | EC P-256 | 3 ans | id-kp-clientAuth |
| **SNISIDOCSPSign** | Répondeur OCSP | EC P-384 | 1 an | id-kp-OCSPSigning |

---

*Document ID : SNISID-PKI-001 v1.0.0 — Mai 2026*  
*Cérémonie des clés Root CA supervisée par : AND, CISO, NDPA, MJ, CNN*  
*Classification : SOUVERAIN / CRYPTOGRAPHIQUE SENSIBLE*
