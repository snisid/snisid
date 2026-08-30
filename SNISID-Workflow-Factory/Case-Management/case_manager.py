# -*- coding: utf-8 -*-
"""
MOTEUR DE GESTION DE DOSSIERS NATIONAUX — SNISID v4.0
Gère le cycle de vie des dossiers d'État Civil, d'Identité, Judiciaires et de Police.
"""

import os
import sys
import uuid
from datetime import datetime
import logging

# Bootstrap path resolution for hyphenated directories
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
sys.path.append(parent_dir)
sys.path.append(os.path.join(parent_dir, "Event-Driven"))

from event_bus import global_event_bus, CloudEvent

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [CASE-MGMT] - %(levelname)s - %(message)s')

class NationalCase:
    def __init__(self, case_type, creator_agency, security_level="CONFIDENTIEL", metadata=None):
        self.case_id = f"CASE-{datetime.now().strftime('%Y%m%d')}-{uuid.uuid4().hex[:8].upper()}"
        self.correlation_id = str(uuid.uuid4())
        self.case_type = case_type  # ex: ETAT_CIVIL_NAISSANCE, IDENTITE_ENROLEMENT, JUSTICE_DOSSIER
        self.status = "CREATED"
        self.creator_agency = creator_agency
        self.assigned_agency = creator_agency
        self.security_level = security_level  # PUBLIC, RESTREINT, CONFIDENTIEL, SECRET-DEFENSE
        self.created_at = datetime.now().isoformat()
        self.updated_at = datetime.now().isoformat()
        self.data_payload = metadata if metadata else {}
        self.audit_trail = []
        
        # Enregistrement de la création
        self._add_audit_entry("CASE_CREATED", f"Initialisation du dossier de type {case_type} par l'agence {creator_agency}")

    def _add_audit_entry(self, action, details, operator="SYSTEM"):
        entry = {
            "timestamp": datetime.now().isoformat(),
            "action": action,
            "operator": operator,
            "details": details,
            "hash": uuid.uuid4().hex # Simulation de chaînage d'audit cryptographique
        }
        self.audit_trail.append(entry)
        logging.info(f"[{self.case_id}] AUDIT: {action} - {details}")

    def transition_status(self, new_status, operator, details=""):
        old_status = self.status
        self.status = new_status
        self.updated_at = datetime.now().isoformat()
        self._add_audit_entry("STATUS_TRANSITION", f"Passage de {old_status} à {new_status}. Réf: {details}", operator)
        
        # Publication de l'événement de changement d'état du dossier
        event = CloudEvent(
            event_type="workflow.case.status_changed",
            subject=self.case_id,
            data={
                "correlationId": self.correlation_id,
                "operatorId": operator,
                "agency": self.assigned_agency,
                "payload": {
                    "caseId": self.case_id,
                    "caseType": self.case_type,
                    "oldStatus": old_status,
                    "newStatus": new_status,
                    "details": details
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)

    def update_payload(self, keys_values: dict, operator, details=""):
        for k, v in keys_values.items():
            self.data_payload[k] = v
        self.updated_at = datetime.now().isoformat()
        self._add_audit_entry("DATA_UPDATE", f"Mise à jour des champs: {', '.join(keys_values.keys())}. Réf: {details}", operator)

    def assign_to_agency(self, target_agency, operator, details=""):
        old_agency = self.assigned_agency
        self.assigned_agency = target_agency
        self.updated_at = datetime.now().isoformat()
        self._add_audit_entry("AGENCY_REASSIGNMENT", f"Réaffecté de {old_agency} à {target_agency}", operator)


class NationalCaseManager:
    def __init__(self):
        self.active_cases = {}

    def create_case(self, case_type, creator_agency, security_level="CONFIDENTIEL", metadata=None) -> NationalCase:
        new_case = NationalCase(case_type, creator_agency, security_level, metadata)
        self.active_cases[new_case.case_id] = new_case
        
        # Notification globale
        event = CloudEvent(
            event_type="workflow.case.created",
            subject=new_case.case_id,
            data={
                "correlationId": new_case.correlation_id,
                "operatorId": "SYSTEM",
                "agency": creator_agency,
                "payload": {
                    "caseId": new_case.case_id,
                    "caseType": case_type,
                    "securityLevel": security_level
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)
        return new_case

    def get_case(self, case_id) -> NationalCase:
        return self.active_cases.get(case_id, None)

# Singleton global de gestion de dossiers
global_case_manager = NationalCaseManager()
