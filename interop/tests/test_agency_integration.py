import pytest
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives import serialization
from interop.security.jws_signing_interceptor import create_jws_detached, verify_jws_detached

@pytest.fixture
def keys():
    private_key = rsa.generate_private_key(public_exponent=65537, key_size=2048)
    public_key = private_key.public_key()
    
    priv_pem = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.NoEncryption()
    )
    
    pub_pem = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo
    )
    return priv_pem, pub_pem

def test_jws_signing_and_verification(keys):
    priv_pem, pub_pem = keys
    payload = {"national_id": "1234567890-X"}
    
    # Agency creates the signature
    jws = create_jws_detached(payload, priv_pem)
    
    # SNISID Gateway verifies it
    is_valid = verify_jws_detached(payload, jws, pub_pem)
    assert is_valid is True

def test_jws_tampering_rejected(keys):
    priv_pem, pub_pem = keys
    payload = {"national_id": "1234567890-X"}
    
    jws = create_jws_detached(payload, priv_pem)
    
    # Attacker modifies the payload over the wire
    tampered_payload = {"national_id": "9999999999-Y"}
    
    is_valid = verify_jws_detached(tampered_payload, jws, pub_pem)
    assert is_valid is False
