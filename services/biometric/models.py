from pydantic import BaseModel, Field
from typing import List, Optional

class BiometricExtractionRequest(BaseModel):
    image_base64: str = Field(..., description="Base64 encoded raw image for extraction.")
    modality: str = Field(..., pattern="^(facial|fingerprint|iris)$", description="Type of biometric data.")
    check_liveness: bool = Field(True, description="Whether to perform Presentation Attack Detection (PAD).")

class BiometricExtractionResponse(BaseModel):
    template_id: str = Field(..., description="UUID of the encrypted template stored in the vault.")
    liveness_score: float = Field(..., description="Score from 0.0 to 1.0 indicating liveness probability.")
    quality_score: float = Field(..., description="NIST-compliant quality assessment score.")

class BiometricVerifyRequest(BaseModel):
    probe_image_base64: str = Field(..., description="Base64 encoded probe image.")
    target_template_id: str = Field(..., description="UUID of the target template to verify against.")

class BiometricVerifyResponse(BaseModel):
    match: bool = Field(..., description="True if FAR < 0.001%.")
    confidence: float = Field(..., description="Cosine similarity confidence score.")

class BiometricIdentifyRequest(BaseModel):
    probe_image_base64: str = Field(..., description="Base64 encoded probe image.")
    top_k: int = Field(5, description="Number of closest matches to return.")

class IdentificationMatch(BaseModel):
    template_id: str
    confidence: float

class BiometricIdentifyResponse(BaseModel):
    matches: List[IdentificationMatch]
