# -*- coding: utf-8 -*-
"""
WORKFLOWS ÉTAT CIVIL — SNISID v4.0
Industrialise les processus de Naissance, Mariage, Divorce, Décès et Adoption selon la norme BPMN 2.0.
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

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [CIVIL-REGISTRY] - %(levelname)s - %(message)s')

class CivilRegistryWorkflowManager:
    @staticmethod
    def start_birth_workflow(birth_type, parent_data, metadata=None):
        """
        Démarre un dossier d'État Civil "Naissance" selon la sous-catégorie BPMN spécifiée.
        Types supportés:
          - SIMPLE: Naissance Simple (SLA standard, validation 1er niveau)
          - RECONNAISSANCE: Naissance par Reconnaissance (SLA standard, double validation + consentement)
          - DECLARATION_TARDIVE: Déclaration Tardive (SLA rapide, processus d'investigation/audition)
          - DECRET: Naissance par Décret (SLA strict, validation par ministère)
          - JUGEMENT: Naissance par Jugement au Rend des Minutes (SLA de fond, validation par un tribunal civil)
        """
        logging.info(f"Démarrage du processus de Naissance: {birth_type}")
        
        # 1. Création du dossier national (Case Management)
        case_data = {
            "birth_type": birth_type,
            "child_last_name": parent_data.get("child_last_name"),
            "child_first_name": parent_data.get("child_first_name"),
            "birth_date": parent_data.get("birth_date"),
            "birth_place": parent_data.get("birth_place"),
            "parents": parent_data.get("parents", []),
            "signatures": {}
        }
        if metadata:
            case_data.update(metadata)

        case = global_case_manager.create_case(
            case_type=f"CIVIL_REG_BIRTH_{birth_type}",
            creator_agency="MINISTERE_INTERIEUR_CIVIL",
            security_level="CONFIDENTIEL",
            metadata=case_data
        )

        # 2. Configuration des SLAs spécifiques au type de naissance
        sla_duration_map = {
            "SIMPLE": 240,            # 4 heures
            "RECONNAISSANCE": 480,   # 8 heures
            "DECLARATION_TARDIVE": 1440, # 24 heures (audition requise)
            "DECRET": 2880,          # 48 heures (secrétariat général)
            "JUGEMENT": 5760         # 96 heures (greffe judiciaire)
        }
        sla_minutes = sla_duration_map.get(birth_type, 240)
        
        # 3. Création de la première tâche humaine de validation
        task_name_map = {
            "SIMPLE": "Validation Naissance Simple",
            "RECONNAISSANCE": "Validation Naissance par Reconnaissance",
            "DECLARATION_TARDIVE": "Audition & Enquête Déclaration Tardive",
            "DECRET": "Contrôle Administratif Décret de Naissance",
            "JUGEMENT": "Enregistrement Greffe - Jugement au Rend des Minutes"
        }
        task_name = task_name_map.get(birth_type, "Validation Naissance")
        
        task = global_human_task_service.create_task(
            task_name=task_name,
            case_id=case.case_id,
            role_required="OFFICIER_ETAT_CIVIL",
            correlation_id=case.correlation_id,
            sla_minutes=sla_minutes,
            security_classification="CONFIDENTIEL"
        )
        
        # Enregistrement de l'échéance dans le moteur de SLA
        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority="P2" if birth_type in ["SIMPLE", "RECONNAISSANCE"] else "P3")

        # 4. Publication de l'événement d'initialisation du workflow
        init_event = CloudEvent(
            event_type="birth.created",
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": "SYSTEM",
                "agency": "MINISTERE_INTERIEUR_CIVIL",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "birthType": birth_type,
                    "childName": f"{case_data['child_first_name']} {case_data['child_last_name']}"
                }
            }
        )
        global_event_bus.publish("civil.registry.events", init_event)
        
        case.transition_status("PENDING_HUMAN_VALIDATION", "SYSTEM", f"Création de la tâche {task.task_id}")
        return case, task

    @staticmethod
    def start_other_civil_workflow(workflow_domain, case_data):
        """
        Démarre un dossier d'État Civil autre (MARIAGE, DIVORCE, DECES, ADOPTION).
        """
        logging.info(f"Démarrage du processus de {workflow_domain}")
        
        case = global_case_manager.create_case(
            case_type=f"CIVIL_REG_{workflow_domain}",
            creator_agency="MINISTERE_INTERIEUR_CIVIL",
            security_level="CONFIDENTIEL",
            metadata=case_data
        )

        task = global_human_task_service.create_task(
            task_name=f"Validation Homologation {workflow_domain.capitalize()}",
            case_id=case.case_id,
            role_required="OFFICIER_ETAT_CIVIL",
            correlation_id=case.correlation_id,
            sla_minutes=480, # 8 Heures
            security_classification="CONFIDENTIEL"
        )
        
        if sla_engine.global_sla_engine:
            sla_engine.global_sla_engine.register_timer(task.task_id, "TASK", task.deadline, priority="P3")

        init_event = CloudEvent(
            event_type=f"{workflow_domain.lower()}.created",
            subject=case.case_id,
            data={
                "correlationId": case.correlation_id,
                "operatorId": "SYSTEM",
                "agency": "MINISTERE_INTERIEUR_CIVIL",
                "payload": {
                    "caseId": case.case_id,
                    "taskId": task.task_id,
                    "domain": workflow_domain
                }
            }
        )
        global_event_bus.publish("civil.registry.events", init_event)
        case.transition_status("PENDING_HUMAN_VALIDATION", "SYSTEM", f"Création de la tâche {task.task_id}")
        return case, task
