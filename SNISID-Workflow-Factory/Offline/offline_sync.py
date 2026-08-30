# -*- coding: utf-8 -*-
"""
MOTEUR DE SYNCHRONISATION HORS-LIGNE (OFFLINE BPMN ENGINE) — SNISID v4.0
Gère la capture locale d'événements, le rejeu (replay) différé et la résolution de conflits transactionnels.
"""

import os
import sys
import json
from datetime import datetime
import logging

# Bootstrap path resolution
current_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(current_dir)
sys.path.append(parent_dir)
sys.path.append(os.path.join(parent_dir, "Event-Driven"))
sys.path.append(os.path.join(parent_dir, "Case-Management"))

from event_bus import global_event_bus, CloudEvent
from case_manager import global_case_manager

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [OFFLINE-SYNC] - %(levelname)s - %(message)s')

class OfflineEvent:
    def __init__(self, action_type, case_id, payload, operator, timestamp=None):
        self.offline_event_id = f"OFF-{datetime.now().strftime('%M%S')}-{operator}"
        self.action_type = action_type # ex: LOCAL_CASE_CREATE, LOCAL_STATUS_TRANSITION, LOCAL_DATA_UPDATE
        self.case_id = case_id
        self.payload = payload
        self.operator = operator
        self.timestamp = timestamp if timestamp else datetime.now().isoformat()


class OfflineSyncEngine:
    def __init__(self):
        self.local_queues = {} # Enregistrement par identifiant de terminal/agent

    def queue_offline_event(self, terminal_id, action_type, case_id, payload, operator):
        """
        Enregistre localement (sur le terminal mobile) un événement pendant que la connexion est coupée.
        """
        if terminal_id not in self.local_queues:
            self.local_queues[terminal_id] = []
            
        event = OfflineEvent(action_type, case_id, payload, operator)
        self.local_queues[terminal_id].append(event)
        logging.info(f"[LOCAL_QUEUE] Terminal [{terminal_id}] a mis en file d'attente l'action '{action_type}' pour le dossier '{case_id}'")
        return event

    def synchronize_terminal(self, terminal_id):
        """
        Simule le retour en ligne et le rejeu des événements locaux (Event Replay) avec résolution de conflits.
        """
        if terminal_id not in self.local_queues or not self.local_queues[terminal_id]:
            logging.info(f"Aucun événement hors-ligne à synchroniser pour le terminal: {terminal_id}")
            return {"status": "SUCCESS", "synced_events_count": 0, "conflicts": []}

        offline_events = self.local_queues[terminal_id]
        logging.info(f"Début de la synchronisation de {len(offline_events)} événements pour le terminal: {terminal_id}")
        
        synced_count = 0
        resolved_conflicts = []

        # Triage chronologique des événements hors-ligne pour un rejeu rigoureux
        offline_events.sort(key=lambda x: x.timestamp)

        for event in offline_events:
            case_id = event.case_id
            server_case = global_case_manager.get_case(case_id)

            if not server_case:
                # Conflit type : Le dossier n'existe pas sur le serveur.
                # Résolution : Si l'action est LOCAL_CASE_CREATE, on le crée, sinon erreur.
                if event.action_type == "LOCAL_CASE_CREATE":
                    logging.info(f"[RESOLVED] Création rétroactive sur le serveur pour le dossier {case_id}")
                    server_case = global_case_manager.create_case(
                        case_type=event.payload.get("case_type"),
                        creator_agency=event.payload.get("creator_agency"),
                        security_level=event.payload.get("security_level", "CONFIDENTIEL"),
                        metadata=event.payload.get("metadata")
                    )
                    
                    # Réindexer dans le gestionnaire de cas sous le bon ID
                    old_generated_id = server_case.case_id
                    server_case.case_id = case_id
                    
                    global_case_manager.active_cases[case_id] = server_case
                    if old_generated_id in global_case_manager.active_cases:
                        del global_case_manager.active_cases[old_generated_id]
                        
                    synced_count += 1
                else:
                    logging.error(f"[CONFLICT] Impossible de rejouer l'action {event.action_type} car le dossier {case_id} n'existe pas.")
                    resolved_conflicts.append({
                        "eventId": event.offline_event_id,
                        "conflictType": "ORPHAN_CASE_ACTION",
                        "resolution": "DISCARDED"
                    })
                continue

            # Cas où le dossier existe déjà sur le serveur : Analyse de conflit
            if event.action_type == "LOCAL_DATA_UPDATE":
                # Stratégie de résolution: Last-Write-Wins (LWW) basée sur le timestamp
                server_last_update = server_case.updated_at
                if event.timestamp > server_last_update:
                    logging.info(f"[RESOLVED - LWW] Mise à jour des données appliquée. Les données hors-ligne sont plus récentes.")
                    server_case.update_payload(
                        keys_values=event.payload.get("data"),
                        operator=event.operator,
                        details=f"Offline Replay LWW - Sync Terminal {terminal_id}"
                    )
                    synced_count += 1
                else:
                    logging.warning(f"[RESOLVED - DISCARD] Conflit LWW détecté. Les données du serveur ({server_last_update}) sont plus récentes que les données hors-ligne ({event.timestamp}). Événement rejeté.")
                    resolved_conflicts.append({
                        "eventId": event.offline_event_id,
                        "conflictType": "LWW_OUTDATED",
                        "resolution": "DISCARDED_IN_FAVOR_OF_SERVER"
                    })

            elif event.action_type == "LOCAL_STATUS_TRANSITION":
                # Conflit d'état: s'assurer que la transition est compatible
                current_server_status = server_case.status
                target_status = event.payload.get("status")
                
                if current_server_status == target_status:
                    logging.info(f"[RESOLVED - NO_OP] État déjà identique sur le serveur.")
                    synced_count += 1
                else:
                    # On applique l'état si l'événement hors-ligne est postérieur à la dernière modification serveur
                    if event.timestamp > server_case.updated_at:
                        logging.info(f"[RESOLVED - OVERWRITE] État forcé selon transition hors-ligne: {target_status}")
                        server_case.transition_status(
                            new_status=target_status,
                            operator=event.operator,
                            details=f"Offline Replay - Sync Terminal {terminal_id}"
                        )
                        synced_count += 1
                    else:
                        logging.warning(f"[RESOLVED - REJECTED] Conflit d'état. L'état actuel du serveur ({current_server_status}) prévaut.")
                        resolved_conflicts.append({
                            "eventId": event.offline_event_id,
                            "conflictType": "STATUS_OUTDATED",
                            "resolution": "SERVER_STATUS_PRESERVED"
                        })

        # Vider la file locale après exécution
        self.local_queues[terminal_id] = []
        logging.info(f"Synchronisation terminée pour [{terminal_id}]. {synced_count} réussis, {len(resolved_conflicts)} conflits gérés.")
        return {
            "status": "SUCCESS",
            "synced_events_count": synced_count,
            "conflicts": resolved_conflicts
        }

# Instance globale
global_offline_sync_engine = OfflineSyncEngine()
