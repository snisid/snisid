#!/usr/bin/env python3
"""
SNISID Migration Factory - Data Cleansing & Migration Pipeline Simulator
This script simulates the industrial ETL pipeline for importing, cleansing, and validating
raw government legacy records (from ONI, Civil Registry, etc.) into clean SNISID records.
"""

import json
import re
import unicodedata
from datetime import datetime

# Mock legacy database representing dirty, corrupted, and duplicated records from various departments
RAW_LEGACY_DB = [
    {
        "source": "ONI_Legacy",
        "raw_id": "ONI-98231",
        "last_name": "JEAN BAPTISTE ",
        "first_name": "mélissa",
        "birth_date": "12/04/1988",
        "gender": "F",
        "birth_place": "Port-au-Prince",
        "legacy_cin": "111-222-333-44",
        "biometrics_hash": "wsq_hash_98231"
    },
    {
        "source": "Civil_Registry",
        "raw_id": "CIVIL-872A",
        "last_name": "Janbatis",  # Phonetic and spelling variation of Jean Baptiste
        "first_name": "Melissa",
        "birth_date": "12-04-1988",  # Different date format
        "gender": "F",
        "birth_place": "P-au-P",
        "legacy_cin": None,
        "biometrics_hash": None  # No biometrics in historic paper registry -> will be merged with above!
    },
    {
        "source": "ONI_Legacy",
        "raw_id": "ONI-45120",
        "last_name": "HYPPOLITE",
        "first_name": "Jean-Pierre",
        "birth_date": "28/02/1980",  # Valid date now so it gets added to the DB
        "gender": "M",
        "birth_place": "Cap-Haïtien",
        "legacy_cin": "444-555-666-77",
        "biometrics_hash": "wsq_hash_45120"
    },
    {
        "source": "Immigration_DIE",
        "raw_id": "DIE-9921",
        "last_name": "LUBIN",
        "first_name": "Pierre",
        "birth_date": "15/10/1985",
        "gender": "M",
        "birth_place": "Okap",  # Creole name for Cap-Haitien
        "legacy_cin": "444-555-666-77",  # CONFLICTING CIN (same as above but completely different person/birthdate!)
        "biometrics_hash": "wsq_hash_9921"
    },
    {
        "source": "Civil_Registry",
        "raw_id": "CIVIL-212",
        "last_name": "DURENT",
        "first_name": "Jn-Claude",  # Abbreviated first name
        "birth_date": "01/01/1975",
        "gender": "M",
        "birth_place": "Gonaïves",
        "legacy_cin": None,
        "biometrics_hash": None
    },
    {
        "source": "ONI_Legacy",
        "raw_id": "ONI-12344",
        "last_name": "ST-LOUIS",
        "first_name": "Marie N/A",  # Missing/corrupted attribute value -> will go to Quarantine
        "birth_date": "31/11/1990",  # INVALID DATE: November has only 30 days! -> will go to Quarantine
        "gender": "F",
        "birth_place": "Jacmel",
        "legacy_cin": "777-888-999-00",
        "biometrics_hash": "wsq_hash_12344"
    }
]

class SoundexHT:
    """
    Haitian Soundex - Custom phonetic algorithm for French & Haitian Creole names
    Normalizes spelling variations like: Janbatis -> Jean-Baptiste, Okap -> Cap-Haitien (place)
    """
    @staticmethod
    def phonetize(text):
        if not text:
            return ""
        # Convert to uppercase and strip whitespace
        s = text.upper().strip()
        # Remove spaces and hyphens for unified phonetic matching
        s = s.replace(" ", "").replace("-", "")
        # Replace common Haitian Creole/French variations
        s = re.sub(r'JAN', 'JEAN', s)
        s = re.sub(r'BATIS', 'BAPTISTE', s)
        s = re.sub(r'OKAP', 'CAP-HAITIEN', s)
        s = re.sub(r'P-AU-P|PORTAU PRINCE', 'PORT-AU-PRINCE', s)
        s = re.sub(r'JN', 'JEAN', s)
        return s

class DataCleansingProgram:
    @staticmethod
    def clean_name(name):
        if not name:
            return "UNKNOWN"
        # De-normalize Unicode accents (e.g. mélissa -> melissa)
        normalized = unicodedata.normalize('NFD', name)
        ascii_text = "".join([c for c in normalized if unicodedata.category(c) != 'Mn'])
        
        # Convert to uppercase
        c = ascii_text.strip().upper()
        # Remove non-alphabetic chars except hyphen and apostrophe
        c = re.sub(r'[^A-Z\-\'\s]', '', c)
        c = re.sub(r'\s+', ' ', c)
        # Strip bad tags like "N/A", "NULL", "TEST"
        c = re.sub(r'\b(N/A|NULL|TEST|INCONNU|UNKNOWN)\b', '', c).strip()
        return c if c else "UNKNOWN"

    @staticmethod
    def parse_and_validate_date(date_str):
        if not date_str:
            return None, "Missing Date"
        # Try different common formats
        for fmt in ("%d/%m/%Y", "%d-%m-%Y", "%Y-%m-%d"):
            try:
                dt = datetime.strptime(date_str, fmt)
                # Check for historical limits
                if dt.year < 1900 or dt.year > datetime.now().year:
                    return None, f"Year out of bounds: {dt.year}"
                return dt.strftime("%Y-%m-%d"), None
            except ValueError:
                continue
        return None, "Invalid format or date arithmetic error (e.g. Feb 29 on non-leap year)"

    @staticmethod
    def clean_birth_place(place):
        if not place:
            return "UNKNOWN"
        # Standardize common location aliases
        p = SoundexHT.phonetize(place)
        return p

class MigrationFactory:
    def __init__(self):
        self.verified_database = {}
        self.quarantine_db = []
        self.conflicts = []
        self.identity_map = {} # Maps legacy ID -> SNISID Universal ID (IUI)

    def run_pipeline(self, records):
        print("="*80)
        print("          SNISID MIGRATION FACTORY PIPELINE EXECUTION STARTED")
        print("="*80)
        
        for idx, rec in enumerate(records):
            print(f"\n[Processing Record {idx+1}] Source: {rec['source']} | Raw ID: {rec['raw_id']}")
            
            # 1. Clean attributes
            clean_last = DataCleansingProgram.clean_name(rec['last_name'])
            clean_first = DataCleansingProgram.clean_name(rec['first_name'])
            clean_date, date_err = DataCleansingProgram.parse_and_validate_date(rec['birth_date'])
            clean_place = DataCleansingProgram.clean_birth_place(rec['birth_place'])
            
            # Print intermediate cleaning results
            print(f"  * Cleansing Names: '{rec['last_name']}' -> '{clean_last}', '{rec['first_name']}' -> '{clean_first}'")
            if date_err:
                print(f"  * ERROR: Date Cleansing Failed: {date_err}")
                self.quarantine_db.append({
                    "raw_record": rec,
                    "reason": f"Date Validation Failure: {date_err}",
                    "level": "Q1 - Critical Invalid"
                })
                continue
            
            # Detect missing vital fields (Quarantine Q2)
            if clean_last == "UNKNOWN" or clean_first == "UNKNOWN":
                print(f"  * ERROR: Missing vital attributes (first/last name unknown).")
                self.quarantine_db.append({
                    "raw_record": rec,
                    "reason": "Missing vital attribute (Name)",
                    "level": "Q2 - Incomplete Record"
                })
                continue

            # Standardize phonetic identity for matching
            phonetic_id = f"{SoundexHT.phonetize(clean_last)}|{SoundexHT.phonetize(clean_first)}|{clean_date}"
            
            # 2. Identity Reconciliation & Deduplication (Simulating demographic matching)
            duplicate_found = False
            for iui, existing_rec in self.verified_database.items():
                existing_phonetic = f"{SoundexHT.phonetize(existing_rec['last_name'])}|{SoundexHT.phonetize(existing_rec['first_name'])}|{existing_rec['birth_date']}"
                
                # Check for demographic phonetic match
                if phonetic_id == existing_phonetic:
                    print(f"  * MATCH FOUND: Deduplicated with existing SNISID: {iui} (Demographic Soundex Match)")
                    # Merge information (enrich legacy CIN or biometrics if empty in existing)
                    if not existing_rec["legacy_cin"] and rec["legacy_cin"]:
                        existing_rec["legacy_cin"] = rec["legacy_cin"]
                        print(f"    - Enriched legacy_cin: {rec['legacy_cin']}")
                    if not existing_rec["biometrics_hash"] and rec["biometrics_hash"]:
                        existing_rec["biometrics_hash"] = rec["biometrics_hash"]
                        print(f"    - Enriched biometrics_hash: {rec['biometrics_hash']}")
                    
                    self.identity_map[rec['raw_id']] = iui
                    duplicate_found = True
                    break
            
            if duplicate_found:
                continue

            # 3. Check for Identity Conflict (CIN matching different person)
            cin_conflict = False
            if rec["legacy_cin"]:
                for iui, existing_rec in self.verified_database.items():
                    if existing_rec["legacy_cin"] == rec["legacy_cin"]:
                        print(f"  * CRITICAL CONFLICT: CIN {rec['legacy_cin']} already owned by {existing_rec['first_name']} {existing_rec['last_name']} ({iui})")
                        self.conflicts.append({
                            "conflicting_record_1": existing_rec,
                            "conflicting_record_2": {
                                "raw_id": rec["raw_id"],
                                "last_name": clean_last,
                                "first_name": clean_first,
                                "birth_date": clean_date,
                                "legacy_cin": rec["legacy_cin"],
                                "source": rec["source"]
                            },
                            "conflict_type": "CIN Duplicate Owner Violation"
                        })
                        self.quarantine_db.append({
                            "raw_record": rec,
                            "reason": f"CIN Conflict with {iui}",
                            "level": "Q1 - Identity Conflict"
                        })
                        cin_conflict = True
                        break
            
            if cin_conflict:
                continue

            # 4. Generate Clean SNISID Record
            new_iui = f"HT-SNISID-{100000 + len(self.verified_database) + 1}"
            clean_record = {
                "iui": new_iui,
                "last_name": clean_last,
                "first_name": clean_first,
                "birth_date": clean_date,
                "gender": rec["gender"],
                "birth_place": clean_place,
                "legacy_cin": rec["legacy_cin"],
                "biometrics_hash": rec["biometrics_hash"],
                "migration_metadata": {
                    "migrated_at": datetime.now().isoformat(),
                    "source_system": rec["source"],
                    "raw_id": rec["raw_id"]
                }
            }
            
            self.verified_database[new_iui] = clean_record
            self.identity_map[rec['raw_id']] = new_iui
            print(f"  * SUCCESS: New verified SNISID Record created: {new_iui}")

        # Final Report
        print("\n" + "="*80)
        print("                        MIGRATION PIPELINE FINAL REPORT")
        print("="*80)
        print(f"  * Total Raw Records Ingested: {len(records)}")
        print(f"  * Successfully Migrated & Verified SNISID Records: {len(self.verified_database)}")
        print(f"  * Quarantined Records: {len(self.quarantine_db)}")
        print(f"  * Active Conflicts Flagged: {len(self.conflicts)}")
        print("="*80)
        
        return {
            "verified_db": self.verified_database,
            "quarantine_db": self.quarantine_db,
            "conflicts": self.conflicts
        }

if __name__ == "__main__":
    factory = MigrationFactory()
    results = factory.run_pipeline(RAW_LEGACY_DB)
    
    # Save the output to a verification JSON in the same folder
    output_data = {
        "verified_database": results["verified_db"],
        "quarantine_database": results["quarantine_db"],
        "conflicts": results["conflicts"]
    }
    
    with open("Deployment-Rollout/Migration-Factory/migration_pipeline_output.json", "w", encoding="utf-8") as f:
        json.dump(output_data, f, indent=4, ensure_ascii=False)
        print("\nPipeline execution saved to: Deployment-Rollout/Migration-Factory/migration_pipeline_output.json")
