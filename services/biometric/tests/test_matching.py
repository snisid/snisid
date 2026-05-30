import pytest
import numpy as np
from services.biometric.matching.engine import MatchingEngine
from services.biometric.security.crypto import BiometricCryptoVault

def test_l2_normalization_and_cosine_similarity():
    engine = MatchingEngine()
    
    # Two identical vectors
    v1 = np.ones(512, dtype=np.float32)
    norm = np.linalg.norm(v1)
    v1_norm = v1 / norm
    
    is_match, score = engine.verify_1_to_1(v1_norm, v1_norm)
    assert is_match is True
    assert score > 0.99  # Should be ~1.0

def test_far_frr_threshold_rejection():
    engine = MatchingEngine()
    
    # Completely orthogonal vectors (cosine sim = 0)
    v1 = np.zeros(512, dtype=np.float32)
    v1[0] = 1.0
    
    v2 = np.zeros(512, dtype=np.float32)
    v2[1] = 1.0
    
    is_match, score = engine.verify_1_to_1(v1, v2)
    assert is_match is False
    assert score < 0.1

def test_crypto_vault_encryption_decryption():
    crypto = BiometricCryptoVault()
    
    original = np.random.rand(512).astype(np.float32)
    ciphertext = crypto.encrypt_template(original)
    
    assert ciphertext.startswith("vault:v1:")
    assert "numpy" not in ciphertext
    
    decrypted = crypto.decrypt_template(ciphertext)
    np.testing.assert_array_almost_equal(original, decrypted)
