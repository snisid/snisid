-- ============================================================
-- SNI-SIDE: ClickHouse — National Analytics Platform
-- Tables materialisées pour l'analyse en temps réel
-- ============================================================

-- ============ 1. FLUX ALPR EN TEMPS RÉEL ============
CREATE TABLE sniside.alpr_stream
(
    read_id UUID,
    plate_text String,
    plate_country String DEFAULT 'HT',
    camera_code String,
    department String,
    latitude Float64,
    longitude Float64,
    read_timestamp DateTime64(3),
    speed_kmh Nullable(Int32),
    vehicle_make Nullable(String),
    vehicle_model Nullable(String),
    vehicle_color Nullable(String),
    ocr_confidence Float32,
    is_wanted UInt8 DEFAULT 0,
    ingestion_time DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(read_timestamp)
ORDER BY (read_timestamp, plate_text, camera_code)
TTL read_timestamp + INTERVAL 1 YEAR DELETE
SETTINGS index_granularity = 8192;

-- ============ 2. ALERTS EN TEMPS RÉEL ============
CREATE TABLE sniside.alerts_stream
(
    alert_id UUID,
    source LowCardinality(String),
    alert_type String,
    severity LowCardinality(String) DEFAULT 'MEDIUM',
    title String,
    description String,
    entity_ids Array(String),
    risk_score Float32 DEFAULT 0.0,
    status LowCardinality(String) DEFAULT 'NEW',
    created_at DateTime64(3),
    acknowledged_at Nullable(DateTime64(3)),
    resolved_at Nullable(DateTime64(3)),
    response_time_sec Nullable(Int32)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, severity, source)
TTL created_at + INTERVAL 2 YEAR DELETE;

-- ============ 3. TRANSACTIONS FINANCIÈRES SUSPECTES ============
CREATE TABLE sniside.suspicious_transactions
(
    transaction_id UUID,
    transaction_ref String,
    transaction_date DateTime64(3),
    transaction_type LowCardinality(String),
    amount Float64,
    currency LowCardinality(String),
    amount_usd Float64,
    source_country LowCardinality(String),
    destination_country LowCardinality(String),
    sender_niu Nullable(String),
    beneficiary_niu Nullable(String),
    risk_score Float32 DEFAULT 0.0,
    mlro_filed UInt8 DEFAULT 0,
    status LowCardinality(String) DEFAULT 'PENDING',
    ingestion_time DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(transaction_date)
ORDER BY (transaction_date, risk_score, source_country, destination_country)
TTL transaction_date + INTERVAL 5 YEAR DELETE
SETTINGS index_granularity = 8192;

-- ============ 4. BORDER CROSSINGS ANALYTICS ============
CREATE TABLE sniside.border_crossings_analytics
(
    crossing_id UUID,
    nationality LowCardinality(String),
    border_point LowCardinality(String),
    crossing_direction LowCardinality(String) DEFAULT 'ENTRY',
    crossing_method LowCardinality(String),
    crossing_date DateTime64(3),
    passport_country LowCardinality(String),
    visa_type Nullable(String),
    risk_score Float32 DEFAULT 0.0,
    alert_triggered UInt8 DEFAULT 0,
    age_years Nullable(Int32),
    gender LowCardinality(String),
    ingestion_time DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(crossing_date)
ORDER BY (crossing_date, border_point, crossing_direction, nationality)
TTL crossing_date + INTERVAL 5 YEAR DELETE;

-- ============ 5. CYBER IOC STREAM ============
CREATE TABLE sniside.cyber_ioc_stream
(
    ioc_id UUID,
    ioc_value String,
    ioc_type LowCardinality(String),
    confidence UInt8 DEFAULT 50,
    severity LowCardinality(String) DEFAULT 'MEDIUM',
    tlp_level LowCardinality(String) DEFAULT 'AMBER',
    threat_actor Nullable(String),
    malware_family Nullable(String),
    first_seen DateTime,
    last_seen DateTime,
    tags Array(String),
    status LowCardinality(String) DEFAULT 'ACTIVE',
    ingestion_time DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(last_seen)
ORDER BY (last_seen, severity, ioc_type, ioc_value);

-- ============ 6. RECHERCHES UNIFIÉES LOG ============
CREATE TABLE sniside.search_log
(
    search_id UUID,
    query String,
    query_type LowCardinality(String),
    query_fingerprint String,
    databases_searched UInt8 DEFAULT 0,
    total_results UInt32 DEFAULT 0,
    search_duration_ms Float64,
    user_id String DEFAULT 'anonymous',
    agency LowCardinality(String),
    user_clearance LowCardinality(String),
    graph_context_used UInt8 DEFAULT 0,
    created_at DateTime DEFAULT now()
)
ENGINE = MergeTree()
ORDER BY (created_at, agency, query_type)
TTL created_at + INTERVAL 1 YEAR DELETE;

-- ============ 7. CRIME STATISTIQUES AGRÉGÉES ============
CREATE TABLE sniside.crime_statistics
(
    date Date,
    department LowCardinality(String),
    crime_type LowCardinality(String),
    total_cases UInt32 DEFAULT 0,
    solved_cases UInt32 DEFAULT 0,
    open_cases UInt32 DEFAULT 0,
    persons_wanted UInt32 DEFAULT 0,
    warrants_issued UInt32 DEFAULT 0,
    warrants_executed UInt32 DEFAULT 0,
    avg_resolution_days Float64 DEFAULT 0.0,
    ingestion_time DateTime DEFAULT now()
)
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(date)
ORDER BY (date, department, crime_type)
TTL date + INTERVAL 10 YEAR DELETE;

-- ============ 8. ALPR VÉHICULES RECHERCHÉS ============
CREATE TABLE sniside.alpr_hits
(
    read_id UUID,
    plate_text String,
    plate_country String,
    camera_code String,
    department String,
    read_timestamp DateTime64(3),
    alert_type LowCardinality(String),
    response_time_sec Nullable(Int32),
    status LowCardinality(String) DEFAULT 'NEW',
    ingestion_time DateTime DEFAULT now()
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(read_timestamp)
ORDER BY (read_timestamp, department, alert_type)
TTL read_timestamp + INTERVAL 2 YEAR DELETE;

-- ============ 9. PERFORMANCE DES MODÈLES IA ============
CREATE TABLE sniside.ai_model_performance
(
    model_name LowCardinality(String),
    model_version String,
    inference_count UInt64 DEFAULT 0,
    avg_latency_ms Float64 DEFAULT 0.0,
    p99_latency_ms Float64 DEFAULT 0.0,
    error_count UInt64 DEFAULT 0,
    error_rate Float32 DEFAULT 0.0,
    gpu_utilization Float32 DEFAULT 0.0,
    memory_usage_mb Float64 DEFAULT 0.0,
    batch_size UInt32 DEFAULT 1,
    observation_time DateTime DEFAULT now()
)
ENGINE = SummingMergeTree()
ORDER BY (observation_time, model_name, model_version);

-- ============ 10. GRAPH INTELLIGENCE EVOLUTION ============
CREATE TABLE sniside.graph_evolution
(
    snapshot_date Date,
    node_count UInt64 DEFAULT 0,
    edge_count UInt64 DEFAULT 0,
    citizen_nodes UInt64 DEFAULT 0,
    vehicle_nodes UInt64 DEFAULT 0,
    phone_nodes UInt64 DEFAULT 0,
    financial_nodes UInt64 DEFAULT 0,
    cyber_nodes UInt64 DEFAULT 0,
    watchlist_nodes UInt64 DEFAULT 0,
    detected_networks UInt32 DEFAULT 0,
    network_density Float32 DEFAULT 0.0,
    centrality_avg Float32 DEFAULT 0.0,
    ingestion_time DateTime DEFAULT now()
)
ENGINE = SummingMergeTree()
ORDER BY (snapshot_date);

-- ============ 11. RECHERCHE TEXTUELLE AVEC FONCTIONS ============
-- ALPR Heatmap par heure
CREATE MATERIALIZED VIEW sniside.alpr_hourly_heatmap
ENGINE = SummingMergeTree()
ORDER BY (hour, department, camera_code)
AS SELECT
    toStartOfHour(read_timestamp) AS hour,
    department,
    camera_code,
    count() AS total_reads,
    countDistinct(plate_text) AS unique_plates,
    avg(speed_kmh) AS avg_speed,
    max(speed_kmh) AS max_speed,
    countIf(is_wanted = 1) AS wanted_hits
FROM sniside.alpr_stream
GROUP BY hour, department, camera_code;

-- Border crossing statistiques journalières
CREATE MATERIALIZED VIEW sniside.border_daily_stats
ENGINE = SummingMergeTree()
ORDER BY (date, border_point)
AS SELECT
    toDate(crossing_date) AS date,
    border_point,
    countIf(crossing_direction = 'ENTRY') AS entries,
    countIf(crossing_direction = 'EXIT') AS exits,
    countDistinct(nationality) AS nationalities,
    avg(risk_score) AS avg_risk,
    countIf(alert_triggered = 1) AS alerts,
    countIf(visa_type = 'TOURIST') AS tourists,
    countIf(visa_type = 'WORK') AS workers
FROM sniside.border_crossings_analytics
GROUP BY date, border_point;

-- Top wanted vehicles (ALPR)
CREATE MATERIALIZED VIEW sniside.top_wanted_vehicles
ENGINE = SummingMergeTree()
ORDER BY (plate_text, hits DESC)
AS SELECT
    plate_text,
    count() AS hits,
    countDistinct(camera_code) AS cameras_seen,
    min(read_timestamp) AS first_seen,
    max(read_timestamp) AS last_seen,
    countDistinct(department) AS departments_seen
FROM sniside.alpr_stream
WHERE is_wanted = 1
GROUP BY plate_text;

-- Alert response time analytics
CREATE MATERIALIZED VIEW sniside.alert_response_analytics
ENGINE = SummingMergeTree()
ORDER BY (date, source, severity)
AS SELECT
    toDate(created_at) AS date,
    source,
    severity,
    count() AS total_alerts,
    avg(response_time_sec) AS avg_response_sec,
    max(response_time_sec) AS max_response_sec,
    countIf(status IN ('RESOLVED', 'ACKNOWLEDGED')) AS handled,
    countIf(status = 'NEW') AS pending
FROM sniside.alerts_stream
GROUP BY date, source, severity;

-- Cross-database entity correlation
CREATE TABLE sniside.entity_correlations
(
    entity_id String,
    entity_type LowCardinality(String),
    databases_present Array(LowCardinality(String)),
    total_references UInt32 DEFAULT 0,
    avg_risk_score Float32 DEFAULT 0.0,
    max_risk_score Float32 DEFAULT 0.0,
    last_seen DateTime,
    created_at DateTime DEFAULT now()
)
ENGINE = ReplacingMergeTree(last_seen)
ORDER BY (entity_id, entity_type);

-- ============ DISTRIBUTED TABLES (across shards) ============
CREATE TABLE sniside.alpr_stream_distributed AS sniside.alpr_stream
ENGINE = Distributed('sniside', 'sniside', 'alpr_stream', rand());

CREATE TABLE sniside.suspicious_transactions_distributed AS sniside.suspicious_transactions
ENGINE = Distributed('sniside', 'sniside', 'suspicious_transactions', cityHash64(transaction_ref));
