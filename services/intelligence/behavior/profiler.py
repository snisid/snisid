from __future__ import annotations
import json
import numpy as np
from dataclasses import dataclass
from typing import Optional
from sklearn.ensemble import IsolationForest
from sklearn.preprocessing import StandardScaler
import redis
import logging

logger = logging.getLogger(__name__)


@dataclass
class BehaviorProfile:
    user_id: str
    anomaly_score: float
    normalized_score: float
    features_used: list[str]
    is_anomaly: bool


class BehavioralProfiler:
    FEATURES = ["transaction_hour", "transaction_day", "amount_normalized",
                 "velocity_24h", "new_location", "device_count_7d"]

    def __init__(self, redis_client: redis.Redis, contamination: float = 0.05):
        self.redis = redis_client
        self.model = IsolationForest(
            n_estimators=100,
            contamination=contamination,
            random_state=42,
            n_jobs=-1
        )
        self.scaler = StandardScaler()
        self._trained = False

    def fit(self, historical_events: list[dict]) -> None:
        if not historical_events:
            logger.warning("behavioral profiler: no historical events, using default model")
            X = np.random.randn(1000, len(self.FEATURES))
        else:
            X = np.array([self._extract_features(e) for e in historical_events])

        X_scaled = self.scaler.fit_transform(X)
        self.model.fit(X_scaled)
        self._trained = True
        logger.info(f"behavioral profiler: trained on {len(X)} samples")

    def profile(self, user_id: str, event: dict) -> BehaviorProfile:
        if not self._trained:
            self.fit([])

        features = self._extract_features(event)
        X = np.array([features])
        X_scaled = self.scaler.transform(X)

        raw_score = self.model.decision_function(X_scaled)[0]
        prediction = self.model.predict(X_scaled)[0]

        normalized = max(0.0, min(1.0, 0.5 - raw_score))

        profile = BehaviorProfile(
            user_id=user_id,
            anomaly_score=float(raw_score),
            normalized_score=normalized,
            features_used=self.FEATURES,
            is_anomaly=(prediction == -1)
        )

        self.redis.setex(
            f"snisid:behavior:{user_id}:profile",
            3600,
            json.dumps({
                "anomaly_score": profile.anomaly_score,
                "normalized_score": profile.normalized_score,
                "is_anomaly": profile.is_anomaly
            })
        )

        return profile

    def _extract_features(self, event: dict) -> list[float]:
        return [
            float(event.get("transaction_hour", 12)),
            float(event.get("transaction_day", 3)),
            float(event.get("amount_normalized", 0.5)),
            float(event.get("velocity_24h", 1.0)),
            float(event.get("new_location", 0.0)),
            float(event.get("device_count_7d", 1.0)),
        ]
