"""
SNISID Retry / Backoff Utilities
=================================
Async retry decorator with exponential backoff for resilient
external service calls.
"""

from __future__ import annotations

import asyncio
import functools
from typing import Any, Callable, Coroutine, TypeVar

import structlog

from shared.logging import get_logger

logger = get_logger(__name__)

T = TypeVar("T")

RetryableFn = Callable[..., Coroutine[Any, Any, T]]


async def async_retry(
    fn: RetryableFn[T],
    max_retries: int = 3,
    base_delay: float = 0.5,
    max_delay: float = 10.0,
    backoff_factor: float = 2.0,
    retryable_exceptions: tuple[type[Exception], ...] = (ConnectionError, TimeoutError, OSError),
) -> T:
    """
    Execute an async callable with exponential backoff retry.

    Args:
        fn: Async callable to execute.
        max_retries: Maximum number of retry attempts.
        base_delay: Initial delay in seconds before first retry.
        max_delay: Maximum delay in seconds between retries.
        backoff_factor: Multiplier applied to delay after each retry.
        retryable_exceptions: Exception types that trigger a retry.

    Returns:
        The return value of *fn*.

    Raises:
        The last exception encountered if all retries are exhausted.
    """
    last_exc: Exception | None = None
    delay = base_delay

    for attempt in range(max_retries + 1):
        try:
            return await fn()
        except retryable_exceptions as exc:
            last_exc = exc
            if attempt < max_retries:
                jitter = delay * (0.5 + asyncio.get_running_loop().time() % 0.5)
                actual_delay = min(jitter, max_delay)
                logger.warning(
                    "retry_attempt",
                    attempt=attempt + 1,
                    max_retries=max_retries,
                    delay_ms=round(actual_delay * 1000),
                    error=str(exc),
                )
                await asyncio.sleep(actual_delay)
                delay = min(delay * backoff_factor, max_delay)
            else:
                logger.error(
                    "retry_exhausted",
                    max_retries=max_retries,
                    error=str(exc),
                )

    if last_exc is not None:
        raise last_exc
    raise RuntimeError("Unexpected: retry loop ended without exception")


def with_retry(
    max_retries: int = 3,
    base_delay: float = 0.5,
    max_delay: float = 10.0,
    backoff_factor: float = 2.0,
    retryable_exceptions: tuple[type[Exception], ...] = (ConnectionError, TimeoutError, OSError),
) -> Callable[[RetryableFn[T]], RetryableFn[T]]:
    """
    Decorator that wraps an async function with exponential backoff retry.

    Usage::

        @with_retry(max_retries=3)
        async def fetch_data(url: str) -> dict:
            ...
    """
    def decorator(fn: RetryableFn[T]) -> RetryableFn[T]:
        @functools.wraps(fn)
        async def wrapper(*args: Any, **kwargs: Any) -> T:
            return await async_retry(
                lambda: fn(*args, **kwargs),
                max_retries=max_retries,
                base_delay=base_delay,
                max_delay=max_delay,
                backoff_factor=backoff_factor,
                retryable_exceptions=retryable_exceptions,
            )
        return wrapper
    return decorator
