# -*- coding: utf-8 -*-
"""
MOTEUR SIMULÉ KAFKA EVENT BUS — SNISID v4.0
Gère la validation des schémas CloudEvents, la publication de messages et la souscription réactive.
"""

import uuid
from datetime import datetime
import json
import logging

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [KAFKA] - %(levelname)s - %(message)s')

class CloudEvent:
    def __init__(self, event_type, subject, data, source="snisid/national/workflow-factory"):
        self.specversion = "1.0"
        self.id = str(uuid.uuid4())
        self.source = source
        self.type = event_type
        self.subject = subject
        self.time = datetime.now().isoformat()
        self.datacontenttype = "application/json"
        
        # Simulation d'une signature cryptographique pour garantir l'intégrité nationale
        self.datacryptosignature = f"SIG_ED25519_{uuid.uuid4().hex[:16].upper()}"
        
        # Données de l'événement standardisées
        self.data = {
            "correlationId": data.get("correlationId", str(uuid.uuid4())),
            "operatorId": data.get("operatorId", "SYSTEM_AUTO"),
            "agency": data.get("agency", "SNISID-CENTRAL"),
            "payload": data.get("payload", {})
        }

    def to_json(self):
        return json.dumps({
            "specversion": self.specversion,
            "id": self.id,
            "source": self.source,
            "type": self.type,
            "subject": self.subject,
            "time": self.time,
            "datacontenttype": self.datacontenttype,
            "datacryptosignature": self.datacryptosignature,
            "data": self.data
        }, indent=2, ensure_ascii=False)


class KafkaEventBus:
    def __init__(self):
        self.topics = {
            "civil.registry.events": [],
            "identity.events": [],
            "justice.events": [],
            "police.events": [],
            "workflow.audit": [],
            "system.alarms": []
        }
        self.subscribers = {}

    def create_topic(self, topic_name):
        if topic_name not in self.topics:
            self.topics[topic_name] = []
            logging.info(f"Topic créé avec succès: {topic_name}")

    def subscribe(self, topic_name, consumer_name, callback):
        if topic_name not in self.subscribers:
            self.subscribers[topic_name] = []
        self.subscribers[topic_name].append({
            "consumer": consumer_name,
            "callback": callback
        })
        logging.info(f"Abonnement enregistré: {consumer_name} sur le topic {topic_name}")

    def publish(self, topic_name, event: CloudEvent):
        if topic_name not in self.topics:
            self.create_topic(topic_name)
        
        # Enregistrement dans le topic
        self.topics[topic_name].append(event)
        logging.info(f"Message publié sur [{topic_name}] ID={event.id} TYPE={event.type}")
        
        # Envoi de l'événement d'audit
        if topic_name != "workflow.audit":
            audit_event = CloudEvent(
                event_type="workflow.audit.event.published",
                subject=event.id,
                data={
                    "correlationId": event.data["correlationId"],
                    "operatorId": event.data["operatorId"],
                    "agency": event.data["agency"],
                    "payload": {
                        "topic": topic_name,
                        "original_type": event.type,
                        "signature": event.datacryptosignature
                    }
                }
            )
            self.topics["workflow.audit"].append(audit_event)

        # Distribution synchrone aux abonnés (simulation asynchrone)
        if topic_name in self.subscribers:
            for sub in self.subscribers[topic_name]:
                try:
                    sub["callback"](event)
                except Exception as e:
                    logging.error(f"Erreur d'exécution de l'abonné [{sub['consumer']}] pour l'événement {event.id}: {str(e)}")
                    # Alarme système en cas d'échec du consommateur
                    self.publish("system.alarms", CloudEvent(
                        event_type="system.consumer.failure",
                        subject=sub["consumer"],
                        data={
                            "correlationId": event.data["correlationId"],
                            "payload": {"error": str(e), "failed_event_id": event.id}
                        }
                    ))


# Instance singleton globale pour la simulation
global_event_bus = KafkaEventBus()
