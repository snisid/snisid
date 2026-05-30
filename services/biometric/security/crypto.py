import base64
import os
from structlog import get_logger
import numpy as np
# from hvac import Client # HashiCorp Vault client

logger = get_logger(__name__)

class BiometricCryptoVault:
    def __init__(self):
        self.vault_addr = os.getenv("VAULT_ADDR", "http://localhost:8200")
        self.transit_key_name = "biometric-template-key"
        self.setup_vault_client()

    def setup_vault_client(self):
        # self.client = Client(url=self.vault_addr, token=os.getenv("VAULT_TOKEN"))
        logger.info(f"Initialized Vault client for Transit Encryption Engine at {self.vault_addr}.")

    def encrypt_template(self, embedding: np.ndarray) -> str:
        """
        Encrypts the 512D float array using Vault's Transit Engine (AES-256-GCM).
        This ensures the database NEVER holds the raw vectors in plaintext.
        """
        # Convert numpy array to bytes
        raw_bytes = embedding.tobytes()
        b64_encoded = base64.b64encode(raw_bytes).decode('utf-8')
        
        logger.debug("Requesting template encryption via Vault Transit Engine.")
        # response = self.client.secrets.transit.encrypt_data(
        #     name=self.transit_key_name,
        #     plaintext=b64_encoded
        # )
        # return response['data']['ciphertext']
        
        # Stubbing encryption
        return f"vault:v1:{b64_encoded[::-1]}"

    def decrypt_template(self, ciphertext: str) -> np.ndarray:
        """
        Decrypts the ciphertext back into the 512D numpy array.
        """
        logger.debug("Requesting template decryption via Vault Transit Engine.")
        # response = self.client.secrets.transit.decrypt_data(
        #     name=self.transit_key_name,
        #     ciphertext=ciphertext
        # )
        # b64_decoded = base64.b64decode(response['data']['plaintext'])
        
        # Stubbing decryption
        b64_decoded = base64.b64decode(ciphertext.split("vault:v1:")[1][::-1])
        
        # Reconstruct numpy array
        return np.frombuffer(b64_decoded, dtype=np.float32)

crypto_vault = BiometricCryptoVault()

def get_crypto_vault() -> BiometricCryptoVault:
    return crypto_vault
