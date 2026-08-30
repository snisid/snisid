-- ============================================================
-- SNI-SIDE: National ALPR Database
-- CockroachDB (Geo-Partitioned)
-- ============================================================

CREATE SCHEMA IF NOT EXISTS snisid_alpr;
SET search_path TO snisid_alpr;

-- ============ ALPR CAMERAS ============
CREATE TABLE alpr_cameras (
    camera_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camera_code VARCHAR(50) UNIQUE NOT NULL,
    camera_type VARCHAR(30) CHECK (camera_type IN ('FIXED','MOBILE','PTZ','HANDHELD','DRONE')),
    location_name VARCHAR(255),
    latitude DECIMAL(10,7) NOT NULL,
    longitude DECIMAL(10,7) NOT NULL,
    address TEXT,
    department VARCHAR(100),
    municipality VARCHAR(100),
    direction VARCHAR(10) CHECK (direction IN ('NORTH','SOUTH','EAST','WEST','BOTH')),
    lane VARCHAR(10),
    speed_limit_kmh INT,
    operator_agency VARCHAR(100) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('ACTIVE','INACTIVE','MAINTENANCE','DECOMMISSIONED')),
    install_date DATE,
    firmware_version VARCHAR(50),
    ip_address INET,
    stream_url VARCHAR(500),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_camera_agency ON alpr_cameras(operator_agency);
CREATE INDEX idx_camera_status ON alpr_cameras(status);
CREATE INDEX idx_camera_location ON alpr_cameras(department, municipality);
CREATE INDEX idx_camera_type ON alpr_cameras(camera_type);

-- ============ ALPR READS ============
CREATE TABLE alpr_reads (
    read_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    camera_id UUID NOT NULL REFERENCES alpr_cameras(camera_id),
    plate_text VARCHAR(50) NOT NULL,
    plate_country VARCHAR(10),
    plate_type VARCHAR(20),
    read_timestamp TIMESTAMPTZ NOT NULL,
    latitude DECIMAL(10,7),
    longitude DECIMAL(10,7),
    speed_kmh INT,
    image_path_front VARCHAR(500),
    image_path_rear VARCHAR(500),
    ocr_confidence DECIMAL(5,2),
    vehicle_make VARCHAR(100),
    vehicle_model VARCHAR(100),
    vehicle_color VARCHAR(50),
    vehicle_year INT,
    vehicle_direction VARCHAR(10),
    vehicle_type VARCHAR(30),
    country VARCHAR(10) DEFAULT 'HT',
    created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY RANGE (read_timestamp);

CREATE INDEX idx_read_plate ON alpr_reads(plate_text);
CREATE INDEX idx_read_camera ON alpr_reads(camera_id);
CREATE INDEX idx_read_timestamp ON alpr_reads(read_timestamp DESC);
CREATE INDEX idx_read_country ON alpr_reads(plate_country);
CREATE INDEX idx_read_confidence ON alpr_reads(ocr_confidence);
CREATE INDEX idx_read_speed ON alpr_reads(speed_kmh);

-- ============ VEHICLE ROUTE HISTORY ============
CREATE TABLE vehicle_routes (
    route_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_text VARCHAR(50) NOT NULL,
    plate_country VARCHAR(10),
    read_sequence JSONB NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL,
    last_seen TIMESTAMPTZ NOT NULL,
    total_reads INT DEFAULT 1,
    distance_km DECIMAL(10,2),
    average_speed DECIMAL(5,2),
    route_pattern JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_route_plate ON vehicle_routes(plate_text);
CREATE INDEX idx_route_first ON vehicle_routes(first_seen DESC);
CREATE INDEX idx_route_reads ON vehicle_routes(total_reads DESC);

-- ============ CROSS BORDER TRACKING ============
CREATE TABLE border_crossings_alpr (
    crossing_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_text VARCHAR(50) NOT NULL,
    plate_country VARCHAR(10),
    exit_camera_id UUID REFERENCES alpr_cameras(camera_id),
    entry_camera_id UUID REFERENCES alpr_cameras(camera_id),
    exit_timestamp TIMESTAMPTZ,
    entry_timestamp TIMESTAMPTZ,
    crossing_duration_min INT,
    border_point VARCHAR(100),
    direction VARCHAR(10) CHECK (direction IN ('EXIT','ENTRY','TRANSIT')),
    vehicle_description JSONB,
    alert_triggered BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_border_plate ON border_crossings_alpr(plate_text);
CREATE INDEX idx_border_entry ON border_crossings_alpr(entry_timestamp DESC);
CREATE INDEX idx_border_point ON border_crossings_alpr(border_point);
CREATE INDEX idx_border_alert ON border_crossings_alpr(alert_triggered);

-- ============ ALPR HEATMAP DATA ============
CREATE TABLE alpr_heatmap_data (
    heatmap_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grid_cell GEOMETRY(POLYGON, 4326) NOT NULL,
    time_bucket TIMESTAMPTZ NOT NULL,
    vehicle_count INT NOT NULL,
    unique_plates INT,
    average_speed DECIMAL(5,2),
    peak_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_heatmap_cell ON alpr_heatmap_data USING gist(grid_cell);
CREATE INDEX idx_heatmap_time ON alpr_heatmap_data(time_bucket DESC);

-- ============ ALERTS ============
CREATE TABLE alpr_alerts (
    alert_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plate_text VARCHAR(50) NOT NULL,
    alert_type VARCHAR(30) CHECK (alert_type IN (
        'WANTED_PERSON','STOLEN_VEHICLE','WANTED_VEHICLE','AMBER_ALERT','SILVER_ALERT',
        'CROSS_BORDER','ANOMALOUS_MOVEMENT','CRIMINAL_CORRELATION'
    )),
    read_id UUID REFERENCES alpr_reads(read_id),
    camera_id UUID REFERENCES alpr_cameras(camera_id),
    detected_at TIMESTAMPTZ NOT NULL,
    risk_score DECIMAL(5,2),
    status VARCHAR(20) CHECK (status IN ('NEW','ACKNOWLEDGED','INVESTIGATING','RESPONDED','FALSE_POSITIVE','CLOSED')),
    acknowledged_by VARCHAR(100),
    acknowledged_at TIMESTAMPTZ,
    response_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
) PARTITION BY LIST (status);

CREATE INDEX idx_alpr_alert_plate ON alpr_alerts(plate_text);
CREATE INDEX idx_alpr_alert_type ON alpr_alerts(alert_type);
CREATE INDEX idx_alpr_alert_time ON alpr_alerts(detected_at DESC);
CREATE INDEX idx_alpr_alert_status ON alpr_alerts(status);

ALTER TABLE alpr_reads ENABLE ROW LEVEL SECURITY;
CREATE POLICY alpr_pnh_select ON alpr_reads FOR SELECT USING (
    current_setting('snisid.agency') IN ('PNH','DCPJ','IMMIGRATION','SNISID_ADMIN')
);
