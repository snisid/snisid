import json
import base64
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives import serialization
from cryptography.exceptions import InvalidSignature

def create_jws_detached(payload: dict, private_key_pem: bytes) -> str:
    """
    Creates a Detached JSON Web Signature (JWS) for the given payload.
    Used by Agencies to sign their requests to SNISID.
    """
    private_key = serialization.load_pem_private_key(private_key_pem, password=None)
    
    # Standard JWS Header for RS256
    header = {"alg": "RS256", "typ": "JWT"}
    header_b64 = base64.urlsafe_b64encode(json.dumps(header).encode()).decode().rstrip("=")
    
    # Payload
    payload_str = json.dumps(payload, separators=(',', ':'))
    payload_b64 = base64.urlsafe_b64encode(payload_str.encode()).decode().rstrip("=")
    
    # Sign: header.payload
    signing_input = f"{header_b64}.{payload_b64}".encode()
    signature = private_key.sign(
        signing_input,
        padding.PKCS1v15(),
        hashes.SHA256()
    )
    sig_b64 = base64.urlsafe_b64encode(signature).decode().rstrip("=")
    
    # Detached JWS format: header..signature (Payload is omitted)
    return f"{header_b64}..{sig_b64}"

def verify_jws_detached(payload: dict, jws: str, public_key_pem: bytes) -> bool:
    """
    Verifies a Detached JWS.
    Used by SNISID Camel routes to verify incoming Agency requests.
    """
    public_key = serialization.load_pem_public_key(public_key_pem)
    
    parts = jws.split('.')
    if len(parts) != 3 or parts[1] != "":
        return False
        
    header_b64 = parts[0]
    sig_b64 = parts[2]
    
    # Reconstruct signing input
    payload_str = json.dumps(payload, separators=(',', ':'))
    payload_b64 = base64.urlsafe_b64encode(payload_str.encode()).decode().rstrip("=")
    signing_input = f"{header_b64}.{payload_b64}".encode()
    
    signature = base64.urlsafe_b64decode(sig_b64 + "=" * (-len(sig_b64) % 4))
    
    try:
        public_key.verify(
            signature,
            signing_input,
            padding.PKCS1v15(),
            hashes.SHA256()
        )
        return True
    except InvalidSignature:
        return False
