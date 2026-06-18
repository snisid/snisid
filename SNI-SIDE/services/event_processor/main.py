import asyncio, json, logging, os, sys, uuid, time
from datetime import datetime, timedelta
from typing import Optional
from collections import defaultdict

import aiokafka
from aiokafka import AIOKafkaConsumer, AIOKafkaProducer
from prometheus_client import Counter, Histogram, Gauge, start_http_server

from graph_updater import Neo4jGraphUpdater
from alert_correlator import AlertCorrelator
from risk_scorer import RiskScorer
from analytics_writer import ClickHouseAnalyticsWriter
from models import *


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(name)s] %(levelname)s: %(message)s",
)
logger = logging.getLogger("sniside.event-processor")

BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
GROUP_ID = os.getenv("KAFKA_GROUP_ID", "sniside-event-processor")
METRICS_PORT = int(os.getenv("METRICS_PORT", "9102"))

TOPICS_ALL = [
    "sniside.ncid.wanted.created", "sniside.ncid.wanted.updated",
    "sniside.ncid.warrant.issued", "sniside.ncid.case.opened",
    "sniside.biometric.match.found", "sniside.biometric.duplicate.detected",
    "sniside.codis.match.positive", "sniside.codis.match.familial",
    "sniside.missing.reported", "sniside.missing.sighting",
    "sniside.vehicle.stolen", "sniside.vehicle.wanted",
    "sniside.alpr.read", "sniside.alpr.wanted.detected",
    "sniside.alpr.anomaly", "sniside.alpr.cross.border",
    "sniside.firearm.ballistic.match",
    "sniside.border.crossing", "sniside.border.wanted.match",
    "sniside.narcotics.seizure", "sniside.narcotics.route.identified",
    "sniside.financial.transaction.suspicious", "sniside.financial.aml.alert",
    "sniside.financial.network.detected",
    "sniside.cyber.ioc.submitted", "sniside.cyber.ioc.critical",
    "sniside.cyber.incident", "sniside.cyber.campaign.detected",
    "sniside.watchlist.match", "sniside.watchlist.alert",
    "sniside.document.fraud.detected",
    "sniside.evidence.face.match", "sniside.evidence.analyzed",
    "sniside.graph.relationship.created", "sniside.graph.network.detected",
    "sniside.ai.fraud.detected", "sniside.ai.predictive.alert",
    "sniside.ai.graph.insight", "sniside.ai.aml.risk.update",
    "sniside.ai.behavioral.alert", "sniside.ai.insider.threat",
]

msg_counter = Counter("sniside_events_total", "Total events processed", ["topic"])
msg_bytes = Counter("sniside_event_bytes_total", "Total bytes processed", ["topic"])
msg_latency = Histogram("sniside_event_latency_seconds", "Event processing latency", ["topic"], buckets=[.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5])
queue_depth = Gauge("sniside_event_queue_depth", "Per-topic queue depth", ["topic"])


class EventProcessor:
    def __init__(self):
        self.graph_updater = Neo4jGraphUpdater()
        self.alert_correlator = AlertCorrelator()
        self.risk_scorer = RiskScorer()
        self.analytics = ClickHouseAnalyticsWriter()
        self.producer: Optional[AIOKafkaProducer] = None
        self.consumer: Optional[AIOKafkaConsumer] = None
        self.running = True

    async def start(self):
        self.graph_updater.start()
        self.alert_correlator.start()
        self.risk_scorer.start()

        self.producer = AIOKafkaProducer(
            bootstrap_servers=BOOTSTRAP_SERVERS,
            acks="all",
            compression_type="zstd",
            max_batch_size=512 * 1024,
            linger_ms=10,
        )
        self.consumer = AIOKafkaConsumer(
            *TOPICS_ALL,
            bootstrap_servers=BOOTSTRAP_SERVERS,
            group_id=GROUP_ID,
            auto_offset_reset="earliest",
            enable_auto_commit=False,
            max_poll_records=500,
            fetch_max_bytes=50 * 1024 * 1024,
            max_partition_fetch_bytes=10 * 1024 * 1024,
            session_timeout_ms=30000,
            heartbeat_interval_ms=3000,
            value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        )
        await self.producer.start()
        await self.consumer.start()
        logger.info(f"Event processor started — consuming {len(TOPICS_ALL)} topics")

    async def stop(self):
        self.running = False
        if self.consumer:
            await self.consumer.stop()
        if self.producer:
            await self.producer.stop()
        self.graph_updater.stop()
        self.alert_correlator.stop()
        self.risk_scorer.stop()

    async def process_event(self, msg, topic: str):
        start = time.monotonic()
        try:
            event = EventEnvelope(topic=topic, key=msg.key.decode() if msg.key else None,
                                  value=msg.value, partition=msg.partition,
                                  offset=msg.offset, timestamp=msg.timestamp)
            await self.route_event(event)
            await self.consumer.commit()
            msg_counter.labels(topic=topic).inc()
            msg_latency.labels(topic=topic).observe(time.monotonic() - start)
        except Exception as e:
            logger.error(f"Failed to process {topic}:{msg.offset}: {str(e)[:200]}", exc_info=True)

    async def route_event(self, event: EventEnvelope):
        topic = event.topic
        value = event.value
        if not value:
            logger.warning(f"Empty event on {topic}, key={event.key}")
            return

        handlers = {
            # Graph updates
            "sniside.ncid.wanted.created": self.handle_wanted_created,
            "sniside.ncid.wanted.updated": self.handle_wanted_updated,
            "sniside.ncid.case.opened": self.handle_case_opened,
            "sniside.ncid.gang.intelligence": self.handle_gang_intel,
            "sniside.biometric.match.found": self.handle_biometric_match,
            "sniside.biometric.enrolled": self.handle_biometric_enrolled,
            "sniside.codis.match.positive": self.handle_codis_match,
            "sniside.codis.match.familial": self.handle_codis_match,
            "sniside.border.crossing": self.handle_border_crossing,
            "sniside.border.wanted.match": self.handle_border_wanted_match,
            "sniside.alpr.read": self.handle_alpr_read,
            "sniside.alpr.wanted.detected": self.handle_alpr_wanted,
            "sniside.alpr.anomaly": self.handle_alpr_anomaly,
            "sniside.alpr.cross.border": self.handle_alpr_cross_border,
            "sniside.vehicle.stolen": self.handle_vehicle_stolen,
            "sniside.vehicle.wanted": self.handle_vehicle_wanted,
            "sniside.vehicle.ownership.changed": self.handle_vehicle_ownership,
            "sniside.missing.reported": self.handle_missing_reported,
            "sniside.missing.sighting": self.handle_missing_sighting,
            "sniside.missing.alert.triggered": self.handle_missing_alert,
            "sniside.firearm.ballistic.match": self.handle_ballistic_match,
            "sniside.financial.transaction.suspicious": self.handle_financial_suspicious,
            "sniside.financial.aml.alert": self.handle_aml_alert,
            "sniside.financial.network.detected": self.handle_financial_network,
            "sniside.cyber.ioc.submitted": self.handle_cyber_ioc,
            "sniside.cyber.ioc.critical": self.handle_cyber_ioc_critical,
            "sniside.cyber.incident": self.handle_cyber_incident,
            "sniside.cyber.campaign.detected": self.handle_cyber_campaign,
            "sniside.watchlist.match": self.handle_watchlist_match,
            "sniside.watchlist.alert": self.handle_watchlist_alert,
            "sniside.document.fraud.detected": self.handle_document_fraud,
            "sniside.evidence.face.match": self.handle_evidence_face_match,
            "sniside.evidence.analyzed": self.handle_evidence_analyzed,
            "sniside.narcotics.seizure": self.handle_narcotics_seizure,
            "sniside.narcotics.route.identified": self.handle_narcotics_route,
            "sniside.ai.fraud.detected": self.handle_ai_fraud,
            "sniside.ai.predictive.alert": self.handle_ai_predictive,
            "sniside.ai.graph.insight": self.handle_ai_graph_insight,
            "sniside.ai.aml.risk.update": self.handle_aml_risk_update,
            "sniside.ai.behavioral.alert": self.handle_behavioral_alert,
            "sniside.ai.insider.threat": self.handle_insider_threat,
            "sniside.graph.relationship.created": self.handle_graph_relationship,
            "sniside.graph.network.detected": self.handle_graph_network,
        }

        handler = handlers.get(topic)
        if handler:
            await handler(event)
        else:
            logger.debug(f"No handler for {topic}")

    # ==================== NCID HANDLERS ====================

    async def handle_wanted_created(self, event: EventEnvelope):
        v = event.value
        niu = v["niu"]
        await self.graph_updater.upsert_citizen(niu, v)
        await self.graph_updater.create_relationship(niu, "HAS_WARRANT", event.key or niu, {"type": "warrant", "count": v.get("warrants_active", 0)})
        for alias in v.get("aliases", []):
            await self.graph_updater.upsert_alias(niu, alias)
            await self.graph_updater.create_relationship(niu, "HAS_ALIAS", f"alias:{alias}", {})
        for gang in v.get("gang_affiliations", []):
            await self.graph_updater.create_relationship(niu, "MEMBER_OF", f"gang:{gang}", {"role": "member"})
        risk = await self.risk_scorer.score_person(niu, v)
        await self.analytics.write_event("ncid", event)

    async def handle_wanted_updated(self, event: EventEnvelope):
        v = event.value
        await self.graph_updater.upsert_citizen(v["niu"], v)
        if v.get("risk_level") == "CRITICAL":
            await self.emit_alert("RISK_ESCALATION", v["niu"], f"Risk escalated to CRITICAL: {v.get('full_name', '')}", "HIGH")
        await self.analytics.write_event("ncid", event)

    async def handle_case_opened(self, event: EventEnvelope):
        v = event.value
        case_id = v["case_id"]
        await self.graph_updater.create_case(case_id, v)
        for niu in v.get("subjects", []):
            await self.graph_updater.create_relationship(niu, "INVOLVED_IN", f"case:{case_id}", {"role": v.get("role", "subject")})

    async def handle_gang_intel(self, event: EventEnvelope):
        v = event.value
        gang_name = v["gang_name"]
        await self.graph_updater.create_gang(gang_name, v)
        for member in v.get("members", []):
            await self.graph_updater.create_relationship(member, "MEMBER_OF", f"gang:{gang_name}", {"role": v.get("membro_role", "member")})
        for rival in v.get("rival_gangs", []):
            await self.graph_updater.create_relationship(f"gang:{gang_name}", "RIVAL", f"gang:{rival}", {})

    # ==================== BIOMETRIC HANDLERS ====================

    async def handle_biometric_match(self, event: EventEnvelope):
        v = event.value
        niu = v.get("niu")
        match_type = v.get("match_type", "face")
        matched_niu = v.get("matched_niu")
        certainty = v.get("confidence", 0.0)
        await self.graph_updater.create_relationship(niu, f"BIOMETRIC_{match_type.upper()}_MATCH", matched_niu, {"confidence": certainty, "timestamp": event.timestamp})
        if certainty > 0.95:
            await self.emit_alert("BIOMETRIC_HIGH_CONFIDENCE", niu, f"High-confidence {match_type} match: {niu} ↔ {matched_niu} ({certainty:.2%})", "HIGH", event)
        await self.analytics.write_event("biometric", event)

    async def handle_biometric_enrolled(self, event: EventEnvelope):
        v = event.value
        niu = v["niu"]
        bio_types = v.get("biometric_types", [])
        for bt in bio_types:
            await self.graph_updater.create_relationship(niu, "HAS_BIOMETRIC", f"{bt}:{niu}", {"type": bt, "enrolled_at": event.timestamp})

    # ==================== CODIS HANDLERS ====================

    async def handle_codis_match(self, event: EventEnvelope):
        v = event.value
        profile_id = v.get("profile_id")
        matched_profile = v.get("matched_profile")
        match_type = v.get("match_type", "direct")
        await self.analytics.write_event("codis", event)
        await self.emit_alert("DNA_MATCH", profile_id, f"DNA {match_type} match: {profile_id} ↔ {matched_profile}",
                              "CRITICAL" if match_type == "direct" else "HIGH", event)

    # ==================== BORDER HANDLERS ====================

    async def handle_border_crossing(self, event: EventEnvelope):
        v = event.value
        niu = v.get("niu")
        await self.graph_updater.create_border_crossing(niu, event.key or niu, v)
        if niu:
            await self.graph_updater.create_relationship(niu, "TRAVELLED_TO", f"country:{v.get('destination_country', 'UNKNOWN')}", {
                "port": v.get("port_of_entry"), "date": event.timestamp})
        await self.analytics.write_event("border", event)

    async def handle_border_wanted_match(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("BORDER_WANTED", v.get("niu", ""), f"Wanted person detected at border: {v.get('full_name', '')} @ {v.get('port_of_entry', '')}",
                              "CRITICAL", event)

    # ==================== ALPR HANDLERS ====================

    async def handle_alpr_read(self, event: EventEnvelope):
        v = event.value
        plate = v.get("plate")
        if plate:
            await self.graph_updater.upsert_vehicle(plate, v)
            if v.get("owner_niu"):
                await self.graph_updater.create_relationship(v["owner_niu"], "OWNS", f"vehicle:{plate}", {"make": v.get("make"), "model": v.get("model"), "year": v.get("year")})
        await self.analytics.write_event("alpr", event)

    async def handle_alpr_wanted(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("ALPR_WANTED_HIT", v.get("plate", ""), f"ALPR hit: wanted vehicle {v.get('plate')} at {v.get('location', '')}",
                              "CRITICAL", event)
        await self.analytics.write_event("alpr", event)

    async def handle_alpr_anomaly(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("ALPR_ANOMALY", v.get("plate", ""), f"ALPR anomaly: {v.get('description', '')}", "MEDIUM", event)

    async def handle_alpr_cross_border(self, event: EventEnvelope):
        v = event.value
        plate = v.get("plate")
        if plate:
            await self.graph_updater.create_relationship(f"vehicle:{plate}", "CROSSED_BORDER", f"country:{v.get('destination', 'UNKNOWN')}", {
                "port": v.get("port"), "timestamp": event.timestamp})

    # ==================== VEHICLE HANDLERS ====================

    async def handle_vehicle_stolen(self, event: EventEnvelope):
        v = event.value
        plate = v.get("plate")
        if plate:
            await self.graph_updater.upsert_vehicle(plate, v)
            await self.graph_updater.create_relationship(f"vehicle:{plate}", "FLAGGED_AS", "stolen", {"theft_date": v.get("theft_date"), "reported_by": v.get("agency")})
        await self.emit_alert("VEHICLE_STOLEN", plate, f"Vehicle reported stolen: {plate}", "HIGH", event)

    async def handle_vehicle_wanted(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("VEHICLE_WANTED", v.get("plate"), f"Vehicle wanted: {v.get('plate')} — {v.get('reason', '')}", "CRITICAL", event)

    async def handle_vehicle_ownership(self, event: EventEnvelope):
        v = event.value
        plate = v.get("plate")
        old_owner = v.get("previous_owner_niu")
        new_owner = v.get("new_owner_niu")
        if plate and old_owner:
            await self.graph_updater.remove_relationship(old_owner, "OWNS", f"vehicle:{plate}")
        if plate and new_owner:
            await self.graph_updater.create_relationship(new_owner, "OWNS", f"vehicle:{plate}", {"from": v.get("transfer_date")})

    # ==================== MISSING PERSONS ====================

    async def handle_missing_reported(self, event: EventEnvelope):
        v = event.value
        niu = v.get("niu")
        if niu:
            await self.graph_updater.upsert_citizen(niu, {"full_name": v.get("full_name"), "status": "MISSING"})
            await self.graph_updater.create_relationship(niu, "REPORTED_MISSING", f"case:{v.get('case_id', '')}", {
                "missing_since": v.get("missing_since"), "reported_by": v.get("reported_by")})
        await self.emit_alert("MISSING_REPORTED", niu, f"Missing person reported: {v.get('full_name', '')}", "HIGH", event)
        await self.analytics.write_event("missing", event)

    async def handle_missing_sighting(self, event: EventEnvelope):
        v = event.value
        niu = v.get("niu")
        location = v.get("location", "")
        await self.emit_alert("MISSING_SIGHTING", niu, f"Missing person sighted: {v.get('full_name', '')} at {location}",
                              "HIGH" if v.get("verified") else "MEDIUM", event)

    async def handle_missing_alert(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("AMBER_SILVER_ALERT", v.get("niu"), v.get("message", ""), "CRITICAL", event)

    # ==================== FIREARMS ====================

    async def handle_ballistic_match(self, event: EventEnvelope):
        v = event.value
        firearm_id = v.get("firearm_id")
        case_id = v.get("case_id")
        if firearm_id and case_id:
            await self.graph_updater.create_relationship(f"firearm:{firearm_id}", "EVIDENCE_IN", f"case:{case_id}", {"match_type": v.get("match_type"), "confidence": v.get("confidence")})
        await self.emit_alert("BALLISTIC_MATCH", firearm_id, f"Ballistic match: {firearm_id} → Case {case_id}", "HIGH", event)

    # ==================== FINANCIAL ====================

    async def handle_financial_suspicious(self, event: EventEnvelope):
        v = event.value
        sender = v.get("sender_niu")
        beneficiary = v.get("beneficiary_niu")
        amount = v.get("amount", 0)
        if sender and beneficiary:
            tx_id = event.key or str(uuid.uuid4())
            await self.graph_updater.create_bank_account(f"acct:{sender}", v.get("sender_account"), v)
            await self.graph_updater.create_bank_account(f"acct:{beneficiary}", v.get("beneficiary_account"), v)
            await self.graph_updater.create_relationship(f"acct:{sender}", "TRANSFERRED_TO", f"acct:{beneficiary}", {
                "amount": amount, "currency": v.get("currency"), "date": event.timestamp, "tx_id": tx_id})
            await self.graph_updater.create_relationship(sender, "CONTROLLED_BY", f"acct:{sender}", {})
            await self.graph_updater.create_relationship(beneficiary, "CONTROLLED_BY", f"acct:{beneficiary}", {})
        await self.analytics.write_event("financial", event)

    async def handle_aml_alert(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("AML_ALERT", v.get("entity_id"), v.get("description", ""), v.get("severity", "HIGH"), event)

    async def handle_financial_network(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("FINANCIAL_NETWORK", v.get("network_id"), v.get("description", ""), "CRITICAL", event)

    # ==================== CYBER ====================

    async def handle_cyber_ioc(self, event: EventEnvelope):
        v = event.value
        ioc_type = v.get("ioc_type", "ip")
        ioc_value = v.get("ioc_value", "")
        if ioc_type == "ip":
            await self.graph_updater.create_ip(ioc_value, v)
        elif ioc_type == "domain":
            await self.graph_updater.create_domain(ioc_value, v)
        elif ioc_type == "wallet":
            await self.graph_updater.create_wallet(ioc_value, v)
        await self.analytics.write_event("cyber", event)

    async def handle_cyber_ioc_critical(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("CRITICAL_IOC", v.get("ioc_value"), f"Critical IOC: {v.get('ioc_type', '')} — {v.get('ioc_value', '')}", "CRITICAL", event)

    async def handle_cyber_incident(self, event: EventEnvelope):
        v = event.value
        incident_id = v.get("incident_id")
        await self.emit_alert("CYBER_INCIDENT", incident_id, f"Cyber incident: {v.get('title', '')} — {v.get('description', '')}",
                              v.get("severity", "HIGH"), event)
        await self.analytics.write_event("cyber", event)

    async def handle_cyber_campaign(self, event: EventEnvelope):
        v = event.value
        await self.graph_updater.create_relationship(v.get("campaign_id"), "TARGETS", v.get("target_sector"), {})
        await self.emit_alert("CYBER_CAMPAIGN", v.get("campaign_id"), f"Campaign detected: {v.get('name', '')}", "CRITICAL", event)

    # ==================== WATCHLIST ====================

    async def handle_watchlist_match(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("WATCHLIST_MATCH", v.get("entity_id"), f"Watchlist match: {v.get('entity_name', '')}", v.get("severity", "HIGH"), event)
        await self.analytics.write_event("watchlist", event)

    async def handle_watchlist_alert(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("WATCHLIST_ALERT", v.get("entity_id"), v.get("message", ""), v.get("severity", "CRITICAL"), event)

    # ==================== DOCUMENTS ====================

    async def handle_document_fraud(self, event: EventEnvelope):
        v = event.value
        doc_id = v.get("document_id")
        niu = v.get("niu")
        if doc_id and niu:
            await self.graph_updater.create_relationship(niu, "USED_DOCUMENT", f"doc:{doc_id}", {"status": "fraud", "fraud_type": v.get("fraud_type")})
        await self.emit_alert("DOCUMENT_FRAUD", doc_id, f"Document fraud detected: {v.get('document_type', '')} — {v.get('fraud_type', '')}", "HIGH", event)

    # ==================== EVIDENCE ====================

    async def handle_evidence_face_match(self, event: EventEnvelope):
        v = event.value
        evidence_id = v.get("evidence_id")
        niu = v.get("niu")
        confidence = v.get("confidence", 0.0)
        if evidence_id and niu:
            await self.graph_updater.create_relationship(niu, "IDENTIFIED_IN", f"evidence:{evidence_id}", {"confidence": confidence, "method": "face_recognition"})
        await self.emit_alert("EVIDENCE_FACE_MATCH", niu, f"Face match in evidence {evidence_id}: confidence {confidence:.2%}", "HIGH" if confidence > 0.9 else "MEDIUM", event)

    async def handle_evidence_analyzed(self, event: EventEnvelope):
        v = event.value
        evidence_id = v.get("evidence_id")
        await self.graph_updater.create_evidence(evidence_id, v)
        for niu in v.get("persons_identified", []):
            await self.graph_updater.create_relationship(niu, "LINKED_TO_EVIDENCE", f"evidence:{evidence_id}", {})

    # ==================== NARCOTICS ====================

    async def handle_narcotics_seizure(self, event: EventEnvelope):
        v = event.value
        seizure_id = v.get("seizure_id")
        await self.emit_alert("NARCOTICS_SEIZURE", seizure_id, f"Seizure: {v.get('drug_type')} — {v.get('quantity_kg', 0)}kg @ {v.get('location', '')}",
                              v.get("severity", "HIGH"), event)
        for niu in v.get("suspects", []):
            await self.graph_updater.create_relationship(niu, "INVOLVED_IN", f"seizure:{seizure_id}", {"role": "suspect"})
        await self.analytics.write_event("narcotics", event)

    async def handle_narcotics_route(self, event: EventEnvelope):
        v = event.value
        route_id = v.get("route_id")
        for loc in v.get("waypoints", []):
            await self.graph_updater.create_relationship(f"route:{route_id}", "PASSES_THROUGH", f"location:{loc}", {})

    # ==================== AI FUSION HANDLERS ====================

    async def handle_ai_fraud(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("AI_FRAUD", v.get("entity_id"), v.get("description", ""), v.get("severity", "CRITICAL"), event)
        if v.get("entities_involved"):
            for ent in v["entities_involved"]:
                if ent.get("entity_type") == "citizen":
                    risk = await self.risk_scorer.apply_ai_score(ent["entity_id"], "fraud", v.get("confidence_score", 0))
                    if risk and risk.get("risk_level") in ("CRITICAL", "HIGH"):
                        await self.emit_alert("RISK_SCORE_CHANGED", ent["entity_id"],
                                              f"Risk score updated: {ent['entity_name']} → {risk['risk_level']} (fraud: {v.get('confidence_score', 0):.2f})",
                                              risk["risk_level"], event)

    async def handle_ai_predictive(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("PREDICTIVE_CRIME", v.get("zone_id"), v.get("description", ""), v.get("severity", "MEDIUM"), event)
        await self.analytics.write_event("ai", event)

    async def handle_ai_graph_insight(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("GRAPH_INSIGHT", v.get("graph_id"), v.get("description", ""), v.get("severity", "HIGH"), event)
        await self.analytics.write_event("ai", event)

    async def handle_aml_risk_update(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("AML_RISK_UPDATE", v.get("entity_id"), v.get("description", ""), v.get("severity", "HIGH"), event)
        await self.analytics.write_event("ai", event)

    async def handle_behavioral_alert(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("BEHAVIORAL_ANOMALY", v.get("entity_id"), v.get("description", ""), v.get("severity", "MEDIUM"), event)

    async def handle_insider_threat(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("INSIDER_THREAT", v.get("user_id"), v.get("description", ""), "CRITICAL", event)

    # ==================== GRAPH HANDLERS ====================

    async def handle_graph_relationship(self, event: EventEnvelope):
        v = event.value
        await self.analytics.write_event("graph", event)

    async def handle_graph_network(self, event: EventEnvelope):
        v = event.value
        await self.emit_alert("NETWORK_DETECTED", v.get("network_id"), v.get("description", ""), v.get("severity", "CRITICAL"), event)

    # ==================== CROSS-CUTTING ====================

    async def emit_alert(self, alert_type: str, entity_id: str, message: str, severity: str, source_event: EventEnvelope = None):
        alert = {
            "event_id": str(uuid.uuid4()),
            "alert_type": alert_type,
            "severity": severity,
            "title": f"{alert_type}: {entity_id}",
            "description": message,
            "entities_involved": [{"entity_type": "unknown", "entity_id": entity_id, "entity_name": entity_id}],
            "confidence_score": 1.0 if severity == "CRITICAL" else 0.8 if severity == "HIGH" else 0.5,
            "ai_model": "event-processor",
            "model_version": "1.0",
            "timestamp": int(time.time() * 1000),
            "correlation_id": str(uuid.uuid4()),
            "source": "sniside-event-processor",
        }
        await self.producer.send_and_wait("sniside.rtcc.incident", key=entity_id.encode(), value=json.dumps(alert).encode("utf-8"))
        logger.info(f"Alert emitted: [{severity}] {alert_type} — {message[:120]}")

        correlated = self.alert_correlator.add_event(alert_type, entity_id, severity, source_event.timestamp if source_event else int(time.time() * 1000))
        if correlated:
            correlation_alert = {
                "event_id": str(uuid.uuid4()),
                "alert_type": "CORRELATION_ALERT",
                "severity": "CRITICAL",
                "title": f"Multi-domain correlation: {entity_id}",
                "description": f"Entity {entity_id} triggered {correlated['alert_count']} alerts across {correlated['domain_count']} domains: {correlated['domains']}",
                "entities_involved": [{"entity_type": "unknown", "entity_id": entity_id, "entity_name": entity_id}],
                "confidence_score": 0.95,
                "ai_model": "alert-correlator",
                "model_version": "1.0",
                "timestamp": int(time.time() * 1000),
                "correlation_id": str(uuid.uuid4()),
                "source": "sniside-alert-correlator",
            }
            await self.producer.send_and_wait("sniside.rtcc.incident", key=entity_id.encode(), value=json.dumps(correlation_alert).encode("utf-8"))
            logger.info(f"CORRELATION ALERT: {entity_id} — {correlated['domains']}")

    async def run(self):
        await self.start()
        try:
            while self.running:
                batch = await self.consumer.getmany(timeout_ms=1000, max_records=500)
                for tp, msgs in batch.items():
                    queue_depth.labels(topic=tp.topic).set(len(msgs))
                    tasks = [self.process_event(m, tp.topic) for m in msgs]
                    await asyncio.gather(*tasks)
        except asyncio.CancelledError:
            pass
        finally:
            await self.stop()


async def main():
    start_http_server(METRICS_PORT)
    logger.info(f"Metrics HTTP server on :{METRICS_PORT}")
    processor = EventProcessor()
    await processor.run()


if __name__ == "__main__":
    asyncio.run(main())
