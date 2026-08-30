"""
SNI-SIDE: Tests d'intégration
Vérification des connexions entre les 15 bases, le Search Engine, et l'AI Fusion.
"""

import asyncio
import json
import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..'))

from ai.national_search_engine import (
    NationalSearchAPI, SearchType, Database, SearchResponse, SearchResultItem
)


class TestSuite:
    """Suite de tests d'intégration SNI-SIDE"""

    def __init__(self):
        self.search_engine = NationalSearchAPI()
        self.passed = 0
        self.failed = 0
        self.errors = []

    async def test_unified_search_by_name(self):
        """Test: Recherche unifiée par nom"""
        result = await self.search_engine.unified_search("Jean Dupont", "NAME")
        assert isinstance(result, SearchResponse), "Doit retourner SearchResponse"
        assert result.query == "Jean Dupont", f"Query mismatch: {result.query}"
        print("  ✓ Recherche unifiée par nom — OK")
        return True

    async def test_unified_search_by_niu(self):
        """Test: Recherche par NIU (format Haïtien 10 chiffres)"""
        result = await self.search_engine.unified_search("0000000001", "NIU")
        assert result.query == "0000000001"
        print("  ✓ Recherche par NIU — OK")
        return True

    async def test_unified_search_by_plate(self):
        """Test: Recherche par plaque d'immatriculation"""
        result = await self.search_engine.unified_search("AA-1234-BB", "PLATE")
        assert result.query == "AA-1234-BB"
        print("  ✓ Recherche par plaque — OK")
        return True

    async def test_unified_search_by_passport(self):
        """Test: Recherche par numéro de passeport"""
        result = await self.search_engine.unified_search("HT1234567", "PASSPORT")
        assert result.query == "HT1234567"
        print("  ✓ Recherche par passeport — OK")
        return True

    async def test_unified_search_by_phone(self):
        """Test: Recherche par téléphone (format Haïtien)"""
        result = await self.search_engine.unified_search("+50934123456", "PHONE")
        assert result.query == "+50934123456"
        print("  ✓ Recherche par téléphone — OK")
        return True

    async def test_unified_search_by_vin(self):
        """Test: Recherche par VIN (17 caractères)"""
        result = await self.search_engine.unified_search("1HGCM82633A004352", "VIN")
        assert result.query == "1HGCM82633A004352"
        print("  ✓ Recherche par VIN — OK")
        return True

    async def test_unified_search_by_wallet(self):
        """Test: Recherche par wallet crypto (Bitcoin)"""
        result = await self.search_engine.unified_search(
            "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "CRYPTO"
        )
        assert result.query == "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
        print("  ✓ Recherche par wallet crypto — OK")
        return True

    async def test_unified_search_by_ip(self):
        """Test: Recherche par adresse IP"""
        result = await self.search_engine.unified_search("192.168.1.1", "IP")
        assert result.query == "192.168.1.1"
        print("  ✓ Recherche par IP — OK")
        return True

    async def test_unified_search_by_domain(self):
        """Test: Recherche par nom de domaine"""
        result = await self.search_engine.unified_search("suspicious-domain.xyz", "DOMAIN")
        assert result.query == "suspicious-domain.xyz"
        print("  ✓ Recherche par domaine — OK")
        return True

    async def test_database_selection_for_niu(self):
        """Test: Sélection des bases appropriées pour NIU"""
        query_type = SearchType.NIU
        expected = {Database.NCID, Database.BIOMETRICS, Database.CODIS,
                    Database.MISSING, Database.BORDER, Database.WATCHLIST,
                    Database.FINANCIAL, Database.DOCUMENTS}
        from ai.national_search_engine import FederatedSearchOrchestrator, SearchQuery
        orchestrator = FederatedSearchOrchestrator()
        query = SearchQuery(raw_query="0000000001", query_type=query_type)
        selected = orchestrator._select_databases(query)
        assert set(selected) == expected, f"Selected: {selected}, Expected: {expected}"
        print(f"  ✓ Sélection bases NIU = {len(selected)} bases — OK")
        return True

    async def test_database_selection_for_plate(self):
        """Test: Sélection des bases pour plaque"""
        query_type = SearchType.PLATE
        expected = {Database.ALPR, Database.VEHICLE, Database.WATCHLIST}
        from ai.national_search_engine import FederatedSearchOrchestrator, SearchQuery
        orchestrator = FederatedSearchOrchestrator()
        query = SearchQuery(raw_query="AA-1234-BB", query_type=query_type)
        selected = orchestrator._select_databases(query)
        assert set(selected) == expected, f"Selected: {selected}, Expected: {expected}"
        print(f"  ✓ Sélection bases Plaque = {len(selected)} bases — OK")
        return True

    async def test_auto_detect_niu(self):
        """Test: Détection automatique du type NIU"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("0000000001")
        assert detected == SearchType.NIU, f"Détecté: {detected}, Attendu: NIU"
        print("  ✓ Détection automatique NIU — OK")

    async def test_auto_detect_plate(self):
        """Test: Détection automatique du type Plaque"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("AA-1234-BB")
        assert detected == SearchType.PLATE, f"Détecté: {detected}, Attendu: PLATE"
        print("  ✓ Détection automatique Plaque — OK")

    async def test_auto_detect_vin(self):
        """Test: Détection automatique du type VIN"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("1HGCM82633A004352")
        assert detected == SearchType.VIN, f"Détecté: {detected}, Attendu: VIN"
        print("  ✓ Détection automatique VIN — OK")

    async def test_auto_detect_phone(self):
        """Test: Détection automatique du type Téléphone"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("+50934123456")
        assert detected == SearchType.PHONE, f"Détecté: {detected}, Attendu: PHONE"
        print("  ✓ Détection automatique Téléphone — OK")

    async def test_auto_detect_email(self):
        """Test: Détection automatique du type Email"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("jean.dupont@email.ht")
        assert detected == SearchType.EMAIL, f"Détecté: {detected}, Attendu: EMAIL"
        print("  ✓ Détection automatique Email — OK")

    async def test_auto_detect_passport(self):
        """Test: Détection automatique du type Passeport"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("HT1234567")
        assert detected == SearchType.PASSPORT, f"Détecté: {detected}, Attendu: PASSPORT"
        print("  ✓ Détection automatique Passeport — OK")

    async def test_auto_detect_ip(self):
        """Test: Détection automatique du type IP"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("192.168.1.1")
        assert detected == SearchType.IP, f"Détecté: {detected}, Attendu: IP"
        print("  ✓ Détection automatique IP — OK")

    async def test_auto_detect_wallet(self):
        """Test: Détection automatique du type Wallet"""
        from ai.national_search_engine import QueryTypeDetector
        detected = QueryTypeDetector.detect("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
        assert detected == SearchType.WALLET, f"Détecté: {detected}, Attendu: WALLET"
        print("  ✓ Détection automatique Wallet — OK")

    async def test_search_engine_health(self):
        """Test: Vérification de santé du moteur de recherche"""
        health = self.search_engine.health()
        assert health["status"] == "healthy"
        assert health["connectors"] > 0
        assert health["databases_supported"] == len(Database)
        assert health["search_types"] == len(SearchType)
        print(f"  ✓ Moteur de recherche sain: {health['connectors']} connecteurs, "
              f"{health['databases_supported']} bases, {health['search_types']} types — OK")

    async def test_cross_reference(self):
        """Test: Référence croisée entre bases"""
        result = await self.search_engine.cross_reference("0000000001", "NIU", depth=2)
        assert "entity_id" in result
        assert "graph" in result
        assert "federated" in result
        print("  ✓ Référence croisée — OK")

    async def test_graph_search(self):
        """Test: Recherche graphique"""
        result = await self.search_engine.graph_search(niu="0000000001", depth=2)
        assert result["entity_id"] == "0000000001"
        print("  ✓ Recherche graphique — OK")

    async def test_search_performance(self):
        """Test: Performance de la recherche (doit être < 5s)"""
        import time
        start = time.monotonic()
        result = await self.search_engine.unified_search("Jean Dupont", "NAME", limit=10)
        duration = (time.monotonic() - start) * 1000
        assert duration < 5000, f"Trop lent: {duration:.0f}ms"
        print(f"  ✓ Performance recherche: {duration:.0f}ms (< 5000ms) — OK")

    # ============ EXÉCUTION ============
    async def run_all(self):
        print(f"\n{'='*60}")
        print("SNI-SIDE — INTEGRATION TESTS")
        print(f"{'='*60}")

        tests = [
            self.test_unified_search_by_name,
            self.test_unified_search_by_niu,
            self.test_unified_search_by_plate,
            self.test_unified_search_by_passport,
            self.test_unified_search_by_phone,
            self.test_unified_search_by_vin,
            self.test_unified_search_by_wallet,
            self.test_unified_search_by_ip,
            self.test_unified_search_by_domain,
            self.test_database_selection_for_niu,
            self.test_database_selection_for_plate,
            self.test_auto_detect_niu,
            self.test_auto_detect_plate,
            self.test_auto_detect_vin,
            self.test_auto_detect_phone,
            self.test_auto_detect_email,
            self.test_auto_detect_passport,
            self.test_auto_detect_ip,
            self.test_auto_detect_wallet,
            self.test_search_engine_health,
            self.test_cross_reference,
            self.test_graph_search,
            self.test_search_performance,
        ]

        for test in tests:
            try:
                await test()
                self.passed += 1
            except AssertionError as e:
                self.failed += 1
                self.errors.append(f"{test.__name__}: {e}")
                print(f"  ✗ {test.__name__} — FAIL: {e}")
            except Exception as e:
                self.failed += 1
                self.errors.append(f"{test.__name__}: {e}")
                print(f"  ✗ {test.__name__} — ERROR: {e}")

        print(f"\n{'='*60}")
        print(f"RÉSULTATS: {self.passed} passed, {self.failed} failed")
        print(f"{'='*60}")

        if self.errors:
            print("\nErreurs:")
            for err in self.errors:
                print(f"  • {err}")

        return self.failed == 0


if __name__ == "__main__":
    suite = TestSuite()
    success = asyncio.run(suite.run_all())
    sys.exit(0 if success else 1)
