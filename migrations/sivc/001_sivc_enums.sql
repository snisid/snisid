BEGIN;

CREATE TYPE sivc_crime_category AS ENUM (
    'VOL_VEHICULE',
    'VOL_PLAQUE',
    'TRAFIC_STUPEFIANTS',
    'ENLEVEMENT',
    'GANG',
    'TRAFIC_ARMES',
    'VEHICULE_ETATIQUE_CLONE',
    'FAUX_POLICIER',
    'TRAITE_PERSONNES',
    'CONTREBANDE',
    'AUTRE_CRIME_GRAVE'
);

CREATE TYPE sivc_alert_level AS ENUM (
    'INFO',
    'CAUTION',
    'WANTED',
    'CRITICAL'
);

CREATE TYPE sivc_alert_status AS ENUM (
    'ACTIVE',
    'SUSPENDED',
    'RESOLVED',
    'EXPIRED',
    'CANCELLED'
);

CREATE TYPE sivc_plate_category AS ENUM (
    'PP',
    'PL',
    'M',
    'TC',
    'SE',
    'CD',
    'OA',
    'AG',
    'MD',
    'TX'
);

CREATE TYPE sivc_stolen_plate_status AS ENUM (
    'STOLEN',
    'RECOVERED',
    'DESTROYED',
    'USED_IN_CRIME'
);

CREATE TYPE sivc_sync_direction AS ENUM (
    'OUTBOUND',
    'INBOUND'
);

CREATE TYPE sivc_sync_status AS ENUM (
    'PENDING',
    'SUCCESS',
    'FAILED',
    'REJECTED'
);

CREATE TYPE sivc_route_type AS ENUM (
    'IMPORT',
    'EXPORT',
    'TRANSIT',
    'DOMESTIC'
);

CREATE TYPE sivc_kidnapping_status AS ENUM (
    'IN_PROGRESS',
    'VICTIM_RESCUED',
    'VICTIM_RELEASED',
    'VICTIM_DECEASED',
    'UNRESOLVED'
);

CREATE TYPE sivc_vehicle_type AS ENUM (
    'BERLINE',
    'SUV',
    'PICKUP',
    'CAMION',
    'MOTO',
    'TAP_TAP',
    'CAMIONNETTE',
    'BUS',
    'QUAD',
    'BATEAU',
    'AUTRE'
);

COMMIT;
