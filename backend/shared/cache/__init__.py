"""
SNISID Cache, Rate Limiting & Session Management
===================================================
Multi-level caching (L1 in-memory + L2 Redis), sliding-window rate
limiting, and Redis-backed session storage for the SNISID platform.

All operations are fully async and use ``orjson`` for high-performance
JSON serialisation.
"""
from __future__ import annotations

import asyncio
import hashlib
import time
import uuid
from dataclasses import dataclass, field
from functools import wraps
from typing import Any, Callable, Coroutine, ParamSpec, TypeVar

import orjson
import redis.asyncio as aioredis

from shared.config import get_settings
from shared.logging import get_logger

logger = get_logger(__name__)

P = ParamSpec("P")
T = TypeVar("T")

__all__ = [
    "RedisCache",
    "cache_aside",
    "RateLimiter",
    "RateLimitResult",
    "SessionStore",
]


# ---------------------------------------------------------------------------
# Internal L1 entry
# ---------------------------------------------------------------------------


@dataclass(slots=True)
class _L1Entry:
    """Single entry in the in-memory L1 cache."""

    value: Any
    expires_at: float


# ---------------------------------------------------------------------------
# RedisCache — multi-level cache
# ---------------------------------------------------------------------------


class RedisCache:
    """
    Two-level cache with stampede protection.

    * **L1** — process-local ``dict`` with per-key TTL.
    * **L2** — Redis, shared across all worker processes.

    A ``setnx``-based lock prevents concurrent callers from regenerating
    the same cache entry at the same time (stampede / thundering-herd
    protection).

    Args:
        redis_client: An ``redis.asyncio.Redis`` instance for L2.
        prefix: Key namespace prefix in Redis.
        default_ttl: Default time-to-live in seconds.
        stampede_ttl: How long the setnx lock is held (seconds).
        stampede_poll: Polling interval while waiting for the lock holder.
    """

    def __init__(
        self,
        redis_client: aioredis.Redis,
        *,
        prefix: str = "cache",
        default_ttl: int = 300,
        stampede_ttl: int = 30,
        stampede_poll: float = 0.1,
    ) -> None:
        self._redis = redis_client
        self._prefix = prefix
        self._default_ttl = default_ttl
        self._stampede_ttl = stampede_ttl
        self._stampede_poll = stampede_poll
        self._l1: dict[str, _L1Entry] = {}
        self._l1_lock = asyncio.Lock()

    # -- helpers --

    def _redis_key(self, key: str) -> str:
        return f"{self._prefix}:{key}"

    def _lock_key(self, key: str) -> str:
        return f"{self._prefix}:lock:{key}"

    def _l1_get(self, key: str) -> Any | None:
        """Return L1 value if present and not expired, else ``None``."""
        entry = self._l1.get(key)
        if entry is None:
            return None
        if time.monotonic() > entry.expires_at:
            self._l1.pop(key, None)
            return None
        return entry.value

    def _l1_set(self, key: str, value: Any, ttl: int) -> None:
        self._l1[key] = _L1Entry(value=value, expires_at=time.monotonic() + ttl)

    def _l1_delete(self, key: str) -> None:
        self._l1.pop(key, None)

    @staticmethod
    def _serialize(value: Any) -> bytes:
        return orjson.dumps(value)

    @staticmethod
    def _deserialize(raw: bytes | None) -> Any | None:
        if raw is None:
            return None
        return orjson.loads(raw)

    # -- public API --

    async def get(self, key: str, default: Any = None) -> Any:
        """
        Retrieve a value by key, checking L1 first then L2.

        Args:
            key: Cache key.
            default: Value returned when the key is absent.

        Returns:
            The cached value or *default*.
        """
        # L1
        value = self._l1_get(key)
        if value is not None:
            logger.debug("cache.hit", level="l1", key=key)
            return value

        # L2
        try:
            raw: bytes | None = await self._redis.get(self._redis_key(key))
        except aioredis.RedisError as exc:
            logger.warning("cache.redis_error", action="get", key=key, error=str(exc))
            return default

        if raw is None:
            logger.debug("cache.miss", key=key)
            return default

        value = self._deserialize(raw)
        # Populate L1 with the remaining TTL from Redis.
        ttl = await self._redis.ttl(self._redis_key(key))
        if ttl and ttl > 0:
            self._l1_set(key, value, ttl)
        logger.debug("cache.hit", level="l2", key=key)
        return value

    async def set(self, key: str, value: Any, ttl: int | None = None) -> None:
        """
        Store a value in both L1 and L2.

        Args:
            key: Cache key.
            value: JSON-serialisable value.
            ttl: Time-to-live in seconds (defaults to ``default_ttl``).
        """
        effective_ttl = ttl if ttl is not None else self._default_ttl

        # L1
        self._l1_set(key, value, effective_ttl)

        # L2
        try:
            await self._redis.set(
                self._redis_key(key),
                self._serialize(value),
                ex=effective_ttl,
            )
        except aioredis.RedisError as exc:
            logger.warning("cache.redis_error", action="set", key=key, error=str(exc))

        logger.debug("cache.set", key=key, ttl=effective_ttl)

    async def delete(self, key: str) -> None:
        """Remove a key from both L1 and L2."""
        self._l1_delete(key)

        try:
            await self._redis.delete(self._redis_key(key))
        except aioredis.RedisError as exc:
            logger.warning(
                "cache.redis_error", action="delete", key=key, error=str(exc)
            )

        logger.debug("cache.delete", key=key)

    async def invalidate_pattern(self, pattern: str) -> int:
        """
        Delete all Redis keys matching *pattern* using ``SCAN`` + ``DEL``.

        The L1 cache is flushed completely to avoid stale entries because
        the in-memory dict cannot be efficiently scanned by glob.

        Args:
            pattern: Redis glob pattern (e.g. ``"user:*"``).

        Returns:
            Number of keys deleted from Redis.
        """
        full_pattern = self._redis_key(pattern)
        deleted = 0

        try:
            cursor: int | bytes = 0
            while True:
                cursor, keys = await self._redis.scan(
                    cursor=cursor, match=full_pattern, count=200
                )
                if keys:
                    deleted += await self._redis.delete(*keys)
                if cursor == 0:
                    break
        except aioredis.RedisError as exc:
            logger.warning(
                "cache.redis_error",
                action="invalidate_pattern",
                pattern=pattern,
                error=str(exc),
            )

        # Flush L1 — no efficient glob on a plain dict.
        async with self._l1_lock:
            self._l1.clear()

        logger.info(
            "cache.invalidate_pattern", pattern=pattern, deleted_count=deleted
        )
        return deleted

    # -- stampede-protected fetch --

    async def get_or_set(
        self,
        key: str,
        factory: Callable[[], Coroutine[Any, Any, T]],
        ttl: int | None = None,
    ) -> T:
        """
        Fetch from cache or call *factory* exactly once (stampede-safe).

        Uses a Redis ``SETNX`` lock so that only one caller regenerates
        the value while all others wait for the result.

        Args:
            key: Cache key.
            factory: Async callable that produces the value on miss.
            ttl: Time-to-live in seconds.

        Returns:
            The cached (or freshly computed) value.
        """
        # Fast path — already cached
        value = await self.get(key)
        if value is not None:
            return value  # type: ignore[return-value]

        effective_ttl = ttl if ttl is not None else self._default_ttl
        lock_key = self._lock_key(key)

        try:
            acquired = await self._redis.set(
                lock_key, b"1", nx=True, ex=self._stampede_ttl
            )
        except aioredis.RedisError:
            # If Redis is down, just compute directly.
            acquired = True

        if acquired:
            try:
                value = await factory()
                await self.set(key, value, effective_ttl)
                return value  # type: ignore[return-value]
            finally:
                try:
                    await self._redis.delete(lock_key)
                except aioredis.RedisError:
                    pass
        else:
            # Wait for the lock holder to populate the cache.
            deadline = time.monotonic() + self._stampede_ttl
            while time.monotonic() < deadline:
                await asyncio.sleep(self._stampede_poll)
                value = await self.get(key)
                if value is not None:
                    return value  # type: ignore[return-value]

            # Fallback: lock holder may have failed; compute ourselves.
            value = await factory()
            await self.set(key, value, effective_ttl)
            return value  # type: ignore[return-value]


# ---------------------------------------------------------------------------
# cache_aside decorator factory
# ---------------------------------------------------------------------------


def _default_key_builder(func: Callable[..., Any], *args: Any, **kwargs: Any) -> str:
    """Build a deterministic cache key from function signature and arguments."""
    parts = [func.__module__, func.__qualname__]
    for arg in args:
        parts.append(repr(arg))
    for k in sorted(kwargs):
        parts.append(f"{k}={repr(kwargs[k])}")
    raw = ":".join(parts)
    return hashlib.sha256(raw.encode()).hexdigest()[:24]


def cache_aside(
    ttl: int = 300,
    key_builder: Callable[..., str] | None = None,
    cache: RedisCache | None = None,
) -> Callable[
    [Callable[P, Coroutine[Any, Any, T]]],
    Callable[P, Coroutine[Any, Any, T]],
]:
    """
    Decorator factory implementing the cache-aside pattern with stampede
    protection.

    The first call computes and caches; subsequent calls are served from
    cache until TTL expires.

    Args:
        ttl: Time-to-live for cached results (seconds).
        key_builder: Optional callable ``(func, *args, **kwargs) → str``.
        cache: Optional pre-configured :class:`RedisCache` instance.  If
            ``None``, a module-level default is lazily constructed from
            settings on first call.

    Usage::

        @cache_aside(ttl=60)
        async def get_citizen(nin: str) -> dict:
            ...
    """
    builder = key_builder or _default_key_builder

    def decorator(
        func: Callable[P, Coroutine[Any, Any, T]],
    ) -> Callable[P, Coroutine[Any, Any, T]]:
        _cache: RedisCache | None = cache

        @wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            nonlocal _cache
            if _cache is None:
                settings = get_settings()
                _redis = aioredis.from_url(
                    settings.redis.get_url(settings.redis.cache_db),
                    decode_responses=False,
                )
                _cache = RedisCache(_redis)

            key = builder(func, *args, **kwargs)
            return await _cache.get_or_set(
                key,
                lambda: func(*args, **kwargs),
                ttl=ttl,
            )

        return wrapper

    return decorator


# ---------------------------------------------------------------------------
# Rate Limiter — Sliding Window
# ---------------------------------------------------------------------------


@dataclass(frozen=True, slots=True)
class RateLimitResult:
    """Outcome of a rate-limit check."""

    allowed: bool
    remaining: int
    reset_at: int
    retry_after: int | None


class RateLimiter:
    """
    Redis-based sliding-window rate limiter.

    Uses a sorted set per key where members are unique request
    identifiers and scores are Unix-epoch timestamps.  Old members
    outside the current window are pruned on each check.

    Supports per-user, per-IP, and per-endpoint granularity by
    constructing the *key* argument appropriately (e.g.
    ``"user:{user_id}:POST:/api/v1/citizens"``).

    Args:
        redis_client: An ``redis.asyncio.Redis`` instance.
        prefix: Namespace prefix for Redis keys.
    """

    def __init__(
        self,
        redis_client: aioredis.Redis,
        prefix: str = "rl",
    ) -> None:
        self._redis = redis_client
        self._prefix = prefix

    def _redis_key(self, key: str) -> str:
        return f"{self._prefix}:{key}"

    async def check(
        self,
        key: str,
        limit: int,
        window_seconds: int,
    ) -> RateLimitResult:
        """
        Check whether the request identified by *key* is within the
        rate limit.

        This method is **atomic**: it prunes, counts, and conditionally
        adds within a Redis pipeline.

        Args:
            key: Composite key identifying the rate-limit bucket
                 (e.g. ``"ip:10.0.0.1:GET:/api/v1/health"``).
            limit: Maximum number of requests allowed in the window.
            window_seconds: Sliding window duration in seconds.

        Returns:
            A :class:`RateLimitResult` describing the outcome.
        """
        rk = self._redis_key(key)
        now = time.time()
        window_start = now - window_seconds
        member = f"{now}:{uuid.uuid4().hex[:8]}"

        try:
            async with self._redis.pipeline(transaction=True) as pipe:
                # Remove expired entries
                pipe.zremrangebyscore(rk, "-inf", window_start)
                # Count remaining entries
                pipe.zcard(rk)
                # Speculatively add the new entry
                pipe.zadd(rk, {member: now})
                # Set key expiry to auto-clean after the window
                pipe.expire(rk, window_seconds + 1)
                results = await pipe.execute()

            current_count: int = results[1]  # zcard before the new add

            reset_at = int(now) + window_seconds

            if current_count >= limit:
                # Over limit — remove the speculatively added member.
                await self._redis.zrem(rk, member)

                # Calculate retry_after from the oldest entry.
                oldest = await self._redis.zrange(rk, 0, 0, withscores=True)
                if oldest:
                    retry_after = int(oldest[0][1]) + window_seconds - int(now)
                    retry_after = max(retry_after, 1)
                else:
                    retry_after = window_seconds

                logger.warning(
                    "rate_limiter.denied",
                    key=key,
                    limit=limit,
                    window=window_seconds,
                )
                return RateLimitResult(
                    allowed=False,
                    remaining=0,
                    reset_at=reset_at,
                    retry_after=retry_after,
                )

            remaining = max(limit - current_count - 1, 0)
            return RateLimitResult(
                allowed=True,
                remaining=remaining,
                reset_at=reset_at,
                retry_after=None,
            )

        except aioredis.RedisError as exc:
            logger.error(
                "rate_limiter.redis_error", key=key, error=str(exc)
            )
            # Fail open: allow the request when Redis is unavailable.
            return RateLimitResult(
                allowed=True,
                remaining=limit,
                reset_at=int(time.time()) + window_seconds,
                retry_after=None,
            )


# ---------------------------------------------------------------------------
# Session Store
# ---------------------------------------------------------------------------


class SessionStore:
    """
    Redis-backed session store with sliding expiration.

    Session data is stored as JSON and keyed by a random session ID.
    Each access can optionally extend the TTL, implementing
    sliding-window session expiration.

    Args:
        redis_client: An ``redis.asyncio.Redis`` instance.
        prefix: Namespace prefix for session keys.
        default_ttl: Default session lifetime in seconds (1 hour).
    """

    def __init__(
        self,
        redis_client: aioredis.Redis,
        *,
        prefix: str = "session",
        default_ttl: int = 3600,
    ) -> None:
        self._redis = redis_client
        self._prefix = prefix
        self._default_ttl = default_ttl

    def _redis_key(self, session_id: str) -> str:
        return f"{self._prefix}:{session_id}"

    async def create(
        self,
        user_id: str,
        data: dict[str, Any] | None = None,
        ttl: int | None = None,
    ) -> str:
        """
        Create a new session.

        Args:
            user_id: The owning user's identifier.
            data: Arbitrary session payload (must be JSON-serialisable).
            ttl: Session lifetime in seconds (defaults to ``default_ttl``).

        Returns:
            A unique session ID.
        """
        session_id = uuid.uuid4().hex
        effective_ttl = ttl if ttl is not None else self._default_ttl

        payload: dict[str, Any] = {
            "user_id": user_id,
            "created_at": time.time(),
            "data": data or {},
        }

        try:
            await self._redis.set(
                self._redis_key(session_id),
                orjson.dumps(payload),
                ex=effective_ttl,
            )
        except aioredis.RedisError as exc:
            logger.error(
                "session.create_failed",
                user_id=user_id,
                error=str(exc),
            )
            raise

        logger.info(
            "session.created",
            session_id=session_id,
            user_id=user_id,
            ttl=effective_ttl,
        )
        return session_id

    async def get(self, session_id: str) -> dict[str, Any] | None:
        """
        Retrieve session data.

        Args:
            session_id: The session identifier.

        Returns:
            Session payload dict or ``None`` if the session does not exist
            or has expired.
        """
        try:
            raw: bytes | None = await self._redis.get(
                self._redis_key(session_id)
            )
        except aioredis.RedisError as exc:
            logger.error(
                "session.get_failed",
                session_id=session_id,
                error=str(exc),
            )
            return None

        if raw is None:
            logger.debug("session.not_found", session_id=session_id)
            return None

        return orjson.loads(raw)  # type: ignore[no-any-return]

    async def delete(self, session_id: str) -> None:
        """
        Delete (invalidate) a session.

        Args:
            session_id: The session identifier.
        """
        try:
            await self._redis.delete(self._redis_key(session_id))
        except aioredis.RedisError as exc:
            logger.error(
                "session.delete_failed",
                session_id=session_id,
                error=str(exc),
            )
            raise

        logger.info("session.deleted", session_id=session_id)

    async def extend(self, session_id: str, ttl: int | None = None) -> bool:
        """
        Extend (slide) the session expiration.

        Args:
            session_id: The session identifier.
            ttl: New TTL in seconds (defaults to ``default_ttl``).

        Returns:
            ``True`` if the session exists and was extended, ``False`` otherwise.
        """
        effective_ttl = ttl if ttl is not None else self._default_ttl
        rk = self._redis_key(session_id)

        try:
            exists = await self._redis.exists(rk)
            if not exists:
                logger.debug(
                    "session.extend_miss", session_id=session_id
                )
                return False

            await self._redis.expire(rk, effective_ttl)
        except aioredis.RedisError as exc:
            logger.error(
                "session.extend_failed",
                session_id=session_id,
                error=str(exc),
            )
            return False

        logger.debug(
            "session.extended",
            session_id=session_id,
            ttl=effective_ttl,
        )
        return True
