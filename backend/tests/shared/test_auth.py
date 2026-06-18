from __future__ import annotations

import time
import uuid
from datetime import datetime, timedelta, timezone
from unittest.mock import patch, MagicMock

import pytest

from shared.auth import JWTHandler, TokenPayload, TokenPair, InvalidTokenError


@pytest.fixture
def jwt_handler():
    return JWTHandler()


class TestJWTTokenGeneration:
    """Test JWT token creation."""

    def test_create_access_token(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN"],
            permissions=["read:citizens", "write:citizens"],
            agency_id="agency-1",
        )
        assert isinstance(token, str)
        assert len(token) > 50
        assert token.count(".") == 2

    def test_create_refresh_token(self, jwt_handler):
        token = jwt_handler.create_refresh_token(
            user_id="user-1",
            roles=["ADMIN"],
            agency_id="agency-1",
        )
        assert isinstance(token, str)
        assert token.count(".") == 2

    def test_create_token_pair(self, jwt_handler):
        pair = jwt_handler.create_token_pair(
            user_id="user-1",
            roles=["REGISTRAR"],
            permissions=["read:citizens"],
            agency_id="agency-1",
        )
        assert isinstance(pair, TokenPair)
        assert pair.token_type == "bearer"
        assert pair.access_token != pair.refresh_token
        assert pair.expires_in > 0
        assert pair.refresh_expires_in > 0


class TestJWTTokenValidation:
    """Test JWT token decoding and validation."""

    def test_decode_valid_token(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN"],
            permissions=["read:citizens"],
            agency_id="agency-1",
        )
        payload = jwt_handler.decode_token(token)
        assert isinstance(payload, TokenPayload)
        assert payload.sub == "user-1"
        assert "ADMIN" in payload.roles
        assert "read:citizens" in payload.permissions
        assert payload.agency_id == "agency-1"
        assert payload.token_type == "access"
        assert payload.iss == "snisid-auth"
        assert payload.aud == "snisid-api"

    def test_decode_invalid_token(self, jwt_handler):
        with pytest.raises(InvalidTokenError):
            jwt_handler.decode_token("invalid.token.here")

    def test_decode_tampered_token(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN"],
        )
        parts = token.split(".")
        tampered = f"{parts[0]}.{parts[1]}modified.{parts[2]}"
        with pytest.raises(InvalidTokenError):
            jwt_handler.decode_token(tampered)

    def test_decode_expired_token(self, jwt_handler):
        with patch.object(jwt_handler._settings, "jwt_access_token_expire_minutes", -1):
            token = jwt_handler.create_access_token(
                user_id="user-1",
                roles=["ADMIN"],
            )
        with pytest.raises(InvalidTokenError):
            jwt_handler.decode_token(token)

    def test_decode_wrong_issuer(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN"],
        )
        payload = jwt_handler.decode_token(token)
        assert payload.iss == "snisid-auth"

    def test_decode_wrong_audience(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN"],
        )
        payload = jwt_handler.decode_token(token)
        assert payload.aud == "snisid-api"


class TestTokenClaims:
    """Test token claims and payload structure."""

    def test_access_token_claims(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN"],
            permissions=["read:citizens"],
            agency_id="agency-1",
            extra_claims={"department": "interior"},
        )
        payload = jwt_handler.decode_token(token)
        assert payload.sub == "user-1"
        assert payload.roles == ["ADMIN"]
        assert payload.permissions == ["read:citizens"]
        assert payload.agency_id == "agency-1"
        assert payload.token_type == "access"

    def test_refresh_token_claims(self, jwt_handler):
        token = jwt_handler.create_refresh_token(
            user_id="user-1",
            roles=["ADMIN"],
            agency_id="agency-1",
        )
        payload = jwt_handler.decode_token(token)
        assert payload.sub == "user-1"
        assert payload.roles == ["ADMIN"]
        assert payload.agency_id == "agency-1"
        assert payload.token_type == "refresh"
        assert payload.permissions == []

    def test_token_has_unique_jti(self, jwt_handler):
        t1 = jwt_handler.create_access_token("u1", ["ADMIN"])
        t2 = jwt_handler.create_access_token("u1", ["ADMIN"])
        p1 = jwt_handler.decode_token(t1)
        p2 = jwt_handler.decode_token(t2)
        assert p1.jti != p2.jti


class TestAgencyAuthorization:
    """Test agency-based authorization in tokens."""

    def test_token_with_agency(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["REGISTRAR"],
            agency_id="agency-dakar",
        )
        payload = jwt_handler.decode_token(token)
        assert payload.agency_id == "agency-dakar"

    def test_token_without_agency(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["SUPER_ADMIN"],
        )
        payload = jwt_handler.decode_token(token)
        assert payload.agency_id is None


class TestJWKS:
    """Test JWKS key exposure."""

    def test_get_jwks(self, jwt_handler):
        jwks = jwt_handler.get_public_key_jwks()
        assert "keys" in jwks
        assert len(jwks["keys"]) == 1
        key = jwks["keys"][0]
        assert key["kty"] == "RSA"
        assert key["alg"] == "RS256"
        assert key["use"] == "sig"
        assert "n" in key
        assert "e" in key


class TestTokenExpiration:
    """Test token expiration behavior."""

    def test_access_token_expiration(self, jwt_handler):
        with patch.object(jwt_handler._settings, "jwt_access_token_expire_minutes", 30):
            token = jwt_handler.create_access_token("u1", ["ADMIN"])
            payload = jwt_handler.decode_token(token)
            now = int(datetime.now(timezone.utc).timestamp())
            expected_exp = now + 30 * 60
            assert abs(payload.exp - expected_exp) < 5

    def test_refresh_token_expiration(self, jwt_handler):
        with patch.object(jwt_handler._settings, "jwt_refresh_token_expire_days", 7):
            token = jwt_handler.create_refresh_token("u1", ["ADMIN"])
            payload = jwt_handler.decode_token(token)
            now = int(datetime.now(timezone.utc).timestamp())
            expected_exp = now + 7 * 86400
            assert abs(payload.exp - expected_exp) < 5


class TestRoleBasedTokens:
    """Test role and permission handling."""

    def test_multiple_roles(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["ADMIN", "REGISTRAR", "AUDITOR"],
        )
        payload = jwt_handler.decode_token(token)
        assert len(payload.roles) == 3
        assert "ADMIN" in payload.roles
        assert "REGISTRAR" in payload.roles

    def test_empty_permissions(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=["VIEWER"],
        )
        payload = jwt_handler.decode_token(token)
        assert payload.permissions == []

    def test_empty_roles(self, jwt_handler):
        token = jwt_handler.create_access_token(
            user_id="user-1",
            roles=[],
        )
        payload = jwt_handler.decode_token(token)
        assert payload.roles == []
