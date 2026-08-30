import logging, os, json, time
from typing import Optional

logger = logging.getLogger("sniside.risk-scorer")

REDIS_HOST = os.getenv("REDIS_HOST", "redis")
REDIS_PORT = int(os.getenv("REDIS_PORT", "6379"))


class RiskScorer:
    def __init__(self):
        self.redis = None

    def start(self):
        try:
            import redis.asyncio as aioredis
            self.redis = aioredis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)
            logger.info("RiskScorer Redis initialized")
        except ImportError:
            logger.warning("redis.asyncio not available — risk scorer runs without cache")

    def stop(self):
        if self.redis:
            self.redis.close()

    async def score_person(self, niu: str, data: dict) -> dict:
        scores = {
            "base": 0.0,
            "warrant_score": 0.0,
            "gang_score": 0.0,
            "interpol_score": 0.0,
        }

        risk_map = {"CRITICAL": 1.0, "HIGH": 0.75, "MEDIUM": 0.5, "LOW": 0.25}
        scores["base"] = risk_map.get(data.get("risk_level", "LOW"), 0.25)
        scores["warrant_score"] = min(data.get("warrants_active", 0) * 0.15, 0.6)
        gang_count = len(data.get("gang_affiliations", []))
        scores["gang_score"] = min(gang_count * 0.1, 0.3)
        interpol_count = len(data.get("interpol_notices", []))
        scores["interpol_score"] = min(interpol_count * 0.2, 0.5)

        total = min(1.0, sum(scores.values()) * 0.4)

        risk_level = "LOW"
        if total >= 0.8:
            risk_level = "CRITICAL"
        elif total >= 0.6:
            risk_level = "HIGH"
        elif total >= 0.35:
            risk_level = "MEDIUM"

        result = {"risk_score": round(total, 3), "risk_level": risk_level, "scores": scores}

        if self.redis:
            await self.redis.setex(f"risk:{niu}", 3600, json.dumps(result))
        return result

    async def apply_ai_score(self, niu: str, model_name: str, score: float) -> Optional[dict]:
        if not self.redis:
            base = {"risk_score": score * 0.5, "risk_level": "MEDIUM", "scores": {}}
            if score >= 0.9:
                base["risk_level"] = "CRITICAL"
            elif score >= 0.7:
                base["risk_level"] = "HIGH"
            return base

        cached = await self.redis.get(f"risk:{niu}")
        if cached:
            current = json.loads(cached)
        else:
            current = {"risk_score": 0.0, "scores": {}}

        current["scores"][model_name] = score
        blended = sum(current["scores"].values()) / max(len(current["scores"]), 1)
        blended = min(1.0, blended)
        current["risk_score"] = round(blended, 3)

        if blended >= 0.8:
            current["risk_level"] = "CRITICAL"
        elif blended >= 0.6:
            current["risk_level"] = "HIGH"
        elif blended >= 0.35:
            current["risk_level"] = "MEDIUM"
        else:
            current["risk_level"] = "LOW"

        await self.redis.setex(f"risk:{niu}", 3600, json.dumps(current))
        return current
