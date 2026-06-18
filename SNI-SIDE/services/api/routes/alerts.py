"""SNI-SIDE API — Routes Alertes"""

from typing import Optional
from fastapi import APIRouter, Depends, HTTPException, Query, Request
from datetime import datetime

from middleware.security import verify_token, create_audit_entry, SecurityContext
from models.schemas import Alert, AcknowledgeAlertResponse, ResolveAlertRequest

router = APIRouter(prefix="/alerts", tags=["National Alerts"])


@router.get("/", response_model=dict)
async def get_alerts(
    request: Request,
    source: Optional[str] = Query(None),
    severity: Optional[str] = Query(None, pattern=r'^(CRITICAL|HIGH|MEDIUM|LOW)$'),
    status: Optional[str] = Query(None, pattern=r'^(NEW|ACKNOWLEDGED|INVESTIGATING|RESOLVED)$'),
    date_from: Optional[str] = Query(None),
    date_to: Optional[str] = Query(None),
    page: int = Query(1, ge=1),
    limit: int = Query(50, ge=1, le=200),
    security: SecurityContext = Depends(verify_token),
):
    """Récupère les alertes avec filtres"""
    # Query PostgreSQL or ClickHouse for alerts
    conditions = []
    params = []
    idx = 1

    if source:
        conditions.append(f"source = ${idx}"); params.append(source); idx += 1
    if severity:
        conditions.append(f"severity = ${idx}"); params.append(severity); idx += 1
    if status:
        conditions.append(f"status = ${idx}"); params.append(status); idx += 1
    if date_from:
        conditions.append(f"created_at >= ${idx}"); params.append(datetime.fromisoformat(date_from)); idx += 1
    if date_to:
        conditions.append(f"created_at <= ${idx}"); params.append(datetime.fromisoformat(date_to)); idx += 1

    where = " AND ".join(conditions) if conditions else "TRUE"
    offset = (page - 1) * limit

    from database import db
    async with db.pg_conn() as conn:
        # Use unified alerts table from sniside schema
        try:
            total = await conn.fetchval(
                f"SELECT COUNT(*) FROM sniside_national.alerts WHERE {where}", *params
            )
            params.extend([limit, offset])
            rows = await conn.fetch(
                f"""SELECT * FROM sniside_national.alerts
                    WHERE {where} ORDER BY created_at DESC
                    LIMIT ${idx} OFFSET ${idx + 1}""",
                *params
            )
        except Exception:
            # Fallback: return empty if table doesn't exist yet
            total = 0
            rows = []

        alerts = [Alert(**dict(r)) for r in rows]

    create_audit_entry(request, security, "SEARCH", "alerts",
                       f"severity={severity} status={status}")
    return {"data": [a.model_dump() for a in alerts], "total": total, "page": page, "limit": limit}


@router.post("/{alert_id}/acknowledge", response_model=AcknowledgeAlertResponse)
async def acknowledge_alert(
    request: Request,
    alert_id: str,
    security: SecurityContext = Depends(verify_token),
):
    """Accuse réception d'une alerte"""
    from database import db
    async with db.pg_conn() as conn:
        result = await conn.execute(
            """UPDATE sniside_national.alerts
               SET status = 'ACKNOWLEDGED',
                   acknowledged_at = NOW(),
                   acknowledged_by = $1
               WHERE alert_id = $2 AND status = 'NEW'""",
            security.user_id, alert_id
        )
        updated = result.split()[-1] != "0"

    if not updated:
        raise HTTPException(status_code=404, detail="Alert not found or already acknowledged")

    create_audit_entry(request, security, "ACKNOWLEDGE", "alerts", alert_id)
    return AcknowledgeAlertResponse(success=True, alert_id=alert_id)


@router.post("/{alert_id}/resolve", response_model=AcknowledgeAlertResponse)
async def resolve_alert(
    request: Request,
    alert_id: str,
    data: ResolveAlertRequest,
    security: SecurityContext = Depends(verify_token),
):
    """Résout une alerte"""
    from database import db
    async with db.pg_conn() as conn:
        result = await conn.execute(
            """UPDATE sniside_national.alerts
               SET status = 'RESOLVED',
                   resolution = $1,
                   updated_at = NOW()
               WHERE alert_id = $2 AND status IN ('ACKNOWLEDGED', 'INVESTIGATING')""",
            data.resolution, alert_id
        )
        updated = result.split()[-1] != "0"

    if not updated:
        raise HTTPException(status_code=404, detail="Alert not found or cannot be resolved")

    create_audit_entry(request, security, "RESOLVE", "alerts", alert_id)
    return AcknowledgeAlertResponse(success=True, alert_id=alert_id)
