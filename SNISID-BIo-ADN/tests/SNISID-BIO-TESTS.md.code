# SNISID-BIO-ADN — Plan de Tests & Validation
**Document ID :** SNISID-BIO-TST-001 | **Version :** 1.0.0

---

## 1. STRATÉGIE DE TESTS

```
Tests Unitaires       ─── Matching engine, calcul LR, qualité score
Tests d'Intégration   ─── LDIS→SDIS→NDIS, Kafka, PostgreSQL
Tests de Performance  ─── LAPI < 200ms, matching < 5s
Tests de Sécurité     ─── Contrôle d'accès RBAC, chiffrement HSM
Tests E2E             ─── Workflow complet : prélèvement → hit → alerte
```

---

## 2. TESTS UNITAIRES — Moteur de Matching

```python
# tests/test_dna_matching.py

import pytest
from unittest.mock import Mock, AsyncMock
from bio_adn.engine.matcher import DNAMatcher, STRProfile, MatchResult
from bio_adn.engine.scoring import calculate_lr_score

# ──────────────────────────────────────────
# Données de test : profils STR standardisés
# ──────────────────────────────────────────

PROFILE_FULL_MATCH = {
    "CSF1PO":   {"value1": 10, "value2": 12},
    "D3S1358":  {"value1": 15, "value2": 17},
    "D5S818":   {"value1": 11, "value2": 13},
    "D7S820":   {"value1": 8,  "value2": 10},
    "D8S1179":  {"value1": 13, "value2": 14},
    "D13S317":  {"value1": 9,  "value2": 11},
    "D16S539":  {"value1": 11, "value2": 12},
    "D18S51":   {"value1": 14, "value2": 16},
    "D21S11":   {"value1": 29, "value2": 30},
    "FGA":      {"value1": 21, "value2": 23},
    "TH01":     {"value1": 7,  "value2": 9},
    "TPOX":     {"value1": 8,  "value2": 11},
    "vWA":      {"value1": 16, "value2": 17},
    "D1S1656":  {"value1": 13, "value2": 15},
    "D2S441":   {"value1": 10, "value2": 11},
    "D2S1338":  {"value1": 19, "value2": 23},
    "D10S1248": {"value1": 13, "value2": 15},
    "D12S391":  {"value1": 18, "value2": 20},
    "D19S433":  {"value1": 13, "value2": 14},
    "D22S1045": {"value1": 15, "value2": 16},
}

PROFILE_PARTIAL_MATCH = {**PROFILE_FULL_MATCH}  # 18/20 loci identiques
PROFILE_PARTIAL_MATCH["D3S1358"] = {"value1": 14, "value2": 16}   # locus diff
PROFILE_PARTIAL_MATCH["FGA"]     = {"value1": 20, "value2": 24}   # locus diff

PROFILE_NO_MATCH = {k: {"value1": v["value1"] + 3, "value2": v["value2"] + 2}
                    for k, v in PROFILE_FULL_MATCH.items()}

# ──────────────────────────────────────────
# Tests de correspondance
# ──────────────────────────────────────────

class TestDNAMatching:

    def test_full_match_detected(self):
        """Un profil identique doit produire un FULL_MATCH (confidence ≥ 0.999)"""
        score = calculate_lr_score(PROFILE_FULL_MATCH, PROFILE_FULL_MATCH)
        assert score.confidence >= 0.999
        assert score.matched_loci == 20
        assert score.match_type == "FULL_MATCH"

    def test_partial_match_detected(self):
        """18/20 loci identiques → PARTIAL (confidence entre 0.85 et 0.999)"""
        score = calculate_lr_score(PROFILE_FULL_MATCH, PROFILE_PARTIAL_MATCH)
        assert 0.85 <= score.confidence < 0.999
        assert score.matched_loci == 18
        assert score.match_type == "PARTIAL"

    def test_no_match_below_threshold(self):
        """Profil non apparenté → confidence < 0.85 → pas de hit"""
        score = calculate_lr_score(PROFILE_FULL_MATCH, PROFILE_NO_MATCH)
        assert score.confidence < 0.85

    def test_incomplete_profile_handling(self):
        """Profil partiel (10/20 loci) doit matcher si loci disponibles concordent"""
        partial = {k: v for i, (k, v) in enumerate(PROFILE_FULL_MATCH.items()) if i < 10}
        score = calculate_lr_score(PROFILE_FULL_MATCH, partial)
        assert score.total_loci == 10
        assert score.matched_loci == 10

    def test_quality_score_threshold_rejection(self):
        """Profil quality_score < 0.60 doit être rejeté (BIO-FSC)"""
        matcher = DNAMatcher(db=Mock(), cache=Mock(), events=Mock())
        profile = STRProfile(
            specimen_number="TEST-001",
            index_type="BIO-FSC",
            loci=PROFILE_FULL_MATCH,
            quality_score=0.55,  # Sous le seuil
            lab_id="LDIS-PAP-001",
            collected_at="2026-06-09"
        )
        with pytest.raises(ValueError, match="quality_score"):
            matcher.validate_profile(profile)


class TestQualityScoring:

    def test_perfect_electropherogram(self):
        """Électrophorégramme parfait → quality_score ≈ 1.0"""
        epg = {locus: {"height": 5000} for locus in [
            "CSF1PO","D3S1358","D5S818","D7S820","D8S1179",
            "D13S317","D16S539","D18S51","D21S11","FGA",
            "TH01","TPOX","vWA","D1S1656","D2S441",
            "D2S1338","D10S1248","D12S391","D19S433","D22S1045"
        ]}
        from bio_adn.quality.scorer import calculate_quality_score
        score = calculate_quality_score(epg)
        assert score >= 0.95

    def test_degraded_sample(self):
        """Échantillon dégradé (8 loci > seuil) → quality_score bas mais accepté pour BIO-RNI"""
        epg = {f"locus_{i}": {"height": 200 if i < 8 else 50} for i in range(20)}
        from bio_adn.quality.scorer import calculate_quality_score
        score = calculate_quality_score(epg)
        assert score >= 0.50  # acceptable pour BIO-RNI
        assert score < 0.60   # pas acceptable pour BIO-FSC


# ──────────────────────────────────────────
# Tests d'intégration
# ──────────────────────────────────────────

class TestDNASubmissionWorkflow:

    @pytest.mark.asyncio
    async def test_submit_and_search_workflow(self, test_client, test_db):
        """Workflow complet : soumettre un profil BIO-CON et le retrouver via BIO-FSC"""
        # 1. Soumettre profil condamné
        response_submit = await test_client.post("/dna/profiles", json={
            "specimen_number": "TEST-CON-001",
            "index_type": "BIO-CON",
            "loci_data": PROFILE_FULL_MATCH,
            "quality_score": 0.98,
            "case_number": "DCPJ-2026-001",
            "collected_date": "2026-06-01",
            "correlation_id": "test-corr-001"
        }, headers={"Authorization": "Bearer test-lab-supervisor-token"})

        assert response_submit.status_code == 201
        sample_id = response_submit.json()["sample_id"]

        # 2. Rechercher avec profil scène de crime (identique)
        response_search = await test_client.post("/dna/search", json={
            "loci_data": PROFILE_FULL_MATCH,
            "index_type": "BIO-FSC",
            "case_number": "DCPJ-2026-002",
            "purpose": "criminal_investigation",
            "min_confidence": 0.85
        }, headers={"Authorization": "Bearer test-ndis-analyst-token"})

        assert response_search.status_code == 200
        result = response_search.json()
        assert result["total_hits"] >= 1
        assert result["hits"][0]["match_type"] == "FULL_MATCH"
        assert result["hits"][0]["confidence"] >= 0.999


class TestLAPIPerformance:

    @pytest.mark.asyncio
    async def test_lapi_response_under_200ms(self, test_client, seeded_db):
        """LAPI doit répondre en < 200ms (P99)"""
        import time

        times = []
        for i in range(50):  # 50 requêtes
            start = time.monotonic()
            response = await test_client.get(
                "/lapi/plate/OE-A-1234",
                headers={"Authorization": "Bearer lapi-service-token"}
            )
            elapsed_ms = (time.monotonic() - start) * 1000
            times.append(elapsed_ms)
            assert response.status_code == 200

        # P99 < 200ms
        times.sort()
        p99 = times[int(len(times) * 0.99)]
        assert p99 < 200, f"P99 LAPI = {p99:.1f}ms > 200ms SLA BREACH"


class TestAccessControl:

    @pytest.mark.asyncio
    async def test_lab_technician_cannot_access_other_lab(self, test_client):
        """Un technicien LDIS-CAP ne peut pas lire les profils de LDIS-PAP"""
        response = await test_client.get(
            "/dna/profiles?lab_id=LDIS-PAP-001",
            headers={"Authorization": "Bearer ldis-cap-technician-token"}
        )
        assert response.status_code == 403

    @pytest.mark.asyncio
    async def test_identity_link_requires_director_role(self, test_client):
        """L'accès aux bio_identity_links requiert le rôle bio.dcpj.director"""
        response = await test_client.get(
            "/dna/identity-links/some-sample-id",
            headers={"Authorization": "Bearer dcpj-investigator-token"}  # Pas director
        )
        assert response.status_code == 403

    @pytest.mark.asyncio
    async def test_expunge_requires_court_order(self, test_client):
        """L'expungement sans court_order_ref doit être rejeté"""
        response = await test_client.post("/dna/profiles/some-id/expunge", json={
            "reason": "acquittement",
            "ordered_by": "juge-001"
            # court_order_ref MANQUANT
        }, headers={"Authorization": "Bearer dcpj-director-token"})
        assert response.status_code == 422
```

---

## 3. TESTS E2E — Scénarios Critiques

### Scénario TC-001 : ADN scène de crime → identification suspect
```
1. Technicien LDIS-PAP-001 soumet BIO-FSC (scène de crime Cité Soleil)
2. Matching SDIS-OUEST : pas de hit
3. Upload NDIS-HT (hebdomadaire)
4. Matching NDIS : hit BIO-CON SDIS-NORD (condamné Cap-Haïtien)
5. Alerte Kafka snisid.bio.hits → DCPJ PAP notifiée
6. Agent DCPJ PAP vérifie hit + dossier condamné
7. RÉSULTAT ATTENDU : Alert niveau CRITIQUE envoyée < 24h après upload NDIS
```

### Scénario TC-002 : Alerte LAPI véhicule volé
```
1. PNH Delmas scanne plaque OE-A-5678 (LAPI MP-16)
2. Requête snisid.bio.lapi.query publiée
3. BIE-VEH Responder trouve hit (véhicule volé DIS-2026-000123)
4. Réponse LAPI envoyée < 200ms avec mco_contact PNH Pétion-Ville
5. Agent LAPI reçoit alerte : VÉHICULE VOLÉ / ALERT HAUTE
6. Contact obligatoire PNH Pétion-Ville avant intervention
7. RÉSULTAT ATTENDU : Réponse < 200ms, contact mco fourni
```

### Scénario TC-003 : Expungement légal BIO-ARR
```
1. Personne arrêtée → profil BIO-ARR créé (specimen_number ARR-2026-001)
2. 90 jours plus tard : classement sans suite
3. Greffier du tribunal envoie notification SNISID-BIO-ADN
4. Expungement déclenché (Directeur DCPJ + court_order_ref obligatoire)
5. Profil supprimé de la base + audit log signé ECDSA
6. RÉSULTAT ATTENDU : Profil non trouvable après expungement
```

---

## 4. COUVERTURE CIBLE

| Module | Cible couverture | Actuel |
|--------|-----------------|--------|
| Matching engine (Go) | 90% | 0% (à créer) |
| Quality scorer (Python) | 85% | 0% (à créer) |
| Index CRUD (Go) | 80% | 0% (à créer) |
| API endpoints (FastAPI) | 85% | 0% (à créer) |
| LAPI responder | 95% | 0% (à créer) |
| Sync LDIS→SDIS→NDIS | 80% | 0% (à créer) |
| Contrôle accès OPA | 95% | 0% (à créer) |
