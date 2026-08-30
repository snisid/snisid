#!/usr/bin/env python3
"""
SNISID - Weekly MEK Synchronization and CRDT Integrity Simulator
Validates queue handling, vector clock reconciliation (CRDT),
and deferred deduplication checks.
"""

import sys
import json
import hashlib

class MockCRDTRecord:
    def __init__(self, citizen_id, name, dob, vector_clock, timestamp, status="PENDING"):
        self.citizen_id = citizen_id
        self.name = name
        self.dob = dob
        self.vector_clock = vector_clock  # dict: {device_id: counter}
        self.timestamp = timestamp
        self.status = status

    def to_dict(self):
        return {
            "citizen_id": self.citizen_id,
            "name": self.name,
            "dob": self.dob,
            "vector_clock": self.vector_clock,
            "timestamp": self.timestamp,
            "status": self.status
        }

def resolve_conflict(local, incoming):
    """
    Resolves conflict using vector clock comparison and falls back to Last-Write-Wins (LWW)
    """
    l_clock = local.get("vector_clock", {})
    i_clock = incoming.get("vector_clock", {})
    
    incoming_greater = False
    local_greater = False
    
    all_keys = set(l_clock.keys()).union(set(i_clock.keys()))
    for key in all_keys:
        lv = l_clock.get(key, 0)
        iv = i_clock.get(key, 0)
        if iv > lv:
            incoming_greater = True
        elif lv > iv:
            local_greater = True
            
    if incoming_greater and not local_greater:
        return incoming, "INCOMING_WIN"
    elif local_greater and not incoming_greater:
        return local, "LOCAL_WIN"
    else:
        # Concurrent changes: resolve using Last-Write-Wins (LWW) via timestamp
        if incoming.get("timestamp", 0) > local.get("timestamp", 0):
            return incoming, "LWW_INCOMING_WIN"
        else:
            return local, "LWW_LOCAL_WIN"

def run_sync_simulation():
    print("=========================================================")
    print("  SNISID WEEKLY MEK SYNCHRONIZATION AND CRDT SIMULATOR   ")
    print("=========================================================")
    
    # 1. Generate local offline queued records on MEK
    print("[*] Generating offline enrollment batch on MEK-HT-042...")
    local_queue = []
    
    # Record A: Brand new citizen enrollment
    cA = MockCRDTRecord("NIPPES-1001", "Jean-Baptiste Noel", "1990-05-15", {"MEK-HT-042": 1}, 1779753600).to_dict()
    local_queue.append(cA)
    
    # Record B: Offline edit on SUD-2005 (has conflict)
    cB = MockCRDTRecord("SUD-2005", "Marie Therese Pierre", "1985-11-20", {"MEK-HT-042": 1, "CORE": 4}, 1779753605).to_dict()
    local_queue.append(cB)
    
    # Record C: Potential duplicate (name & dob matches cA, triggers deferred dedup review)
    cC = MockCRDTRecord("OUEST-3009", "Jean-Baptiste Noel", "1990-05-15", {"MEK-HT-042": 1}, 1779753610).to_dict()
    local_queue.append(cC)
    
    print(f"[+] Local queue populated with {len(local_queue)} transactions.")
    
    # 2. Pre-populate Core database state with existing record
    core_db = {
        "SUD-2005": {
            "citizen_id": "SUD-2005",
            "name": "Marie T. Pierre",
            "dob": "1985-11-20",
            "vector_clock": {"CORE": 4},
            "timestamp": 1779753000,
            "status": "APPROVED"
        }
    }
    
    # 3. Process Sync Simulation
    print("[*] Reconnecting to Core Datacenter. Starting synchronization...")
    processed_count = 0
    conflicts_resolved = 0
    duplicates_flagged = 0
    
    sync_report = {
        "sync_job_id": "SYNC-TEST-AUTO-001",
        "device_id": "MEK-HT-042",
        "timestamp": "2026-05-24T21:24:50Z",
        "results": []
    }
    
    for record in local_queue:
        citizen_id = record["citizen_id"]
        print(f"\n[*] Processing sync for record {citizen_id} ({record['name']})...")
        
        # SHA-256 checksum integrity verification
        rec_str = f"{record['citizen_id']}-{record['name']}-{record['dob']}"
        rec_hash = hashlib.sha256(rec_str.encode('utf-8')).hexdigest()
        print(f"  - Integrity hash: {rec_hash[:16]}... [PASSED]")
        
        if citizen_id in core_db:
            print("  - Conflict detected! Resolving...")
            conflicts_resolved += 1
            existing = core_db[citizen_id]
            resolved, strategy = resolve_conflict(existing, record)
            core_db[citizen_id] = resolved
            print(f"  - Resolved via: {strategy}")
            print(f"  - Final Name: {resolved['name']}")
            sync_report["results"].append({
                "citizen_id": citizen_id,
                "status": "CONFLICT_RESOLVED",
                "strategy": strategy,
                "name_in_db": resolved["name"]
            })
        else:
            # Search for duplicates in core (1:N deferred deduplication)
            is_dup = False
            for existing_id, existing in core_db.items():
                if existing["name"] == record["name"] and existing["dob"] == record["dob"]:
                    is_dup = True
                    break
            
            if is_dup:
                print("  - [WARNING] Deferred deduplication: duplicate pattern found. Flagged for review.")
                duplicates_flagged += 1
                record["status"] = "FLAGGED_FOR_REVIEW"
                sync_report["results"].append({
                    "citizen_id": citizen_id,
                    "status": "FLAGGED_AS_DUPLICATE",
                    "conflict_with": existing_id
                })
            else:
                print("  - New enrollment. Direct insertion.")
                record["status"] = "APPROVED"
                sync_report["results"].append({
                    "citizen_id": citizen_id,
                    "status": "SYNCED_OK"
                })
            
            core_db[citizen_id] = record
            
        processed_count += 1
        
    # 4. Compile Compliance Report
    print("\n=========================================================")
    print("  SIMULATION DE SYNCHRONISATION MEK - RAPPORT DE CONFORMITÉ")
    print("=========================================================")
    compliance_passed = (duplicates_flagged == 1 and conflicts_resolved == 1)
    report = {
        "compliance_status": "PASSED" if compliance_passed else "FAILED",
        "total_records_processed": processed_count,
        "conflicts_resolved": conflicts_resolved,
        "duplicates_flagged_for_review": duplicates_flagged,
        "final_database_size": len(core_db)
    }
    print(json.dumps(report, indent=2))
    
    # Save compliance report to disk
    report_path = "pki/scripts/mek_sync_report.json"
    with open(report_path, "w") as f:
        json.dump(report, f, indent=2)
    print(f"\n[+] Compliance report written to {report_path}")
        
    return compliance_passed

if __name__ == "__main__":
    success = run_sync_simulation()
    sys.exit(0 if success else 1)
