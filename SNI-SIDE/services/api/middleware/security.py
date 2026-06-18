import time, hashlib, uuid, json
from typing import Optional, List
from datetime import datetime

import jwt
from fastapi import Request, HTTPException, Depends
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

from config import settings

security_scheme = HTTPBearer(auto_error=False)


class SecurityContext:
    def __init__(self, user_id: str, agency: str, roles: List[str],
                 clearance: str, session_id: str):
        self.user_id = user_id
        self.agency = agency
        self.roles = roles
        self.clearance = clearance
        self.session_id = session_id

    def has_role(self, role: str) -> bool:
        return role in self.roles

    def has_any_role(self, roles: List[str]) -> bool:
        return any(r in self.roles for r in roles)

    def has_clearance(self, required: str) -> bool:
        levels = {"UNCLASSIFIED": 0, "RESTRICTED": 1, "CONFIDENTIAL": 2,
                  "SECRET": 3, "TOP_SECRET": 4}
        return levels.get(self.clearance, 0) >= levels.get(required, 0)

    def __repr__(self):
        return f"SecurityContext(user={self.user_id}, agency={self.agency}, clearance={self.clearance})"


async def verify_token(credentials: HTTPAuthorizationCredentials = Depends(security_scheme)) -> SecurityContext:
    if credentials is None:
        raise HTTPException(status_code=401, detail="Authentication required")
    try:
        payload = jwt.decode(
            credentials.credentials,
            settings.jwt_secret_key,
            algorithms=[settings.jwt_algorithm],
        )
        return SecurityContext(
            user_id=payload.get("sub", "unknown"),
            agency=payload.get("agency", "UNKNOWN"),
            roles=payload.get("roles", []),
            clearance=payload.get("clearance", "CONFIDENTIAL"),
            session_id=payload.get("jti", str(uuid.uuid4())),
        )
    except jwt.ExpiredSignatureError:
        raise HTTPException(status_code=401, detail="Token expired")
    except jwt.InvalidTokenError:
        raise HTTPException(status_code=401, detail="Invalid token")


def verify_agency(allowed_agencies: List[str]):
    async def dependency(security: SecurityContext = Depends(verify_token)) -> SecurityContext:
        if security.agency not in allowed_agencies:
            raise HTTPException(status_code=403, detail=f"Agency '{security.agency}' not authorized")
        return security
    return dependency


def require_clearance(required: str):
    async def dependency(security: SecurityContext = Depends(verify_token)) -> SecurityContext:
        if not security.has_clearance(required):
            raise HTTPException(status_code=403, detail="Insufficient clearance")
        return security
    return dependency


async def create_audit_entry(request: Request, security: SecurityContext,
                              action: str, resource: str, resource_id: str,
                              status: str = "SUCCESS") -> dict:
    entry = {
        "audit_id": str(uuid.uuid4()),
        "timestamp": datetime.utcnow().isoformat(),
        "user_id": security.user_id,
        "agency": security.agency,
        "session_id": security.session_id,
        "action": action,
        "resource": resource,
        "resource_id": resource_id,
        "ip_address": request.client.host if request.client else "unknown",
        "user_agent": request.headers.get("user-agent", ""),
        "status": status,
        "request_id": request.headers.get("X-Request-ID", ""),
    }
    return entry
