# MP-50 — RVIN-HT
## Registre National des Victimes Non Identifiées
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-50 | Code : RVIN-HT
Dépendances      : SNISID-BIO-ADN, AFIS-HT (MP-19), DIPE-HT (MP-43), SIGDC-HT (MP-49)
Normes           : INTERPOL Disaster Victim Identification (DVI), PCAST DNA Standards
Acteurs          : MSPP (Médecine Légale), PNH, Parquet, Croix-Rouge, MSF
```

---

## 1. CONTEXTE

En Haïti, des milliers de corps non identifiés s'accumulent dans les morgues et les
charniers après les événements de masse (séismes, massacres, inondations). Ce module
crée le registre forensique permettant l'identification via ADN, AFIS et odontologie.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE rvin_source AS ENUM (
    'CRIME_SCENE', 'DISASTER_SITE', 'MASS_GRAVE',
    'RIVER', 'STREET', 'HOSPITAL_DOA', 'OTHER'
);

CREATE TYPE rvin_status AS ENUM (
    'UNIDENTIFIED', 'TENTATIVE_MATCH', 'CONFIRMED_IDENTIFIED', 'CLAIMED'
);

CREATE TABLE rvin_unidentified_remains (
    remains_id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_rvin_id    VARCHAR(25) UNIQUE NOT NULL,  -- RVIN-HT-AAAA-NNNNNN
    discovery_date      TIMESTAMPTZ NOT NULL,
    discovery_location  VARCHAR(300) NOT NULL,
    dept_code           CHAR(2) NOT NULL,
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    discovery_source    rvin_source NOT NULL,
    status              rvin_status NOT NULL DEFAULT 'UNIDENTIFIED',

    -- Description physique
    estimated_sex       VARCHAR(10),
    estimated_age_min   SMALLINT,
    estimated_age_max   SMALLINT,
    estimated_height_cm SMALLINT,
    skin_tone           VARCHAR(30),
    hair_type           VARCHAR(50),
    clothing_description TEXT,
    distinguishing_marks TEXT,
    decomposition_level SMALLINT CHECK (decomposition_level BETWEEN 1 AND 5),

    -- Biometrie forensique
    afis_latent_id      UUID,
    dna_sample_taken    BOOLEAN DEFAULT FALSE,
    dna_sample_ref      VARCHAR(100),
    dna_profile_id      UUID,
    dental_chart_ref    VARCHAR(200),
    photo_refs          TEXT[] DEFAULT '{}',
    xray_refs           TEXT[] DEFAULT '{}',

    -- Morgue / conservation
    morgue_location     VARCHAR(200),
    morgue_ref          VARCHAR(50),
    storage_date        TIMESTAMPTZ,
    estimated_death_date TIMESTAMPTZ,

    -- Liens contextuels
    disaster_id         UUID,
    mass_incident_id    UUID,
    gang_id             UUID,
    case_reference      VARCHAR(100),
    interpol_dvi_ref    VARCHAR(50),

    -- Identification (si resolue)
    matched_dipe_case_id UUID,
    matched_snisid_id    UUID,
    identification_method VARCHAR(50),
    identification_date  TIMESTAMPTZ,
    identified_by        UUID,

    examiner_id         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rvin_dna_comparisons (
    comparison_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    remains_id          UUID NOT NULL REFERENCES rvin_unidentified_remains(remains_id),
    dipe_case_id        UUID,
    reference_dna_ref   VARCHAR(100),
    comparison_date     TIMESTAMPTZ NOT NULL,
    match_probability   DECIMAL(10,8),
    is_match            BOOLEAN DEFAULT FALSE,
    lab_reference       VARCHAR(100),
    examiner_id         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rvin_status    ON rvin_unidentified_remains(status);
CREATE INDEX idx_rvin_dept      ON rvin_unidentified_remains(dept_code, discovery_date DESC);
CREATE INDEX idx_rvin_disaster  ON rvin_unidentified_remains(disaster_id) WHERE disaster_id IS NOT NULL;

COMMIT;
```

---

## 3. API REST

| Méthode | Endpoint                           | Rôle          | Description                      |
|---------|------------------------------------|---------------|----------------------------------|
| `POST`  | `/api/v1/rvin/remains`             | MSPP, PNH     | Enregistrer restes non identifiés|
| `GET`   | `/api/v1/rvin/remains/:id`         | MSPP, PARQUET | Détail fiche                     |
| `POST`  | `/api/v1/rvin/remains/:id/dna`     | LAB_FORENSIC  | Soumettre résultat ADN           |
| `GET`   | `/api/v1/rvin/match/dipe/:case_id` | PNH, MSPP     | Chercher correspondance DIPE-HT  |
| `GET`   | `/api/v1/rvin/unidentified`        | MSPP           | Tous les restes non identifiés   |
| `GET`   | `/api/v1/rvin/stats/by-source`     | MSPP_ADMIN    | Stats par source découverte      |

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
RVIN_DB_HOST=localhost
RVIN_DB_NAME=snisid_rvin
RVIN_AFIS_SERVICE_URL=http://afis-svc:8091
RVIN_BIO_ADN_SERVICE_URL=http://bio-adn-svc:8080
RVIN_DIPE_SERVICE_URL=http://dipe-svc:8118
RVIN_INTERPOL_DVI_URL=https://i247-gateway.pnh.gov.ht/dvi
RVIN_SERVICE_PORT=8120
```

---
*MP-50 — RVIN-HT — Victimes Non Identifiées — SNISID — République d'Haïti*

---
---

# MP-51 — MVSM-HT
## Surveillance des Rassemblements de Masse et Mobilisations Populaires
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-51 | Code : MVSM-HT
Dépendances      : SIGEO-HT (MP-48), TERR-HT (MP-26), SIVC-HT (MP-18), GANG-HT (MP-24)
Normes           : Principes ONU usage force rassemblements, OSCE Guidelines
Acteurs          : PNH BRI, GIPNH, MSP, Préfectures, Mairies
```

---

## 1. CONTEXTE

Haïti connaît régulièrement des manifestations, barrages routiers (peyi lòk) et
rassemblements qui peuvent dégénérer en violence. Ce module permet la planification
opérationnelle, le suivi en temps réel et l'évaluation post-événement des rassemblements.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE mvsm_event_type AS ENUM (
    'POLITICAL_PROTEST', 'LABOR_STRIKE', 'COMMUNITY_ACTION',
    'RELIGIOUS_GATHERING', 'CULTURAL_EVENT', 'PEYI_LOK_BARRICADE',
    'GANG_MOBILIZATION', 'SPONTANEOUS_UNREST', 'OTHER'
);

CREATE TYPE mvsm_risk_level AS ENUM (
    'LOW', 'MODERATE', 'HIGH', 'CRITICAL'
);

CREATE TABLE mvsm_events (
    event_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_mvsm_id    VARCHAR(25) UNIQUE NOT NULL,
    event_type          mvsm_event_type NOT NULL,
    event_name          VARCHAR(200),
    risk_level          mvsm_risk_level NOT NULL DEFAULT 'LOW',
    status              VARCHAR(20) DEFAULT 'PLANNED',
    organizer_name      VARCHAR(200),
    organizer_snisid_id UUID,
    gang_id             UUID,                -- Si mobilisation gang
    scheduled_date      TIMESTAMPTZ NOT NULL,
    actual_start        TIMESTAMPTZ,
    actual_end          TIMESTAMPTZ,
    location_desc       VARCHAR(300),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7),
    estimated_crowd     INTEGER,
    peak_crowd          INTEGER,
    deployed_units      TEXT[] DEFAULT '{}',
    incidents_during    INTEGER DEFAULT 0,
    casualties          INTEGER DEFAULT 0,
    arrests_made        INTEGER DEFAULT 0,
    weapons_found       INTEGER DEFAULT 0,
    vehicles_involved   INTEGER DEFAULT 0,
    sivc_alert_ids      UUID[] DEFAULT '{}',
    post_event_notes    TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE mvsm_real_time_updates (
    update_id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id            UUID NOT NULL REFERENCES mvsm_events(event_id),
    update_time         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_crowd_est   INTEGER,
    situation           TEXT NOT NULL,
    risk_change         mvsm_risk_level,
    action_taken        TEXT,
    reported_by         UUID NOT NULL,
    lat                 DECIMAL(10,7),
    lng                 DECIMAL(10,7)
);

CREATE INDEX idx_mvsm_events_date   ON mvsm_events(scheduled_date DESC);
CREATE INDEX idx_mvsm_events_dept   ON mvsm_events(dept_code, scheduled_date DESC);
CREATE INDEX idx_mvsm_events_risk   ON mvsm_events(risk_level) WHERE status IN ('PLANNED','ACTIVE');
CREATE INDEX idx_mvsm_updates_event ON mvsm_real_time_updates(event_id, update_time DESC);

COMMIT;
```

## 3. API REST

| Méthode | Endpoint                               | Rôle         | Description                   |
|---------|----------------------------------------|--------------|-------------------------------|
| `POST`  | `/api/v1/mvsm/events`                  | PNH_ADMIN    | Créer événement rassemblement |
| `GET`   | `/api/v1/mvsm/events/upcoming`         | PNH, BRI     | Prochains événements          |
| `POST`  | `/api/v1/mvsm/events/:id/updates`      | PNH_OFFICER  | Mise à jour temps réel        |
| `GET`   | `/api/v1/mvsm/events/active`           | PNH          | Événements en cours           |
| `GET`   | `/api/v1/mvsm/events/:id/debrief`      | BRI_CMD      | Rapport post-événement        |
| `PATCH` | `/api/v1/mvsm/events/:id/risk`         | BRI_SUPERVISOR| Changer niveau de risque      |

## 4. VARIABLES D'ENVIRONNEMENT

```dotenv
MVSM_DB_HOST=localhost
MVSM_DB_NAME=snisid_mvsm
MVSM_SIVC_SERVICE_URL=http://sivc-svc:8090
MVSM_SIGEO_SERVICE_URL=http://sigeo-svc:8125
MVSM_KAFKA_BROKERS=kafka:9092
MVSM_SERVICE_PORT=8127
```

---
*MP-51 — MVSM-HT — Surveillance Rassemblements — SNISID — République d'Haïti*

---
---

# MP-52 — SISAL-HT
## Système National d'Alerte Précoce Multi-Risques
### Master Prompt d'Implémentation — SNISID

```
Classification   : RESTREINT / USAGE OFFICIEL
Module           : MP-52 | Code : SISAL-HT
Dépendances      : SIGDC-HT (MP-49), SIGEO-HT (MP-48), MVSM-HT (MP-51), GANG-HT (MP-24)
Normes           : Cadre Sendai, OCHA Early Warning, WMO Multi-Hazard Early Warning
Acteurs          : CSPAN, SNGRD, MSP, BRH (météo), UHM (séismes), OCHA
```

---

## 1. CONTEXTE

SISAL-HT est le système d'orchestration des alertes multi-risques : il fusionne les
signaux des systèmes de surveillance sismique (UHM), météorologique (BRH), sécuritaire
(SNISID) et humanitaire (OCHA) pour produire des alertes consolidées diffusées
via SMS, radio, application mobile et sirènes.

---

## 2. BASE DE DONNÉES

```sql
BEGIN;

CREATE TYPE sisal_hazard_type AS ENUM (
    'EARTHQUAKE', 'HURRICANE', 'FLOOD', 'TSUNAMI', 'LANDSLIDE',
    'SECURITY_GANG', 'SECURITY_MASS_CASUALTY', 'EPIDEMIC',
    'INDUSTRIAL', 'COMPOSITE'
);

CREATE TYPE sisal_severity AS ENUM (
    'ADVISORY', 'WATCH', 'WARNING', 'EMERGENCY', 'CATASTROPHE'
);

CREATE TYPE sisal_channel AS ENUM (
    'SMS_MASS', 'PUSH_NOTIFICATION', 'RADIO_BROADCAST',
    'SIRENE', 'SOCIAL_MEDIA', 'OFFICIAL_AGENCIES', 'LOUD_SPEAKER'
);

CREATE TABLE sisal_alerts (
    alert_id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    national_sisal_id   VARCHAR(25) UNIQUE NOT NULL,  -- SISAL-HT-AAAA-NNNNNN
    hazard_type         sisal_hazard_type NOT NULL,
    severity            sisal_severity NOT NULL,
    title               VARCHAR(200) NOT NULL,
    message_fr          TEXT NOT NULL,              -- Message en français
    message_ht          TEXT NOT NULL,              -- Message en créole haïtien
    affected_depts      CHAR(2)[] DEFAULT '{}',
    affected_communes   TEXT[] DEFAULT '{}',
    affected_pop_est    INTEGER,
    issued_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until         TIMESTAMPTZ,
    source_agency       VARCHAR(100) NOT NULL,
    source_event_id     UUID,                       -- Lien SIGDC ou autre module
    recommended_actions TEXT[],
    channels_used       sisal_channel[] DEFAULT '{}',
    sms_count_sent      INTEGER DEFAULT 0,
    push_count_sent     INTEGER DEFAULT 0,
    is_cancelled        BOOLEAN DEFAULT FALSE,
    cancelled_at        TIMESTAMPTZ,
    cancel_reason       TEXT,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sisal_data_feeds (
    feed_id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    feed_name           VARCHAR(100) NOT NULL,      -- UHM_SEISMIC, BRH_WEATHER, USGS, NHC
    feed_url            VARCHAR(500),
    feed_type           VARCHAR(30),                -- REST_API, RSS, WEBSOCKET
    hazard_types        sisal_hazard_type[] DEFAULT '{}',
    polling_interval_sec INTEGER DEFAULT 60,
    is_active           BOOLEAN DEFAULT TRUE,
    last_poll           TIMESTAMPTZ,
    last_alert_generated TIMESTAMPTZ,
    alert_thresholds    JSONB,                      -- {"magnitude": 5.0, "wind_speed": 120}
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE sisal_subscriptions (
    sub_id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    snisid_person_id    UUID,
    phone_number        VARCHAR(30),
    email               VARCHAR(200),
    dept_code           CHAR(2),
    commune             VARCHAR(100),
    hazard_types        sisal_hazard_type[] DEFAULT '{}',
    min_severity        sisal_severity DEFAULT 'WARNING',
    channels            sisal_channel[] DEFAULT ARRAY['SMS_MASS'::sisal_channel],
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sisal_alerts_dept     ON sisal_alerts USING gin(affected_depts);
CREATE INDEX idx_sisal_alerts_severity ON sisal_alerts(severity, issued_at DESC);
CREATE INDEX idx_sisal_alerts_hazard   ON sisal_alerts(hazard_type, issued_at DESC);
CREATE INDEX idx_sisal_subs_dept       ON sisal_subscriptions(dept_code) WHERE is_active = TRUE;

COMMIT;
```

---

## 3. SERVICE GO — ORCHESTRATEUR D'ALERTES

```go
package service

import (
    "context"
    "sync"
    "time"
    "github.com/snisid/sisal-svc/internal/domain"
)

// BroadcastAlert diffuse une alerte sur tous les canaux en parallele
func (s *AlertService) BroadcastAlert(
    ctx context.Context,
    alert *domain.SISALAlert,
) *domain.BroadcastResult {
    result := &domain.BroadcastResult{AlertID: alert.AlertID}
    var wg sync.WaitGroup

    channels := []struct {
        name    string
        handler func(context.Context, *domain.SISALAlert) (int, error)
    }{
        {"SMS", s.smsGateway.SendMassSMS},
        {"PUSH", s.pushService.BroadcastPush},
        {"RADIO", s.radioClient.Broadcast},
    }

    for _, ch := range channels {
        wg.Add(1)
        go func(c struct {
            name    string
            handler func(context.Context, *domain.SISALAlert) (int, error)
        }) {
            defer wg.Done()
            sent, err := c.handler(ctx, alert)
            if err == nil {
                result.AddChannel(c.name, sent)
            }
        }(ch)
    }

    wg.Wait()
    result.CompletedAt = time.Now()

    // Publier sur Kafka pour historisation
    _ = s.kafka.Publish(ctx, "sisal.alert.broadcast.completed", result)
    return result
}
```

---

## 4. API REST

| Méthode | Endpoint                          | Rôle         | Description                     |
|---------|-----------------------------------|--------------|---------------------------------|
| `POST`  | `/api/v1/sisal/alerts`            | CSPAN, SNGRD | Émettre alerte multi-risques    |
| `GET`   | `/api/v1/sisal/alerts/active`     | PUBLIC       | Alertes actives en cours        |
| `GET`   | `/api/v1/sisal/alerts/history`    | CSPAN_ADMIN  | Historique des alertes          |
| `POST`  | `/api/v1/sisal/alerts/:id/cancel` | CSPAN        | Annuler une alerte              |
| `POST`  | `/api/v1/sisal/subscribe`         | PUBLIC       | S'abonner aux alertes SMS/Push  |
| `GET`   | `/api/v1/sisal/feeds/status`      | CSPAN_ADMIN  | Statut des flux de données      |

## 5. VARIABLES D'ENVIRONNEMENT

```dotenv
SISAL_DB_HOST=localhost
SISAL_DB_NAME=snisid_sisal
SISAL_SMS_GATEWAY_URL=http://sms-gw:8080
SISAL_PUSH_FCM_URL=https://fcm.googleapis.com/fcm/send
SISAL_PUSH_FCM_KEY=<VAULT:sisal/fcm_key>
SISAL_RADIO_API_URL=http://radio-broadcast:8081
SISAL_UHM_FEED_URL=https://uhm.gov.ht/api/seismic
SISAL_BRH_WEATHER_URL=https://brh.gouv.ht/api/weather
SISAL_USGS_FEED_URL=https://earthquake.usgs.gov/fdsnws/event/1
SISAL_NHC_FEED_URL=https://www.nhc.noaa.gov/CurrentStorms.json
SISAL_MIN_EARTHQUAKE_MAG=4.5
SISAL_KAFKA_BROKERS=kafka:9092
SISAL_SERVICE_PORT=8128
```

---
*MP-52 — SISAL-HT — Alerte Précoce Multi-Risques — SNISID — République d'Haïti*
