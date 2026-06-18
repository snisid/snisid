"""
Tests for AlertCorrelator
"""
import pytest
import time
from services.event_processor.alert_correlator import AlertCorrelator


@pytest.fixture
def correlator():
    c = AlertCorrelator()
    c.start()
    return c


class TestAlertCorrelator:
    def test_single_event_no_correlation(self, correlator):
        result = correlator.add_event("ALPR_WANTED_HIT", "HT12345678", "HIGH", int(time.time() * 1000))
        assert result is None

    def test_multi_event_same_domain_no_correlation(self, correlator):
        ts = int(time.time() * 1000)
        for _ in range(5):
            correlator.add_event("ALPR_WANTED_HIT", "HT12345678", "HIGH", ts)
        result = correlator.add_event("ALPR_ANOMALY", "HT12345678", "MEDIUM", ts)
        assert result is None  # 6 events, 1 domain → no correlation

    def test_multi_domain_correlation(self, correlator):
        ts = int(time.time() * 1000)
        correlator.add_event("ALPR_WANTED_HIT", "HT12345678", "HIGH", ts)
        correlator.add_event("BORDER_WANTED", "HT12345678", "CRITICAL", ts)
        correlator.add_event("WATCHLIST_MATCH", "HT12345678", "HIGH", ts)
        result = correlator.add_event("NCID_WANTED", "HT12345678", "HIGH", ts)
        assert result is not None
        assert result["entity_id"] == "HT12345678"
        assert result["domain_count"] >= 2
        assert result["alert_count"] >= 3

    def test_different_entities_no_correlation(self, correlator):
        ts = int(time.time() * 1000)
        correlator.add_event("ALPR_WANTED_HIT", "HT11111111", "HIGH", ts)
        correlator.add_event("BORDER_WANTED", "HT22222222", "HIGH", ts)
        correlator.add_event("WATCHLIST_MATCH", "HT33333333", "HIGH", ts)
        result = correlator.add_event("NCID_WANTED", "HT44444444", "HIGH", ts)
        assert result is None

    def test_window_expiry(self, correlator):
        old_ts = int(time.time() * 1000) - 600000  # 10 min ago (window=5min)
        correlator.add_event("ALPR_WANTED_HIT", "HT12345678", "HIGH", old_ts)
        correlator.add_event("BORDER_WANTED", "HT12345678", "HIGH", old_ts)
        correlator.add_event("WATCHLIST_MATCH", "HT12345678", "HIGH", old_ts)
        result = correlator.add_event("NCID_WANTED", "HT12345678", "HIGH", old_ts)
        assert result is not None  # everything in same old window
