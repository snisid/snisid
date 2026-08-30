import pytest
from api_nire import NIREResolver


@pytest.fixture
def resolver():
    return NIREResolver()


class TestNIREValidation:
    def test_valid_nire(self, resolver):
        assert resolver.validate_nire("HT2026000010001") is True

    def test_empty_nire(self, resolver):
        assert resolver.validate_nire("") is False

    def test_short_nire(self, resolver):
        assert resolver.validate_nire("HT12") is False

    def test_long_nire(self, resolver):
        assert resolver.validate_nire("HT202600001000100") is False

    def test_invalid_prefix(self, resolver):
        assert resolver.validate_nire("122026000010001") is False


class TestNIREFormat:
    def test_format_correct(self, resolver):
        result = resolver.format_nire(2026, "HT", 10001)
        assert result == "HT2026000010001"
        assert len(result) == 15

    def test_format_roundtrip(self, resolver):
        formatted = resolver.format_nire(2026, "HT", 10001)
        assert resolver.validate_nire(formatted) is True
