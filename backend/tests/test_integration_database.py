from __future__ import annotations

import uuid
from datetime import date

import pytest
import pytest_asyncio
from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from services.agency.commands import AgencyCommandHandler, CreateAgencyCommand
from services.agency.models import Agency, AgencyStatus, AgencyType
from services.identity.models import Citizen, CitizenReadModel, IdentityStatus


# ── Agency Integration Tests ──────────────────────────────────────────────


@pytest_asyncio.fixture
async def agency_engine():
    engine = create_async_engine("sqlite+aiosqlite://", echo=False)
    async with engine.begin() as conn:
        await conn.run_sync(Agency.metadata.create_all)
    yield engine
    await engine.dispose()


@pytest_asyncio.fixture
async def agency_session(agency_engine):
    factory = async_sessionmaker(bind=agency_engine, class_=AsyncSession, expire_on_commit=False)
    async with factory() as sess:
        yield sess
        await sess.rollback()


@pytest.mark.asyncio
class TestAgencyIntegration:
    async def test_create_and_retrieve_agency(self, agency_session):
        handler = AgencyCommandHandler(agency_session)
        cmd = CreateAgencyCommand(
            name="Office National d'Identification",
            code="ONI",
            agency_type=AgencyType.CENTRAL,
            department="Interior",
            city="Port-au-Prince",
        )
        result = await handler.handle_create(cmd)
        agency_id = result["agency_id"]
        assert agency_id is not None

        stmt = select(Agency).where(Agency.id == agency_id)
        result = await agency_session.execute(stmt)
        agency = result.scalar_one_or_none()
        assert agency is not None
        assert agency.name == "Office National d'Identification"
        assert agency.code == "ONI"

    async def test_agency_default_status(self, agency_session):
        handler = AgencyCommandHandler(agency_session)
        cmd = CreateAgencyCommand(name="Test", code="TST", agency_type=AgencyType.LOCAL)
        create_result = await handler.handle_create(cmd)

        stmt = select(Agency).where(Agency.id == create_result["agency_id"])
        query_result = await agency_session.execute(stmt)
        agency = query_result.scalar_one_or_none()
        assert agency is not None
        assert agency.status == AgencyStatus.ACTIVE


# ── Identity Integration Tests ────────────────────────────────────────────


@pytest_asyncio.fixture
async def identity_engine():
    engine = create_async_engine("sqlite+aiosqlite://", echo=False)
    async with engine.begin() as conn:
        await conn.run_sync(Citizen.metadata.create_all)
        await conn.run_sync(CitizenReadModel.metadata.create_all)
    yield engine
    await engine.dispose()


@pytest_asyncio.fixture
async def identity_session(identity_engine):
    factory = async_sessionmaker(bind=identity_engine, class_=AsyncSession, expire_on_commit=False)
    async with factory() as sess:
        yield sess
        await sess.rollback()


@pytest_asyncio.fixture
async def seed_citizen(identity_session):
    citizen = Citizen(
        id=str(uuid.uuid4()),
        national_id="TEST-000001",
        first_name="Jean",
        last_name="Dupont",
        date_of_birth=date(1990, 1, 15),
        place_of_birth="Port-au-Prince",
        gender="male",
        nationality="HTI",
        email="jean.dupont@example.com",
        phone="+50912345678",
        status=IdentityStatus.ACTIVE,
        agency_id=str(uuid.uuid4()),
        created_by="system",
    )
    identity_session.add(citizen)
    await identity_session.flush()
    return citizen


@pytest_asyncio.fixture
async def seed_read_model(identity_session):
    agency_id = str(uuid.uuid4())
    entries = [
        CitizenReadModel(
            id=str(uuid.uuid4()), national_id="R-000001",
            full_name="Jean Dupont", first_name="Jean", last_name="Dupont",
            date_of_birth=date(1990, 1, 15), gender="male", nationality="HTI",
            status="active", agency_id=agency_id, document_count=1, has_biometrics=True, verified=True,
        ),
        CitizenReadModel(
            id=str(uuid.uuid4()), national_id="R-000002",
            full_name="Marie Durand", first_name="Marie", last_name="Durand",
            date_of_birth=date(1992, 5, 20), gender="female", nationality="HTI",
            status="pending", agency_id=agency_id, document_count=0, has_biometrics=False, verified=False,
        ),
    ]
    for e in entries:
        identity_session.add(e)
    await identity_session.flush()
    return entries


@pytest.mark.asyncio
class TestIdentityReadIntegration:
    async def test_get_citizen_by_id(self, identity_session, seed_citizen):
        stmt = select(Citizen).where(Citizen.id == seed_citizen.id)
        result = await identity_session.execute(stmt)
        citizen = result.scalar_one_or_none()
        assert citizen is not None
        assert citizen.national_id == "TEST-000001"
        assert citizen.first_name == "Jean"
        assert citizen.last_name == "Dupont"
        assert citizen.status == IdentityStatus.ACTIVE

    async def test_get_citizen_by_national_id(self, identity_session, seed_citizen):
        stmt = select(Citizen).where(Citizen.national_id == "TEST-000001")
        result = await identity_session.execute(stmt)
        citizen = result.scalar_one_or_none()
        assert citizen is not None
        assert citizen.id == seed_citizen.id

    async def test_search_read_model(self, identity_session, seed_read_model):
        stmt = (
            select(CitizenReadModel)
            .where(CitizenReadModel.national_id.ilike("%R-000001%"))
        )
        result = await identity_session.execute(stmt)
        items = result.scalars().all()
        assert len(items) == 1
        assert items[0].full_name == "Jean Dupont"

    async def test_stats(self, identity_session, seed_read_model):
        total = await identity_session.execute(select(func.count()).select_from(CitizenReadModel))
        assert total.scalar_one() == 2

        by_status = await identity_session.execute(
            select(CitizenReadModel.status, func.count())
            .group_by(CitizenReadModel.status)
        )
        status_map = {row[0]: row[1] for row in by_status.all()}
        assert status_map.get("active") == 1
        assert status_map.get("pending") == 1

    async def test_read_model_with_agency_filter(self, identity_session, seed_read_model):
        agency_id = seed_read_model[0].agency_id
        stmt = (
            select(func.count())
            .select_from(CitizenReadModel)
            .where(CitizenReadModel.agency_id == agency_id)
        )
        result = await identity_session.execute(stmt)
        assert result.scalar_one() == 2


# ── Projection Integration Tests ─────────────────────────────────────────


@pytest_asyncio.fixture
async def projection_engine():
    engine = create_async_engine("sqlite+aiosqlite://", echo=False)
    async with engine.begin() as conn:
        await conn.run_sync(CitizenReadModel.metadata.create_all)
    yield engine
    await engine.dispose()


@pytest_asyncio.fixture
async def projection_session(projection_engine):
    factory = async_sessionmaker(bind=projection_engine, class_=AsyncSession, expire_on_commit=False)
    async with factory() as sess:
        yield sess
        await sess.rollback()


@pytest.mark.asyncio
class TestCitizenProjectionIntegration:
    async def test_project_creates_read_model(self, projection_session):
        from services.identity.projections.citizen_projector import CitizenProjector
        projector = CitizenProjector(projection_session)
        aggregate_id = str(uuid.uuid4())

        await projector.project("IdentityCreated", aggregate_id, {
            "national_id": "PRJ-000001",
            "first_name": "Alice",
            "last_name": "Martin",
            "date_of_birth": "1985-06-12",
            "gender": "female",
            "nationality": "HTI",
            "agency_id": str(uuid.uuid4()),
        })

        stmt = select(CitizenReadModel).where(CitizenReadModel.id == aggregate_id)
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is not None
        assert model.full_name == "Alice Martin"
        assert model.status == "pending"
        assert model.verified is False
        assert model.document_count == 0
        assert model.has_biometrics is False

    async def test_project_verify_updates_status(self, projection_session):
        from services.identity.projections.citizen_projector import CitizenProjector
        projector = CitizenProjector(projection_session)
        aggregate_id = str(uuid.uuid4())

        await projector.project("IdentityCreated", aggregate_id, {
            "national_id": "PRJ-000002",
            "first_name": "Bob",
            "last_name": "Pierre",
            "date_of_birth": "1990-01-01",
            "gender": "male",
            "nationality": "HTI",
            "agency_id": str(uuid.uuid4()),
        })
        await projector.project("IdentityVerified", aggregate_id, {
            "verification_method": "biometric",
            "verifier_id": "verifier-1",
        })

        stmt = select(CitizenReadModel).where(CitizenReadModel.id == aggregate_id)
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is not None
        assert model.status == "active"
        assert model.verified is True

    async def test_project_suspend_and_revoke(self, projection_session):
        from services.identity.projections.citizen_projector import CitizenProjector
        projector = CitizenProjector(projection_session)
        aggregate_id = str(uuid.uuid4())

        await projector.project("IdentityCreated", aggregate_id, {
            "national_id": "PRJ-000003",
            "first_name": "Claire",
            "last_name": "Desir",
            "date_of_birth": "1995-03-15",
            "gender": "female",
            "nationality": "HTI",
            "agency_id": str(uuid.uuid4()),
        })

        await projector.project("IdentitySuspended", aggregate_id, {
            "reason": "Fraud investigation",
        })
        stmt = select(CitizenReadModel).where(CitizenReadModel.id == aggregate_id)
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is not None
        assert model.status == "suspended"

        await projector.project("IdentityRevoked", aggregate_id, {
            "reason": "Confirmed fraud",
        })
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is None

    async def test_project_biometric_enrolled_updates_flag(self, projection_session):
        from services.identity.projections.citizen_projector import CitizenProjector
        projector = CitizenProjector(projection_session)
        aggregate_id = str(uuid.uuid4())

        await projector.project("IdentityCreated", aggregate_id, {
            "national_id": "PRJ-000004",
            "first_name": "David",
            "last_name": "Lubin",
            "date_of_birth": "2000-07-22",
            "gender": "male",
            "nationality": "HTI",
            "agency_id": str(uuid.uuid4()),
        })
        await projector.project("BiometricEnrolled", aggregate_id, {
            "biometric_type": "fingerprint",
        })

        stmt = select(CitizenReadModel).where(CitizenReadModel.id == aggregate_id)
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is not None
        assert model.has_biometrics is True

    async def test_project_document_issued_increments_count(self, projection_session):
        from services.identity.projections.citizen_projector import CitizenProjector
        projector = CitizenProjector(projection_session)
        aggregate_id = str(uuid.uuid4())

        await projector.project("IdentityCreated", aggregate_id, {
            "national_id": "PRJ-000005",
            "first_name": "Emma",
            "last_name": "Jean",
            "date_of_birth": "1992-11-30",
            "gender": "female",
            "nationality": "HTI",
            "agency_id": str(uuid.uuid4()),
        })
        await projector.project("DocumentIssued", aggregate_id, {
            "document_type": "national_id",
            "document_number": "NID-000001",
        })

        stmt = select(CitizenReadModel).where(CitizenReadModel.id == aggregate_id)
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is not None
        assert model.document_count == 1

    async def test_project_unknown_event_is_ignored(self, projection_session):
        from services.identity.projections.citizen_projector import CitizenProjector
        projector = CitizenProjector(projection_session)
        aggregate_id = str(uuid.uuid4())

        await projector.project("UnknownEvent", aggregate_id, {})
        stmt = select(CitizenReadModel).where(CitizenReadModel.id == aggregate_id)
        result = await projection_session.execute(stmt)
        model = result.scalar_one_or_none()
        assert model is None
