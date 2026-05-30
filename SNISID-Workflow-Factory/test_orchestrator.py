# -*- coding: utf-8 -*-
"""
SCRIPT DE TEST D'INTÉGRATION ET DE VÉRIFICATION — SNISID v4.0
Vérifie la syntaxe de tous les modules de la Workflow Factory et exécute un scénario complet d'orchestration.
"""

import os
import sys
import logging

logging.basicConfig(level=logging.INFO, format='%(asctime)s - [TEST-INTEGRATION] - %(levelname)s - %(message)s')

def run_integration_test():
    try:
        logging.info("======================================================================")
        logging.info("DÉBUT DU TEST D'INTÉGRATION ET DE ROBUSTESSE DU MOTEUR DE WORKFLOW SNISID")
        logging.info("======================================================================")

        # Append paths
        base_dir = os.path.dirname(os.path.abspath(__file__))
        sys.path.append(base_dir)
        sys.path.append(os.path.join(base_dir, "Event-Driven"))
        sys.path.append(os.path.join(base_dir, "Case-Management"))
        sys.path.append(os.path.join(base_dir, "Human-Tasks"))
        sys.path.append(os.path.join(base_dir, "SLA"))
        sys.path.append(os.path.join(base_dir, "Civil-Registry"))
        sys.path.append(os.path.join(base_dir, "Identity"))
        sys.path.append(os.path.join(base_dir, "Justice"))
        sys.path.append(os.path.join(base_dir, "Police"))
        sys.path.append(os.path.join(base_dir, "Offline"))

        # Import modules using resolved paths
        from event_bus import global_event_bus, CloudEvent
        from case_manager import global_case_manager
        from human_task_service import global_human_task_service
        from civil_registry_workflows import CivilRegistryWorkflowManager
        from identity_workflows import IdentityWorkflowManager
        from justice_workflows import JusticeWorkflowManager
        from police_workflows import PoliceWorkflowManager
        from offline_sync import global_offline_sync_engine
        from orchestrator import global_workflow_factory
        import sla_engine

        logging.info("✓ Tous les modules de la Workflow-Factory ont été importés avec succès (Aucune erreur de syntaxe).")

        # -------------------------------------------------------------
        # SCÉNARIO 1 : ÉTAT CIVIL — Naissance Simple
        # -------------------------------------------------------------
        logging.info("\n--- SCÉNARIO 1 : Démarrage d'un processus de Naissance Simple ---")
        parent_data = {
            "child_last_name": "VALÉRY",
            "child_first_name": "Magalie",
            "birth_date": "2026-05-25",
            "birth_place": "Maternité de Delmas",
            "parents": ["VALÉRY Joseph", "JEAN Marie"]
        }
        
        # Lancer le workflow
        case, task = CivilRegistryWorkflowManager.start_birth_workflow("SIMPLE", parent_data)
        logging.info(f"✓ Dossier d'État Civil créé: {case.case_id} (Statut: {case.status})")
        logging.info(f"✓ Tâche humaine générée: {task.task_id} (Rôle requis: {task.role_required}, Deadline: {task.deadline})")

        # Validation de la tâche par l'officier
        logging.info("\n--- Approbation de la tâche par l'officier d'État Civil ---")
        task.complete(decision="APPROVED", comments="Dossier de naissance entièrement conforme, pièces médicales validées.", operator="OFFICIER_MARC_01")
        
        # Vérification de l'état final du dossier
        logging.info(f"✓ Statut final du dossier: {case.status}")
        assert case.status == "COMPLETED_APPROVED", f"Erreur: Statut attendu 'COMPLETED_APPROVED', reçu '{case.status}'"
        logging.info("✓ Scénario 1 validé avec succès (Dossier complété et persistant).")

        # -------------------------------------------------------------
        # SCÉNARIO 2 : HORS-LIGNE — Enregistrement terrain en zone blanche
        # -------------------------------------------------------------
        logging.info("\n--- SCÉNARIO 2 : Enregistrement Hors-Ligne & Rejoueur d'Événements ---")
        terminal_id = "TERM-MOBILE-OPJ-88"
        offline_case_id = "CASE-OFFLINE-999"
        
        # Simulation d'actions hors-ligne accumulées dans le terminal mobile de l'agent
        logging.info("1. Enregistrement d'une création de dossier locale (Offline)")
        global_offline_sync_engine.queue_offline_event(
            terminal_id=terminal_id,
            action_type="LOCAL_CASE_CREATE",
            case_id=offline_case_id,
            payload={
                "case_type": "POLICE_DCPJ_INVESTIGATION",
                "creator_agency": "DCPJ_POLICE_NATIONALE",
                "security_level": "SECRET-DEFENSE",
                "metadata": {"case_title": "Enquête Fraude Nationale", "chief_investigator": "COMMISSAIRE_PIERRE"}
            },
            operator="AGENT_TERRAIN_01"
        )
        
        logging.info("2. Enregistrement d'une mise à jour de données locale (Offline)")
        global_offline_sync_engine.queue_offline_event(
            terminal_id=terminal_id,
            action_type="LOCAL_DATA_UPDATE",
            case_id=offline_case_id,
            payload={
                "data": {"witnesses": ["Témoin A", "Témoin B"], "status_field": "IN_PROGRESS"}
            },
            operator="AGENT_TERRAIN_01"
        )

        # Simulation du retour réseau et de la synchronisation différée (Deferred Sync with Replay)
        logging.info("3. Rétablissement réseau. Lancement de la synchronisation et de la résolution de conflits...")
        sync_result = global_offline_sync_engine.synchronize_terminal(terminal_id)
        
        logging.info(f"✓ Rapport de Synchronisation: {sync_result}")
        assert sync_result["status"] == "SUCCESS", "Erreur: Échec de la synchronisation hors-ligne"
        
        synced_case = global_case_manager.get_case(offline_case_id)
        logging.info(f"✓ Dossier récréé rétroactivement sur le serveur central: {synced_case.case_id}")
        logging.info(f"✓ Contenu synchronisé des données: {synced_case.data_payload}")
        logging.info("✓ Scénario 2 (Offline Support) validé avec succès.")

        # -------------------------------------------------------------
        # SCÉNARIO 3 : DÉPASSEMENT DE SLA & ALERTE CRITIQUE
        # -------------------------------------------------------------
        logging.info("\n--- SCÉNARIO 3 : Simulation d'un dépassement de SLA et Alarme Système ---")
        suspect = {"first_name": "Jean", "last_name": "VALJEAN"}
        infraction = {"title": "Vol qualifié de miche de pain"}
        
        # Ouverture d'un dossier criminel
        court_case, court_task = JusticeWorkflowManager.open_criminal_case(suspect, infraction, "COURT_CENTRAL", "JUGE_JACOB")
        
        # Forcer l'expiration temporelle et déclencher l'évaluation SLA
        logging.info("Forçage artificiel du dépassement de la date limite...")
        # On remplace l'heure limite du timer associé pour la mettre dans le passé
        timer_id = f"SLA-{court_task.task_id}"
        if timer_id in sla_engine.global_sla_engine.active_timers:
            import datetime
            sla_engine.global_sla_engine.active_timers[timer_id].limit_time = datetime.datetime.now() - datetime.timedelta(minutes=1)
            
        logging.info("Évaluation des SLA par le démon SLA Engine...")
        breaches = sla_engine.global_sla_engine.evaluate_all_slas()
        
        logging.info(f"✓ Violations de SLA détectées: {len(breaches)}")
        logging.info(f"✓ Nouveau statut de la tâche suite à l'escalade: {court_task.status}")
        assert court_task.status == "ESCALATED", f"Erreur: La tâche aurait dû être escaladée, statut actuel: {court_task.status}"
        logging.info("✓ Scénario 3 (SLA Engine & Escalation) validé avec succès.")

        logging.info("\n======================================================================")
        logging.info("TOUS LES TESTS D'INTÉGRATION ONT ÉTÉ PASSÉS AVEC SUCCÈS ! FORENSIC-READY ✓")
        logging.info("======================================================================")
        return True

    except Exception as e:
        logging.error(f"ECHEC CRITIQUE DU TEST D'INTEGRATION: {str(e)}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    success = run_integration_test()
    sys.exit(0 if success else 1)
