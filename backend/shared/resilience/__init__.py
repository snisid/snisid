"""
SNISID Resilience Patterns
============================
Production-grade resilience primitives for distributed microservice
communication: circuit breaker, retry with backoff, bulkhead isolation,
and timeout enforcement.

All patterns are fully async, emit structured logs on state transitions,
and expose both class-based and decorator-factory APIs.
"""
from __future__ import annotations

import asyncio
import enum
import functools
import random
import time
from dataclasses import dataclass, field
from typing import Any, Callable, Concatenate, Coroutine, ParamSpec, TypeVar

from shared.logging import get_logger

logger = get_logger(__name__)

P = ParamSpec("P")
T = TypeVar("T")

__all__ = [
    "CircuitState",
    "CircuitBreaker",
    "circuit_breaker",
    "RetryPolicy",
    "retry",
    "Bulkhead",
    "BulkheadFullError",
    "bulkhead",
    "with_timeout",
    "timeout",
]


# ---------------------------------------------------------------------------
# Exceptions
# ---------------------------------------------------------------------------


class CircuitOpenError(Exception):
    """Raised when the circuit breaker is in the OPEN state."""

    def __init__(self, name: str, recovery_in: float) -> None:
        self.name = name
        self.recovery_in = recovery_in
        super().__init__(
            f"Circuit '{name}' is OPEN. Recovery attempt in {recovery_in:.1f}s."
        )


class BulkheadFullError(Exception):
    """Raised when the bulkhead concurrency limit is reached and the wait times out."""

    def __init__(self, name: str, max_concurrent: int, max_wait: float) -> None:
        self.name = name
        self.max_concurrent = max_concurrent
        self.max_wait = max_wait
        super().__init__(
            f"Bulkhead '{name}' full ({max_concurrent} concurrent). "
            f"Timed out after {max_wait:.1f}s."
        )


# ---------------------------------------------------------------------------
# Circuit Breaker
# ---------------------------------------------------------------------------


class CircuitState(enum.Enum):
    """States for the circuit breaker finite-state machine."""

    CLOSED = "closed"
    OPEN = "open"
    HALF_OPEN = "half_open"


@dataclass
class CircuitBreaker:
    """
    Three-state circuit breaker protecting downstream calls.

    State machine:
        CLOSED → OPEN       when ``failure_count`` reaches ``failure_threshold``
        OPEN → HALF_OPEN    after ``recovery_timeout`` seconds elapse
        HALF_OPEN → CLOSED  when a trial call succeeds
        HALF_OPEN → OPEN    when a trial call fails

    Args:
        name: Human-readable identifier for logging / metrics.
        failure_threshold: Consecutive failures before opening the circuit.
        recovery_timeout: Seconds to wait in OPEN before probing.
        half_open_max_calls: Max concurrent trial calls in HALF_OPEN.
    """

    name: str = "default"
    failure_threshold: int = 5
    recovery_timeout: float = 30.0
    half_open_max_calls: int = 3

    # ---- internal mutable state (not part of __init__ signature) ----
    _state: CircuitState = field(default=CircuitState.CLOSED, init=False, repr=False)
    _failure_count: int = field(default=0, init=False, repr=False)
    _half_open_calls: int = field(default=0, init=False, repr=False)
    _last_failure_time: float = field(default=0.0, init=False, repr=False)
    _lock: asyncio.Lock = field(default_factory=asyncio.Lock, init=False, repr=False)

    # -- public properties --

    @property
    def state(self) -> CircuitState:
        """Return the current (possibly auto-transitioned) state."""
        if self._state is CircuitState.OPEN and self._recovery_elapsed():
            # Don't mutate here; the transition is performed inside `call`.
            return CircuitState.HALF_OPEN
        return self._state

    @property
    def failure_count(self) -> int:
        return self._failure_count

    @property
    def last_failure_time(self) -> float:
        return self._last_failure_time

    # -- helpers --

    def _recovery_elapsed(self) -> bool:
        return (time.monotonic() - self._last_failure_time) >= self.recovery_timeout

    def _transition(self, new_state: CircuitState) -> None:
        old_state = self._state
        self._state = new_state
        logger.info(
            "circuit_breaker.state_transition",
            circuit=self.name,
            from_state=old_state.value,
            to_state=new_state.value,
            failure_count=self._failure_count,
        )

    # -- core API --

    async def call(
        self,
        func: Callable[..., Coroutine[Any, Any, T]],
        *args: Any,
        **kwargs: Any,
    ) -> T:
        """
        Execute *func* with circuit breaker protection.

        Raises:
            CircuitOpenError: When the circuit is OPEN and the recovery
                timeout has not yet elapsed.
        """
        async with self._lock:
            # --- OPEN → HALF_OPEN probe ---
            if self._state is CircuitState.OPEN:
                if self._recovery_elapsed():
                    self._half_open_calls = 0
                    self._transition(CircuitState.HALF_OPEN)
                else:
                    remaining = self.recovery_timeout - (
                        time.monotonic() - self._last_failure_time
                    )
                    raise CircuitOpenError(self.name, recovery_in=max(remaining, 0.0))

            # --- HALF_OPEN gate ---
            if self._state is CircuitState.HALF_OPEN:
                if self._half_open_calls >= self.half_open_max_calls:
                    raise CircuitOpenError(
                        self.name,
                        recovery_in=0.0,
                    )
                self._half_open_calls += 1

        # Execute outside the lock so calls are not serialised.
        try:
            result: T = await func(*args, **kwargs)
        except Exception:
            await self._on_failure()
            raise
        else:
            await self._on_success()
            return result

    async def _on_success(self) -> None:
        async with self._lock:
            if self._state is CircuitState.HALF_OPEN:
                self._failure_count = 0
                self._transition(CircuitState.CLOSED)
            elif self._state is CircuitState.CLOSED:
                # Reset consecutive failure counter on any success.
                self._failure_count = 0

    async def _on_failure(self) -> None:
        async with self._lock:
            self._failure_count += 1
            self._last_failure_time = time.monotonic()

            if self._state is CircuitState.HALF_OPEN:
                self._transition(CircuitState.OPEN)
            elif (
                self._state is CircuitState.CLOSED
                and self._failure_count >= self.failure_threshold
            ):
                self._transition(CircuitState.OPEN)

            logger.warning(
                "circuit_breaker.failure_recorded",
                circuit=self.name,
                failure_count=self._failure_count,
                state=self._state.value,
            )

    def reset(self) -> None:
        """Manually reset the circuit breaker to CLOSED."""
        self._state = CircuitState.CLOSED
        self._failure_count = 0
        self._half_open_calls = 0
        self._last_failure_time = 0.0
        logger.info("circuit_breaker.manual_reset", circuit=self.name)


def circuit_breaker(
    name: str = "default",
    *,
    failure_threshold: int = 5,
    recovery_timeout: float = 30.0,
    half_open_max_calls: int = 3,
) -> Callable[
    [Callable[P, Coroutine[Any, Any, T]]],
    Callable[P, Coroutine[Any, Any, T]],
]:
    """
    Decorator factory that wraps an async function with a :class:`CircuitBreaker`.

    Usage::

        @circuit_breaker("payment-gateway", failure_threshold=3)
        async def charge(amount: Decimal) -> Receipt:
            ...
    """
    cb = CircuitBreaker(
        name=name,
        failure_threshold=failure_threshold,
        recovery_timeout=recovery_timeout,
        half_open_max_calls=half_open_max_calls,
    )

    def decorator(
        func: Callable[P, Coroutine[Any, Any, T]],
    ) -> Callable[P, Coroutine[Any, Any, T]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            return await cb.call(func, *args, **kwargs)

        # Expose the breaker instance for introspection / testing.
        wrapper.circuit_breaker = cb  # type: ignore[attr-defined]
        return wrapper

    return decorator


# ---------------------------------------------------------------------------
# Retry Policy
# ---------------------------------------------------------------------------


@dataclass
class RetryPolicy:
    """
    Configurable retry policy with exponential back-off and optional jitter.

    Args:
        max_retries: Maximum number of retry attempts (0 = no retries).
        base_delay: Initial delay in seconds before the first retry.
        max_delay: Upper cap on the computed delay.
        exponential_base: Multiplicative base for exponential growth.
        jitter: If ``True``, adds uniform random jitter ``[0, delay)``.
        retryable_exceptions: Tuple of exception classes eligible for retry.
            Defaults to ``(Exception,)`` (retry on any error).
    """

    max_retries: int = 3
    base_delay: float = 1.0
    max_delay: float = 60.0
    exponential_base: float = 2.0
    jitter: bool = True
    retryable_exceptions: tuple[type[Exception], ...] = (Exception,)

    def _compute_delay(self, attempt: int) -> float:
        """Compute the sleep duration for a given attempt (0-indexed)."""
        delay = self.base_delay * (self.exponential_base ** attempt)
        delay = min(delay, self.max_delay)
        if self.jitter:
            delay = random.uniform(0.0, delay)  # noqa: S311
        return delay

    async def execute(
        self,
        func: Callable[..., Coroutine[Any, Any, T]],
        *args: Any,
        **kwargs: Any,
    ) -> T:
        """
        Execute *func* with the configured retry policy.

        Returns:
            The return value of *func* on success.

        Raises:
            The last exception raised by *func* after all retries are exhausted.
        """
        last_exc: Exception | None = None

        for attempt in range(self.max_retries + 1):
            try:
                return await func(*args, **kwargs)
            except self.retryable_exceptions as exc:
                last_exc = exc
                if attempt >= self.max_retries:
                    logger.error(
                        "retry_policy.exhausted",
                        function=getattr(func, "__qualname__", str(func)),
                        attempts=attempt + 1,
                        error=str(exc),
                    )
                    raise

                delay = self._compute_delay(attempt)
                logger.warning(
                    "retry_policy.retrying",
                    function=getattr(func, "__qualname__", str(func)),
                    attempt=attempt + 1,
                    max_retries=self.max_retries,
                    delay_seconds=round(delay, 3),
                    error=str(exc),
                )
                await asyncio.sleep(delay)
            except Exception:
                # Non-retryable exception — fail immediately.
                raise

        # Should be unreachable, but satisfy the type checker.
        assert last_exc is not None  # noqa: S101
        raise last_exc


def retry(
    max_retries: int = 3,
    *,
    base_delay: float = 1.0,
    max_delay: float = 60.0,
    exponential_base: float = 2.0,
    jitter: bool = True,
    retryable_exceptions: tuple[type[Exception], ...] = (Exception,),
) -> Callable[
    [Callable[P, Coroutine[Any, Any, T]]],
    Callable[P, Coroutine[Any, Any, T]],
]:
    """
    Decorator factory that wraps an async function with a :class:`RetryPolicy`.

    Usage::

        @retry(max_retries=5, retryable_exceptions=(httpx.TransportError,))
        async def fetch_remote(url: str) -> bytes:
            ...
    """
    policy = RetryPolicy(
        max_retries=max_retries,
        base_delay=base_delay,
        max_delay=max_delay,
        exponential_base=exponential_base,
        jitter=jitter,
        retryable_exceptions=retryable_exceptions,
    )

    def decorator(
        func: Callable[P, Coroutine[Any, Any, T]],
    ) -> Callable[P, Coroutine[Any, Any, T]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            return await policy.execute(func, *args, **kwargs)

        wrapper.retry_policy = policy  # type: ignore[attr-defined]
        return wrapper

    return decorator


# ---------------------------------------------------------------------------
# Bulkhead (Semaphore-based Concurrency Limiter)
# ---------------------------------------------------------------------------


@dataclass
class Bulkhead:
    """
    Semaphore-based concurrency limiter (bulkhead pattern).

    Limits how many concurrent invocations of a protected resource can
    execute simultaneously.  Callers that cannot acquire a slot within
    ``max_wait`` seconds receive a :class:`BulkheadFullError`.

    Args:
        name: Human-readable identifier for logging / metrics.
        max_concurrent: Maximum number of concurrent executions.
        max_wait: Maximum seconds to wait for a slot before raising.
    """

    name: str = "default"
    max_concurrent: int = 10
    max_wait: float = 5.0

    _semaphore: asyncio.Semaphore = field(init=False, repr=False)

    def __post_init__(self) -> None:
        self._semaphore = asyncio.Semaphore(self.max_concurrent)

    async def execute(
        self,
        func: Callable[..., Coroutine[Any, Any, T]],
        *args: Any,
        **kwargs: Any,
    ) -> T:
        """
        Execute *func* within the concurrency limit.

        Raises:
            BulkheadFullError: If a slot cannot be acquired within ``max_wait``.
        """
        try:
            acquired = await asyncio.wait_for(
                self._acquire(),
                timeout=self.max_wait,
            )
        except asyncio.TimeoutError:
            logger.warning(
                "bulkhead.full",
                bulkhead=self.name,
                max_concurrent=self.max_concurrent,
                max_wait=self.max_wait,
            )
            raise BulkheadFullError(self.name, self.max_concurrent, self.max_wait)

        try:
            return await func(*args, **kwargs)
        finally:
            self._semaphore.release()

    async def _acquire(self) -> bool:
        """Coroutine wrapper around ``Semaphore.acquire`` for ``wait_for``."""
        await self._semaphore.acquire()
        return True

    @property
    def active(self) -> int:
        """Number of slots currently in use."""
        return self.max_concurrent - self._semaphore._value  # noqa: SLF001

    @property
    def available(self) -> int:
        """Number of slots still available."""
        return self._semaphore._value  # noqa: SLF001


def bulkhead(
    name: str = "default",
    max_concurrent: int = 10,
    max_wait: float = 5.0,
) -> Callable[
    [Callable[P, Coroutine[Any, Any, T]]],
    Callable[P, Coroutine[Any, Any, T]],
]:
    """
    Decorator factory that wraps an async function with a :class:`Bulkhead`.

    Usage::

        @bulkhead("db-pool", max_concurrent=20, max_wait=3.0)
        async def query_db(sql: str) -> list[Row]:
            ...
    """
    bh = Bulkhead(name=name, max_concurrent=max_concurrent, max_wait=max_wait)

    def decorator(
        func: Callable[P, Coroutine[Any, Any, T]],
    ) -> Callable[P, Coroutine[Any, Any, T]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            return await bh.execute(func, *args, **kwargs)

        wrapper.bulkhead = bh  # type: ignore[attr-defined]
        return wrapper

    return decorator


# ---------------------------------------------------------------------------
# Timeout Helpers
# ---------------------------------------------------------------------------


async def with_timeout(
    coro: Coroutine[Any, Any, T],
    timeout_seconds: float,
) -> T:
    """
    Execute a coroutine with a strict timeout.

    Args:
        coro: The coroutine to execute.
        timeout_seconds: Maximum wall-clock seconds to wait.

    Returns:
        The coroutine's return value.

    Raises:
        asyncio.TimeoutError: If the coroutine does not complete in time.
    """
    try:
        return await asyncio.wait_for(coro, timeout=timeout_seconds)
    except asyncio.TimeoutError:
        logger.warning(
            "timeout.exceeded",
            timeout_seconds=timeout_seconds,
        )
        raise


def timeout(
    seconds: float,
) -> Callable[
    [Callable[P, Coroutine[Any, Any, T]]],
    Callable[P, Coroutine[Any, Any, T]],
]:
    """
    Decorator that enforces a wall-clock timeout on an async function.

    Usage::

        @timeout(10.0)
        async def slow_operation() -> Result:
            ...
    """

    def decorator(
        func: Callable[P, Coroutine[Any, Any, T]],
    ) -> Callable[P, Coroutine[Any, Any, T]]:
        @functools.wraps(func)
        async def wrapper(*args: P.args, **kwargs: P.kwargs) -> T:
            return await with_timeout(func(*args, **kwargs), seconds)

        return wrapper

    return decorator
