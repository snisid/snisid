import uuid
from dataclasses import dataclass, field
from typing import Any, Optional


@dataclass
class EventEnvelope:
    topic: str
    key: Optional[str]
    value: dict[str, Any]
    partition: int
    offset: int
    timestamp: int
    event_id: str = field(default_factory=lambda: str(uuid.uuid4()))


@dataclass
class AlertEvent:
    alert_id: str
    alert_type: str
    entity_id: str
    severity: str
    title: str
    description: str
    timestamp: int
    source: str

    def to_dict(self) -> dict:
        return {
            "event_id": self.alert_id,
            "alert_type": self.alert_type,
            "severity": self.severity,
            "title": self.title,
            "description": self.description,
            "timestamp": self.timestamp,
            "source": self.source,
        }


@dataclass
class CorrelationWindow:
    entity_id: str
    alerts: dict[str, list[dict]] = field(default_factory=dict)
    first_seen: int = 0
    last_seen: int = 0

    @property
    def alert_count(self) -> int:
        return sum(len(v) for v in self.alerts.values())

    @property
    def domains(self) -> list[str]:
        return list(self.alerts.keys())

    @property
    def domain_count(self) -> int:
        return len(self.alerts)
