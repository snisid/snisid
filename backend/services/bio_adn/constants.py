HAITI_DEPARTMENTS: list[str] = [
    "OUEST", "NORD", "NORD-EST", "NORD-OUEST",
    "ARTIBONITE", "CENTRE", "SUD", "SUD-EST",
    "GRAND-ANSE", "NIPPES",
]

HAITI_PORTS: list[str] = [
    "PORT-AU-PRINCE", "CAP-HAITIEN", "JACMEL",
    "LES-CAYES", "JEREMIE", "SAINT-MARC", "GONAIVES",
    "PORT-DE-PAIX", "FORT-LIBERTE", "MIRAGOANE",
]

ALERT_LEVELS: dict[str, str] = {
    "FULL_MATCH": "CRITICAL",
    "PARTIAL": "HIGH",
    "FAMILIAL": "MEDIUM",
    "VEHICLE_STOLEN": "HIGH",
    "PERSON_WANTED": "HIGH",
    "PERSON_WANTED_ARMED": "CRITICAL",
}

QUALITY_THRESHOLDS: dict[str, dict[str, float | int]] = {
    "BIO-CON": {"min_score": 0.95, "min_loci": 20},
    "BIO-ARR": {"min_score": 0.90, "min_loci": 18},
    "BIO-FSC": {"min_score": 0.60, "min_loci": 10},
    "BIO-DIS": {"min_score": 0.85, "min_loci": 15},
    "BIO-RNI": {"min_score": 0.50, "min_loci": 8},
}

RETENTION_DAYS: dict[str, int | None] = {
    "BIO-ARR": 3 * 365,
    "BIO-CON": None,
    "PER-REC": None,
    "BIE-VEH": 5 * 365,
}
