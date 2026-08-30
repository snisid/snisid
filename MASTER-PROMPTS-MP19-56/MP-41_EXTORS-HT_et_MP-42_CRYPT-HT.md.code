# MP-41 — EXTORS-HT
## Registre National des Extorsions, Rançons et Économie Criminelle des Gangs
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-41 | Code : EXTORS-HT
Dépendances      : GANG-HT (MP-24), SIVC-HT (MP-18), UCREF-INT (MP-39), BLAN-HT (MP-40)
Normes           : GAFI/FATF Rec.5 (financement terrorisme), INTERPOL Financial Crime
Acteurs          : DCPJ CAE, UCREF, MJSP, BRH (Banque de la République d'Haïti)
```

---

## 1. CONTEXTE

L'économie des gangs haïtiens repose sur trois piliers documentés :
- **Rançons kidnapping** : 10,000–5,000,000 USD par victime selon le profil
- **Péages illicites** : Barrières sur RN1, RN2, RN3 — taxation systématique véhicules et marchands
- **Extorsion commerce** : « Taxe » mensuelle imposée aux commerçants dans zones contrôlées

Ce module trace tous les flux d'extorsion identifiés, les victimes, les montants et
les canaux de paiement pour alimenter les dossiers UCREF et les poursuites judiciaires.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE extors_type AS ENUM (
    'KIDNAPPING_RANSOM',
    'ROAD_TOLL_ILLEGAL',
    'BUSINESS_PROTECTION_RACKET',
    'REAL_ESTATE_EXTORTION',
    'PUBLIC_SERVANT_EXTORTION',
    'NGO_EXTORTION',
    'FUEL_TRUCK_HIJACK',
    'OTHER'
);

CREATE TYPE extors_payment_channel AS ENUM (
    'MONCASH', 'NATCASH', 'DIGICEL_MONEY',
    'WIRE_TRANSFER', 'CASH_DROP',
    'CRYPTOCURRENCY', 'INTERMEDIARY', 'UNKNOWN'
);

CREATE TYPE extors_status AS ENUM (
    'ACTIVE','PAID','REFUSED','NEGOTIATING',
    'LAW_ENFORCEMENT_INVOLVED','RESOLVED','VICTIM_HARMED'
);

CREATE TABLE extors_cases (
    case_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_extors_id  VARCHAR(25) UNIQUE NOT NULL,  -- EXTORS-HT-AAAA-NNNNNN
    extors_type         extors_type NOT NULL,
    status              extors_status NOT NULL DEFAULT 'ACTIVE',

    -- Perpetrateurs
    gang_id             UUID,
    gang_name           VARCHAR(150),
    perpetrator_ids     UUID[] DEFAULT '{}',       -- CHEF-HT member IDs
    chef_member_ids     UUID[] DEFAULT '{}',

    -- Victimes
    victim_count        SMALLINT DEFAULT 1,
    victim_snisid_ids   UUID[] DEFAULT '{}',
    victim_types        TEXT[] DEFAULT '{}',       -- INDIVIDUAL, BUSINESS, NGO, GOVERNMENT
    victim_nationality  CHAR(3)[] DEFAULT '{}',
    is_foreigner_victim BOOLEAN DEFAULT FALSE,

    -- Localisation
    incident_location   VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    route_number        VARCHAR(10),               -- RN1, RN2 etc pour les peages

    -- Montants
    demanded_amount     DECIMAL(15,2),
    demanded_currency   CHAR(3) DEFAULT 'USD',
    paid_amount         DECIMAL(15,2),
    paid_currency       CHAR(3),
    payment_channel     extors_payment_channel,
    payment_ref         VARCHAR(200),              -- Numero de transaction MonCash etc.
    payment_date        TIMESTAMPTZ,

    -- Chronologie
    first_contact_date  TIMESTAMPTZ NOT NULL,
    resolution_date     TIMESTAMPTZ,

    -- Enquête
    case_reference      VARCHAR(100),
    investigating_unit  VARCHAR(50),
    ucref_str_id        UUID,                      -- Lien STR si rançon tracée
    blan_case_id        UUID,                      -- Lien dossier blanchiment

    notes               TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE extors_road_toll_points (
    toll_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gang_id             UUID NOT NULL,
    location_desc       VARCHAR(300) NOT NULL,
    route_number        VARCHAR(10),
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    daily_revenue_usd   DECIMAL(10,2),
    vehicle_types_taxed TEXT[] DEFAULT '{}',
    toll_rates          JSONB,                     -- {moto: 50, voiture: 200, camion: 500}
    active_since        DATE,
    is_active           BOOLEAN DEFAULT TRUE,
    source_intel        TEXT,
    last_confirmed_at   TIMESTAMPTZ,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE extors_negotiations (
    neg_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id             UUID NOT NULL REFERENCES extors_cases(case_id),
    negotiation_date    TIMESTAMPTZ NOT NULL,
    contact_method      VARCHAR(50),               -- PHONE, INTERMEDIARY, DROP_NOTE
    contact_number      VARCHAR(30),
    demand_updated      DECIMAL(15,2),
    demand_currency     CHAR(3),
    position_update     TEXT,
    recorded_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_extors_cases_gang    ON extors_cases(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_extors_cases_type    ON extors_cases(extors_type, status);
CREATE INDEX idx_extors_cases_dept    ON extors_cases(dept_code, first_contact_date DESC);
CREATE INDEX idx_extors_cases_channel ON extors_cases(payment_channel) WHERE paid_amount IS NOT NULL;
CREATE INDEX idx_extors_tolls_route   ON extors_road_toll_points(route_number) WHERE is_active = TRUE;
CREATE INDEX idx_extors_tolls_dept    ON extors_road_toll_points(dept_code) WHERE is_active = TRUE;

COMMIT;
```

---

## 3. SERVICE GO CLÉ — ANALYSE DES PÉAGES ILLICITES

```go
package service

import (
    "context"
    "github.com/snisid/extors-svc/internal/domain"
)

// ComputeGangRevenue calcule le revenu estime d un gang via extorsion
func (s *ExtorsService) ComputeGangRevenue(
    ctx context.Context,
    gangID string,
) (*domain.GangRevenueReport, error) {
    report := &domain.GangRevenueReport{GangID: gangID}

    // 1. Revenus péages (sources documentees)
    tolls, _ := s.repo.FindActiveTollsByGang(ctx, gangID)
    for _, t := range tolls {
        report.TollRevenueDaily += t.DailyRevenueUSD
    }
    report.TollRevenueMonthly = report.TollRevenueDaily * 30

    // 2. Revenus rançons (cas documentes)
    ransoms, _ := s.repo.FindPaidRansomsByGang(ctx, gangID)
    for _, r := range ransoms {
        report.RansomRevenue += r.PaidAmountUSD
    }

    // 3. Extorsions commerciales (estimations)
    rackets, _ := s.repo.FindActiveRacketsByGang(ctx, gangID)
    for _, r := range rackets {
        report.RacketRevenueMonthly += r.EstimatedMonthlyUSD
    }

    report.TotalMonthlyEstimateUSD = report.TollRevenueMonthly +
        report.RacketRevenueMonthly +
        (report.RansomRevenue / 12) // annualise

    // Publier vers UCREF pour STR automatique si > 100k USD/mois
    if report.TotalMonthlyEstimateUSD > 100000 {
        _ = s.kafka.Publish(ctx, "extors.high.revenue.gang", report)
    }
    return report, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                              | Rôle          | Description                       |
|---------|---------------------------------------|---------------|-----------------------------------|
| `POST`  | `/api/v1/extors/cases`                | DCPJ, CAE     | Ouvrir dossier extorsion          |
| `GET`   | `/api/v1/extors/cases/:id`            | DCPJ, CAE     | Détail dossier                    |
| `POST`  | `/api/v1/extors/cases/:id/negotiations`| CAE_OFFICER  | Enregistrer négociation           |
| `POST`  | `/api/v1/extors/toll-points`          | DCPJ_INTEL    | Documenter péage illicite         |
| `GET`   | `/api/v1/extors/toll-points/map`      | DCPJ, BRI     | Carte GeoJSON péages actifs       |
| `GET`   | `/api/v1/extors/gang/:id/revenue`     | UCREF, DCPJ   | Revenu estimé d'un gang           |
| `GET`   | `/api/v1/extors/stats/by-type`        | DCPJ_ADMIN    | Statistiques par type             |
| `GET`   | `/api/v1/extors/moncash/patterns`     | UCREF         | Patterns MonCash suspects         |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
EXTORS_DB_HOST=localhost
EXTORS_DB_NAME=snisid_extors
EXTORS_GANG_SERVICE_URL=http://gang-svc:8095
EXTORS_UCREF_SERVICE_URL=http://ucref-svc:8112
EXTORS_KAFKA_BROKERS=kafka:9092
EXTORS_HIGH_REVENUE_THRESHOLD_USD=100000
EXTORS_SERVICE_PORT=8116
```

---
*MP-41 — EXTORS-HT — Extorsions et Rançons — SNISID — République d'Haïti*

---
---

# MP-42 — CRYPT-HT
## Surveillance des Cryptomonnaies à Usage Criminel
### Master Prompt d'Implémentation — SNISID

```
Classification   : SECRET / USAGE OFFICIEL EXCLUSIF
Module           : MP-42 | Code : CRYPT-HT
Dépendances      : UCREF-INT (MP-39), BLAN-HT (MP-40), EXTORS-HT (MP-41), SANC-HT (MP-27)
Normes           : GAFI/FATF Rec.15 (Virtual Assets), FinCEN Virtual Currency Guidance
Acteurs          : UCREF, BRH, Cellule cybercriminalité DCPJ, DEA Financial
```

---

## 1. CONTEXTE

L'usage des cryptomonnaies dans les crimes haïtiens est documenté :
- Réception de rançons en Bitcoin / USDT pour éviter la traçabilité MonCash
- Blanchiment via services de mixing (Tornado Cash, Wasabi Wallet)
- Paiements pour armes via marchés dark web
- Transferts entre membres de gangs déportés USA ↔ Haiti

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE crypt_asset_type AS ENUM (
    'BITCOIN', 'ETHEREUM', 'USDT', 'USDC',
    'MONERO', 'ZCASH', 'LITECOIN', 'OTHER_ERC20', 'UNKNOWN'
);

CREATE TYPE crypt_suspicion_type AS ENUM (
    'RANSOM_RECEIPT', 'SANCTIONS_EVASION', 'DARKWEB_PAYMENT',
    'MIXER_SERVICE', 'PEER_TO_PEER_UNREGULATED',
    'EXCHANGE_HIGH_RISK', 'GANG_PAYMENT', 'UNKNOWN'
);

CREATE TABLE crypt_flagged_wallets (
    wallet_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_crypt_id   VARCHAR(25) UNIQUE NOT NULL,  -- CRYPT-HT-NNNNNN
    wallet_address      VARCHAR(200) NOT NULL,
    asset_type          crypt_asset_type NOT NULL,
    blockchain_network  VARCHAR(50),                  -- Bitcoin, Ethereum, Tron, etc.
    suspicion_type      crypt_suspicion_type NOT NULL,
    snisid_person_id    UUID,
    gang_id             UUID,
    estimated_balance_usd DECIMAL(18,2),
    total_received_usd  DECIMAL(18,2),
    total_sent_usd      DECIMAL(18,2),
    first_tx_date       TIMESTAMPTZ,
    last_tx_date        TIMESTAMPTZ,
    is_sanctioned       BOOLEAN DEFAULT FALSE,
    ofac_sdn_ref        VARCHAR(50),
    chainalysis_ref     VARCHAR(100),                 -- Ref rapport Chainalysis
    elliptic_ref        VARCHAR(100),                 -- Ref rapport Elliptic
    source_intel        TEXT,
    linked_cases        UUID[] DEFAULT '{}',
    is_frozen           BOOLEAN DEFAULT FALSE,
    freeze_jurisdiction VARCHAR(50),
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE crypt_transactions (
    tx_id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id           UUID REFERENCES crypt_flagged_wallets(wallet_id),
    tx_hash             VARCHAR(100) NOT NULL,
    asset_type          crypt_asset_type NOT NULL,
    direction           VARCHAR(10) NOT NULL,          -- INCOMING, OUTGOING
    from_address        VARCHAR(200),
    to_address          VARCHAR(200),
    amount_crypto       DECIMAL(30,18),
    amount_usd_at_tx    DECIMAL(18,2),
    tx_timestamp        TIMESTAMPTZ NOT NULL,
    block_number        BIGINT,
    is_mixer_involved   BOOLEAN DEFAULT FALSE,
    mixer_service       VARCHAR(100),
    risk_score          SMALLINT,
    suspicion_flags     TEXT[] DEFAULT '{}',
    extors_case_id      UUID,                          -- Lien rançon si applicable
    ucref_str_id        UUID,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE crypt_exchange_accounts (
    exchange_id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id    UUID,
    exchange_name       VARCHAR(100) NOT NULL,         -- Binance, Coinbase, LocalBitcoins
    exchange_country    CHAR(3),
    account_ref         VARCHAR(200),                  -- Partiellement masque
    kyc_level           VARCHAR(20),                   -- NONE, BASIC, FULL
    total_volume_usd    DECIMAL(18,2),
    is_flagged          BOOLEAN DEFAULT FALSE,
    flagging_reason     TEXT,
    legal_hold_request  BOOLEAN DEFAULT FALSE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_crypt_wallets_address ON crypt_flagged_wallets(wallet_address);
CREATE INDEX idx_crypt_wallets_gang    ON crypt_flagged_wallets(gang_id) WHERE gang_id IS NOT NULL;
CREATE INDEX idx_crypt_wallets_sanctioned ON crypt_flagged_wallets(is_sanctioned) WHERE is_sanctioned = TRUE;
CREATE INDEX idx_crypt_tx_wallet       ON crypt_transactions(wallet_id, tx_timestamp DESC);
CREATE INDEX idx_crypt_tx_hash         ON crypt_transactions(tx_hash);
CREATE INDEX idx_crypt_tx_mixer        ON crypt_transactions(is_mixer_involved) WHERE is_mixer_involved = TRUE;

COMMIT;
```

---

## 3. SERVICE GO CLÉ — ANALYSE BLOCKCHAIN

```go
package service

import (
    "context"
    "github.com/snisid/crypt-svc/internal/domain"
)

// AnalyzeWalletRisk evalue le risque d un wallet via API Chainalysis
func (s *CryptService) AnalyzeWalletRisk(
    ctx context.Context,
    walletAddress, assetType string,
) (*domain.WalletRiskReport, error) {
    // 1. Verif contre wallets flagges locaux
    local, _ := s.repo.FindByAddress(ctx, walletAddress)
    if local != nil {
        return &domain.WalletRiskReport{
            WalletAddress: walletAddress,
            RiskScore:     95,
            IsKnownCriminal: true,
            LocalRecord:   local,
        }, nil
    }

    // 2. API Chainalysis KYT (Know Your Transaction)
    chainalysis, err := s.chainalysisClient.CheckWallet(ctx, walletAddress, assetType)
    if err != nil {
        return nil, err
    }

    // 3. Verif sanctions OFAC (wallets SDN)
    sancCheck, _ := s.sancClient.CheckCryptoWallet(ctx, walletAddress)

    report := &domain.WalletRiskReport{
        WalletAddress:       walletAddress,
        RiskScore:           chainalysis.RiskScore,
        ExposureToSanctions: sancCheck != nil && sancCheck.IsMatch,
        MixerExposure:       chainalysis.MixerExposure,
        DarkWebExposure:     chainalysis.DarkWebExposure,
        Source:              "CHAINALYSIS",
    }

    // Si tres haut risque -> creer automatiquement fiche flaggee
    if report.RiskScore >= 80 {
        _ = s.kafka.Publish(ctx, "crypt.high.risk.wallet.detected", report)
    }
    return report, nil
}
```

---

## 4. API REST

| Méthode | Endpoint                               | Rôle          | Description                      |
|---------|----------------------------------------|---------------|----------------------------------|
| `GET`   | `/api/v1/crypt/check/:address`         | UCREF, DCPJ   | Analyser risque d'un wallet      |
| `POST`  | `/api/v1/crypt/wallets`                | UCREF_ANALYST | Enregistrer wallet flagué        |
| `POST`  | `/api/v1/crypt/wallets/:id/transactions`| UCREF_ANALYST| Ajouter transactions suspectes   |
| `GET`   | `/api/v1/crypt/wallets/sanctioned`     | UCREF, MJSP   | Wallets sous sanctions OFAC      |
| `GET`   | `/api/v1/crypt/wallets/gang/:id`       | UCREF, DCPJ   | Wallets liés à un gang           |
| `GET`   | `/api/v1/crypt/stats/by-asset`         | UCREF_ADMIN   | Stats par type de crypto         |

---

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
CRYPT_DB_HOST=localhost
CRYPT_DB_NAME=snisid_crypt
CRYPT_CHAINALYSIS_API_URL=https://api.chainalysis.com/api/kyt
CRYPT_CHAINALYSIS_API_KEY=<VAULT:crypt/chainalysis_api_key>
CRYPT_ELLIPTIC_API_URL=https://aml.elliptic.co/v2
CRYPT_ELLIPTIC_API_KEY=<VAULT:crypt/elliptic_api_key>
CRYPT_SANC_SERVICE_URL=http://sanc-svc:8100
CRYPT_UCREF_SERVICE_URL=http://ucref-svc:8112
CRYPT_HIGH_RISK_THRESHOLD=80
CRYPT_SERVICE_PORT=8117
```

---
*MP-42 — CRYPT-HT — Cryptomonnaies Criminelles — SNISID — République d'Haïti*
