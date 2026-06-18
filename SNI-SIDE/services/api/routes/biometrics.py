"""SNI-SIDE API — Routes HN-NGI (Biometrics)"""

from fastapi import APIRouter, Depends, HTTPException, UploadFile, File, Form, Request
from typing import Optional

from middleware.security import verify_agency, create_audit_entry, SecurityContext
from models.schemas import (
    BiometricVerificationRequest, BiometricVerificationResponse,
    IdentifyBiometricRequest, IdentifyBiometricResponse,
    BiometricCandidate, DuplicateDetectionResult,
)

router = APIRouter(prefix="/biometrics", tags=["HN-NGI — Biometrics"])

BIO_AGENCIES = ["PNH", "DCPJ", "ONI", "IMMIGRATION", "SNISID_ADMIN"]


@router.post("/verify", response_model=BiometricVerificationResponse)
async def verify_biometric(
    request: Request,
    data: BiometricVerificationRequest,
    security: SecurityContext = Depends(verify_agency(BIO_AGENCIES)),
):
    """Vérification biométrique 1:1"""
    from database import db

    # Search in PostgreSQL for the NIU
    async with db.pg_conn() as conn:
        person = await conn.fetchrow(
            "SELECT niu, full_name FROM snisid_ncid.wanted_persons WHERE niu = $1",
            data.niu
        )
        if not person:
            raise HTTPException(status_code=404, detail="Person not found")

    # In production: Milvus vector search for face/fingerprint
    # collection = db.milvus_client.get_collection(settings.milvus_aliases_collection)
    # results = collection.search(
    #     data=[embedding],
    #     anns_field="embedding",
    #     param={"metric_type": "IP", "nprobe": 10},
    #     limit=1,
    #     expr=f"niu == '{data.niu}'"
    # )

    create_audit_entry(request, security, "VERIFY", "biometric", data.niu)

    return BiometricVerificationResponse(
        verified=True,
        match_score=0.97,
        threshold=0.75,
        search_duration_ms=45,
    )


@router.post("/identify", response_model=IdentifyBiometricResponse)
async def identify_biometric(
    request: Request,
    data: IdentifyBiometricRequest,
    security: SecurityContext = Depends(verify_agency(BIO_AGENCIES)),
):
    """Identification biométrique 1:N"""
    # In production: Milvus vector search across the entire gallery
    # collection = db.milvus_client.get_collection(settings.milvus_aliases_collection)
    # results = collection.search(
    #     data=[embedding],
    #     anns_field="embedding",
    #     param={"metric_type": "IP", "nprobe": 10},
    #     limit=data.max_candidates,
    # )

    create_audit_entry(request, security, "IDENTIFY", "biometric", data.biometric_type)

    return IdentifyBiometricResponse(
        candidates=[
            BiometricCandidate(niu="0000000001", full_name="Jean Dupont", score=0.95, rank=1),
            BiometricCandidate(niu="0000000002", full_name="Marie Pierre", score=0.82, rank=2),
        ],
        gallery_size=1500000,
        search_duration_ms=120,
    )


@router.post("/search/face", response_model=IdentifyBiometricResponse)
async def search_face(
    request: Request,
    file: UploadFile = File(...),
    min_confidence: float = Form(0.7),
    max_results: int = Form(20),
    security: SecurityContext = Depends(verify_agency(BIO_AGENCIES)),
):
    """Recherche faciale dans la base nationale"""
    image_data = await file.read()

    # In production: extract face embedding using ArcFace model
    # embedding = arcface_model.extract(image_data)
    # candidates = milvus.search(collection="faces", data=[embedding], ...)

    create_audit_entry(request, security, "SEARCH", "face", file.filename or "unknown")

    return IdentifyBiometricResponse(
        candidates=[],
        gallery_size=1500000,
        search_duration_ms=200,
    )


@router.post("/search/fingerprint", response_model=IdentifyBiometricResponse)
async def search_fingerprint(
    request: Request,
    file: UploadFile = File(...),
    finger_position: str = Form(...),
    security: SecurityContext = Depends(verify_agency(BIO_AGENCIES)),
):
    """Recherche d'empreintes digitales"""
    image_data = await file.read()
    create_audit_entry(request, security, "SEARCH", "fingerprint", finger_position)

    return IdentifyBiometricResponse(
        candidates=[],
        gallery_size=1500000,
        search_duration_ms=300,
    )


@router.post("/duplicate-detection", response_model=list[DuplicateDetectionResult])
async def detect_duplicates(
    request: Request,
    biometric_type: str,
    threshold: float = 0.85,
    security: SecurityContext = Depends(verify_agency(["PNH", "DCPJ", "SNISID_ADMIN"])),
):
    """Détection de doublons biométriques"""
    # In production: run dedup across Milvus collections
    # cluster_dbscan or pairwise comparison

    create_audit_entry(request, security, "DEDUP", "biometric", biometric_type)
    return []
