"""
SNISID Database — Async SQLAlchemy Engine & Session Management
===============================================================
Provides async database engine, session factory, base model, and
FastAPI dependency injection for transactional database access.
"""
from __future__ import annotations

import uuid
from contextlib import asynccontextmanager
from datetime import datetime, timezone
from typing import Any, AsyncGenerator

from sqlalchemy import MetaData, event, text
from sqlalchemy.ext.asyncio import (
    AsyncAttrs,
    AsyncEngine,
    AsyncSession,
    async_sessionmaker,
    create_async_engine,
)
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column
from sqlalchemy import DateTime, String, Boolean

from shared.config import get_settings
from shared.logging import get_logger

logger = get_logger(__name__)

# Naming convention for constraints (required for Alembic autogenerate)
NAMING_CONVENTION: dict[str, str] = {
    "ix": "ix_%(column_0_label)s",
    "uq": "uq_%(table_name)s_%(column_0_name)s",
    "ck": "ck_%(table_name)s_%(constraint_name)s",
    "fk": "fk_%(table_name)s_%(column_0_name)s_%(referred_table_name)s",
    "pk": "pk_%(table_name)s",
}


class Base(AsyncAttrs, DeclarativeBase):
    """
    Base declarative model for all SNISID entities.
    Provides automatic id, timestamps, and soft-delete support.
    """

    metadata = MetaData(naming_convention=NAMING_CONVENTION)

    id: Mapped[str] = mapped_column(
        String(36),
        primary_key=True,
        default=lambda: str(uuid.uuid4()),
        doc="UUID primary key",
    )
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: datetime.now(timezone.utc),
        nullable=False,
        doc="Record creation timestamp (UTC)",
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: datetime.now(timezone.utc),
        onupdate=lambda: datetime.now(timezone.utc),
        nullable=False,
        doc="Last update timestamp (UTC)",
    )
    is_deleted: Mapped[bool] = mapped_column(
        Boolean,
        default=False,
        nullable=False,
        index=True,
        doc="Soft delete flag",
    )
    deleted_at: Mapped[datetime | None] = mapped_column(
        DateTime(timezone=True),
        nullable=True,
        default=None,
        doc="Soft delete timestamp (UTC)",
    )

    def soft_delete(self) -> None:
        """Mark this record as soft-deleted."""
        self.is_deleted = True
        self.deleted_at = datetime.now(timezone.utc)

    def __repr__(self) -> str:
        return f"<{self.__class__.__name__}(id={self.id!r})>"


# ── Engine & Session Factory ──────────────────────────────────────────

_engine: AsyncEngine | None = None
_session_factory: async_sessionmaker[AsyncSession] | None = None


async def init_database(database_url: str | None = None) -> AsyncEngine:
    """
    Initialize the async database engine and session factory.

    Args:
        database_url: Override database URL. If None, uses settings.

    Returns:
        The created AsyncEngine.
    """
    global _engine, _session_factory

    if database_url is None:
        settings = get_settings()
        database_url = settings.database.async_url

    _engine = create_async_engine(
        database_url,
        pool_size=get_settings().database.pool_size,
        max_overflow=get_settings().database.max_overflow,
        pool_timeout=get_settings().database.pool_timeout,
        pool_recycle=get_settings().database.pool_recycle,
        echo=get_settings().database.echo,
        pool_pre_ping=True,
        json_serializer=_json_serializer,
        json_deserializer=_json_deserializer,
    )

    _session_factory = async_sessionmaker(
        bind=_engine,
        class_=AsyncSession,
        expire_on_commit=False,
        autoflush=False,
    )

    logger.info("database_initialized", url=_mask_password(database_url))
    return _engine


async def close_database() -> None:
    """Close the database engine and release connections."""
    global _engine, _session_factory
    if _engine is not None:
        await _engine.dispose()
        logger.info("database_closed")
        _engine = None
        _session_factory = None


@asynccontextmanager
async def get_session() -> AsyncGenerator[AsyncSession, None]:
    """
    Async context manager providing a transactional database session.
    Commits on success, rolls back on exception.
    """
    if _session_factory is None:
        raise RuntimeError("Database not initialized. Call init_database() first.")

    session = _session_factory()
    try:
        yield session
        await session.commit()
    except Exception:
        await session.rollback()
        raise
    finally:
        await session.close()


async def get_db_session() -> AsyncGenerator[AsyncSession, None]:
    """
    FastAPI dependency that provides a database session.
    Usage: db: AsyncSession = Depends(get_db_session)
    """
    async with get_session() as session:
        yield session


async def check_database_health() -> bool:
    """Check if the database connection is healthy."""
    if _engine is None:
        return False
    try:
        async with _engine.connect() as conn:
            await conn.execute(text("SELECT 1"))
        return True
    except Exception as e:
        logger.error("database_health_check_failed", error=str(e))
        return False


# ── Helpers ───────────────────────────────────────────────────────────

def _json_serializer(obj: Any) -> str:
    """Fast JSON serializer using orjson."""
    import orjson
    return orjson.dumps(obj).decode("utf-8")


def _json_deserializer(s: str) -> Any:
    """Fast JSON deserializer using orjson."""
    import orjson
    return orjson.loads(s)


def _mask_password(url: str) -> str:
    """Mask password in database URL for logging."""
    if "@" in url and ":" in url.split("@")[0]:
        parts = url.split("@")
        creds = parts[0].rsplit(":", 1)
        return f"{creds[0]}:****@{parts[1]}"
    return url
