from __future__ import annotations

import asyncio
import time
from unittest.mock import AsyncMock, patch

import pytest
import pytest_asyncio

from shared.cache import RedisCache, RateLimiter, SessionStore, cache_aside


class TestRedisCache:
    """Test multi-level cache get/set operations."""

    @pytest.mark.asyncio
    async def test_set_and_get(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("key1", {"name": "test"})
        value = await cache.get("key1")
        assert value == {"name": "test"}

    @pytest.mark.asyncio
    async def test_get_missing_key(self, redis_client):
        cache = RedisCache(redis_client)
        value = await cache.get("nonexistent")
        assert value is None

    @pytest.mark.asyncio
    async def test_get_with_default(self, redis_client):
        cache = RedisCache(redis_client)
        value = await cache.get("nonexistent", default="fallback")
        assert value == "fallback"

    @pytest.mark.asyncio
    async def test_set_with_custom_ttl(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("ttl_key", "value", ttl=10)
        ttl = await redis_client.ttl("cache:ttl_key")
        assert 0 < ttl <= 10

    @pytest.mark.asyncio
    async def test_delete(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("del_key", "value")
        await cache.delete("del_key")
        value = await cache.get("del_key")
        assert value is None

    @pytest.mark.asyncio
    async def test_delete_nonexistent(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.delete("nonexistent")
        assert True  # Should not raise

    @pytest.mark.asyncio
    async def test_overwrite_value(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("overwrite", "old")
        await cache.set("overwrite", "new")
        value = await cache.get("overwrite")
        assert value == "new"

    @pytest.mark.asyncio
    async def test_l1_hit(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("l1_test", "value")
        await redis_client.delete("cache:l1_test")
        value = await cache.get("l1_test")
        assert value == "value"

    @pytest.mark.asyncio
    async def test_l2_fallback_on_l1_miss(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("l2_test", "l2_value")
        cache._l1_delete("l2_test")
        value = await cache.get("l2_test")
        assert value == "l2_value"


class TestCacheStampedeProtection:
    """Test stampede protection in get_or_set."""

    @pytest.mark.asyncio
    async def test_get_or_set_miss(self, redis_client):
        cache = RedisCache(redis_client)
        factory_call_count = 0

        async def factory():
            nonlocal factory_call_count
            factory_call_count += 1
            return {"computed": 42}

        result = await cache.get_or_set("stampede", factory)
        assert result == {"computed": 42}
        assert factory_call_count == 1

    @pytest.mark.asyncio
    async def test_get_or_set_hit(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("hit_test", {"cached": True})
        factory_call_count = 0

        async def factory():
            nonlocal factory_call_count
            factory_call_count += 1
            return {"computed": 42}

        result = await cache.get_or_set("hit_test", factory)
        assert result == {"cached": True}
        assert factory_call_count == 0

    @pytest.mark.asyncio
    async def test_factory_only_called_once(self, redis_client):
        cache = RedisCache(redis_client)
        call_count = 0

        async def factory():
            nonlocal call_count
            call_count += 1
            await asyncio.sleep(0.05)
            return {"result": call_count}

        async def concurrent_call():
            return await cache.get_or_set("concurrent", factory)

        r1, r2 = await asyncio.gather(concurrent_call(), concurrent_call())
        assert r1 == {"result": 1}
        assert r2 == {"result": 1}
        assert call_count >= 1

    @pytest.mark.asyncio
    async def test_stampede_lock_released(self, redis_client):
        cache = RedisCache(redis_client)

        async def factory():
            return {"data": "new"}

        await cache.get_or_set("lock_release", factory)
        lock_exists = await redis_client.exists("cache:lock:lock_release")
        assert lock_exists == 0

    @pytest.mark.asyncio
    async def test_fallback_when_lock_holder_fails(self, redis_client):
        cache = RedisCache(redis_client)

        class FailingFactory:
            def __init__(self):
                self.call_count = 0

            async def __call__(self):
                self.call_count += 1
                if self.call_count == 1:
                    raise RuntimeError("Factory failed")
                return {"recovered": True}

        factory = FailingFactory()
        with pytest.raises(RuntimeError):
            await cache.get_or_set("failing", factory)


class TestCacheInvalidation:
    """Test cache invalidation by key and pattern."""

    @pytest.mark.asyncio
    async def test_invalidate_pattern(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("user:1", "data1")
        await cache.set("user:2", "data2")
        await cache.set("other:1", "data3")

        deleted = await cache.invalidate_pattern("user:*")
        assert deleted >= 2
        assert await cache.get("user:1") is None
        assert await cache.get("user:2") is None
        assert await cache.get("other:1") == "data3"

    @pytest.mark.asyncio
    async def test_invalidate_pattern_clears_l1(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("bulk:1", "l1_value")
        await cache.invalidate_pattern("bulk:*")
        assert cache._l1_get("bulk:1") is None


class TestCacheErrorHandling:
    """Test cache behavior when Redis is down."""

    @pytest.mark.asyncio
    async def test_get_returns_default_on_redis_error(self, redis_client):
        cache = RedisCache(redis_client)
        await redis_client.aclose()
        value = await cache.get("key", default="fallback")
        assert value == "fallback"

    @pytest.mark.asyncio
    async def test_set_does_not_raise_on_redis_error(self, redis_client):
        cache = RedisCache(redis_client)
        await redis_client.aclose()
        await cache.set("key", "value")
        assert cache._l1_get("key") == "value"

    @pytest.mark.asyncio
    async def test_delete_does_not_raise_on_redis_error(self, redis_client):
        cache = RedisCache(redis_client)
        await cache.set("key", "value")
        await redis_client.aclose()
        await cache.delete("key")

    @pytest.mark.asyncio
    async def test_get_or_set_fallback_on_redis_error(self, redis_client):
        cache = RedisCache(redis_client)

        async def factory():
            return {"computed": True}

        await redis_client.aclose()
        result = await cache.get_or_set("redis_down", factory)
        assert result == {"computed": True}


class TestTTLExpiration:
    """Test TTL-based cache expiration."""

    @pytest.mark.asyncio
    async def test_l1_expiration(self, redis_client):
        cache = RedisCache(redis_client, default_ttl=1)
        await cache.set("expire_key", "value")
        assert cache._l1_get("expire_key") == "value"
        time.sleep(1.1)
        assert cache._l1_get("expire_key") is None

    @pytest.mark.asyncio
    async def test_l2_expiration(self, redis_client):
        cache = RedisCache(redis_client, default_ttl=1)
        await cache.set("ttl_test", "value")
        ttl = await redis_client.ttl("cache:ttl_test")
        assert 0 < ttl <= 1


class TestCacheAsideDecorator:
    """Test cache_aside decorator."""

    @pytest.mark.asyncio
    async def test_cache_aside_hit(self, redis_client):
        cache = RedisCache(redis_client)
        call_count = 0

        @cache_aside(ttl=60, cache=cache)
        async def get_data(key: str):
            nonlocal call_count
            call_count += 1
            return {"key": key, "value": 42}

        r1 = await get_data("test")
        assert call_count == 1

        r2 = await get_data("test")
        assert call_count == 1
        assert r2 == {"key": "test", "value": 42}

    @pytest.mark.asyncio
    async def test_cache_aside_different_args(self, redis_client):
        cache = RedisCache(redis_client)
        call_count = 0

        @cache_aside(ttl=60, cache=cache)
        async def fetch(id: int):
            nonlocal call_count
            call_count += 1
            return {"id": id}

        await fetch(1)
        await fetch(2)
        assert call_count == 2

    @pytest.mark.asyncio
    async def test_cache_aside_no_cache_instance(self):
        @cache_aside(ttl=60)
        async def fetch_data():
            return {"result": "ok"}

        result = await fetch_data()
        assert result == {"result": "ok"}


class TestRateLimiter:
    """Test rate limiter functionality."""

    @pytest.mark.asyncio
    async def test_allowed(self, redis_client):
        limiter = RateLimiter(redis_client)
        result = await limiter.check("user:1", limit=5, window_seconds=60)
        assert result.allowed is True
        assert result.remaining == 4
        assert result.retry_after is None

    @pytest.mark.asyncio
    async def test_denied(self, redis_client):
        limiter = RateLimiter(redis_client)
        for _ in range(5):
            await limiter.check("user:burst", limit=5, window_seconds=60)
        result = await limiter.check("user:burst", limit=5, window_seconds=60)
        assert result.allowed is False
        assert result.remaining == 0
        assert result.retry_after is not None

    @pytest.mark.asyncio
    async def test_independent_keys(self, redis_client):
        limiter = RateLimiter(redis_client)
        r1 = await limiter.check("ip:1", limit=2, window_seconds=60)
        r2 = await limiter.check("ip:2", limit=2, window_seconds=60)
        assert r1.allowed
        assert r2.allowed

    @pytest.mark.asyncio
    async def test_fail_open_on_redis_error(self, redis_client):
        limiter = RateLimiter(redis_client)
        await redis_client.aclose()
        result = await limiter.check("user:err", limit=5, window_seconds=60)
        assert result.allowed is True


class TestSessionStore:
    """Test session store operations."""

    @pytest.mark.asyncio
    async def test_create_session(self, redis_client):
        store = SessionStore(redis_client)
        session_id = await store.create("user-1", data={"role": "admin"})
        assert session_id is not None
        assert len(session_id) > 0

    @pytest.mark.asyncio
    async def test_get_session(self, redis_client):
        store = SessionStore(redis_client)
        session_id = await store.create("user-1", data={"theme": "dark"})
        session = await store.get(session_id)
        assert session is not None
        assert session["user_id"] == "user-1"
        assert session["data"]["theme"] == "dark"

    @pytest.mark.asyncio
    async def test_get_nonexistent_session(self, redis_client):
        store = SessionStore(redis_client)
        session = await store.get("nonexistent")
        assert session is None

    @pytest.mark.asyncio
    async def test_delete_session(self, redis_client):
        store = SessionStore(redis_client)
        session_id = await store.create("user-1")
        await store.delete(session_id)
        session = await store.get(session_id)
        assert session is None

    @pytest.mark.asyncio
    async def test_extend_session(self, redis_client):
        store = SessionStore(redis_client)
        session_id = await store.create("user-1", ttl=60)
        result = await store.extend(session_id, ttl=120)
        assert result is True

    @pytest.mark.asyncio
    async def test_extend_nonexistent(self, redis_client):
        store = SessionStore(redis_client)
        result = await store.extend("nonexistent", ttl=60)
        assert result is False
