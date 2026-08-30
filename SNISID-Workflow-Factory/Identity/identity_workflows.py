# -*- coding: utf-8 -*-
"""
WORKFLOWS IDENTITÉ NATIONALE — SNISID v4.0
Industrialise l'enrôlement biométrique, la vérification, la correction, la révocation et la récupération d'identités.
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

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [IDENTITY-WORKFLOW] - %(levelname)s - %(message)s')

class IdentityWorkflowManager:
    @staticmethod
    def start_identity_workflow(action_type, identity_data, operator_id="AGENT_ENROLMENT"):
        """
        Démarre un dossier de gestion de l'Identité Nationale.
        Types:
          - ENROLLMENT: Création, enregistrement biométrique.
          - VERIFICATION: Processus de vérification approfondie ou audit de doublon.
          - CORRECTION: Demande de correction administrative (changement de nom/statut).
          - REVOCATION: Révocation ou suspension de l'identité ou de la carte d'identité (perte/vol/décès).
          - RECOVERY: Récupération de l'identité, réémission de carte.
        """
        logging.info(f"Démarrage du processus Identité: {action_type}")
        
        case = global_case_manager.create_case(
            case_type=f"IDENTITY_{action_type}",
            creator_agency="OFFICE_NATIONAL_IDENTITE_ONI",
            security_level="CONFIDENTIEL",
            metadata=identity_data
        )

        # Les workflows d'identité ont des SLAs différents
        sla_map = {
            "ENROLLMENT": 480,   # 8 heures (ONI validation)
            "VERIFICATION": 120, # 2 heures (critique pour les frontières/banques)
            "CORRECTION": 1440,  # 24 heures (correction administrative)
            "REVOCATION": 30,    # 30 minutes (URGENT - P1)
            "RECOVERY": 240      # 4 heures
        }
        sla_minutes = sla_map.get(action_type, 240)
        priority = "P1" if action_type == "REVOCATION" else ("P2" if action_type == "VERIFICATION" else "P3")

        # Création de la tâche de validation requise
        task_name_map = {
            "ENROLLMENT": "Validation Biométrique & Enrôlement ONI",
            "VERIFICATION": "Vérification des Doublons & Correspondance",
            "CORRECTION": "Approbation de la Modification Civile",
            "REVOCATION": "Révocation Immédiate Titre Identité",
            "RECOVERY": "Authentification & Récupération Profil"
        }
        task_name = task_name_map.get(action_type, "Validation Identité")

        task = global_human_task_service.create_task(
            task_name=task_name,
            case_id=case.case_id,
            role_required="SUPERVISEUR_IDENTITE_ONI",
            correlation_id=case.correlation_id,
            sla_minutes=sla_minutes,
            security_classification="CONFIDENTIEL"
        )

        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority=priority)

        # Publication de l'événement d'identité sur le bus
        event_type_map = {
            "ENROLLMENT": "identity.enrollment.started",
            "VERIFICATION": "identity.verified", # Déclenche de manière autonome les vérifications
            "CORRECTION": "identity.correction.requested",
            "REVOCATION": "identity.revocation.initiated",
            "RECOVERY": "identity.recovery.requested"
        }
        
        init_event = CloudEvent(
            event_type=event_type_map.get(action_type, "identity.workflow.started"),
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": operator_id,
                "agency": "OFFICE_NATIONAL_IDENTITE_ONI",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "nationalId": identity_data.get("national_id", "Nouveau"),
                    "action": action_type
                }
            }
        )
        global_event_bus.publish("identity.events", init_event)
        case.transition_status("PENDING_BIOMETRIC_VALIDATION", operator_id, f"Création de la tâche {task.task_id}")
        return case, task
