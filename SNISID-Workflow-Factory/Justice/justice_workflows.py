# -*- coding: utf-8 -*-
"""
WORKFLOWS JUSTICE NATIONALE — SNISID v4.0
Industrialise la gestion des dossiers judiciaires (Criminal Case), des Mandats de Justice et de l'administration pénitentiaire.
"""

import os
import sys
import logging

# Bootstrap path resolution
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
sys.path.append(parent_dir)
sys.path.append(os.path.join(parent_dir, "Event-Driven"))
sys.path.append(os.path.join(parent_dir, "Case-Management"))
sys.path.append(os.path.join(parent_dir, "Human-Tasks"))
sys.path.append(os.path.join(parent_dir, "SLA"))

from event_bus import global_event_bus, CloudEvent
from case_manager import global_case_manager
from human_task_service import global_human_task_service
import sla_engine

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [JUSTICE-WORKFLOW] - %(levelname)s - %(message)s')

class JusticeWorkflowManager:
    @staticmethod
    def open_criminal_case(suspect_identity, infraction_details, court_id, judge_operator):
        """
        Démarre un dossier criminel judiciaire (Criminal Case Record).
        Classification de sécurité maximale: SECRET-DEFENSE
        """
        logging.info(f"Ouverture d'un dossier criminel judiciaire pour le suspect: {suspect_identity.get('first_name')} {suspect_identity.get('last_name')}")
        
        case_data = {
            "suspect": suspect_identity,
            "infraction": infraction_details,
            "court_id": court_id,
            "judge_operator": judge_operator,
            "status_history": ["INITIATED"]
        }
        
        case = global_case_manager.create_case(
            case_type="JUSTICE_CRIMINAL_CASE",
            creator_agency="MINISTERE_JUSTICE_TRIBUNAL",
            security_level="SECRET-DEFENSE",
            metadata=case_data
        )

        # Les dossiers criminels requièrent des vérifications sous 24h
        task = global_human_task_service.create_task(
            task_name="Instruction Initiale Dossier Criminel",
            case_id=case.case_id,
            role_required="JUGE_D_INSTRUCTION",
            correlation_id=case.correlation_id,
            sla_minutes=1440, # 24 heures
            security_classification="SECRET-DEFENSE"
        )
        
        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority="P2")

        # Notification sur le topic justice
        event = CloudEvent(
            event_type="judicial.case.opened",
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": judge_operator,
                "agency": "MINISTERE_JUSTICE_TRIBUNAL",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "suspectName": f"{suspect_identity.get('first_name')} {suspect_identity.get('last_name')}",
                    "infraction": infraction_details.get("title")
                }
            }
        )
        global_event_bus.publish("justice.events", event)
        case.transition_status("PENDING_JUDICIAL_INSTRUCTION", judge_operator, f"Ouverture par le Juge {judge_operator}")
        return case, task

    @staticmethod
    def issue_warrant(warrant_type, suspect_id, warrant_reason, issuing_judge):
        """
        Gère l'émission d'un mandat judiciaire (Mandat d'arrêt, de dépôt, d'amener, de perquisition).
        Niveau de priorité maximal (P1 / Urgence Nationale) pour les mandats d'arrêt.
        """
        logging.info(f"Émission d'un mandat: {warrant_type} pour suspect {suspect_id}")
        
        case_data = {
            "warrant_type": warrant_type,
            "suspect_id": suspect_id,
            "reason": warrant_reason,
            "issuing_judge": issuing_judge
        }
        
        case = global_case_manager.create_case(
            case_type=f"JUSTICE_WARRANT_{warrant_type}",
            creator_agency="MINISTERE_JUSTICE_TRIBUNAL",
            security_level="SECRET-DEFENSE",
            metadata=case_data
        )

        # Validation immédiate requise par le greffe
        task = global_human_task_service.create_task(
            task_name=f"Homologation Mandat {warrant_type}",
            case_id=case.case_id,
            role_required="GREFFIER_TRIBUNAL",
            correlation_id=case.correlation_id,
            sla_minutes=30, # 30 minutes de validation pour mandats d'urgence
            security_classification="SECRET-DEFENSE"
        )

        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority="P1")

        event = CloudEvent(
            event_type="judicial.warrant.issued",
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": issuing_judge,
                "agency": "MINISTERE_JUSTICE_TRIBUNAL",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "warrantType": warrant_type,
                    "suspectId": suspect_id
                }
            }
        )
        global_event_bus.publish("justice.events", event)
        case.transition_status("PENDING_REGISTRY_SIGNATURE", issuing_judge, f"Validation par le greffier de service")
        return case, task
