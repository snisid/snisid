# -*- coding: utf-8 -*-
"""
WORKFLOWS POLICE NATIONALE — SNISID v4.0
Industrialise les procédures d'arrestation, d'audition de garde à vue et d'enquêtes criminelles DCPJ.
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

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [POLICE-WORKFLOW] - %(levelname)s - %(message)s')

class PoliceWorkflowManager:
    @staticmethod
    def log_arrestation(officer_badge, suspect_details, arrest_location, charges_summary):
        """
        Démarre un processus d'arrestation et de garde à vue (GAV).
        SLA de GAV légal strict: 24h (1440 minutes). Dépasser cette limite sans validation du juge est illégal!
        """
        logging.info(f"Arrestation enregistrée par l'officier {officer_badge} à {arrest_location}")
        
        case_data = {
            "officer_badge": officer_badge,
            "suspect": suspect_details,
            "arrest_location": arrest_location,
            "charges_summary": charges_summary,
            "detention_start_time": None
        }
        
        case = global_case_manager.create_case(
            case_type="POLICE_ARREST_GAV",
            creator_agency="DCPJ_POLICE_NATIONALE",
            security_level="SECRET-DEFENSE",
            metadata=case_data
        )

        # SLA légal strict de 24 Heures pour la Garde à Vue
        # Une tâche de validation de libération ou de présentation au juge doit être résolue
        task = global_human_task_service.create_task(
            task_name="Contrôle Constitutionnel & Légalité Garde à Vue",
            case_id=case.case_id,
            role_required="COMMISSAIRE_DCPJ",
            correlation_id=case.correlation_id,
            sla_minutes=1440, # 24 heures strictes
            security_classification="SECRET-DEFENSE"
        )

        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority="P1") # Priorité critique

        event = CloudEvent(
            event_type="police.arrestation.created",
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": officer_badge,
                "agency": "DCPJ_POLICE_NATIONALE",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "suspectName": f"{suspect_details.get('first_name')} {suspect_details.get('last_name')}",
                    "charges": charges_summary
                }
            }
        )
        global_event_bus.publish("police.events", event)
        case.transition_status("IN_DETENTION_GAV", officer_badge, f"Début de la garde à vue à {arrest_location}")
        return case, task

    @staticmethod
    def start_dcpj_investigation(case_title, target_individual, chief_investigator):
        """
        Démarre une enquête criminelle de la Direction Centrale de la Police Judiciaire (DCPJ).
        """
        logging.info(f"Ouverture d'enquête DCPJ: '{case_title}' menée par {chief_investigator}")
        
        case_data = {
            "case_title": case_title,
            "target": target_individual,
            "chief_investigator": chief_investigator,
            "evidence_list": [],
            "witnesses": []
        }
        
        case = global_case_manager.create_case(
            case_type="POLICE_DCPJ_INVESTIGATION",
            creator_agency="DCPJ_POLICE_NATIONALE",
            security_level="SECRET-DEFENSE",
            metadata=case_data
        )

        task = global_human_task_service.create_task(
            task_name="Rapport d'Enquête Préliminaire DCPJ",
            case_id=case.case_id,
            role_required="OFFICIER_POLICE_JUDICIAIRE",
            correlation_id=case.correlation_id,
            sla_minutes=2880, # 48 Heures pour la première validation de piste
            security_classification="SECRET-DEFENSE"
        )

        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority="P2")

        event = CloudEvent(
            event_type="police.investigation.opened",
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": chief_investigator,
                "agency": "DCPJ_POLICE_NATIONALE",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "caseTitle": case_title
                }
            }
        )
        global_event_bus.publish("police.events", event)
        case.transition_status("UNDER_INVESTIGATION", chief_investigator, f"Initialisation de l'enquête DCPJ")
        return case, task
