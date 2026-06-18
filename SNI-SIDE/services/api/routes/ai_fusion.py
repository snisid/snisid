"""SNI-SIDE API — Routes AI Fusion Center"""

from fastapi import APIRouter, Depends, HTTPException, Request, UploadFile, File

from middleware.security import verify_agency, verify_clearance, create_audit_entry, SecurityContext
from models.schemas import (
    FraudAnalysisRequest, FraudAnalysisResponse,
    GraphRAGQuery, GraphRAGResponse, GraphNode, GraphEdge,
    AMLAnalysisResponse, AMLRiskFactor,
)

router = APIRouter(prefix="/ai", tags=["National AI Fusion Center"])

AI_AGENCIES = ["PNH", "DCPJ", "INTELLIGENCE", "FIU", "SOC", "SNISID_ADMIN"]


@router.post("/fraud-analysis", response_model=FraudAnalysisResponse)
async def analyze_fraud(
    request: Request,
    data: FraudAnalysisRequest,
    security: SecurityContext = Depends(verify_agency(AI_AGENCIES)),
):
    """Analyse de fraude multi-domaine via GNN"""
    from database import db

    # 1. Query Neo4j for the person's subgraph
    async with await db.neo4j_session() as session:
        result = await session.run(
            """
            MATCH path = (c:Citizen {niu: $niu})-[*1..$depth]-(connected)
            RETURN c, collect(path) as paths
            """,
            niu=data.entity_niu, depth=data.graph_depth
        )
        record = await result.single()

    # 2. Run fraud GNN inference
    # risk_score = fraud_gnn_model.predict(subgraph_features)
    risk_score = 0.72

    # 3. Generate explanation
    risk_level = "HIGH" if risk_score >= 0.7 else "MEDIUM" if risk_score >= 0.4 else "LOW"
    indicators = []
    if risk_score > 0.5:
        indicators.append("Multiple associated entities with high risk scores")
    if risk_score > 0.7:
        indicators.append("Financial transactions to high-risk jurisdictions")
        indicators.append("Connected to known criminal organization")

    create_audit_entry(request, security, "AI_ANALYSIS", "fraud", data.entity_niu)

    return FraudAnalysisResponse(
        risk_score=risk_score,
        risk_level=risk_level,
        fraud_indicators=indicators,
        feature_importance={"graph_centrality": 0.45, "transaction_anomaly": 0.30, "associate_risk": 0.25},
        model_version="fraud-gnn-v1.0.0",
        network_context={"nodes": 12, "edges": 18, "depth": data.graph_depth},
    )


@router.post("/graphrag", response_model=GraphRAGResponse)
async def graphrag_query(
    request: Request,
    data: GraphRAGQuery,
    security: SecurityContext = Depends(verify_clearance("SECRET")),
):
    """Requête GraphRAG — Intelligence augmentée par le graphe national"""
    from database import db

    # 1. Retrieve context from Neo4j
    context_nodes = []
    context_edges = []

    if data.seed_entity_id:
        async with await db.neo4j_session() as session:
            result = await session.run(
                """
                MATCH (start)
                WHERE start.niu = $entity_id OR start.phone = $entity_id
                      OR start.plate = $entity_id
                CALL apoc.path.subgraph(start, {
                    maxLevel: $depth,
                    relationshipFilter: 'OWNS|USES|ASSOCIATED_WITH|FINANCED_BY|LINKED_TO|TRAVELLED_WITH'
                })
                YIELD path
                RETURN path
                """,
                entity_id=data.seed_entity_id, depth=data.graph_depth
            )
            async for record in result:
                # Extract nodes and relationships from path
                path = record["path"]
                # Process into context_nodes and context_edges
                pass

    # 2. Generate analysis using LLM with graph context
    # analysis = llm.generate(prompt_with_graph_context)
    analysis = (
        f"GraphRAG analysis for query: '{data.query}'. "
        f"Seed entity: {data.seed_entity_id} ({data.seed_entity_type}). "
        f"Discovered {len(context_nodes)} nodes and {len(context_edges)} edges "
        f"at depth {data.graph_depth}. "
        f"Risk assessment and intelligence report generated."
    )

    create_audit_entry(request, security, "GRAPHRAG", "query", data.query[:50])

    return GraphRAGResponse(
        analysis=analysis,
        context_nodes=context_nodes,
        context_edges=context_edges,
        overall_risk_score=0.65,
        key_findings=[
            "Entity connected to 3 persons of interest",
            "Financial flow pattern matches known ML typology",
            "Phone contact with watchlisted individual",
        ],
        confidence=0.78,
    )


@router.post("/aml-risk", response_model=AMLAnalysisResponse)
async def aml_risk_analysis(
    request: Request,
    entity_niu: str,
    security: SecurityContext = Depends(verify_agency(["FIU", "INTELLIGENCE", "SNISID_ADMIN"])),
):
    """Analyse AML (Anti-Money Laundering)"""
    from database import db

    # Query financial transactions + beneficial owners + PEP status
    async with db.pg_conn() as conn:
        txn_count = await conn.fetchval(
            "SELECT COUNT(*) FROM snisid_financial.suspicious_transactions WHERE sender_niu = $1 OR beneficiary_niu = $1",
            entity_niu
        )
        is_pep = await conn.fetchval(
            "SELECT COUNT(*) FROM snisid_financial.politically_exposed_persons WHERE niu = $1",
            entity_niu
        )
        bo_count = await conn.fetchval(
            "SELECT COUNT(*) FROM snisid_financial.beneficial_owners WHERE niu = $1",
            entity_niu
        )

    risk_score = min(1.0, (txn_count * 0.1 + is_pep * 0.3 + bo_count * 0.2))
    risk_factors = []
    if is_pep:
        risk_factors.append(AMLRiskFactor(factor="PEP", score=0.8, description="Politically Exposed Person"))
    if txn_count > 5:
        risk_factors.append(AMLRiskFactor(factor="HIGH_VOLUME", score=0.6, description=f"{txn_count} suspicious transactions"))
    if bo_count > 3:
        risk_factors.append(AMLRiskFactor(factor="COMPLEX_STRUCTURE", score=0.7, description="Complex beneficial ownership structure"))

    return AMLAnalysisResponse(
        overall_risk_score=risk_score,
        risk_level="HIGH" if risk_score > 0.7 else "MEDIUM",
        risk_factors=risk_factors,
        recommendations=["Enhanced due diligence required", "File STR to FIU", "Monitor all related accounts"],
    )


@router.post("/detect-deepfake")
async def detect_deepfake(
    request: Request,
    file: UploadFile = File(...),
    security: SecurityContext = Depends(verify_agency(AI_AGENCIES)),
):
    """Détection de deepfake sur média"""
    media_data = await file.read()

    # In production: run through EfficientNet deepfake model
    # result = deepfake_model.predict(media_data)

    create_audit_entry(request, security, "AI_DEEPFAKE", "detect", file.filename or "unknown")

    return {
        "is_fake": False,
        "confidence": 0.03,
        "artifact_scores": {"face_warping": 0.02, "blinking": 0.01, "lighting": 0.05},
        "analysis_method": "efficientnet-b4",
    }


@router.post("/predict-crime")
async def predict_crime(
    region: str,
    days_ahead: int = 30,
    security: SecurityContext = Depends(verify_agency(["PNH", "INTELLIGENCE", "SNISID_ADMIN"])),
):
    """Prédiction de criminalité par région"""
    # In production: spatio-temporal crime prediction model
    return {
        "region": region,
        "predictions": [
            {"crime_type": "VOL", "probability": 0.45, "trend": "increasing", "hotspot_zones": ["zone_a", "zone_b"]},
            {"crime_type": "AGRESSION", "probability": 0.32, "trend": "stable", "hotspot_zones": ["zone_c"]},
            {"crime_type": "NARCOTICS", "probability": 0.28, "trend": "decreasing", "hotspot_zones": ["zone_d"]},
        ],
        "model_confidence": 0.73,
        "period": f"{days_ahead} days",
    }
