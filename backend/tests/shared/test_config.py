"""
SNISID Configuration Tests
==========================
Tests for Pydantic-based configuration management with Vault integration.
"""
from __future__ import annotations

import os
from unittest.mock import patch, MagicMock

import pytest

from shared.config import (
    Settings,
    DatabaseConfig,
    RedisConfig,
    KafkaConfig,
    AuthConfig,
    VaultConfig,
    Environment,
    get_settings,
)


class TestDatabaseConfig:
    """Test DatabaseConfig settings."""

    def test_default_values(self):
        config = DatabaseConfig()
        assert config.host == "localhost"
        assert config.port == 5432
        assert config.name == "snisid"
        assert config.user == "snisid"
        assert config.pool_size == 20
        assert config.max_overflow == 10

    def test_async_url_property(self):
        config = DatabaseConfig(password="secret")
        url = config.async_url
        assert "postgresql+asyncpg" in url
        assert "snisid:secret" in url
        assert "localhost:5432" in url

    def test_sync_url_property(self):
        config = DatabaseConfig(password="secret")
        url = config.sync_url
        assert "postgresql+psycopg" in url
        assert "snisid:secret" in url

    def test_env_prefix_override(self):
        with patch.dict(os.environ, {"DB_HOST": "db-prod.internal"}):
            config = DatabaseConfig()
            assert config.host == "db-prod.internal"

    def test_ssl_mode_default(self):
        config = DatabaseConfig()
        assert config.ssl_mode == "prefer"


class TestRedisConfig:
    """Test RedisConfig settings."""

    def test_default_values(self):
        config = RedisConfig()
        assert config.host == "localhost"
        assert config.port == 6379
        assert config.db == 0
        assert config.max_connections == 50

    def test_url_no_password(self):
        config = RedisConfig()
        assert config.url == "redis://localhost:6379/0"

    def test_url_with_password(self):
        config = RedisConfig(password="redispass")
        assert ":redispass@" in config.url
        assert config.url.startswith("redis://")

    def test_url_ssl(self):
        config = RedisConfig(ssl=True, password="pass")
        assert config.url.startswith("rediss://")

    def test_get_url_custom_db(self):
        config = RedisConfig()
        url = config.get_url(db=5)
        assert url.endswith("/5")

    def test_dedicated_dbs(self):
        config = RedisConfig()
        assert config.cache_db == 0
        assert config.session_db == 1
        assert config.rate_limit_db == 2
        assert config.celery_db == 3


class TestKafkaConfig:
    """Test KafkaConfig settings."""

    def test_default_values(self):
        config = KafkaConfig()
        assert config.bootstrap_servers == "localhost:9092"
        assert config.producer_acks == "all"
        assert config.producer_retries == 3

    def test_bootstrap_list(self):
        config = KafkaConfig(bootstrap_servers="kafka1:9092,kafka2:9092")
        assert len(config.bootstrap_list) == 2
        assert "kafka1:9092" in config.bootstrap_list

    def test_single_bootstrap(self):
        config = KafkaConfig()
        assert len(config.bootstrap_list) == 1

    def test_consumer_defaults(self):
        config = KafkaConfig()
        assert config.consumer_group_prefix == "snisid"
        assert config.consumer_auto_offset_reset == "earliest"
        assert config.consumer_max_poll_records == 500

    def test_topic_prefix(self):
        config = KafkaConfig()
        assert config.topic_prefix == "snisid"


class TestAuthConfig:
    """Test AuthConfig settings."""

    def test_default_values(self):
        config = AuthConfig()
        assert config.jwt_algorithm == "RS256"
        assert config.jwt_access_token_expire_minutes == 30
        assert config.jwt_refresh_token_expire_days == 7

    def test_mfa_required_roles(self):
        config = AuthConfig()
        assert "SUPER_ADMIN" in config.mfa_required_roles
        assert "ADMIN" in config.mfa_required_roles

    def test_account_security_defaults(self):
        config = AuthConfig()
        assert config.max_failed_attempts == 5
        assert config.lockout_duration_minutes == 30
        assert config.password_min_length == 12


class TestVaultConfig:
    """Test VaultConfig settings."""

    def test_default_disabled(self):
        config = VaultConfig()
        assert config.enabled is False
        assert config.url == "http://localhost:8200"

    def test_path_prefix(self):
        config = VaultConfig()
        assert config.path_prefix == "snisid"


class TestSettings:
    """Test master Settings aggregation."""

    def test_default_environment(self):
        settings = Settings()
        assert settings.environment == Environment.DEVELOPMENT

    def test_production_environment(self):
        with patch.dict(os.environ, {"SNISID_ENVIRONMENT": "production"}):
            settings = Settings()
            assert settings.environment == Environment.PRODUCTION

    def test_staging_environment(self):
        with patch.dict(os.environ, {"SNISID_ENVIRONMENT": "staging"}):
            settings = Settings()
            assert settings.environment == Environment.STAGING

    def test_invalid_environment_raises(self):
        with patch.dict(os.environ, {"SNISID_ENVIRONMENT": "invalid"}):
            with pytest.raises(ValueError):
                Settings()

    def test_sub_configs_initialized(self):
        settings = Settings()
        assert isinstance(settings.database, DatabaseConfig)
        assert isinstance(settings.redis, RedisConfig)
        assert isinstance(settings.kafka, KafkaConfig)
        assert isinstance(settings.auth, AuthConfig)
        assert isinstance(settings.vault, VaultConfig)

    def test_service_defaults(self):
        settings = Settings()
        assert settings.service_name == "snisid"
        assert settings.service_version == "1.0.0"
        assert settings.port == 8000
        assert settings.workers == 4

    def test_cors_origins_default(self):
        settings = Settings()
        assert settings.cors_origins == ["*"]

    def test_env_override_service_name(self):
        with patch.dict(os.environ, {"SNISID_SERVICE_NAME": "identity-service"}):
            settings = Settings()
            assert settings.service_name == "identity-service"

    def test_env_override_port(self):
        with patch.dict(os.environ, {"SNISID_PORT": "9090"}):
            settings = Settings()
            assert settings.port == 9090


class TestSettingsVault:
    """Test Vault secret loading."""

    def test_vault_disabled_skips_loading(self):
        settings = Settings()
        settings.vault.enabled = False
        # Should not raise
        settings.load_secrets_from_vault()

    def test_vault_enabled_without_hvac_raises(self):
        settings = Settings()
        settings.vault.enabled = True
        with patch.dict("sys.modules", {"hvac": None}):
            with pytest.raises(RuntimeError, match="hvac package required"):
                settings.load_secrets_from_vault()

    def test_vault_enabled_imports_hvac(self):
        settings = Settings()
        settings.vault.enabled = True
        settings.vault.url = "http://vault:8200"
        settings.vault.token = "root-token"

        mock_client = MagicMock()
        mock_client.is_authenticated.return_value = True

        # Mock successful secret reads
        mock_client.secrets.kv.v2.read_secret_version.side_effect = [
            {"data": {"data": {"password": "db-pass", "username": "db-user"}}},
            {"data": {"data": {"password": "redis-pass"}}},
            {"data": {"data": {"private_key": "jwt-priv", "public_key": "jwt-pub"}}},
            {"data": {"data": {"username": "kafka-user", "password": "kafka-pass"}}},
        ]

        with patch("hvac.Client", return_value=mock_client):
            settings.load_secrets_from_vault()

        assert settings.database.password == "db-pass"
        assert settings.database.user == "db-user"
        assert settings.redis.password == "redis-pass"
        assert settings.auth.jwt_private_key == "jwt-priv"
        assert settings.auth.jwt_public_key == "jwt-pub"
        assert settings.kafka.sasl_username == "kafka-user"
        assert settings.kafka.sasl_password == "kafka-pass"

    def test_vault_authentication_failure(self):
        settings = Settings()
        settings.vault.enabled = True

        mock_client = MagicMock()
        mock_client.is_authenticated.return_value = False

        with patch("hvac.Client", return_value=mock_client):
            with pytest.raises(RuntimeError, match="Vault authentication failed"):
                settings.load_secrets_from_vault()

    def test_vault_partial_secrets(self):
        """Test that missing secret paths don't crash."""
        settings = Settings()
        settings.vault.enabled = True

        mock_client = MagicMock()
        mock_client.is_authenticated.return_value = True
        mock_client.secrets.kv.v2.read_secret_version.side_effect = [
            {"data": {"data": {"password": "only-db-pass"}}},
            Exception("Redis path not found"),
            {"data": {"data": {}}},
            {"data": {"data": {"username": "kafka-user"}}},
        ]

        with patch("hvac.Client", return_value=mock_client):
            with pytest.raises(RuntimeError, match="Failed to load secrets"):
                settings.load_secrets_from_vault()


class TestGetSettings:
    """Test cached settings singleton."""

    def test_get_settings_returns_settings_instance(self):
        settings = get_settings()
        assert isinstance(settings, Settings)

    def test_get_settings_cached(self):
        s1 = get_settings()
        s2 = get_settings()
        assert s1 is s2

    def test_get_settings_has_all_configs(self):
        settings = get_settings()
        assert hasattr(settings, "database")
        assert hasattr(settings, "redis")
        assert hasattr(settings, "kafka")
        assert hasattr(settings, "auth")
        assert hasattr(settings, "vault")
        assert hasattr(settings, "observability")
