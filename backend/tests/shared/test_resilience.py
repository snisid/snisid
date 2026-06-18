from __future__ import annotations

import asyncio
import time
from unittest.mock import AsyncMock

import pytest

from shared.resilience import (
    CircuitBreaker,
    CircuitOpenError,
    CircuitState,
    RetryPolicy,
    Bulkhead,
    BulkheadFullError,
    with_timeout,
    circuit_breaker,
    retry,
    bulkhead,
    timeout,
)


class TestCircuitBreaker:
    """Test circuit breaker states: closed, open, half-open."""

    @pytest.mark.asyncio
    async def test_closed_state_initial(self):
        cb = CircuitBreaker(name="test", failure_threshold=3, recovery_timeout=60)
        assert cb.state == CircuitState.CLOSED
        assert cb.failure_count == 0

    @pytest.mark.asyncio
    async def test_opens_after_threshold_failures(self):
        cb = CircuitBreaker(
            name="test",
            failure_threshold=3,
            recovery_timeout=60,
        )
        func = AsyncMock(side_effect=ValueError("fail"))

        for _ in range(3):
            with pytest.raises(ValueError):
                await cb.call(func)

        assert cb.state == CircuitState.OPEN
        assert cb.failure_count == 3

    @pytest.mark.asyncio
    async def test_open_raises_circuit_open_error(self):
        cb = CircuitBreaker(
            name="test",
            failure_threshold=1,
            recovery_timeout=60,
        )
        func = AsyncMock(side_effect=ValueError("fail"))

        with pytest.raises(ValueError):
            await cb.call(func)

        with pytest.raises(CircuitOpenError):
            await cb.call(AsyncMock(return_value="ok"))

    @pytest.mark.asyncio
    async def test_transition_to_half_open_after_recovery(self):
        cb = CircuitBreaker(
            name="test",
            failure_threshold=1,
            recovery_timeout=0.05,
        )

        func = AsyncMock(side_effect=ValueError("fail"))
        with pytest.raises(ValueError):
            await cb.call(func)

        assert cb.state == CircuitState.OPEN

        await asyncio.sleep(0.06)
        assert cb.state == CircuitState.HALF_OPEN

    @pytest.mark.asyncio
    async def test_half_open_success_closes_circuit(self):
        cb = CircuitBreaker(
            name="test",
            failure_threshold=1,
            recovery_timeout=0.05,
        )

        fail_func = AsyncMock(side_effect=ValueError("fail"))
        with pytest.raises(ValueError):
            await cb.call(fail_func)

        await asyncio.sleep(0.06)

        success_func = AsyncMock(return_value="recovered")
        result = await cb.call(success_func)
        assert result == "recovered"
        assert cb.state == CircuitState.CLOSED
        assert cb.failure_count == 0

    @pytest.mark.asyncio
    async def test_half_open_failure_reopens(self):
        cb = CircuitBreaker(
            name="test",
            failure_threshold=2,
            recovery_timeout=0.05,
        )

        fail_func = AsyncMock(side_effect=ValueError("fail"))
        with pytest.raises(ValueError):
            await cb.call(fail_func)
        with pytest.raises(ValueError):
            await cb.call(fail_func)

        assert cb.state == CircuitState.OPEN
        await asyncio.sleep(0.06)

        with pytest.raises(ValueError):
            await cb.call(fail_func)

        assert cb.state == CircuitState.OPEN

    @pytest.mark.asyncio
    async def test_success_resets_failure_count(self):
        cb = CircuitBreaker(name="test", failure_threshold=5)

        func = AsyncMock(side_effect=[ValueError("fail"), "ok"])
        with pytest.raises(ValueError):
            await cb.call(func)

        assert cb.failure_count == 1

        result = await cb.call(func)
        assert result == "ok"
        assert cb.failure_count == 0

    @pytest.mark.asyncio
    async def test_manual_reset(self):
        cb = CircuitBreaker(name="test", failure_threshold=1)
        func = AsyncMock(side_effect=ValueError("fail"))

        with pytest.raises(ValueError):
            await cb.call(func)

        assert cb.state == CircuitState.OPEN
        cb.reset()
        assert cb.state == CircuitState.CLOSED
        assert cb.failure_count == 0

    @pytest.mark.asyncio
    async def test_half_open_max_calls(self):
        cb = CircuitBreaker(
            name="test",
            failure_threshold=1,
            recovery_timeout=0.05,
            half_open_max_calls=1,
        )

        with pytest.raises(ValueError):
            await cb.call(AsyncMock(side_effect=ValueError("fail")))

        await asyncio.sleep(0.06)

        # Use an event to make the first probe hang so the second call
        # enters the HALF_OPEN gate before the probe completes
        probe_entered = asyncio.Event()
        probe_continue = asyncio.Event()

        async def slow_probe():
            probe_entered.set()
            await probe_continue.wait()
            return "ok"

        # Call 1 starts, enters HALF_OPEN gate, starts executing slow_probe
        t1 = asyncio.create_task(cb.call(slow_probe))
        await probe_entered.wait()

        # Call 2 tries to enter HALF_OPEN gate — should be rejected
        with pytest.raises(CircuitOpenError):
            await cb.call(AsyncMock(return_value="ok"))

        # Let probe finish
        probe_continue.set()
        result = await t1
        assert result == "ok"

    @pytest.mark.asyncio
    async def test_circuit_breaker_decorator(self):
        @circuit_breaker("decorator_test", failure_threshold=2, recovery_timeout=60)
        async def failing_func():
            raise ValueError("fail")

        with pytest.raises(ValueError):
            await failing_func()
        assert failing_func.circuit_breaker.failure_count == 1


class TestRetryPolicy:
    """Test retry with backoff."""

    @pytest.mark.asyncio
    async def test_success_without_retry(self):
        policy = RetryPolicy(max_retries=3)
        func = AsyncMock(return_value="success")
        result = await policy.execute(func)
        assert result == "success"
        func.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_retry_on_failure_then_success(self):
        policy = RetryPolicy(max_retries=3, base_delay=0.01)
        func = AsyncMock(side_effect=[ValueError("fail"), ValueError("fail"), "ok"])
        result = await policy.execute(func)
        assert result == "ok"
        assert func.await_count == 3

    @pytest.mark.asyncio
    async def test_exhaust_retries(self):
        policy = RetryPolicy(max_retries=2, base_delay=0.01)
        func = AsyncMock(side_effect=ValueError("persistent"))
        with pytest.raises(ValueError, match="persistent"):
            await policy.execute(func)
        assert func.await_count == 3

    @pytest.mark.asyncio
    async def test_no_retry_on_success(self):
        policy = RetryPolicy(max_retries=3, base_delay=0.01)
        func = AsyncMock(return_value="immediate")
        result = await policy.execute(func)
        assert result == "immediate"
        func.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_zero_retries(self):
        policy = RetryPolicy(max_retries=0)
        func = AsyncMock(side_effect=ValueError("fail"))
        with pytest.raises(ValueError):
            await policy.execute(func)
        func.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_retryable_exceptions_filter(self):
        policy = RetryPolicy(
            max_retries=2,
            base_delay=0.01,
            retryable_exceptions=(ValueError,),
        )
        func = AsyncMock(side_effect=TypeError("non-retryable"))
        with pytest.raises(TypeError):
            await policy.execute(func)
        func.assert_awaited_once()

    @pytest.mark.asyncio
    async def test_retry_decorator(self):
        call_count = 0

        @retry(max_retries=2, base_delay=0.01)
        async def flaky():
            nonlocal call_count
            call_count += 1
            if call_count < 3:
                raise ValueError("not yet")
            return "done"

        result = await flaky()
        assert result == "done"
        assert call_count == 3

    @pytest.mark.asyncio
    async def test_compute_delay(self):
        policy = RetryPolicy(
            max_retries=3,
            base_delay=1.0,
            exponential_base=2.0,
            max_delay=60.0,
            jitter=False,
        )
        assert policy._compute_delay(0) == 1.0
        assert policy._compute_delay(1) == 2.0
        assert policy._compute_delay(2) == 4.0

    @pytest.mark.asyncio
    async def test_max_delay_cap(self):
        policy = RetryPolicy(
            max_retries=10,
            base_delay=1.0,
            exponential_base=2.0,
            max_delay=10.0,
            jitter=False,
        )
        assert policy._compute_delay(10) == 10.0


class TestBulkhead:
    """Test bulkhead concurrency isolation."""

    @pytest.mark.asyncio
    async def test_execute_under_limit(self):
        bh = Bulkhead(name="test", max_concurrent=5, max_wait=5.0)
        func = AsyncMock(return_value="ok")
        result = await bh.execute(func)
        assert result == "ok"

    @pytest.mark.asyncio
    async def test_bulkhead_rejects_when_full(self):
        bh = Bulkhead(name="test", max_concurrent=1, max_wait=0.1)

        async def slow():
            await asyncio.sleep(1)

        task = asyncio.create_task(bh.execute(slow))
        await asyncio.sleep(0.05)

        with pytest.raises(BulkheadFullError):
            await bh.execute(AsyncMock(return_value="fast"))

        task.cancel()

    @pytest.mark.asyncio
    async def test_bulkhead_releases_slot(self):
        bh = Bulkhead(name="test", max_concurrent=1, max_wait=5.0)
        func = AsyncMock(return_value="ok")
        await bh.execute(func)
        result = await bh.execute(func)
        assert result == "ok"

    @pytest.mark.asyncio
    async def test_concurrent_count(self):
        bh = Bulkhead(name="test", max_concurrent=10)

        async def occupy():
            await asyncio.sleep(0.1)
            return "done"

        tasks = [asyncio.create_task(bh.execute(occupy)) for _ in range(3)]
        await asyncio.sleep(0.05)
        assert bh.active == 3
        assert bh.available == 7

        await asyncio.gather(*tasks)
        assert bh.active == 0
        assert bh.available == 10

    @pytest.mark.asyncio
    async def test_bulkhead_decorator(self):
        @bulkhead("decorator", max_concurrent=2, max_wait=5.0)
        async def limited():
            return "done"

        result = await limited()
        assert result == "done"

    @pytest.mark.asyncio
    async def test_bulkhead_full_error_message(self):
        bh = Bulkhead(name="pool", max_concurrent=1, max_wait=0.05)

        async def blocker():
            await asyncio.sleep(1)

        task = asyncio.create_task(bh.execute(blocker))
        await asyncio.sleep(0.1)

        with pytest.raises(BulkheadFullError) as exc_info:
            await bh.execute(AsyncMock(return_value="x"))
        assert "pool" in str(exc_info.value)
        assert "1 concurrent" in str(exc_info.value)

        task.cancel()


class TestTimeout:
    """Test timeout enforcement."""

    @pytest.mark.asyncio
    async def test_within_timeout(self):
        async def fast():
            await asyncio.sleep(0.01)
            return "done"

        result = await with_timeout(fast(), timeout_seconds=5)
        assert result == "done"

    @pytest.mark.asyncio
    async def test_exceeds_timeout(self):
        async def slow():
            await asyncio.sleep(10)

        with pytest.raises(asyncio.TimeoutError):
            await with_timeout(slow(), timeout_seconds=0.05)

    @pytest.mark.asyncio
    async def test_timeout_decorator(self):
        @timeout(seconds=0.05)
        async def slow_func():
            await asyncio.sleep(10)
            return "never"

        with pytest.raises(asyncio.TimeoutError):
            await slow_func()

    @pytest.mark.asyncio
    async def test_timeout_decorator_within_limit(self):
        @timeout(seconds=5)
        async def fast_func():
            return "fast"

        result = await fast_func()
        assert result == "fast"
