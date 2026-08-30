# -*- coding: utf-8 -*-
"""
ORCHESTRATEUR CENTRAL DE WORKFLOWS SNISID — SNISID v4.0
Coordonne le Bus d'Événements, le Gestionnaire de Dossiers, le Service de Tâches Humaines et le Moteur de SLA.
"""

import os
import sys
import logging

# Bootstrap path resolution
current_dir = os.path.dirname(os.path.abspath(__file__))
sys.path.append(current_dir)
sys.path.append(os.path.join(current_dir, "Event-Driven"))
sys.path.append(os.path.join(current_dir, "Case-Management"))
sys.path.append(os.path.join(current_dir, "Human-Tasks"))
sys.path.append(os.path.join(current_dir, "SLA"))
sys.path.append(os.path.join(current_dir, "Civil-Registry"))
sys.path.append(os.path.join(current_dir, "Identity"))
sys.path.append(os.path.join(current_dir, "Justice"))
sys.path.append(os.path.join(current_dir, "Police"))
sys.path.append(os.path.join(current_dir, "Offline"))

from event_bus import global_event_bus, CloudEvent
from case_manager import global_case_manager
from human_task_service import global_human_task_service
from sla_engine import init_sla_engine, global_sla_engine
from civil_registry_workflows import CivilRegistryWorkflowManager
from identity_workflows import IdentityWorkflowManager
from justice_workflows import JusticeWorkflowManager
from police_workflows import PoliceWorkflowManager
from offline_sync import global_offline_sync_engine

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [ORCHESTRATOR] - %(levelname)s - %(message)s')

class NationalWorkflowFactory:
    def __init__(self):
        # Initialisation du moteur de SLA lié au gestionnaire de cas et de tâches
        self.sla_engine = init_sla_engine(global_human_task_service, global_case_manager)
        self._setup_event_handlers()
        logging.info("National Workflow Factory initialisée et prête à orchestrer les flux étatiques.")

    def _setup_event_handlers(self):
        """
        Configure les écouteurs d'événements Kafka pour réagir automatiquement
        aux changements de statut des tâches humaines, des dossiers et des alertes.
        """
        # S'abonner aux événements d'approbation et rejet de tâches
        global_event_bus.subscribe(
            topic_name="workflow.audit",
            consumer_name="National_Orchestrator_Audit_Consumer",
            callback=self._handle_audit_events
        )
        
        global_event_bus.subscribe(
            topic_name="system.alarms",
            consumer_name="National_Orchestrator_Alarms_Consumer",
            callback=self._handle_alarm_events
        )

    def _handle_audit_events(self, event: CloudEvent):
        """
        Gère les transitions d'état de manière automatisée pour simuler les routes BPMN.
        """
        logging.info(f"[ROUTING_ENGINE] Traitement de l'événement d'audit: {event.type}")
        
        # Réaction à l'approbation d'un dossier par validation humaine
        if event.type == "workflow.approved":
            case_id = event.data["payload"]["caseId"]
            task_id = event.data["payload"]["taskId"]
            comments = event.data["payload"]["comments"]
            
            case = global_case_manager.get_case(case_id)
            if case:
                # Désactiver le timer SLA de la tâche résolue
                self.sla_engine.deactivate_timer(task_id)
                
                # Transitionner l'état du dossier vers l'état approuvé final
                case.transition_status("COMPLETED_APPROVED", "SYSTEM_ORCHESTRATOR", f"Validé par l'officier. Commentaire: {comments}")
                
                # Simuler la production du livrable national final
                creation_event = CloudEvent(
                    event_type="workflow.state.persisted_registry",
                    subject=case_id,
                    data={
                        "correlationId": case.correlation_id,
                        "operatorId": "SYSTEM_ORCHESTRATOR",
                        "agency": case.creator_agency,
                        "payload": {
                            "caseId": case_id,
                            "finalStatus": "REGISTERED_IN_NATIONAL_LEDGER",
                            "digitalSignature": f"SECURE_CERT_HASH_{case_id}"
                        }
                    }
                )
                global_event_bus.publish("workflow.audit", creation_event)

        # Réaction au rejet d'un dossier
        elif event.type == "workflow.rejected":
            case_id = event.data["payload"]["caseId"]
            task_id = event.data["payload"]["taskId"]
            comments = event.data["payload"]["comments"]
            
            case = global_case_manager.get_case(case_id)
            if case:
                self.sla_engine.deactivate_timer(task_id)
                case.transition_status("REJECTED_TERMINATED", "SYSTEM_ORCHESTRATOR", f"Rejeté par l'officier. Motif: {comments}")

    def _handle_alarm_events(self, event: CloudEvent):
        logging.critical(f"[ALERT_MONITOR] ALARME CRITIQUE DÉTECTÉE - TYPE: {event.type}")
        # Traitement d'urgence pour les escalades nationales
        if event.type == "system.consumer.failure":
            logging.error(f"[SYSTEM_FATAL] Échec critique d'un consommateur Kafka: {event.subject}")
        elif event.type == "workflow.emergency.triggered":
            logging.critical(f"[RED_ALERT] Intervention manuelle requise! Dossier: {event.subject}. Raison: {event.data['payload']['reason']}")


# Instance globale d'orchestration
global_workflow_factory = NationalWorkflowFactory()
