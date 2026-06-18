"""SNI-SIDE API — NCID Service (National Criminal Intelligence Database)"""

from typing import Optional, List, Tuple
from datetime import date, datetime
import json
import uuid

from database import db
from models.schemas import (
    WantedPerson, WantedPersonCreate, Warrant, CriminalCase,
    Gang, CriminalOrganization, InterpolNotice,
)


class NCIDService:
    """Service pour la base NCID — Criminal Intelligence"""

    # ============ WANTED PERSONS ============
    async def search_wanted_persons(
        self, name: Optional[str] = None, alias: Optional[str] = None,
        niu: Optional[str] = None, nationality: Optional[str] = None,
        risk_level: Optional[str] = None, status: Optional[str] = None,
        page: int = 1, limit: int = 20
    ) -> Tuple[List[WantedPerson], int]:
        """Recherche de personnes recherchées avec filtres"""
        conditions = []
        params = []
        param_idx = 1

        if name:
            conditions.append(f"full_name ILIKE ${param_idx}")
            params.append(f"%{name}%")
            param_idx += 1
        if alias:
            conditions.append(f"alias ILIKE ${param_idx}")
            params.append(f"%{alias}%")
            param_idx += 1
        if niu:
            conditions.append(f"niu = ${param_idx}")
            params.append(niu)
            param_idx += 1
        if nationality:
            conditions.append(f"nationality = ${param_idx}")
            params.append(nationality)
            param_idx += 1
        if risk_level:
            conditions.append(f"risk_level = ${param_idx}")
            params.append(risk_level)
            param_idx += 1
        if status:
            conditions.append(f"status = ${param_idx}")
            params.append(status)
            param_idx += 1

        where_clause = " AND ".join(conditions) if conditions else "TRUE"
        offset = (page - 1) * limit

        async with db.pg_conn() as conn:
            # Total count
            count_query = f"SELECT COUNT(*) FROM snisid_ncid.wanted_persons WHERE {where_clause}"
            total = await conn.fetchval(count_query, *params)

            # Results
            query = f"""
                SELECT niu, full_name, alias, date_of_birth, place_of_birth,
                       gender, nationality, height_cm, weight_kg, eye_color,
                       hair_color, skin_tone, scars_marks, last_known_address,
                       occupation, risk_level, status, photos, biometric_references,
                       created_at, updated_at
                FROM snisid_ncid.wanted_persons
                WHERE {where_clause}
                ORDER BY risk_level DESC, updated_at DESC
                LIMIT ${param_idx} OFFSET ${param_idx + 1}
            """
            params.extend([limit, offset])
            rows = await conn.fetch(query, *params)

            persons = [WantedPerson(**dict(r)) for r in rows]
            return persons, total

    async def get_wanted_person(self, niu: str) -> Optional[WantedPerson]:
        """Récupère une personne recherchée par NIU"""
        async with db.pg_conn() as conn:
            row = await conn.fetchrow(
                "SELECT * FROM snisid_ncid.wanted_persons WHERE niu = $1", niu
            )
            return WantedPerson(**dict(row)) if row else None

    async def create_wanted_person(self, data: WantedPersonCreate) -> WantedPerson:
        """Crée une nouvelle personne recherchée"""
        async with db.pg_conn() as conn:
            row = await conn.fetchrow(
                """INSERT INTO snisid_ncid.wanted_persons
                   (niu, full_name, alias, date_of_birth, place_of_birth,
                    gender, nationality, height_cm, weight_kg, eye_color,
                    hair_color, skin_tone, scars_marks, last_known_address,
                    occupation, risk_level, photos)
                   VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
                   RETURNING *""",
                data.niu, data.full_name, data.alias, data.date_of_birth,
                data.place_of_birth, data.gender, data.nationality,
                data.height_cm, data.weight_kg, data.eye_color,
                data.hair_color, data.skin_tone, data.scars_marks,
                data.last_known_address, data.occupation, data.risk_level,
                json.dumps(data.photos),
            )
            return WantedPerson(**dict(row))

    # ============ WARRANTS ============
    async def get_person_warrants(self, niu: str) -> List[Warrant]:
        """Récupère les mandats d'arrêt d'une personne"""
        async with db.pg_conn() as conn:
            rows = await conn.fetch(
                "SELECT * FROM snisid_ncid.arrest_warrants WHERE person_niu = $1 ORDER BY issued_date DESC",
                niu
            )
            return [Warrant(**dict(r)) for r in rows]

    # ============ CRIMINAL CASES ============
    async def search_criminal_cases(
        self, case_number: Optional[str] = None, case_type: Optional[str] = None,
        status: Optional[str] = None, agency: Optional[str] = None,
        date_from: Optional[date] = None, date_to: Optional[date] = None,
        page: int = 1, limit: int = 20
    ) -> Tuple[List[CriminalCase], int]:
        """Recherche de cas criminels"""
        conditions = []
        params = []
        idx = 1

        if case_number:
            conditions.append(f"case_number ILIKE ${idx}"); params.append(f"%{case_number}%"); idx += 1
        if case_type:
            conditions.append(f"case_type = ${idx}"); params.append(case_type); idx += 1
        if status:
            conditions.append(f"status = ${idx}"); params.append(status); idx += 1
        if agency:
            conditions.append(f"lead_agency ILIKE ${idx}"); params.append(f"%{agency}%"); idx += 1
        if date_from:
            conditions.append(f"incident_date >= ${idx}"); params.append(date_from); idx += 1
        if date_to:
            conditions.append(f"incident_date <= ${idx}"); params.append(date_to); idx += 1

        where = " AND ".join(conditions) if conditions else "TRUE"
        offset = (page - 1) * limit

        async with db.pg_conn() as conn:
            total = await conn.fetchval(
                f"SELECT COUNT(*) FROM snisid_ncid.criminal_cases WHERE {where}", *params
            )
            params.extend([limit, offset])
            rows = await conn.fetch(
                f"""SELECT * FROM snisid_ncid.criminal_cases
                    WHERE {where} ORDER BY incident_date DESC NULLS LAST
                    LIMIT ${idx} OFFSET ${idx + 1}""",
                *params
            )
            return [CriminalCase(**dict(r)) for r in rows], total

    async def get_criminal_case(self, case_id: str) -> Optional[CriminalCase]:
        """Récupère un cas criminel"""
        async with db.pg_conn() as conn:
            row = await conn.fetchrow(
                "SELECT * FROM snisid_ncid.criminal_cases WHERE case_id = $1", case_id
            )
            return CriminalCase(**dict(row)) if row else None

    # ============ GANGS ============
    async def search_gangs(self, name: Optional[str] = None,
                           territory: Optional[str] = None,
                           risk_level: Optional[str] = None) -> List[Gang]:
        conditions, params = [], []
        idx = 1
        if name:
            conditions.append(f"name ILIKE ${idx}"); params.append(f"%{name}%"); idx += 1
        if territory:
            conditions.append(f"territory ILIKE ${idx}"); params.append(f"%{territory}%"); idx += 1
        if risk_level:
            conditions.append(f"risk_level = ${idx}"); params.append(risk_level); idx += 1
        where = " AND ".join(conditions) if conditions else "TRUE"

        async with db.pg_conn() as conn:
            rows = await conn.fetch(
                f"SELECT * FROM snisid_ncid.gangs WHERE {where} ORDER BY risk_level DESC, name", *params
            )
            return [Gang(**dict(r)) for r in rows]

    # ============ ORGANIZATIONS ============
    async def search_organizations(self, name: Optional[str] = None,
                                   type_: Optional[str] = None,
                                   country: Optional[str] = None) -> List[CriminalOrganization]:
        conditions, params = [], []
        idx = 1
        if name:
            conditions.append(f"name ILIKE ${idx}"); params.append(f"%{name}%"); idx += 1
        if type_:
            conditions.append(f"type = ${idx}"); params.append(type_); idx += 1
        if country:
            conditions.append(f"geographic_reach ILIKE ${idx}"); params.append(f"%{country}%"); idx += 1
        where = " AND ".join(conditions) if conditions else "TRUE"

        async with db.pg_conn() as conn:
            rows = await conn.fetch(
                f"SELECT * FROM snisid_ncid.criminal_organizations WHERE {where} ORDER BY risk_level DESC", *params
            )
            return [CriminalOrganization(**dict(r)) for r in rows]

    # ============ INTERPOL NOTICES ============
    async def search_interpol_notices(self, notice_type: Optional[str] = None,
                                      name: Optional[str] = None,
                                      nationality: Optional[str] = None,
                                      status: Optional[str] = None) -> List[InterpolNotice]:
        conditions, params = [], []
        idx = 1
        if notice_type:
            conditions.append(f"notice_type = ${idx}"); params.append(notice_type); idx += 1
        if name:
            conditions.append(f"person_name ILIKE ${idx}"); params.append(f"%{name}%"); idx += 1
        if nationality:
            conditions.append(f"nationality = ${idx}"); params.append(nationality); idx += 1
        if status:
            conditions.append(f"status = ${idx}"); params.append(status); idx += 1
        where = " AND ".join(conditions) if conditions else "TRUE"

        async with db.pg_conn() as conn:
            rows = await conn.fetch(
                f"SELECT * FROM snisid_ncid.interpol_notices WHERE {where} ORDER BY issued_date DESC", *params
            )
            return [InterpolNotice(**dict(r)) for r in rows]


# Instance singleton
ncid_service = NCIDService()
