import asyncio, json, logging, os, uuid
from datetime import datetime
from typing import Optional

import uvicorn
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel, Field

from graphrag_engine import GraphRAGEngine, ReportType

logger = logging.getLogger("sniside.graphrag-api")
app = FastAPI(title="SNI-SIDE GraphRAG Intelligence Engine", version="2.1")

engine = GraphRAGEngine()


class ReportRequest(BaseModel):
    entity_id: str = Field(..., description="NIU, plate, phone, etc.")
    entity_type: str = Field("Citizen", description="Citizen, Vehicle, Phone, Gang, etc.")
    entity_label: str = Field(None, description="Display name")
    report_type: ReportType = Field(ReportType.ENTITY_PROFILE)
    depth: int = Field(2, ge=1, le=6)
    since: int = Field(None, description="Timestamp in ms for temporal analysis")
    entity_id2: str = Field(None, description="Second entity for CROSS_ENTITY")


class ReportResponse(BaseModel):
    report_id: str
    report_type: str
    generated_at: str
    target_entity: dict
    executive_summary: str
    key_findings: list
    risk_assessment: dict
    graph_context: dict
    recommendations: list
    confidence_score: float
    model_used: str


class CrossSearchRequest(BaseModel):
    query: str = Field(..., description="Natural language query")
    max_results: int = Field(20, le=100)
    include_graph: bool = Field(True)


class CrossSearchResult(BaseModel):
    query: str
    results: list
    graph_context: Optional[dict]
    execution_time_ms: int


@app.on_event("startup")
async def startup():
    engine.start()
    logger.info("GraphRAG API started")


@app.on_event("shutdown")
async def shutdown():
    engine.stop()


@app.get("/health")
async def health():
    return {"status": "ok", "service": "graphrag-engine", "version": "2.1"}


@app.post("/intelligence/report", response_model=ReportResponse)
async def generate_report(req: ReportRequest):
    try:
        report = await engine.generate_report(
            entity_id=req.entity_id,
            entity_type=req.entity_type,
            entity_label=req.entity_label or req.entity_id,
            report_type=req.report_type,
            depth=req.depth,
            since=req.since,
            entity_id2=req.entity_id2,
        )
        return report
    except Exception as e:
        logger.error(f"Report generation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/intelligence/report/{report_type}/{entity_id}", response_model=ReportResponse)
async def get_report(
    report_type: ReportType,
    entity_id: str,
    entity_type: str = Query("Citizen"),
    entity_label: str = Query(None),
    depth: int = Query(2),
):
    try:
        report = await engine.generate_report(
            entity_id=entity_id,
            entity_type=entity_type,
            entity_label=entity_label or entity_id,
            report_type=report_type,
            depth=depth,
        )
        return report
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/intelligence/cross-search", response_model=CrossSearchResult)
async def cross_search(req: CrossSearchRequest):
    import time
    start = time.monotonic()
    report = await engine.generate_report(
        entity_id=req.query,
        entity_type="Query",
        entity_label=req.query,
        report_type=ReportType.LINK_ANALYSIS,
    )
    elapsed = int((time.monotonic() - start) * 1000)
    return CrossSearchResult(
        query=req.query,
        results=[report] if report else [],
        graph_context=report.get("graph_context") if req.include_graph else None,
        execution_time_ms=elapsed,
    )


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8080)
