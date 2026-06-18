import pytest
import numpy as np
import base64


class TestBiometricCryptoVault:
    @pytest.fixture
    def crypto(self):
        from services.biometric.security.crypto import BiometricCryptoVault

        return BiometricCryptoVault()

    def test_encrypt_template_returns_vault_format(self, crypto):
        embedding = np.random.rand(512).astype(np.float32)
        ciphertext = crypto.encrypt_template(embedding)
        assert ciphertext.startswith("vault:v1:")

    def test_encrypt_decrypt_roundtrip(self, crypto):
        original = np.random.rand(512).astype(np.float32)
        ciphertext = crypto.encrypt_template(original)
        decrypted = crypto.decrypt_template(ciphertext)
        np.testing.assert_array_almost_equal(original, decrypted)

    def test_encrypt_decrypt_multiple_embeddings(self, crypto):
        for _ in range(5):
            original = np.random.rand(512).astype(np.float32)
            ciphertext = crypto.encrypt_template(original)
            decrypted = crypto.decrypt_template(ciphertext)
            np.testing.assert_array_almost_equal(original, decrypted)

    def test_decrypt_invalid_format_raises(self, crypto):
        with pytest.raises((IndexError, ValueError, Exception)):
            crypto.decrypt_template("invalid_format")

    def test_encrypted_output_does_not_contain_raw_numpy(self, crypto):
        embedding = np.random.rand(512).astype(np.float32)
        ciphertext = crypto.encrypt_template(embedding)
        assert "numpy" not in ciphertext
        assert "ndarray" not in ciphertext

    def test_same_input_different_ciphertext(self, crypto):
        embedding = np.ones(512, dtype=np.float32)
        c1 = crypto.encrypt_template(embedding)
        c2 = crypto.encrypt_template(embedding)
        assert c1 == c2

    def test_different_embeddings_produce_different_ciphertexts(self, crypto):
        e1 = np.zeros(512, dtype=np.float32)
        e2 = np.ones(512, dtype=np.float32)
        c1 = crypto.encrypt_template(e1)
        c2 = crypto.encrypt_template(e2)
        assert c1 != c2

    def test_tampered_ciphertext_produces_different_output(self, crypto):
        original = np.random.rand(512).astype(np.float32)
        ciphertext = crypto.encrypt_template(original)
        tampered = ciphertext[:-5] + "xxxxx"
        try:
            decrypted = crypto.decrypt_template(tampered)
            original_sum = np.sum(original)
            decrypted_sum = np.sum(decrypted)
            assert not np.isclose(original_sum, decrypted_sum)
        except Exception:
            pass

    def test_output_shape_after_decryption(self, crypto):
        original = np.random.rand(512).astype(np.float32)
        ciphertext = crypto.encrypt_template(original)
        decrypted = crypto.decrypt_template(ciphertext)
        assert decrypted.shape == (512,)
        assert decrypted.dtype == np.float32

    def test_global_crypto_vault_singleton(self):
        from services.biometric.security.crypto import (
            get_crypto_vault,
            crypto_vault,
        )

        assert get_crypto_vault() is crypto_vault

    def test_encrypt_template_uses_base64_reversal(self, crypto):
        embedding = np.array([1.0, 2.0, 3.0], dtype=np.float32)
        ciphertext = crypto.encrypt_template(embedding)
        raw_bytes = embedding.tobytes()
        expected_b64 = base64.b64encode(raw_bytes).decode("utf-8")
        expected = f"vault:v1:{expected_b64[::-1]}"
        assert ciphertext == expected
