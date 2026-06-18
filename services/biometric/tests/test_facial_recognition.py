import pytest
import numpy as np
from unittest.mock import patch, MagicMock


class TestFacialInferenceEngine:
    @pytest.fixture
    def engine(self):
        from services.biometric.inference.facial import FacialInferenceEngine

        eng = FacialInferenceEngine()
        eng.load_model()
        return eng

    def test_extract_embedding_output_shape(self, engine):
        image_bytes = b"fake_image_data_12345"
        embedding = engine.extract_embedding(image_bytes)
        assert isinstance(embedding, np.ndarray)
        assert embedding.shape == (512,)
        assert embedding.dtype == np.float32

    def test_extract_embedding_is_l2_normalized(self, engine):
        image_bytes = b"test_normalization"
        embedding = engine.extract_embedding(image_bytes)
        norm = np.linalg.norm(embedding)
        assert abs(norm - 1.0) < 1e-5

    def test_extract_embedding_deterministic_for_same_input(self, engine):
        image_bytes = b"deterministic_test"
        emb1 = engine.extract_embedding(image_bytes)
        emb2 = engine.extract_embedding(image_bytes)
        np.testing.assert_array_equal(emb1, emb2)

    def test_extract_embedding_raises_when_model_not_loaded(self):
        from services.biometric.inference.facial import FacialInferenceEngine

        engine = FacialInferenceEngine()
        with pytest.raises(RuntimeError, match="Model not loaded"):
            engine.extract_embedding(b"test")

    def test_assess_quality_returns_float(self, engine):
        quality = engine.assess_quality(b"some_image_bytes")
        assert isinstance(quality, float)
        assert 0.0 <= quality <= 1.0

    def test_embedding_differs_for_different_inputs(self, engine):
        emb_a = engine.extract_embedding(b"a" * 50)
        emb_b = engine.extract_embedding(b"b" * 100)
        assert not np.allclose(emb_a, emb_b)

    def test_load_model_sets_flag(self):
        from services.biometric.inference.facial import FacialInferenceEngine

        engine = FacialInferenceEngine()
        assert engine.model_loaded is False
        engine.load_model()
        assert engine.model_loaded is True

    def test_global_engine_singleton(self):
        from services.biometric.inference.facial import (
            get_facial_engine,
            init_facial_model,
        )

        engine_before = get_facial_engine()
        assert engine_before.model_loaded is False
        init_facial_model()
        engine_after = get_facial_engine()
        assert engine_after.model_loaded is True
        assert engine_before is engine_after

    def test_extract_embedding_handles_large_input(self, engine):
        large_bytes = b"x" * 100_000
        embedding = engine.extract_embedding(large_bytes)
        assert embedding.shape == (512,)
        assert np.isclose(np.linalg.norm(embedding), 1.0)

    def test_quality_score_stable(self, engine):
        q1 = engine.assess_quality(b"img_data")
        q2 = engine.assess_quality(b"img_data")
        assert q1 == q2
