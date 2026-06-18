"""
Tests for RiskScorer
"""
import pytest
import json
from tests.mocks import MockRedis
from services.event_processor.risk_scorer import RiskScorer


@pytest.fixture
def scorer():
    rs = RiskScorer()
    rs.redis = MockRedis()
    return rs


class TestRiskScorer:
    @pytest.mark.asyncio
    async def test_score_low_risk(self, scorer):
        result = await scorer.score_person("HT00000001", {
            "risk_level": "LOW",
            "warrants_active": 0,
            "gang_affiliations": [],
            "interpol_notices": [],
        })
        assert result["risk_level"] == "LOW"
        assert result["risk_score"] <= 0.3

    @pytest.mark.asyncio
    async def test_score_critical_risk(self, scorer):
        result = await scorer.score_person("HT99999999", {
            "risk_level": "CRITICAL",
            "warrants_active": 5,
            "gang_affiliations": ["Gang A", "Gang B", "Gang C"],
            "interpol_notices": ["RED", "BLUE"],
        })
        assert result["risk_level"] in ("CRITICAL", "HIGH")
        assert result["risk_score"] >= 0.3

    @pytest.mark.asyncio
    async def test_score_medium_risk(self, scorer):
        result = await scorer.score_person("HT55555555", {
            "risk_level": "MEDIUM",
            "warrants_active": 1,
            "gang_affiliations": ["Gang A"],
            "interpol_notices": [],
        })
        assert result["risk_level"] in ("MEDIUM", "LOW")
        assert 0.1 <= result["risk_score"] <= 0.8

    @pytest.mark.asyncio
    async def test_apply_ai_score_new(self, scorer):
        result = await scorer.apply_ai_score("HT00000001", "fraud", 0.85)
        assert result is not None
        assert result["risk_level"] == "HIGH"

    @pytest.mark.asyncio
    async def test_apply_ai_score_existing(self, scorer):
        await scorer.redis.setex("risk:HT00000001", 3600, json.dumps({
            "risk_score": 0.5, "risk_level": "MEDIUM", "scores": {"base": 0.5},
        }))
        result = await scorer.apply_ai_score("HT00000001", "fraud", 0.9)
        assert result is not None
        blended = (0.5 + 0.9) / 2
        assert abs(result["risk_score"] - round(blended, 3)) < 0.001
