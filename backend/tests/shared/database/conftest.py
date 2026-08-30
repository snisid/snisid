from __future__ import annotations

import pytest
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from shared.database import Base

_SQLITE_URL = "sqlite+aiosqlite://"

@pytest.fixture(scope="session")
def engine():
    engine = create_async_engine(_SQLITE_URL, echo=False)
    return engine


@pytest.fixture(scope="function")
async def db_session(engine) -> AsyncSession:
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    factory = async_sessionmaker(bind=engine, expire_on_commit=False)
    async with factory() as session:
        yield session
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
