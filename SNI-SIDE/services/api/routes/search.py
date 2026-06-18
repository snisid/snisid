"""SNI-SIDE API — Routes du Moteur de Recherche National"""

from fastapi import APIRouter, Depends, HTTPException, Query, Request
from typing import Optional

from middleware.security import verify_token, verify_agency, create_audit_entry, SecurityContext
from models.schemas import (
    UnifiedSearchQuery, UnifiedSearchResponse, SearchResultItem,
    GraphSearchQuery, GraphSearchResponse, GraphNode, GraphEdge,
)

router = APIRouter(prefix="/search", tags=["National Sovereign Search Engine"])


@router.get("/unified", response_model=UnifiedSearchResponse)
async def unified_search(
    request: Request,
    q: str = Query(..., min_length=1, max_length=500),
    type: str = Query("ALL", pattern=r'^(ALL|NAME|NIU|PHOTO|FINGERPRINT|DNA|PHONE|EMAIL|ADDRESS|PLATE|VIN|PASSPORT|CASE)$'),
    databases: Optional[str] = Query(None),
    page: int = Query(1, ge=1),
    limit: int = Query(20, ge=1, le=100),
    fuzzy: bool = Query(True),
    security: SecurityContext = Depends(verify_token),
):
    """Recherche unifiée à travers toutes les bases nationales"""
    from ai.national_search_engine import search_engine

    db_list = databases.split(",") if databases else []

    result = await search_engine.unified_search(
        query_str=q, search_type=type, databases=db_list,
        page=page, limit=limit, fuzzy=fuzzy,
    )

    # Convert to response model
    response = UnifiedSearchResponse(
        query=result.query,
        total_results=result.total_results,
        page=result.page,
        limit=result.limit,
        search_duration_ms=result.search_duration_ms,
        databases_searched=result.databases_searched,
        results={
            db.value: [
                SearchResultItem(
                    database=item.database.value,
                    result_type=item.result_type,
                    id=item.id,
                    title=item.title,
                    description=item.description,
                    score=item.score,
                    confidence=item.match_confidence,
                    risk_score=item.risk_score,
                    metadata=item.metadata,
                ) for item in items
            ] for db, items in result.results.items()
        },
        graph_context=result.graph_context,
        suggested_queries=result.suggested_queries,
    )

    create_audit_entry(request, security, "SEARCH", "unified_search", f"q={q}")
    return response


@router.get("/graph", response_model=GraphSearchResponse)
async def graph_search(
    request: Request,
    niu: Optional[str] = Query(None),
    phone: Optional[str] = Query(None),
    plate: Optional[str] = Query(None),
    depth: int = Query(2, ge=1, le=5),
    relationship_types: Optional[str] = Query(None),
    security: SecurityContext = Depends(verify_agency(["PNH", "DCPJ", "INTELLIGENCE", "SNISID_ADMIN"])),
):
    """Recherche graphique Neo4j"""
    from ai.national_search_engine import search_engine

    rel_types = relationship_types.split(",") if relationship_types else None

    result = await search_engine.graph_search(
        niu=niu, phone=phone, plate=plate,
        depth=depth, relationship_types=rel_types,
    )

    response = GraphSearchResponse(
        nodes=[GraphNode(**n) for n in result.get("nodes", [])],
        edges=[GraphEdge(**e) for e in result.get("edges", [])],
        centrality_score=result.get("centrality", 0),
        network_size=result.get("network_size", 0),
        detected_patterns=result.get("detected_patterns", []),
        analysis_summary=result.get("analysis_summary"),
    )

    create_audit_entry(request, security, "SEARCH", "graph_search",
                       f"niu={niu} phone={phone} plate={plate}")
    return response


@router.get("/cross-reference", response_model=dict)
async def cross_reference(
    request: Request,
    entity_id: str = Query(..., min_length=1),
    entity_type: str = Query(..., pattern=r'^(NIU|PHONE|PLATE|PASSPORT|VIN|EMAIL|IP)$'),
    depth: int = Query(2, ge=1, le=5),
    security: SecurityContext = Depends(verify_agency(["PNH", "DCPJ", "INTELLIGENCE", "SNISID_ADMIN"])),
):
    """Référence croisée à travers toutes les bases"""
    from ai.national_search_engine import search_engine

    result = await search_engine.cross_reference(entity_id, entity_type, depth)

    create_audit_entry(request, security, "CROSS_REFERENCE", entity_type, entity_id)
    return result
