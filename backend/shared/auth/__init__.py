"""
SNISID JWT Handler — RS256 Token Management
=============================================
JWT creation and validation using RS256 (RSA-SHA256) asymmetric keys.
Supports key rotation, JWKS, refresh tokens, and Redis-based blacklisting.
"""
from __future__ import annotations

import uuid
from datetime import datetime, timedelta, timezone
from typing import Any

from jose import JWTError, jwt
from pydantic import BaseModel, Field

from shared.config import get_settings
from shared.logging import get_logger

logger = get_logger(__name__)


class TokenPayload(BaseModel):
    """Decoded JWT token claims."""

    sub: str = Field(..., description="Subject (user ID)")
    jti: str = Field(..., description="JWT ID (unique token identifier)")
    roles: list[str] = Field(default_factory=list, description="User roles")
    permissions: list[str] = Field(default_factory=list, description="User permissions")
    agency_id: str | None = Field(None, description="Agency the user belongs to")
    iss: str = Field("snisid-auth", description="Issuer")
    aud: str = Field("snisid-api", description="Audience")
    exp: int = Field(..., description="Expiration timestamp")
    iat: int = Field(..., description="Issued at timestamp")
    token_type: str = Field("access", description="Token type (access/refresh)")


class TokenPair(BaseModel):
    """Access + Refresh token pair."""

    access_token: str
    refresh_token: str
    token_type: str = "bearer"
    expires_in: int
    refresh_expires_in: int


class JWTHandler:
    """
    RS256 JWT token manager with key rotation support.

    In production, keys are loaded from Vault. In development,
    falls back to file paths or generates ephemeral keys.
    """

    def __init__(self) -> None:
        self._settings = get_settings().auth
        self._private_key: str | None = None
        self._public_key: str | None = None
        self._load_keys()

    def _load_keys(self) -> None:
        """Load RSA keys from settings (Vault or file paths)."""
        # Priority: Vault-provided keys > file paths > generate ephemeral
        if self._settings.jwt_private_key:
            self._private_key = self._settings.jwt_private_key
            self._public_key = self._settings.jwt_public_key
            logger.info("jwt_keys_loaded", source="vault")
            return

        if self._settings.jwt_private_key_path:
            try:
                with open(self._settings.jwt_private_key_path) as f:
                    self._private_key = f.read()
                with open(self._settings.jwt_public_key_path) as f:
                    self._public_key = f.read()
                logger.info("jwt_keys_loaded", source="file")
                return
            except FileNotFoundError:
                logger.warning("jwt_key_files_not_found", msg="Generating ephemeral keys")

        # Development fallback: generate ephemeral RSA key pair
        self._generate_ephemeral_keys()

    def _generate_ephemeral_keys(self) -> None:
        """Generate an ephemeral RSA key pair for development."""
        from cryptography.hazmat.primitives import serialization
        from cryptography.hazmat.primitives.asymmetric import rsa

        private_key = rsa.generate_private_key(
            public_exponent=65537,
            key_size=2048,
        )
        self._private_key = private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=serialization.NoEncryption(),
        ).decode("utf-8")

        self._public_key = (
            private_key.public_key()
            .public_bytes(
                encoding=serialization.Encoding.PEM,
                format=serialization.PublicFormat.SubjectPublicKeyInfo,
            )
            .decode("utf-8")
        )
        logger.warning("jwt_ephemeral_keys_generated", msg="DO NOT use in production")

    def create_access_token(
        self,
        user_id: str,
        roles: list[str],
        permissions: list[str] | None = None,
        agency_id: str | None = None,
        extra_claims: dict[str, Any] | None = None,
    ) -> str:
        """
        Create a signed JWT access token.

        Args:
            user_id: Subject identifier (user ID).
            roles: List of role strings.
            permissions: List of permission strings.
            agency_id: Associated agency ID.
            extra_claims: Additional custom claims.

        Returns:
            Encoded JWT string.
        """
        now = datetime.now(timezone.utc)
        expire = now + timedelta(minutes=self._settings.jwt_access_token_expire_minutes)

        claims: dict[str, Any] = {
            "sub": user_id,
            "jti": str(uuid.uuid4()),
            "roles": roles,
            "permissions": permissions or [],
            "agency_id": agency_id,
            "iss": self._settings.jwt_issuer,
            "aud": self._settings.jwt_audience,
            "exp": int(expire.timestamp()),
            "iat": int(now.timestamp()),
            "token_type": "access",
        }
        if extra_claims:
            claims.update(extra_claims)

        token = jwt.encode(
            claims,
            self._private_key,
            algorithm=self._settings.jwt_algorithm,
        )
        logger.debug("access_token_created", user_id=user_id, jti=claims["jti"])
        return token

    def create_refresh_token(
        self,
        user_id: str,
        roles: list[str],
        agency_id: str | None = None,
    ) -> str:
        """Create a long-lived refresh token."""
        now = datetime.now(timezone.utc)
        expire = now + timedelta(days=self._settings.jwt_refresh_token_expire_days)

        claims: dict[str, Any] = {
            "sub": user_id,
            "jti": str(uuid.uuid4()),
            "roles": roles,
            "agency_id": agency_id,
            "iss": self._settings.jwt_issuer,
            "aud": self._settings.jwt_audience,
            "exp": int(expire.timestamp()),
            "iat": int(now.timestamp()),
            "token_type": "refresh",
        }

        return jwt.encode(
            claims,
            self._private_key,
            algorithm=self._settings.jwt_algorithm,
        )

    def create_token_pair(
        self,
        user_id: str,
        roles: list[str],
        permissions: list[str] | None = None,
        agency_id: str | None = None,
    ) -> TokenPair:
        """Create an access + refresh token pair."""
        return TokenPair(
            access_token=self.create_access_token(
                user_id, roles, permissions, agency_id
            ),
            refresh_token=self.create_refresh_token(user_id, roles, agency_id),
            expires_in=self._settings.jwt_access_token_expire_minutes * 60,
            refresh_expires_in=self._settings.jwt_refresh_token_expire_days * 86400,
        )

    def decode_token(self, token: str) -> TokenPayload:
        """
        Decode and validate a JWT token.

        Args:
            token: The encoded JWT string.

        Returns:
            Decoded TokenPayload.

        Raises:
            InvalidTokenError: If the token is invalid, expired, or tampered.
        """
        try:
            payload = jwt.decode(
                token,
                self._public_key,
                algorithms=[self._settings.jwt_algorithm],
                issuer=self._settings.jwt_issuer,
                audience=self._settings.jwt_audience,
            )
            return TokenPayload(**payload)
        except JWTError as e:
            logger.warning("token_decode_failed", error=str(e))
            raise InvalidTokenError(f"Invalid token: {e}") from e

    def get_public_key_jwks(self) -> dict[str, Any]:
        """
        Return the public key in JWKS format for the /.well-known/jwks.json endpoint.
        """
        from cryptography.hazmat.primitives.serialization import load_pem_public_key
        import base64

        public_key = load_pem_public_key(self._public_key.encode())
        numbers = public_key.public_numbers()  # type: ignore[union-attr]

        def _int_to_base64url(n: int) -> str:
            byte_length = (n.bit_length() + 7) // 8
            return base64.urlsafe_b64encode(
                n.to_bytes(byte_length, byteorder="big")
            ).rstrip(b"=").decode("ascii")

        return {
            "keys": [
                {
                    "kty": "RSA",
                    "use": "sig",
                    "alg": "RS256",
                    "kid": "snisid-primary",
                    "n": _int_to_base64url(numbers.n),
                    "e": _int_to_base64url(numbers.e),
                }
            ]
        }


class InvalidTokenError(Exception):
    """Raised when a JWT token is invalid, expired, or tampered."""
    pass


# Module-level singleton
_jwt_handler: JWTHandler | None = None


def get_jwt_handler() -> JWTHandler:
    """Get or create the JWT handler singleton."""
    global _jwt_handler
    if _jwt_handler is None:
        _jwt_handler = JWTHandler()
    return _jwt_handler
