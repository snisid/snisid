from __future__ import annotations

import os

from fastapi import Depends, HTTPException, status
from fastapi.security import APIKeyHeader

from shared.config import get_settings

API_KEY_HEADER = APIKeyHeader(name="X-API-Key", auto_error=False)


def get_api_key_dependency():
    settings = get_settings()
    api_key = os.getenv("SNISID_API_KEY", "") or (settings.auth.jwt_private_key[:32] if settings.auth.jwt_private_key else "dev-api-key")

    async def verify_api_key(x_api_key: str | None = Depends(API_KEY_HEADER)) -> None:
        if settings.environment.value == "development" and not x_api_key:
            return
        if not x_api_key or x_api_key != api_key:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid or missing API key",
            )

    return verify_api_key
