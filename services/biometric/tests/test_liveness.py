import pytest
from unittest.mock import patch, MagicMock


class TestLivenessEngine:
    @pytest.fixture
    def engine(self):
        from services.biometric.inference.liveness import LivenessEngine

        eng = LivenessEngine()
        eng.load_model()
        return eng

    def test_detect_returns_score_between_zero_and_one(self, engine):
        score = engine.detect(b"some_image_data")
        assert isinstance(score, float)
        assert 0.0 <= score <= 1.0

    def test_large_image_returns_high_score(self, engine):
        score = engine.detect(b"x" * 50_000)
        assert score >= 0.9

    def test_small_image_returns_low_score_spoof_detected(self, engine):
        score = engine.detect(b"small")
        assert score < 0.5

    def test_default_threshold_is_high(self, engine):
        assert engine.threshold == 0.995

    def test_detect_below_threshold_for_small_image(self, engine):
        score = engine.detect(b"tiny")
        assert score < engine.threshold

    def test_detect_raises_when_model_not_loaded(self):
        from services.biometric.inference.liveness import LivenessEngine

        engine = LivenessEngine()
        with pytest.raises(RuntimeError, match="Liveness Model not loaded"):
            engine.detect(b"test")

    def test_load_model_sets_flag(self):
        from services.biometric.inference.liveness import LivenessEngine

        engine = LivenessEngine()
        assert engine.model_loaded is False
        engine.load_model()
        assert engine.model_loaded is True

    def test_global_liveness_engine_singleton(self):
        from services.biometric.inference.liveness import (
            get_liveness_engine,
            init_liveness_model,
        )

        engine_before = get_liveness_engine()
        assert engine_before.model_loaded is False
        init_liveness_model()
        engine_after = get_liveness_engine()
        assert engine_after.model_loaded is True
        assert engine_before is engine_after

    def test_edge_case_exactly_boundary_size(self, engine):
        score_under = engine.detect(b"a" * 9_999)
        score_over = engine.detect(b"a" * 10_000)
        assert score_under < score_over

    def test_detect_is_deterministic(self, engine):
        data = b"deterministic_check"
        s1 = engine.detect(data)
        s2 = engine.detect(data)
        assert s1 == s2
