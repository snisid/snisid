# SNISID OMERTA CENTRAL DASHBOARD
# Interface de consultation des alertes de corruption
# Accès réservé au Directeur Central uniquement

import curses
import json
import datetime
import hashlib
from cryptography.hazmat.primitives.ciphers.aead import AESGCM

class OmertaDashboard:
    def __init__(self):
        self.central_vault = "/var/snisid/central/vault"
        self.alert_queue = []
        self.master_code_hash = None
        self.session_active = False
        
    def authenticate_director(self):
        """Authentification par code maître + biométrie"""
        print("╔════════════════════════════════════════════════╗")
        print("║  ACCÈS BUREAU CENTRAL - OMERTA DASHBOARD       ║")
        print("║  TOP SECRET / NOFORN                           ║")
        print("╚════════════════════════════════════════════════╝")
        
        # Code maître (hashé)
        master_code = input("Code Maître Directeur: ")
        master_hash = hashlib.sha512(master_code.encode()).hexdigest()
        
        # Vérification du hash
        with open("/etc/snisid/secrets/director_master_hash", "r") as f:
            stored_hash = f.read().strip()
            
        if master_hash != stored_hash:
            print("❌ CODE INCORRECT - TENTATIVE LOGUÉE")
            return False
            
        # Biométrie supplémentaire (empreinte)
        bio_result = subprocess.run(
            ["snisid-cli", "biometric", "verify", "--level=L3"],
            capture_output=True
        )
        
        if bio_result.returncode != 0:
            print("❌ BIOMÉTRIE ÉCHOUÉE - ACCÈS REFUSÉ")
            return False
            
        self.session_active = True
        print("✅ AUTHENTIFICATION RÉUSSIE - BIENVENUE DIRECTEUR")
        return True
    
    def fetch_alerts(self):
        """Récupérer les alertes non acquittées du vault central"""
        vault_path = f"{self.central_vault}/pending_alerts"
        
        for alert_file in os.listdir(vault_path):
            if alert_file.endswith(".enc"):
                # Déchiffrer avec clé centrale
                with open(f"{vault_path}/{alert_file}", "rb") as f:
                    encrypted_data = f.read()
                    
                aesgcm = AESGCM(load_central_key())
                decrypted = aesgcm.decrypt(nonce, encrypted_data, None)
                
                alert_data = json.loads(decrypted.decode('utf-8'))
                self.alert_queue.append(alert_data)
                
        # Trier par criticité et date
        self.alert_queue.sort(
            key=lambda x: (x['priority'], x['timestamp']), 
            reverse=True
        )
        
    def display_alert(self, alert):
        """Afficher une alerte détaillée"""
        print("\n" + "="*70)
        print(f"🚨 ALERTE JUDAS #{alert['alert_id']}")
        print("="*70)
        print(f"⏰ Timestamp: {alert['timestamp']}")
        print(f"📍 Région: {alert['culprit']['assigned_region']}")
        print(f"\n👤 IDENTITÉ DU CONTREVENANT:")
        print(f"   Nom: {alert['culprit']['nom']}")
        print(f"   Prénom: {alert['culprit']['prenom']}")
        print(f"   NIF: {alert['culprit']['nif']}")
        print(f"   Rôle: {alert['culprit']['role']}")
        print(f"   Photo: [AFFICHÉE DANS FENÊTRE SÉCURISÉE]")
        
        print(f"\n📋 VIOLATIONS DÉTECTÉES:")
        if alert['violations']['privilege_escalation']:
            print(f"   ⚠️  Élévation de privilèges: DECRYPTED_DATA")
        if alert['violations']['jurisdiction_violation']:
            print(f"   ⚠️  Accès hors juridiction: DECRYPTED_DATA")
        if alert['violations']['bulk_export_attempt']:
            print(f"   ⚠️  Export massif de données: {alert['violations']['bulk_export_attempt']} fichiers")
        if alert['violations']['log_tampering']:
            print(f"   ⚠️  Tentative modification logs: CONFIRMÉE")
            
        print(f"\n📁 PREUVES DISPONIBLES:")
        for evidence in alert['evidence_files']:
            print(f"   📄 {evidence}")
            
        print(f"\n⚖️  BASE LÉGALE: {alert['legal_basis']}")
        print(f"🔄 ALERTES ENVOYÉES: Toutes les 15 minutes depuis {alert['timestamp']}")
        
    def acknowledge_alert(self, alert_id):
        """Acquitter une alerte et clore le cycle"""
        print("\n[A] Acquitter et Clôturer")
        print("[E] Envoyer à la Justice")
        print("[S] Surveillance Renforcée")
        print("[Q] Quitter sans action")
        
        choice = input("Action: ").upper()
        
        if choice == "A":
            # Marquer comme acquitté
            status = "ACKNOWLEDGED_BY_DIRECTOR"
            timestamp = datetime.datetime.utcnow().isoformat()
            
            update_alert_status(alert_id, status, timestamp)
            print("✅ Alert acquittée. Dossier archivé.")
            
        elif choice == "E":
            # Transmission automatique au parquet
            transmit_to_prosecutor(alert_id)
            print("⚖️  Dossier transmis au Bureau du Procureur")
            
        elif choice == "S":
            # Activer surveillance renforcée sur l'agent
            activate_enhanced_surveillance(alert_id)
            print("👁️  Surveillance renforcée activée sur l'agent")
            
        return choice
    
    def run_dashboard(self):
        """Boucle principale du dashboard"""
        if not self.authenticate_director():
            return
            
        while self.session_active:
            self.fetch_alerts()
            
            if not self.alert_queue:
                print("\n✅ AUCUNE ALERTE EN ATTENTE - SYSTÈME SÉCURISÉ")
            else:
                print(f"\n🔴 {len(self.alert_queue)} ALERTES EN ATTENTE")
                
                for i, alert in enumerate(self.alert_queue, 1):
                    print(f"\n[{i}] Alerte #{alert['alert_id']} - {alert['culprit']['nom']}")
                    self.display_alert(alert)
                    
                    action = self.acknowledge_alert(alert['alert_id'])
                    
                    if action in ['A', 'E']:
                        self.alert_queue.remove(alert)
                        
            print("\n[R] Rafraîchir | [Q] Quitter")
            cmd = input("Commande: ").upper()
            
            if cmd == 'Q':
                self.session_active = False
                print("🔒 Session fermée. Journalisation complète.")
                
if __name__ == "__main__":
    dashboard = OmertaDashboard()
    dashboard.run_dashboard()
