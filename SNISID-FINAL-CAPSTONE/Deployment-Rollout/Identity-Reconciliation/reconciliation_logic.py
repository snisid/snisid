#!/usr/bin/env python3
"""
SNISID National Identity Reconciliation Engine (NIRE) - Core Engine
Implements the hybrid demographic & biometric reconciliation logic, duplicate resolution,
and identity fraud detection model.
"""

import math

def levenshtein_similarity(s1, s2):
    """
    Computes Levenshtein similarity score between 0.0 and 1.0.
    """
    s1 = s1.upper().strip()
    s2 = s2.upper().strip()
    if s1 == s2:
        return 1.0
    if not s1 or not s2:
        return 0.0
    
    # Initialize matrix
    rows = len(s1) + 1
    cols = len(s2) + 1
    dist = [[0 for _ in range(cols)] for _ in range(rows)]
    
    for i in range(1, rows):
        dist[i][0] = i
    for j in range(1, cols):
        dist[0][j] = j
        
    for col in range(1, cols):
        for row in range(1, rows):
            if s1[row-1] == s2[col-1]:
                cost = 0
            else:
                cost = 1
            dist[row][col] = min(
                dist[row-1][col] + 1,      # deletion
                dist[row][col-1] + 1,      # insertion
                dist[row-1][col-1] + cost  # substitution
            )
            
    lev_dist = dist[-1][-1]
    max_len = max(len(s1), len(s2))
    return 1.0 - (lev_dist / max_len)

class NationalIdentityReconciliationEngine:
    @staticmethod
    def calculate_demographic_score(rec1, rec2):
        """
        Computes a weighted demographic similarity score between 0.0 and 1.0.
        Weights:
          - Last Name: 25%
          - First Name: 20%
          - Birth Date: 25%
          - Birth Place: 15%
          - Mother's First Name: 15%
        """
        # Compare names
        score_last = levenshtein_similarity(rec1.get("last_name", ""), rec2.get("last_name", ""))
        score_first = levenshtein_similarity(rec1.get("first_name", ""), rec2.get("first_name", ""))
        
        # Compare Birth Date (exact string match = 1.0, else 0.0, but with partial match for typos)
        dob1 = rec1.get("birth_date", "")
        dob2 = rec2.get("birth_date", "")
        if dob1 == dob2:
            score_dob = 1.0
        else:
            # Check Levenshtein for 1-2 character typos (e.g. 1988-04-12 vs 1988-04-21)
            dob_lev = levenshtein_similarity(dob1, dob2)
            score_dob = dob_lev if dob_lev >= 0.8 else 0.0
            
        # Compare Birth Place
        score_place = levenshtein_similarity(rec1.get("birth_place", ""), rec2.get("birth_place", ""))
        
        # Compare Mother's name (if present, default to 1.0 if not provided on either to avoid penalizing)
        m1 = rec1.get("mother_first_name", "")
        m2 = rec2.get("mother_first_name", "")
        if not m1 or not m2:
            score_mother = 1.0 # Neutral
        else:
            score_mother = levenshtein_similarity(m1, m2)
            
        # Weighted calculation
        weighted_score = (
            (score_last * 0.25) +
            (score_first * 0.20) +
            (score_dob * 0.25) +
            (score_place * 0.15) +
            (score_mother * 0.15)
        )
        return round(weighted_score, 4)

    @staticmethod
    def calculate_biometric_score(rec1, rec2):
        """
        Simulates Biometric Fusion Score from fingerprints, face, and iris scan.
        If biometric hashes/signatures match perfectly, return 1.0.
        If they do not match at all, return 0.0.
        In real production, this parses actual minutiae vectors and face embeddings.
        """
        b1 = rec1.get("biometrics_hash", "")
        b2 = rec2.get("biometrics_hash", "")
        if not b1 or not b2:
            return 0.0  # No biometrics available for cross-match
        
        if b1 == b2:
            return 1.0
            
        # Cosine similarity simulation for neural face embeddings
        f_embed1 = rec1.get("facial_embedding", [])
        f_embed2 = rec2.get("facial_embedding", [])
        if f_embed1 and f_embed2 and len(f_embed1) == len(f_embed2):
            # Calculate cosine similarity
            dot_product = sum(a*b for a,b in zip(f_embed1, f_embed2))
            norm_a = math.sqrt(sum(a*a for a in f_embed1))
            norm_b = math.sqrt(sum(b*b for b in f_embed2))
            if norm_a == 0 or norm_b == 0:
                return 0.0
            return round(dot_product / (norm_a * norm_b), 4)
            
        return 0.0

    @classmethod
    def reconcile(cls, new_record, existing_record):
        """
        Reconciles a new registration against an existing record in the master base.
        Returns a tuple: (Decision, Details)
        Possible Decisions:
          - 'IDEMPOTENT_MERGE': Identical person, merge and update.
          - 'FRAUD_SUSPICION': Same physical person (biometrics match) but totally different identity text!
          - 'MANUAL_AUDIT': Ambiguous score, needs human verification.
          - 'NEW_UNIQUE_IDENTITY': No match, completely new unique citizen registration.
        """
        dem_score = cls.calculate_demographic_score(new_record, existing_record)
        bio_score = cls.calculate_biometric_score(new_record, existing_record)
        
        details = {
            "demographic_similarity_score": dem_score,
            "biometric_fusion_score": bio_score,
            "existing_iui": existing_record.get("iui"),
            "existing_name": f"{existing_record.get('first_name')} {existing_record.get('last_name')}"
        }
        
        # Scenario 1: Perfect Biometric Match AND High Demographic Match -> Merge
        if bio_score >= 0.88 and dem_score >= 0.80:
            return "IDEMPOTENT_MERGE", details
            
        # Scenario 2: Biometrics match perfectly, BUT text is completely different -> Fraud USURPATION!
        if bio_score >= 0.88 and dem_score < 0.50:
            return "FRAUD_SUSPICION", details
            
        # Scenario 3: Ambiguous region
        if (0.60 <= bio_score < 0.88) or (0.70 <= dem_score < 0.80 and bio_score > 0.0):
            return "MANUAL_AUDIT", details
            
        # Scenario 4: No Match
        return "NEW_UNIQUE_IDENTITY", details

# Quick self-test inside script if executed directly
if __name__ == "__main__":
    # Test cases
    rec_original = {
        "iui": "HT-SNISID-100001",
        "last_name": "JEAN-BAPTISTE",
        "first_name": "Melissa",
        "birth_date": "1988-04-12",
        "birth_place": "Port-au-Prince",
        "mother_first_name": "Marie",
        "biometrics_hash": "wsq_hash_98231"
    }
    
    # Typos and spelling variations (Legitimate citizen coming back to update or re-register)
    rec_typo = {
        "last_name": "Janbatis",
        "first_name": "Mélissa",
        "birth_date": "1988-04-12",
        "birth_place": "P-au-P",
        "mother_first_name": "Mari",
        "biometrics_hash": "wsq_hash_98231"
    }
    
    # Usurper: same biometrics but completely different demographic info! (FRAUD)
    rec_usurper = {
        "last_name": "LUBIN",
        "first_name": "Pierre",
        "birth_date": "1975-08-25",
        "birth_place": "Cap-Haïtien",
        "mother_first_name": "Rose",
        "biometrics_hash": "wsq_hash_98231" # Stole biometrics!
    }
    
    print("Testing Levenshtein Similarity:")
    print(f"  - 'MELISSA' vs 'MELISSA': {levenshtein_similarity('MELISSA', 'MELISSA')}")
    print(f"  - 'JEAN-BAPTISTE' vs 'JANBATIS': {levenshtein_similarity('JEAN-BAPTISTE', 'JANBATIS'):.4f}")
    
    print("\nTesting NIRE Engine:")
    decision, det = NationalIdentityReconciliationEngine.reconcile(rec_typo, rec_original)
    print(f"  * Case 1 (Legitimate Typo Update): Decision={decision} | DemoScore={det['demographic_similarity_score']} | BioScore={det['biometric_fusion_score']}")
    
    decision_fraud, det_fraud = NationalIdentityReconciliationEngine.reconcile(rec_usurper, rec_original)
    print(f"  * Case 2 (Identity Fraud Attempt): Decision={decision_fraud} | DemoScore={det_fraud['demographic_similarity_score']} | BioScore={det_fraud['biometric_fusion_score']}")
