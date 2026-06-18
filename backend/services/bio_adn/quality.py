from __future__ import annotations

from services.bio_adn.constants import QUALITY_THRESHOLDS as INDEX_QUALITY_THRESHOLDS

CODIS_LOCI: list[str] = [
    "CSF1PO", "D3S1358", "D5S818", "D7S820", "D8S1179",
    "D13S317", "D16S539", "D18S51", "D21S11", "FGA",
    "TH01", "TPOX", "vWA", "D1S1656", "D2S441",
    "D2S1338", "D10S1248", "D12S391", "D19S433", "D22S1045",
]


def calculate_quality_score(electropherogram: dict[str, dict]) -> float:
    valid_loci = 0
    intensity_scores: list[float] = []
    noise_penalties = 0

    for locus in CODIS_LOCI:
        peak = electropherogram.get(locus)
        if peak is None:
            noise_penalties += 1
            continue

        height = peak.get("height", 0)
        if height < 150:
            noise_penalties += 1
            continue
        if height > 30000:
            noise_penalties += 0.5
            continue

        valid_loci += 1
        intensity_scores.append(min(max(height / 5000, 0.5), 1.0))

    total_loci = len(CODIS_LOCI)
    base_score = valid_loci / total_loci
    intensity_factor = sum(intensity_scores) / max(len(intensity_scores), 1)
    noise_factor = max(0, 1 - (noise_penalties * 0.05))

    return round(base_score * 0.6 + intensity_factor * 0.3 + noise_factor * 0.1, 3)


def validate_profile(index_type: str, quality_score: float, loci_count: int) -> list[str]:
    errors: list[str] = []
    thresholds = INDEX_QUALITY_THRESHOLDS.get(index_type)
    if thresholds is None:
        errors.append(f"Type d'index inconnu: {index_type}")
        return errors

    if quality_score < thresholds["min_score"]:
        errors.append(
            f"Quality score {quality_score} < seuil requis {thresholds['min_score']} pour {index_type}"
        )
    if loci_count < thresholds["min_loci"]:
        errors.append(
            f"Loci {loci_count}/{thresholds['min_loci']} insuffisants pour {index_type}"
        )
    return errors
