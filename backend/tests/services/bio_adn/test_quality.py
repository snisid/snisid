import pytest

from services.bio_adn.quality import (
    CODIS_LOCI,
    INDEX_QUALITY_THRESHOLDS,
    calculate_quality_score,
    validate_profile,
)


class TestCalculateQualityScore:

    def test_perfect_electropherogram(self):
        epg = {locus: {"height": 5000} for locus in CODIS_LOCI}
        score = calculate_quality_score(epg)
        assert score >= 0.95

    def test_all_loci_present_min_height(self):
        epg = {locus: {"height": 150} for locus in CODIS_LOCI}
        score = calculate_quality_score(epg)
        assert score > 0.80

    def test_degraded_sample(self):
        epg = {}
        for i, locus in enumerate(CODIS_LOCI):
            epg[locus] = {"height": 5000 if i < 8 else 50}
        score = calculate_quality_score(epg)
        assert 0.40 <= score < 0.95

    def test_all_loci_missing(self):
        score = calculate_quality_score({})
        assert score < 0.40

    def test_saturated_peak_penalty(self):
        epg = {locus: {"height": 35000} for locus in CODIS_LOCI}
        score = calculate_quality_score(epg)
        assert score < 0.95

    def test_mixed_quality(self):
        epg = {}
        for i, locus in enumerate(CODIS_LOCI):
            if i < 10:
                epg[locus] = {"height": 3000}
            elif i < 15:
                epg[locus] = {"height": 120}
            else:
                epg[locus] = {"height": 0}
        score = calculate_quality_score(epg)
        assert 0.30 <= score <= 0.70

    def test_fluctuating_heights(self):
        epg = {}
        for i, locus in enumerate(CODIS_LOCI):
            heights = [200, 500, 1000, 3000, 8000, 15000]
            epg[locus] = {"height": heights[i % len(heights)]}
        score = calculate_quality_score(epg)
        assert 0.50 <= score <= 1.0


class TestValidateProfile:

    def test_bio_con_valid(self):
        errors = validate_profile("BIO-CON", 0.98, 20)
        assert errors == []

    def test_bio_con_below_threshold(self):
        errors = validate_profile("BIO-CON", 0.90, 20)
        assert len(errors) == 1
        assert "seuil" in errors[0]

    def test_bio_con_insufficient_loci(self):
        errors = validate_profile("BIO-CON", 0.98, 15)
        assert len(errors) == 1
        assert "Loci" in errors[0]

    def test_bio_con_both_failures(self):
        errors = validate_profile("BIO-CON", 0.50, 8)
        assert len(errors) == 2

    def test_bio_fsc_low_threshold(self):
        errors = validate_profile("BIO-FSC", 0.60, 10)
        assert errors == []

    def test_bio_fsc_below_threshold(self):
        errors = validate_profile("BIO-FSC", 0.50, 8)
        assert len(errors) == 2

    def test_bio_rni_minimum_acceptable(self):
        errors = validate_profile("BIO-RNI", 0.50, 8)
        assert errors == []

    def test_bio_arr_valid(self):
        errors = validate_profile("BIO-ARR", 0.90, 18)
        assert errors == []

    def test_bio_dis_valid(self):
        errors = validate_profile("BIO-DIS", 0.85, 15)
        assert errors == []

    def test_unknown_index_type(self):
        errors = validate_profile("BIO-UNKNOWN", 0.95, 20)
        assert len(errors) == 1
        assert "inconnu" in errors[0]
