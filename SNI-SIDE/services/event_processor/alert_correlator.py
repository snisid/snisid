import logging, os, time
from collections import defaultdict
from typing import Optional
from models import CorrelationWindow

logger = logging.getLogger("sniside.alert-correlator")

CORRELATION_WINDOW_MS = int(os.getenv("CORRELATION_WINDOW_MS", "300000"))  # 5 minutes
MIN_ALERTS_FOR_CORRELATION = int(os.getenv("MIN_ALERTS_FOR_CORRELATION", "3"))
MIN_DOMAINS_FOR_CORRELATION = int(os.getenv("MIN_DOMAINS_FOR_CORRELATION", "2"))


class AlertCorrelator:
    def __init__(self):
        self.windows: dict[str, CorrelationWindow] = {}
        self.domain_topic_map = {
            "WANTED_MATCH": "NCID", "BIOMETRIC_MATCH": "BIOMETRICS",
            "DNA_MATCH": "CODIS", "MISSING": "MISSING",
            "ALPR_WANTED_HIT": "ALPR", "ALPR_ANOMALY": "ALPR",
            "VEHICLE_STOLEN": "VEHICLE", "BORDER_WANTED": "BORDER",
            "BALLISTIC_MATCH": "FIREARMS", "AML_ALERT": "FINANCIAL",
            "AML_RISK_UPDATE": "FINANCIAL", "FINANCIAL_NETWORK": "FINANCIAL",
            "AI_FRAUD": "AI", "PREDICTIVE_CRIME": "AI",
            "GRAPH_INSIGHT": "AI", "BEHAVIORAL_ANOMALY": "AI",
            "INSIDER_THREAT": "AI", "CRITICAL_IOC": "CYBER",
            "CYBER_INCIDENT": "CYBER", "CYBER_CAMPAIGN": "CYBER",
            "WATCHLIST_MATCH": "WATCHLIST", "WATCHLIST_ALERT": "WATCHLIST",
            "DOCUMENT_FRAUD": "DOCUMENT", "EVIDENCE_FACE_MATCH": "EVIDENCE",
            "NARCOTICS_SEIZURE": "NARCOTICS",
            "NETWORK_DETECTED": "GRAPH", "RISK_ESCALATION": "NCID",
            "MISSING_SIGHTING": "MISSING", "AMBER_SILVER_ALERT": "MISSING",
        }

    def start(self):
        logger.info(f"AlertCorrelator started — window={CORRELATION_WINDOW_MS}ms, min_alerts={MIN_ALERTS_FOR_CORRELATION}, min_domains={MIN_DOMAINS_FOR_CORRELATION}")

    def stop(self):
        pass

    def add_event(self, alert_type: str, entity_id: str, severity: str, timestamp: int) -> Optional[dict]:
        now = time.time() * 1000
        self._purge(now)

        if entity_id not in self.windows:
            self.windows[entity_id] = CorrelationWindow(entity_id=entity_id)

        window = self.windows[entity_id]
        domain = self.domain_topic_map.get(alert_type, "OTHER")

        if domain not in window.alerts:
            window.alerts[domain] = []
        window.alerts[domain].append({
            "alert_type": alert_type,
            "severity": severity,
            "timestamp": timestamp,
        })
        window.last_seen = now
        if window.first_seen == 0:
            window.first_seen = now

        if window.alert_count >= MIN_ALERTS_FOR_CORRELATION and window.domain_count >= MIN_DOMAINS_FOR_CORRELATION:
            return {
                "entity_id": entity_id,
                "alert_count": window.alert_count,
                "domain_count": window.domain_count,
                "domains": ",".join(window.domains),
                "first_seen": window.first_seen,
                "last_seen": window.last_seen,
            }
        return None

    def _purge(self, now: int):
        expired = [eid for eid, w in self.windows.items()
                   if now - w.last_seen > CORRELATION_WINDOW_MS]
        for eid in expired:
            del self.windows[eid]
        if expired:
            logger.debug(f"Purged {len(expired)} correlation windows")
