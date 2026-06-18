"""
SNI-SIDE: National Sovereign API Server
========================================
FastAPI application connecting all 15 national databases,
the Search Engine, AI Fusion Center, and infrastructure.

Services:
- NCID Criminal Intelligence
- HN-NGI Biometrics
- HN-CODIS DNA
- Missing Persons
- Vehicle Intelligence
- National ALPR
- Firearms Intelligence
- Border Intelligence
- Counter Narcotics
- Financial Crime
- Cybercrime
- National Watchlist
- Document Fraud
- GEOINT
- Digital Evidence
- National Sovereign Search Engine
- National AI Fusion Center
"""

import time
import uuid
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request, Response
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import JSONResponse

from config import settings
from database import db


# ============ LIFESPAN ============
@asynccontextmanager
async def lifespan(app: FastAPI):
    """Initialisation et nettoyage des ressources"""
    print(f"[SNI-SIDE] Starting API Server — {settings.app_name}")
    print(f"[SNI-SIDE] Environment: {settings.environment.value}")

    # Initialize database connections
    await db.initialize()
    print("[SNI-SIDE] Database connections initialized")

    # Verify connections
    async with db.pg_conn() as conn:
        version = await conn.fetchval("SELECT version()")
        print(f"[SNI-SIDE] PostgreSQL: {version[:50]}...")

    async with db.cockroach_conn() as conn:
        version = await conn.fetchval("SELECT version()")
        print(f"[SNI-SIDE] CockroachDB: {version[:50]}...")

    print(f"[SNI-SIDE] Neo4j: {settings.neo4j_uri}")
    print(f"[SNI-SIDE] Milvus: {settings.milvus_host}:{settings.milvus_port}")
    print(f"[SNI-SIDE] ClickHouse: {settings.clickhouse_host}:{settings.clickhouse_port}")
    print(f"[SNI-SIDE] Kafka: {settings.kafka_bootstrap_servers}")
    print(f"[SNI-SIDE] Redis: {settings.redis_host}:{settings.redis_port}")
    print(f"[SNI-SIDE] MinIO: {settings.minio_endpoint}")

    yield

    # Cleanup
    await db.close()
    print("[SNI-SIDE] Database connections closed")


# ============ APPLICATION ============
app = FastAPI(
    title=settings.app_name,
    description="SNISID National Intelligence, Security, Investigation and Sovereign Data Ecosystem",
    version="1.0.0",
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.cors_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


# ============ MIDDLEWARE ============
@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    """Ajoute le temps de traitement aux réponses"""
    start_time = time.time()
    response = await call_next(request)
    process_time = (time.time() - start_time) * 1000
    response.headers["X-Process-Time-Ms"] = str(int(process_time))
    response.headers["X-Request-ID"] = request.headers.get("X-Request-ID", str(uuid.uuid4()))
    return response


@app.middleware("http")
async def catch_exceptions(request: Request, call_next):
    """Capture globale des exceptions"""
    try:
        return await call_next(request)
    except Exception as e:
        print(f"[ERROR] {type(e).__name__}: {e}")
        return JSONResponse(
            status_code=500,
            content={"detail": "Internal server error", "type": type(e).__name__},
        )


# ============ HEALTH ============
@app.get("/health", tags=["System"])
async def health():
    """Vérification de santé du service"""
    checks = {
        "status": "healthy",
        "app": settings.app_name,
        "version": "1.0.0",
        "environment": settings.environment.value,
    }

    # Database checks
    try:
        async with db.pg_conn() as conn:
            await conn.fetchval("SELECT 1")
            checks["postgresql"] = "connected"
    except Exception as e:
        checks["postgresql"] = f"error: {e}"

    try:
        async with db.cockroach_conn() as conn:
            await conn.fetchval("SELECT 1")
            checks["cockroachdb"] = "connected"
    except Exception as e:
        checks["cockroachdb"] = f"error: {e}"

    try:
        await db.redis_client.ping()
        checks["redis"] = "connected"
    except Exception as e:
        checks["redis"] = f"error: {e}"

    try:
        async with await db.neo4j_session() as session:
            await session.run("RETURN 1")
            checks["neo4j"] = "connected"
    except Exception as e:
        checks["neo4j"] = f"error: {e}"

    return checks


@app.get("/ready", tags=["System"])
async def ready():
    """Prêt à recevoir du trafic"""
    return {"status": "ready"}


# ============ ROUTES ============

# NCID — Criminal Intelligence
from routes.ncid import router as ncid_router
app.include_router(ncid_router, prefix="/intelligence/v1")

# HN-NGI — Biometrics
from routes.biometrics import router as biometrics_router
app.include_router(biometrics_router, prefix="/intelligence/v1")

# National Sovereign Search Engine
from routes.search import router as search_router
app.include_router(search_router, prefix="/intelligence/v1")

# National Alerts
from routes.alerts import router as alerts_router
app.include_router(alerts_router, prefix="/intelligence/v1")

# National AI Fusion Center
from routes.ai_fusion import router as ai_fusion_router
app.include_router(ai_fusion_router, prefix="/intelligence/v1")


# ============ MÉTADONNÉES ============
@app.get("/", tags=["System"])
async def root():
    """Racine de l'API — informations générales"""
    return {
        "name": settings.app_name,
        "version": "1.0.0",
        "environment": settings.environment.value,
        "databases": [
            "NCID", "HN-NGI", "HN-CODIS", "Missing Persons",
            "Vehicle Intelligence", "National ALPR", "Firearms Intelligence",
            "Border Intelligence", "Counter Narcotics", "Financial Crime",
            "Cybercrime Intelligence", "National Watchlist", "Document Fraud",
            "GEOINT", "Digital Evidence",
        ],
        "capabilities": [
            "Unified Search", "Graph Intelligence", "AI Fusion",
            "Real-Time Alerts", "Federated Queries",
        ],
        "documentation": "/docs",
    }


# ============ ENTRY POINT ============
if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "main:app",
        host=settings.host,
        port=settings.port,
        workers=settings.workers,
        reload=settings.reload,
    )
