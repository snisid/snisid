"""SNI-SIDE API — Routes NCID (Criminal Intelligence)"""

from typing import Optional
from fastapi import APIRouter, Depends, HTTPException, Query, Request

from middleware.security import (
    verify_token, verify_agency, create_audit_entry, SecurityContext
)
from services.ncid import ncid_service
from models.schemas import (
    WantedPerson, WantedPersonCreate, Warrant, CriminalCase,
    Gang, CriminalOrganization, InterpolNotice,
)

router = APIRouter(prefix="/ncid", tags=["NCID — Criminal Intelligence"])

ALLOWED_AGENCIES = ["PNH", "DCPJ", "INTELLIGENCE", "SNISID_ADMIN"]


@router.get("/wanted-persons", response_model=dict)
async def search_wanted_persons(
    request: Request,
    name: Optional[str] = Query(None),
    alias: Optional[str] = Query(None),
    niu: Optional[str] = Query(None),
    nationality: Optional[str] = Query(None),
    risk_level: Optional[str] = Query(None, pattern=r'^(CRITICAL|HIGH|MEDIUM|LOW)$'),
    status: Optional[str] = Query(None, pattern=r'^(ACTIVE|CAPTURED|DECEASED|INACTIVE)$'),
    page: int = Query(1, ge=1),
    limit: int = Query(20, ge=1, le=100),
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Recherche de personnes recherchées"""
    persons, total = await ncid_service.search_wanted_persons(
        name, alias, niu, nationality, risk_level, status, page, limit
    )
    audit = create_audit_entry(request, security, "SEARCH", "wanted_persons",
                                f"query={name or niu or alias}")
    return {"data": [p.model_dump() for p in persons], "total": total, "page": page, "limit": limit}


@router.get("/wanted-persons/{niu}", response_model=WantedPerson)
async def get_wanted_person(
    request: Request,
    niu: str,
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Détails d'une personne recherchée"""
    person = await ncid_service.get_wanted_person(niu)
    if not person:
        raise HTTPException(status_code=404, detail="Wanted person not found")
    create_audit_entry(request, security, "READ", "wanted_persons", niu)
    return person


@router.post("/wanted-persons", response_model=WantedPerson, status_code=201)
async def create_wanted_person(
    request: Request,
    data: WantedPersonCreate,
    security: SecurityContext = Depends(verify_agency(["PNH", "DCPJ", "SNISID_ADMIN"])),
):
    """Crée une nouvelle personne recherchée"""
    person = await ncid_service.create_wanted_person(data)
    create_audit_entry(request, security, "CREATE", "wanted_persons", data.niu)
    return person


@router.get("/wanted-persons/{niu}/warrants", response_model=list[Warrant])
async def get_person_warrants(
    request: Request,
    niu: str,
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Mandats d'arrêt d'une personne"""
    warrants = await ncid_service.get_person_warrants(niu)
    create_audit_entry(request, security, "READ", "warrants", niu)
    return warrants


@router.get("/cases", response_model=dict)
async def search_criminal_cases(
    request: Request,
    case_number: Optional[str] = Query(None),
    case_type: Optional[str] = Query(None),
    status: Optional[str] = Query(None, pattern=r'^(OPEN|UNDER_INVESTIGATION|CLOSED|COLD)$'),
    agency: Optional[str] = Query(None),
    date_from: Optional[str] = Query(None),
    date_to: Optional[str] = Query(None),
    page: int = Query(1, ge=1),
    limit: int = Query(20, ge=1, le=100),
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Recherche de cas criminels"""
    from datetime import date
    d_from = date.fromisoformat(date_from) if date_from else None
    d_to = date.fromisoformat(date_to) if date_to else None
    cases, total = await ncid_service.search_criminal_cases(
        case_number, case_type, status, agency, d_from, d_to, page, limit
    )
    create_audit_entry(request, security, "SEARCH", "criminal_cases", f"type={case_type}")
    return {"data": [c.model_dump() for c in cases], "total": total, "page": page, "limit": limit}


@router.get("/cases/{case_id}", response_model=CriminalCase)
async def get_criminal_case(
    request: Request,
    case_id: str,
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Détails d'un cas criminel"""
    case = await ncid_service.get_criminal_case(case_id)
    if not case:
        raise HTTPException(status_code=404, detail="Case not found")
    create_audit_entry(request, security, "READ", "criminal_cases", case_id)
    return case


@router.get("/gangs", response_model=list[Gang])
async def search_gangs(
    request: Request,
    name: Optional[str] = Query(None),
    territory: Optional[str] = Query(None),
    risk_level: Optional[str] = Query(None),
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Recherche de gangs"""
    gangs = await ncid_service.search_gangs(name, territory, risk_level)
    create_audit_entry(request, security, "SEARCH", "gangs", "")
    return gangs


@router.get("/organizations", response_model=list[CriminalOrganization])
async def search_organizations(
    request: Request,
    name: Optional[str] = Query(None),
    org_type: Optional[str] = Query(None, alias="type"),
    country: Optional[str] = Query(None),
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Recherche d'organisations criminelles"""
    orgs = await ncid_service.search_organizations(name, org_type, country)
    create_audit_entry(request, security, "SEARCH", "organizations", "")
    return orgs


@router.get("/interpol-notices", response_model=list[InterpolNotice])
async def search_interpol_notices(
    request: Request,
    notice_type: Optional[str] = Query(None, pattern=r'^(RED|BLUE|GREEN|YELLOW|BLACK|ORANGE|PURPLE)$'),
    name: Optional[str] = Query(None),
    nationality: Optional[str] = Query(None),
    status: Optional[str] = Query(None, pattern=r'^(ACTIVE|WITHDRAWN|EXECUTED)$'),
    security: SecurityContext = Depends(verify_agency(ALLOWED_AGENCIES)),
):
    """Recherche de notices Interpol"""
    notices = await ncid_service.search_interpol_notices(notice_type, name, nationality, status)
    create_audit_entry(request, security, "SEARCH", "interpol_notices", "")
    return notices
