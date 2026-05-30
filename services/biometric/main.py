from fastapi import FastAPI
from structlog import get_logger
from shared.middleware import setup_middleware
from services.biometric.api import routes

logger = get_logger(__name__)

app = FastAPI(
    title="SNISID Biometric Service",
    description="NIST-compliant Biometric Matching, Liveness Detection, and Extraction Engine.",
    version="1.0.0",
)

# Setup Zero Trust Middlewares (Logging, Security Headers, Audit, Input Sanitization)
setup_middleware(app)

# Include API routes
app.include_router(routes.router, prefix="/api/v1/biometrics", tags=["Biometrics"])

@app.on_event("startup")
async def startup_event():
    logger.info("Biometric Service starting. Loading AI models and FAISS indexes.")
    # Initialize engines
    from services.biometric.inference.facial import init_facial_model
    from services.biometric.inference.liveness import init_liveness_model
    from services.biometric.matching.engine import init_faiss_index
    
    init_facial_model()
    init_liveness_model()
    init_faiss_index()

@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Biometric Service shutting down. Clearing memory.")
