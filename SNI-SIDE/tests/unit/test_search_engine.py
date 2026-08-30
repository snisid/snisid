"""
Tests for the National Sovereign Search Engine — type detection
"""
import pytest
import json


class TestQueryTypeDetection:
    @pytest.mark.parametrize("query,expected_type", [
        ("HT12345678", "NIU"),
        ("AB-123-CD", "PLATE"),
        ("1HGCM82633A004352", "VIN"),
        ("HT123456", "PASSPORT"),
        ("+50912345678", "PHONE"),
        ("jean.dupont@email.com", "EMAIL"),
        ("192.168.1.1", "IP"),
        ("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "WALLET"),
        ("malware-site.com", "DOMAIN"),
        ("Case-2024-001", "CASE"),
        ("Jean Dupont", "NAME"),
        ("", "UNKNOWN"),
        ("12345", "UNKNOWN"),
    ])
    def test_detect_type(self, query, expected_type):
        from ai.national_search_engine import detect_query_type
        assert detect_query_type(query) == expected_type


class TestSearchIndex:
    def test_niu_pattern(self):
        from ai.national_search_engine import detect_query_type
        examples = ["HT00000001", "HT99999999", "HTABCDEFGH"]
        for ex in examples:
            assert detect_query_type(ex) == "NIU", f"Failed: {ex}"

    def test_plate_patterns(self):
        from ai.national_search_engine import detect_query_type
        plates = ["AB-123-CD", "AA-000-BB", "ZZ-999-XX", "HT-001-AA"]
        for p in plates:
            assert detect_query_type(p) == "PLATE", f"Failed: {p}"

    def test_phone_patterns(self):
        from ai.national_search_engine import detect_query_type
        phones = ["+50912345678", "+50987654321", "+50900000000"]
        for p in phones:
            assert detect_query_type(p) == "PHONE", f"Failed: {p}"

    def test_email_patterns(self):
        from ai.national_search_engine import detect_query_type
        emails = ["test@sniside.ht", "jean.dupont@pnh.ht", "admin@gov.ht"]
        for e in emails:
            assert detect_query_type(e) == "EMAIL", f"Failed: {e}"

    def test_ip_patterns(self):
        from ai.national_search_engine import detect_query_type
        ips = ["10.0.0.1", "192.168.1.100", "185.220.101.0"]
        for ip in ips:
            assert detect_query_type(ip) == "IP", f"Failed: {ip}"

    def test_wallet_patterns(self):
        from ai.national_search_engine import detect_query_type
        wallets = [
            "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
            "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq",
            "3J98t1WpEZ73CNmQviecrnyiWrnqRhWNLy",
        ]
        for w in wallets:
            assert detect_query_type(w) == "WALLET", f"Failed: {w}"
