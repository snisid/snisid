import uuid
import base64
from fastapi import APIRouter, HTTPException, Depends, status
from structlog import get_logger

from services.biometric.models import (
    BiometricExtractionRequest, BiometricExtractionResponse,
    BiometricVerifyRequest, BiometricVerifyResponse,
    BiometricIdentifyRequest, BiometricIdentifyResponse
)
from services.biometric.inference.facial import get_facial_engine, FacialInferenceEngine
from services.biometric.inference.liveness import get_liveness_engine, LivenessEngine
from services.biometric.matching.engine import get_matching_engine, MatchingEngine
from services.biometric.security.crypto import get_crypto_vault, BiometricCryptoVault

logger = get_logger(__name__)
router = APIRouter()

# In a real system, templates would be stored in Postgres and pulled into Faiss memory on startup.
# We stub a tiny memory store for testing the routes.
TEMPLATE_DB = {}

@router.post("/extract", response_model=BiometricExtractionResponse, status_code=status.HTTP_201_CREATED)
async def extract_template(
    req: BiometricExtractionRequest,
    facial_engine: FacialInferenceEngine = Depends(get_facial_engine),
    liveness_engine: LivenessEngine = Depends(get_liveness_engine),
    crypto: BiometricCryptoVault = Depends(get_crypto_vault),
    matching: MatchingEngine = Depends(get_matching_engine)
):
    try:
        image_bytes = base64.b64decode(req.image_base64)
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid base64 image data.")

    # 1. Quality Check
    quality = facial_engine.assess_quality(image_bytes)
    if quality < 0.8: # ISO/NIST minimum quality threshold
        raise HTTPException(status_code=422, detail="Image quality too low for extraction.")

    # 2. Liveness Check (PAD)
    liveness_score = 1.0
    if req.check_liveness:
        liveness_score = liveness_engine.detect(image_bytes)
        if liveness_score < liveness_engine.threshold:
            logger.warning("Liveness detection failed. Possible presentation attack.", liveness_score=liveness_score)
            raise HTTPException(status_code=403, detail="Presentation attack detected.")

    # 3. Extraction
    embedding = facial_engine.extract_embedding(image_bytes)

    # 4. Encryption & Storage
    template_id = str(uuid.uuid4())
    ciphertext = crypto.encrypt_template(embedding)
    
    # Save to DB (mock) and add to in-memory Faiss index
    TEMPLATE_DB[template_id] = ciphertext
    matching.add_template(template_id, embedding)
    
    # IMPORTANT: The image_bytes variable is strictly locally scoped and will be garbage collected.
    # No raw images are written to the container filesystem.

    logger.info("Biometric template extracted and encrypted successfully.", template_id=template_id)
    return BiometricExtractionResponse(
        template_id=template_id,
        liveness_score=liveness_score,
        quality_score=quality
    )

@router.post("/verify", response_model=BiometricVerifyResponse)
async def verify_1_to_1(
    req: BiometricVerifyRequest,
    facial_engine: FacialInferenceEngine = Depends(get_facial_engine),
    crypto: BiometricCryptoVault = Depends(get_crypto_vault),
    matching: MatchingEngine = Depends(get_matching_engine)
):
    # Retrieve target from DB
    if req.target_template_id not in TEMPLATE_DB:
        raise HTTPException(status_code=404, detail="Target template not found.")
        
    target_ciphertext = TEMPLATE_DB[req.target_template_id]
    target_embedding = crypto.decrypt_template(target_ciphertext)

    try:
        probe_bytes = base64.b64decode(req.probe_image_base64)
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid base64 image data.")

    probe_embedding = facial_engine.extract_embedding(probe_bytes)
    
    is_match, conf = matching.verify_1_to_1(probe_embedding, target_embedding)
    
    return BiometricVerifyResponse(match=is_match, confidence=conf)

@router.delete("/{template_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_template(
    template_id: str,
    matching: MatchingEngine = Depends(get_matching_engine)
):
    """
    RGPD Right to Erasure.
    Removes the template from the database and the FAISS index.
    """
    if template_id in TEMPLATE_DB:
        del TEMPLATE_DB[template_id]
        matching.remove_template(template_id)
        logger.info(f"Template {template_id} cryptographically erased.")
    return None
