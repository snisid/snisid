# -*- coding: utf-8 -*-
"""
MOTEUR DE DELAIS ADMINISTRATIFS ET SLA — SNISID v4.0
Surveille les limites de temps des dossiers/tâches et déclenche des actions correctives automatiques.
"""

import os
import sys
from datetime import datetime
import logging

# Bootstrap path resolution
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
sys.path.append(parent_dir)
sys.path.append(os.path.join(parent_dir, "Event-Driven"))

from event_bus import global_event_bus, CloudEvent

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [SLA-ENGINE] - %(levelname)s - %(message)s')

class SLATimer:
    def __init__(self, timer_id, target_id, target_type, limit_time_iso, priority="P3"):
        self.timer_id = timer_id
        self.target_id = target_id      # ID de la tâche ou du dossier
        self.target_type = target_type  # TASK ou CASE
        self.limit_time = datetime.fromisoformat(limit_time_iso)
        self.priority = priority        # P1, P2, P3, P4
        self.is_active = True
        self.triggered_at = None
        self.breached = False

    def check_breach(self) -> bool:
        if not self.is_active:
            return False
        
        if datetime.now() > self.limit_time:
            self.breached = True
            self.is_active = False
            self.triggered_at = datetime.now().isoformat()
            logging.warning(f"[BREACH] SLA Dépassé pour {self.target_type} {self.target_id}! Priorité: {self.priority}")
            return True
        return False


class NationalSLAEngine:
    def __init__(self, human_task_service, case_manager):
        self.active_timers = {}
        self.task_service = human_task_service
        self.case_manager = case_manager

    def register_timer(self, target_id, target_type, limit_time_iso, priority="P3") -> SLATimer:
        timer_id = f"SLA-{target_id}"
        timer = SLATimer(timer_id, target_id, target_type, limit_time_iso, priority)
        self.active_timers[timer_id] = timer
        logging.info(f"Nouveau Timer SLA enregistré: {timer_id} (Type: {target_type}, Limite: {limit_time_iso})")
        return timer

    def deactivate_timer(self, target_id):
        timer_id = f"SLA-{target_id}"
        if timer_id in self.active_timers:
            self.active_timers[timer_id].is_active = False
            logging.info(f"Timer SLA désactivé pour {target_id}")

    def evaluate_all_slas(self):
        """
        Vérifie toutes les limites de temps et déclenche les alertes et les actions d'escalade.
        """
        breaches_detected = []
        for timer_id, timer in list(self.active_timers.items()):
            if timer.check_breach():
                breaches_detected.append(timer)
                self._handle_sla_breach(timer)
        return breaches_detected

    def _handle_sla_breach(self, timer: SLATimer):
        # Publier l'événement de violation de SLA
        event = CloudEvent(
            event_type="workflow.sla.breach",
            subject=timer.target_id,
            data={
                "correlationId": "SLA-SYSTEM",
                "operatorId": "SLA_ENGINE",
                "agency": "SNISID-CENTRAL",
                "payload": {
                    "timerId": timer.timer_id,
                    "targetId": timer.target_id,
                    "targetType": timer.target_type,
                    "priority": timer.priority,
                    "limitTime": timer.limit_time.isoformat(),
                    "breachTime": timer.triggered_at
                }
            }
        )
        global_event_bus.publish("workflow.audit", event)

        # Actions d'escalade automatiques selon la cible
        if timer.target_type == "TASK":
            task = self.task_service.get_task(timer.target_id)
            if task and task.status in ["PENDING", "ASSIGNED", "DELEGATED"]:
                escalation_reason = f"Violation de SLA ({timer.priority}) - Temps limite dépassé ({timer.limit_time})"
                task.escalate(escalation_reason)
                
                # Déclenchement d'un workflow d'urgence si Priorité P1 (Urgence Nationale)
                if timer.priority == "P1":
                    self._trigger_emergency_workflow(task)

    def _trigger_emergency_workflow(self, task):
        logging.critical(f"[CRITICAL_SLA] Déclenchement du flux d'urgence nationale pour la tâche {task.task_id}")
        emergency_event = CloudEvent(
            event_type="workflow.emergency.triggered",
            subject=task.case_id,
            data={
                "correlationId": task.correlation_id,
                "operatorId": "SLA_ENGINE_CRITICAL",
                "agency": "SNISID-CENTRAL",
                "payload": {
                    "taskId": task.task_id,
                    "reason": "SLA Breach on P1 task",
                    "escalationTarget": "DIRECTEUR_NATIONAL_AGENCE"
                }
            }
        )
        global_event_bus.publish("system.alarms", emergency_event)

# Création d'une instance globale pour l'initialisation ultérieure
global_sla_engine = None
def init_sla_engine(task_service, case_manager):
    global global_sla_engine
    global_sla_engine = NationalSLAEngine(task_service, case_manager)
    return global_sla_engine
