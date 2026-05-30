#!/usr/bin/env python3
"""
SNISID National Identity Reconciliation Engine (NIRE) - Test Suite
Verifies demographic and biometric reconciliation edge cases.
"""

import unittest
import sys
import os

# Append the current directory of the script to sys.path so it can find reconciliation_logic
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from reconciliation_logic import NationalIdentityReconciliationEngine, levenshtein_similarity

class TestIdentityReconciliation(unittest.TestCase):
    
    def setUp(self):
        # Base master record of a citizen
        self.master_citizen = {
            "iui": "HT-SNISID-100001",
            "last_name": "JEAN-BAPTISTE",
            "first_name": "Melissa",
            "birth_date": "1988-04-12",
            "birth_place": "Port-au-Prince",
            "mother_first_name": "Marie-Thérèse",
            "biometrics_hash": "fingerprint_wsq_98122",
            "facial_embedding": [0.15, -0.22, 0.88, 0.04, -0.52]
        }

    def test_levenshtein_similarity(self):
        self.assertEqual(levenshtein_similarity("MELISSA", "MELISSA"), 1.0)
        self.assertEqual(levenshtein_similarity("MELISSA", ""), 0.0)
        # 1 insertion/substitution away
        sim = levenshtein_similarity("MELISSA", "MELISA")
        self.assertGreater(sim, 0.8)

    def test_exact_match_reconciliation(self):
        # Reconciling exact same record (Idempotency test)
        decision, details = NationalIdentityReconciliationEngine.reconcile(self.master_citizen, self.master_citizen)
        self.assertEqual(decision, "IDEMPOTENT_MERGE")
        self.assertEqual(details["demographic_similarity_score"], 1.0)
        self.assertEqual(details["biometric_fusion_score"], 1.0)

    def test_phonetic_and_typo_match(self):
        # Legitimate citizen with slight variations (e.g., Jean-Baptiste without hyphen, spelling typo in mother's name)
        re_enrolling_citizen = {
            "last_name": "JEAN BAPTISTE",
            "first_name": "Mélissa",
            "birth_date": "1988-04-12",
            "birth_place": "Port-au-Prince",
            "mother_first_name": "Marie Therese",
            "biometrics_hash": "fingerprint_wsq_98122", # Biometrics match perfectly
            "facial_embedding": [0.15, -0.22, 0.88, 0.04, -0.52]
        }
        
        decision, details = NationalIdentityReconciliationEngine.reconcile(re_enrolling_citizen, self.master_citizen)
        self.assertEqual(decision, "IDEMPOTENT_MERGE")
        self.assertGreaterEqual(details["demographic_similarity_score"], 0.85)
        self.assertEqual(details["biometric_fusion_score"], 1.0)

    def test_identity_fraud_detection(self):
        # USURPER: same biometrics but completely different names (e.g. attempting multiple registration under fake names)
        usurper_citizen = {
            "last_name": "LUBIN",
            "first_name": "Pierre-Paul",
            "birth_date": "1972-11-04",
            "birth_place": "Gonaïves",
            "mother_first_name": "Evelyne",
            "biometrics_hash": "fingerprint_wsq_98122", # STOLE the biometric profile of Melissa!
            "facial_embedding": [0.15, -0.22, 0.88, 0.04, -0.52]
        }
        
        decision, details = NationalIdentityReconciliationEngine.reconcile(usurper_citizen, self.master_citizen)
        self.assertEqual(decision, "FRAUD_SUSPICION")
        self.assertLess(details["demographic_similarity_score"], 0.50)
        self.assertEqual(details["biometric_fusion_score"], 1.0)

    def test_facial_recognition_similarity(self):
        # Citizen with slightly different facial lighting (cosine similarity close but not perfect 1.0)
        # Note: We must change the biometrics_hash so it falls back to facial embedding comparison!
        subtle_photo_change = {
            "last_name": "JEAN-BAPTISTE",
            "first_name": "Melissa",
            "birth_date": "1988-04-12",
            "birth_place": "Port-au-Prince",
            "mother_first_name": "Marie-Thérèse",
            "biometrics_hash": "different_fingerprint_wsq_1234",
            "facial_embedding": [0.14, -0.21, 0.89, 0.05, -0.51] # 99% similar facial embedding
        }
        
        bio_score = NationalIdentityReconciliationEngine.calculate_biometric_score(subtle_photo_change, self.master_citizen)
        self.assertGreater(bio_score, 0.98)
        self.assertLess(bio_score, 1.0) # Cosine sim is slightly lower than 1.0

    def test_new_unique_identity(self):
        # A completely different citizen (No biometric match and different demographics)
        new_citizen = {
            "last_name": "CHARLES",
            "first_name": "Stevenson",
            "birth_date": "1994-07-22",
            "birth_place": "Cap-Haïtien",
            "mother_first_name": "Chantal",
            "biometrics_hash": "fingerprint_wsq_12099",
            "facial_embedding": [-0.55, 0.12, 0.01, 0.94, 0.33]
        }
        
        decision, details = NationalIdentityReconciliationEngine.reconcile(new_citizen, self.master_citizen)
        self.assertEqual(decision, "NEW_UNIQUE_IDENTITY")
        self.assertLess(details["demographic_similarity_score"], 0.40)
        # For completely different faces, the cosine similarity is very low (or negative)
        self.assertLess(details["biometric_fusion_score"], 0.30)

if __name__ == "__main__":
    unittest.main()
