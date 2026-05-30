# -*- coding: utf-8 -*-
"""
SERVICE DE GESTION DES TÂCHES HUMAINES (HUMAN TASK MANAGER) — SNISID v4.0
Gère le cycle de vie des validations humaines, des chaînes d'approbation et des délégations.
"""

import os
import sys
import uuid
from datetime import datetime, timedelta
import logging

# Bootstrap path resolution
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
sys.path.append(parent_dir)
sys.path.append(os.path.join(parent_dir, "Event-Driven"))

from event_bus import global_event_bus, CloudEvent

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [HUMAN-TASK] - %(levelname)s - %(message)s')

class HumanTask:
    def __init__(self, task_name, case_id, role_required, correlation_id, sla_minutes=240, security_classification="CONFIDENTIEL"):
        self.task_id = f"TASK-{uuid.uuid4().hex[:8].upper()}"
        self.task_name = task_name
        self.case_id = case_id
        self.correlation_id = correlation_id
        self.role_required = role_required # ex: OFFICIER_ETAT_CIVIL, SUPERVISEUR_DCPJ
        self.security_classification = security_classification
        self.status = "PENDING"  # PENDING, ASSIGNED, COMPLETED, ESCALATED, DELEGATED, CANCELLED
        
        self.assigned_user = None
        self.original_user = None # Pour garder trace en cas de délégation
        self.created_at = datetime.now().isoformat()
        self.assigned_at = None
        self.completed_at = None
        
        # Gestion des délais
        self.sla_duration = timedelta(minutes=sla_minutes)
        self.deadline = (datetime.now() + self.sla_duration).isoformat()
        
        self.approval_chain = [] # Historique des décisions
        self.audit_log = []
        
        self._log_audit("TASK_INITIATED", f"Tâche créée. Rôle requis: {role_required}. Limite de résolution: {self.deadline}")

    def _log_audit(self, action, details, operator="SYSTEM"):
        entry = {
            "timestamp": datetime.now().isoformat(),
            "action": action,
            "operator": operator,
            "details": details
        }
        self.audit_log.append(entry)
        logging.info(f"[{self.task_id}] HUMAN AUDIT: {action} - {details} (Op: {operator})")

    def assign_to(self, user_id, operator="SYSTEM"):
        self.assigned_user = user_id
        if not self.original_user:
            self.original_user = user_id
        self.status = "ASSIGNED"
        self.assigned_at = datetime.now().isoformat()
        self._log_audit("TASK_ASSIGNED", f"Assignée à l'agent {user_id}", operator)
        
        # Notification d'assignation
        event = CloudEvent(
            event_type="workflow.task.assigned",
            subject=self.task_id,
            data={
                "correlationId": self.correlation_id,
                "operatorId": operator,
                "agency": "SNISID-CENTRAL",
                "payload": {
                    "taskId": self.task_id,
                    "caseId": self.case_id,
                    "assignedUser": user_id,
                    "deadline": self.deadline
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)

    def delegate_to(self, delegate_user_id, reason, operator):
        if self.status not in ["PENDING", "ASSIGNED"]:
            raise ValueError(f"Impossible de déléguer une tâche dans l'état {self.status}")
        
        previous_user = self.assigned_user
        self.assigned_user = delegate_user_id
        self.status = "DELEGATED"
        self._log_audit("TASK_DELEGATED", f"Délégation de {previous_user} à {delegate_user_id}. Motif: {reason}", operator)
        
        event = CloudEvent(
            event_type="workflow.task.delegated",
            subject=self.task_id,
            data={
                "correlationId": self.correlation_id,
                "operatorId": operator,
                "agency": "SNISID-CENTRAL",
                "payload": {
                    "taskId": self.task_id,
                    "fromUser": previous_user,
                    "toUser": delegate_user_id,
                    "reason": reason
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)

    def complete(self, decision, comments, operator):
        """
        decision: APPROVED ou REJECTED
        """
        if self.status not in ["PENDING", "ASSIGNED", "DELEGATED"]:
            raise ValueError(f"Impossible de compléter une tâche dans l'état {self.status}")
            
        self.status = "COMPLETED"
        self.completed_at = datetime.now().isoformat()
        self.approval_chain.append({
            "step": self.task_name,
            "operator": operator,
            "decision": decision,
            "comments": comments,
            "timestamp": self.completed_at
        })
        self._log_audit("TASK_COMPLETED", f"Décision: {decision}. Commentaires: {comments}", operator)
        
        # Événement de complétion (workflow.approved ou workflow.rejected)
        event_type = "workflow.approved" if decision == "APPROVED" else "workflow.rejected"
        event = CloudEvent(
            event_type=event_type,
            subject=self.case_id,
            data={
                "correlationId": self.correlation_id,
                "operatorId": operator,
                "agency": "SNISID-CENTRAL",
                "payload": {
                    "taskId": self.task_id,
                    "caseId": self.case_id,
                    "decision": decision,
                    "comments": comments
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)

    def escalate(self, escalation_reason):
        self.status = "ESCALATED"
        self._log_audit("TASK_ESCALATED", f"Déclenchement d'escalade automatique. Raison: {escalation_reason}")
        
        event = CloudEvent(
            event_type="workflow.task.escalated",
            subject=self.task_id,
            data={
                "correlationId": self.correlation_id,
                "operatorId": "SLA_ENGINE",
                "agency": "SNISID-CENTRAL",
                "payload": {
                    "taskId": self.task_id,
                    "caseId": self.case_id,
                    "reason": escalation_reason,
                    "roleRequired": self.role_required
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)


class HumanTaskService:
    def __init__(self):
        self.tasks = {}

    def create_task(self, task_name, case_id, role_required, correlation_id, sla_minutes=240, security_classification="CONFIDENTIEL") -> HumanTask:
        task = HumanTask(task_name, case_id, role_required, correlation_id, sla_minutes, security_classification)
        self.tasks[task.task_id] = task
        return task

    def get_task(self, task_id) -> HumanTask:
        return self.tasks.get(task_id, None)

    def list_pending_tasks(self):
        return [t for t in self.tasks.values() if t.status in ["PENDING", "ASSIGNED", "DELEGATED"]]

# Singleton de gestion de tâches humaines
global_human_task_service = HumanTaskService()
