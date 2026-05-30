#!/usr/bin/env python3
"""
SNISID Offline-First Synchronization Protocol Simulator
Simulates the Delayed Sync Protocol from a Local Edge Node (LEN) in a remote Haitian department
(e.g., Grand'Anse or Sud) back to the Sovereign Central Datacenter in Port-au-Prince.
"""

import json
import hashlib
from datetime import datetime

class LocalEdgeNode:
    """
    Simulates a secure Local Edge Node running in an offline office (e.g., Jeremie BLC).
    """
    def __init__(self, node_id, region):
        self.node_id = node_id
        self.region = region
        self.transaction_queue = []
        self.local_cache = {}

    def local_enroll(self, last_name, first_name, birth_date, gender, biometrics_hash):
        """
        Enrolls a citizen locally without network.
        Generates a signed temporary IUI (IUI-T) and secures the transaction in the queue.
        """
        temp_id = f"HT-T-{self.node_id}-{len(self.transaction_queue) + 1001}"
        timestamp = datetime.now().isoformat()
        
        # Structure the payload
        payload = {
            "iui_temp": temp_id,
            "last_name": last_name.upper().strip(),
            "first_name": first_name.strip(),
            "birth_date": birth_date,
            "gender": gender,
            "birth_place": self.region,
            "biometrics_hash": biometrics_hash,
            "enrolled_at": timestamp,
            "node_id": self.node_id
        }
        
        # Simulate local cryptographic signing of the transaction payload
        payload_bytes = json.dumps(payload, sort_keys=True).encode()
        local_signature = hashlib.sha256(payload_bytes + f"SECRET_KEY_FOR_{self.node_id}".encode()).hexdigest()
        
        transaction = {
            "payload": payload,
            "signature": local_signature,
            "hash": hashlib.sha256(payload_bytes).hexdigest()
        }
        
        self.transaction_queue.append(transaction)
        self.local_cache[temp_id] = payload
        print(f"[{self.node_id}] Local enrollment succeeded! Temp ID: {temp_id} created for {first_name} {last_name}")
        return temp_id

class CentralDatacenter:
    """
    Simulates the Sovereign Central Datacenter running the NIRE (National Identity Reconciliation Engine).
    """
    def __init__(self):
        # Already registered citizens in the central production database
        self.central_db = {
            "HT-SNISID-100001": {
                "iui": "HT-SNISID-100001",
                "last_name": "LUBIN",
                "first_name": "Joseph",
                "birth_date": "1983-05-18",
                "gender": "M",
                "biometrics_hash": "wsq_hash_joseph_lubin"
            }
        }
        self.sync_logs = []
        self.quarantine_db = []

    def verify_node_signature(self, tx):
        """
        Verifies that the transaction was signed by a registered and trusted Local Edge Node key.
        """
        payload = tx["payload"]
        node_id = payload["node_id"]
        payload_bytes = json.dumps(payload, sort_keys=True).encode()
        expected_sig = hashlib.sha256(payload_bytes + f"SECRET_KEY_FOR_{node_id}".encode()).hexdigest()
        return tx["signature"] == expected_sig

    def process_sync_queue(self, node_id, queue):
        """
        Processes a batch of transactions uploaded by a Local Edge Node.
        Applies cryptographic checks, biometric deduplication, and conflict resolution.
        """
        print(f"\n[Central DC] Initiating Sync Session with Node: {node_id}...")
        success_count = 0
        error_count = 0
        
        for tx in queue:
            payload = tx["payload"]
            temp_id = payload["iui_temp"]
            
            # 1. Cryptographic Security Check
            if not self.verify_node_signature(tx):
                print(f"  * REJECTED: Cryptographic signature mismatch on transaction {temp_id}!")
                self.quarantine_db.append({
                    "transaction": tx,
                    "reason": "Cryptographic Signature Fraud/Corruption"
                })
                error_count += 1
                continue
            
            # 2. Biometric Deduplication Check
            biometrics = payload["biometrics_hash"]
            duplicate_found = False
            
            for central_iui, citizen in self.central_db.items():
                if citizen["biometrics_hash"] == biometrics:
                    print(f"  * CONFLICT DETECTED: Biometrics on {temp_id} matches existing central citizen {central_iui}!")
                    # Check if demographics match (Merge/Idempotent update) or mismatch (FRAUD!)
                    if citizen["last_name"] == payload["last_name"] and citizen["birth_date"] == payload["birth_date"]:
                        print(f"    - Resolution: Legitimate duplicate registration. Merging logs.")
                        self.sync_logs.append({
                            "action": "MERGE_RECORD",
                            "iui": central_iui,
                            "temp_id": temp_id,
                            "node_id": node_id
                        })
                    else:
                        print(f"    - CRITICAL WAR ROOM ALERTE: Identity Fraud suspicion. Placing {temp_id} in Quarantine.")
                        self.quarantine_db.append({
                            "transaction": tx,
                            "reason": "Biometric Identity Hijacking / Usurpation Detected",
                            "conflict_with_iui": central_iui
                        })
                        error_count += 1
                    duplicate_found = True
                    break
                    
            if duplicate_found:
                continue
                
            # 3. Successful Unique Citizen Registration
            new_iui = f"HT-SNISID-{100000 + len(self.central_db) + 1}"
            self.central_db[new_iui] = {
                "iui": new_iui,
                "last_name": payload["last_name"],
                "first_name": payload["first_name"],
                "birth_date": payload["birth_date"],
                "gender": payload["gender"],
                "biometrics_hash": biometrics
            }
            self.sync_logs.append({
                "action": "CREATE_RECORD",
                "iui": new_iui,
                "temp_id": temp_id,
                "node_id": node_id
            })
            print(f"  * SUCCESS: Synchronized and Promoted {temp_id} -> {new_iui}")
            success_count += 1
            
        print(f"[Central DC] Sync complete. Node {node_id} processed: Success={success_count}, Rejected/Quarantined={error_count}\n")
        return success_count, error_count

if __name__ == "__main__":
    # Initialize components
    dc = CentralDatacenter()
    jeremie_node = LocalEdgeNode(node_id="LEN-JEREMIE", region="Grand'Anse")
    
    # 1. Enrolls 3 citizens offline at Jeremie
    # Citizen A: Completely new unique citizen
    jeremie_node.local_enroll(
        last_name="PHILIPPE", first_name="Jean", birth_date="1992-07-15", gender="M",
        biometrics_hash="wsq_hash_jean_philippe"
    )
    
    # Citizen B: Already registered central citizen attempting duplicate registration (same demographics & biometrics)
    jeremie_node.local_enroll(
        last_name="LUBIN", first_name="Joseph", birth_date="1983-05-18", gender="M",
        biometrics_hash="wsq_hash_joseph_lubin" # Already in central_db
    )
    
    # Citizen C: Fraudster stealing Joseph Lubin's biometrics but registering as a different person
    jeremie_node.local_enroll(
        last_name="METELLUS", first_name="Pierre", birth_date="1965-12-01", gender="M",
        biometrics_hash="wsq_hash_joseph_lubin" # Stealing biometrics!
    )
    
    # 2. Simulate internet reconnection and uploading the queue to central Datacenter
    success, rejected = dc.process_sync_queue(jeremie_node.node_id, jeremie_node.transaction_queue)
    
    # Print state of Databases after sync
    print("="*60)
    print("                 POST-SYNC CENTRAL STATE REPORT")
    print("="*60)
    print(f"Active Central Citizens Count: {len(dc.central_db)}")
    print(f"Quarantined/Failed Registrations Count: {len(dc.quarantine_db)}")
    print(f"Sync Logs Recorded: {len(dc.sync_logs)}")
    print("="*60)
    
    # Write output log to file
    sync_report = {
        "sync_logs": dc.sync_logs,
        "quarantine_db": dc.quarantine_db,
        "central_citizens_count": len(dc.central_db)
    }
    with open("Deployment-Rollout/Offline-Deployment/offline_sync_report.json", "w") as f:
        json.dump(sync_report, f, indent=4)
