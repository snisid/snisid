"""
SNISID Auth Dependencies — FastAPI Zero Trust Enforcement
==========================================================
FastAPI dependencies for authentication and authorization.
Every endpoint MUST use these — Zero Trust by default.
"""
from __future__ import annotations

from typing import Annotated, Any

from fastapi import Depends, HTTPException, Request, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
from pydantic import BaseModel, Field

from shared.auth import InvalidTokenError, TokenPayload, get_jwt_handler
from shared.logging import get_logger, set_log_context

logger = get_logger(__name__)

# FastAPI security scheme
_bearer_scheme = HTTPBearer(
    scheme_name="JWT Bearer",
    description="RS256 JWT token from the SNISID Auth Service",
    auto_error=True,
)


class CurrentUser(BaseModel):
    """Authenticated user context available in every request."""

    id: str = Field(..., description="User ID (from JWT sub)")
    roles: list[str] = Field(default_factory=list, description="User roles")
    permissions: list[str] = Field(default_factory=list, description="User permissions")
    agency_id: str | None = Field(None, description="Agency ID")
    token_jti: str = Field(..., description="JWT ID for revocation tracking")
    token_type: str = Field("access", description="Token type")

    def has_role(self, role: str) -> bool:
        """Check if user has a specific role."""
        return role in self.roles or "SUPER_ADMIN" in self.roles

    def has_any_role(self, *roles: str) -> bool:
        """Check if user has any of the specified roles."""
        return any(self.has_role(r) for r in roles)

    def has_permission(self, permission: str) -> bool:
        """Check if user has a specific permission."""
        return permission in self.permissions or "SUPER_ADMIN" in self.roles

    def belongs_to_agency(self, agency_id: str) -> bool:
        """Check if user belongs to a specific agency."""
        return self.agency_id == agency_id or "SUPER_ADMIN" in self.roles


async def _extract_token(
    credentials: HTTPAuthorizationCredentials = Depends(_bearer_scheme),
) -> TokenPayload:
    """Extract and validate JWT token from Authorization header."""
    try:
        handler = get_jwt_handler()
        payload = handler.decode_token(credentials.credentials)

        if payload.token_type != "access":
            raise InvalidTokenError("Expected access token, got refresh token")

        return payload
    except InvalidTokenError as e:
        logger.warning("auth_failed", error=str(e))
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid or expired token",
            headers={"WWW-Authenticate": "Bearer"},
        )


async def get_current_user(
    request: Request,
    token: TokenPayload = Depends(_extract_token),
) -> CurrentUser:
    """
    FastAPI dependency that returns the authenticated user.
    Sets logging context with user info.

    Usage:
        @router.get("/endpoint")
        async def endpoint(user: CurrentUser = Depends(get_current_user)):
            ...
    """
    user = CurrentUser(
        id=token.sub,
        roles=token.roles,
        permissions=token.permissions,
        agency_id=token.agency_id,
        token_jti=token.jti,
        token_type=token.token_type,
    )

    # Set logging context for all downstream log calls
    set_log_context(
        user_id=user.id,
        trace_id=request.headers.get("X-Trace-ID", ""),
        correlation_id=request.headers.get("X-Correlation-ID", ""),
    )

    return user


# Type alias for cleaner route signatures
AuthenticatedUser = Annotated[CurrentUser, Depends(get_current_user)]


def require_role(*roles: str):
    """
    Dependency factory that enforces role-based access.

    Usage:
        @router.post("/admin-action")
        async def admin_action(
            user: CurrentUser = Depends(require_role("ADMIN", "SUPER_ADMIN"))
        ):
            ...
    """

    async def _check_role(user: CurrentUser = Depends(get_current_user)) -> CurrentUser:
        if not user.has_any_role(*roles):
            logger.warning(
                "authorization_denied",
                user_id=user.id,
                required_roles=list(roles),
                user_roles=user.roles,
            )
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail=f"Requires one of roles: {', '.join(roles)}",
            )
        return user

    return _check_role


def require_permission(*permissions: str):
    """
    Dependency factory that enforces permission-based access.

    Usage:
        @router.delete("/resource/{id}")
        async def delete_resource(
            user: CurrentUser = Depends(require_permission("resource:delete"))
        ):
            ...
    """

    async def _check_permission(
        user: CurrentUser = Depends(get_current_user),
    ) -> CurrentUser:
        for perm in permissions:
            if not user.has_permission(perm):
                logger.warning(
                    "permission_denied",
                    user_id=user.id,
                    required_permissions=list(permissions),
                )
                raise HTTPException(
                    status_code=status.HTTP_403_FORBIDDEN,
                    detail=f"Missing permission: {perm}",
                )
        return user

    return _check_permission


def require_agency(agency_id_param: str = "agency_id"):
    """
    Dependency factory that enforces agency-level access control.
    The user must belong to the agency specified in the path/query parameter.

    Usage:
        @router.get("/agencies/{agency_id}/resources")
        async def get_resources(
            agency_id: str,
            user: CurrentUser = Depends(require_agency("agency_id"))
        ):
            ...
    """

    async def _check_agency(
        request: Request,
        user: CurrentUser = Depends(get_current_user),
    ) -> CurrentUser:
        target_agency = request.path_params.get(agency_id_param) or request.query_params.get(
            agency_id_param
        )
        if target_agency and not user.belongs_to_agency(target_agency):
            logger.warning(
                "agency_access_denied",
                user_id=user.id,
                user_agency=user.agency_id,
                target_agency=target_agency,
            )
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Access denied for this agency",
            )
        return user

    return _check_agency
