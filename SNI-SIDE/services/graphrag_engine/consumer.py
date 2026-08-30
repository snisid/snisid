import asyncio, json, logging, os, uuid
from datetime import datetime
from typing import Optional

import aiokafka
from graphrag_engine import GraphRAGEngine, ReportType

logger = logging.getLogger("sniside.graphrag-consumer")

BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
GROUP_ID = os.getenv("KAFKA_GROUP_ID", "sniside-graphrag")
TRIGGER_TOPICS = [
    "sniside.graph.relationship.created",
    "sniside.graph.network.detected",
    "sniside.codis.match.positive",
    "sniside.codis.match.familial",
    "sniside.financial.network.detected",
    "sniside.cyber.campaign.detected",
    "sniside.ai.aml.risk.update",
    "sniside.fusion.cell.intelligence",
    "sniside.narcotics.route.identified",
]


class IntelligenceTrigger:
    def __init__(self, engine: GraphRAGEngine):
        self.engine = engine
        self.producer = None
        self.consumer = None
        self.running = True

    async def start(self):
        self.producer = aiokafka.AIOKafkaProducer(
            bootstrap_servers=BOOTSTRAP_SERVERS,
            acks="all", compression_type="zstd",
        )
        self.consumer = aiokafka.AIOKafkaConsumer(
            *TRIGGER_TOPICS,
            bootstrap_servers=BOOTSTRAP_SERVERS,
            group_id=GROUP_ID,
            auto_offset_reset="earliest",
            enable_auto_commit=False,
            value_deserializer=lambda v: json.loads(v.decode("utf-8")),
        )
        await self.producer.start()
        await self.consumer.start()
        logger.info(f"GraphRAG consumer started — {len(TRIGGER_TOPICS)} trigger topics")

    async def stop(self):
        self.running = False
        if self.consumer:
            await self.consumer.stop()
        if self.producer:
            await self.producer.stop()

    async def run(self):
        await self.start()
        try:
            while self.running:
                batch = await self.consumer.getmany(timeout_ms=1000, max_records=10)
                for tp, msgs in batch.items():
                    for msg in msgs:
                        await self.handle_trigger(tp.topic, msg)
                if batch:
                    await self.consumer.commit()
        except asyncio.CancelledError:
            pass
        finally:
            await self.stop()

    async def handle_trigger(self, topic: str, msg):
        try:
            v = msg.value
            entity_id = v.get("niu") or v.get("entity_id") or v.get("network_id") or v.get("profile_id")
            if not entity_id:
                logger.debug(f"No entity_id in event from {topic}")
                return

            report_type = ReportType.LINK_ANALYSIS
            if "network" in topic:
                report_type = ReportType.NETWORK_MAP
            elif "financial" in topic or "aml" in topic:
                report_type = ReportType.FINANCIAL_FLOW
            elif "match" in topic:
                report_type = ReportType.ENTITY_PROFILE

            report = await self.engine.generate_report(
                entity_id=entity_id,
                entity_type="Citizen",
                entity_label=v.get("full_name", v.get("entity_name", entity_id)),
                report_type=report_type,
            )

            await self.producer.send_and_wait(
                "sniside.ai.graph.insight",
                key=entity_id.encode(),
                value=json.dumps({
                    "event_id": str(uuid.uuid4()),
                    "graph_id": report["report_id"],
                    "entity_id": entity_id,
                    "report_type": report_type.value,
                    "executive_summary": report.get("executive_summary", ""),
                    "risk_level": report.get("risk_assessment", {}).get("overall_risk", "MEDIUM"),
                    "confidence_score": report.get("confidence_score", 0.7),
                    "description": report.get("executive_summary", "")[:500],
                    "severity": report.get("risk_assessment", {}).get("overall_risk", "MEDIUM"),
                    "timestamp": int(datetime.utcnow().timestamp() * 1000),
                    "source": "sniside-graphrag-engine",
                }).encode("utf-8"),
            )
            logger.info(f"GraphRAG insight generated for {entity_id} ({report_type.value})")
        except Exception as e:
            logger.error(f"Trigger failed for {topic}:{msg.offset}: {e}")
