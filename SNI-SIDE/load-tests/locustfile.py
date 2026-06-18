"""
SNI-SIDE: Load Testing Framework

Usage:
    locust -f locustfile.py --host=https://api.sniside.ht --users=100 --spawn-rate=10
    python locustfile.py --headless --users=500 --spawn-rate=20 --run-time=10m
"""

import json, random, uuid
from locust import HttpUser, task, between, tag

QUERIES = [
    "HT12345678", "HT87654321", "HT55555555", "HT99999999",
    "AB-123-CD", "AA-000-BB", "ZZ-999-XX",
    "+50912345678", "+50987654321",
    "jean.dupont@sniside.ht",
    "185.220.101.42", "10.0.0.1",
    "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    "Jean Dupont", "Marie Jean", "Pierre Paul",
]

PLATES = [f"{random.choice('ABCDEFGH')}{random.choice('ABCDEFGH')}-{random.randint(100,999)}-{random.choice('ABCDEFGH')}{random.choice('ABCDEFGH')}" for _ in range(50)]
NIUS = [f"HT{''.join(random.choices('0123456789ABCDEFGHJKLMNPQRSTUVWXYZ', k=8))}" for _ in range(100)]


class SNISIDAnalystUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        self.client.headers.update({
            "X-API-Key": "sniside-api-key-change-in-production",
            "Content-Type": "application/json",
        })

    @tag("search")
    @task(30)
    def unified_search(self):
        query = random.choice(QUERIES)
        with self.client.get(
            f"/intelligence/v1/search/unified?q={query}&limit=20",
            name="/search/unified",
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                resp.success()
            elif resp.status_code == 429:
                resp.failure("Rate limited")
            else:
                resp.failure(f"Status: {resp.status_code}")

    @tag("ncid")
    @task(15)
    def get_wanted(self):
        niu = random.choice(NIUS)
        with self.client.get(
            f"/intelligence/v1/ncid/wanted-persons/{niu}",
            name="/ncid/wanted-persons/[niu]",
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                resp.success()
            elif resp.status_code == 404:
                resp.success()  # 404 is valid for missing persons
            else:
                resp.failure(f"Status: {resp.status_code}")

    @tag("ncid")
    @task(10)
    def list_wanted(self):
        risk = random.choice(["CRITICAL", "HIGH", "MEDIUM", "LOW"])
        with self.client.get(
            f"/intelligence/v1/ncid/wanted-persons?risk_level={risk}&limit=20",
            name="/ncid/wanted-persons",
        ) as resp:
            pass

    @tag("alpr")
    @task(20)
    def search_alpr(self):
        plate = random.choice(PLATES)
        with self.client.get(
            f"/intelligence/v1/alpr/reads?plate={plate}&limit=10",
            name="/alpr/reads",
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                resp.success()
            else:
                resp.failure(f"Status: {resp.status_code}")

    @tag("alpr")
    @task(10)
    def ingest_alpr(self):
        read = {
            "plate": random.choice(PLATES),
            "make": random.choice(["TOYOTA", "HONDA", "NISSAN", "FORD", "MITSUBISHI"]),
            "model": random.choice(["HILUX", "CIVIC", "SENTRA", "RANGER", "PAJERO"]),
            "year": random.randint(2000, 2025),
            "color": random.choice(["BLACK", "WHITE", "SILVER", "BLUE", "RED"]),
            "location": f"{18.4 + random.random():.4f},{-72.3 + random.random():.4f}",
            "timestamp": int(random.random() * 1000000000000) + 1700000000000,
            "camera_id": f"CAM-PAP-{random.randint(1,500):03d}",
        }
        with self.client.post(
            "/intelligence/v1/alpr/ingest",
            json=read,
            name="/alpr/ingest",
        ) as resp:
            pass

    @tag("graphrag")
    @task(5)
    def generate_intelligence(self):
        niu = random.choice(NIUS)
        report = {
            "entity_id": niu,
            "entity_type": "Citizen",
            "report_type": random.choice(["ENTITY_PROFILE", "LINK_ANALYSIS"]),
        }
        with self.client.post(
            "/intelligence/v1/ai/report",
            json=report,
            name="/ai/report",
            catch_response=True,
        ) as resp:
            if resp.status_code in (200, 201):
                resp.success()
            else:
                resp.failure(f"Status: {resp.status_code}")

    @tag("alerts")
    @task(10)
    def list_alerts(self):
        severity = random.choice(["CRITICAL", "HIGH", "MEDIUM", None])
        params = {"limit": 20}
        if severity:
            params["severity"] = severity
        with self.client.get(
            "/intelligence/v1/alerts",
            params=params,
            name="/alerts",
        ) as resp:
            pass

    @tag("financial")
    @task(5)
    def search_financial(self):
        with self.client.get(
            "/intelligence/v1/financial/suspicious?limit=20",
            name="/financial/suspicious",
        ) as resp:
            pass

    @tag("watchlist")
    @task(8)
    def search_watchlist(self):
        category = random.choice(["TERRORISM", "ORGANIZED_CRIME", "NARCOTICS", None])
        params = {"limit": 20}
        if category:
            params["category"] = category
        with self.client.get(
            "/intelligence/v1/watchlist/entries",
            params=params,
            name="/watchlist/entries",
        ) as resp:
            pass

    @tag("biometrics")
    @task(3)
    def biometric_identify(self):
        # Simulate 1:N identification (no real image data in load test)
        with self.client.post(
            "/intelligence/v1/biometrics/identify",
            json={"face_image": "base64_placeholder_for_load_test"},
            name="/biometrics/identify",
            catch_response=True,
        ) as resp:
            if resp.status_code in (200, 400, 422):
                resp.success()
            else:
                resp.failure(f"Status: {resp.status_code}")

    @tag("cyber")
    @task(5)
    def search_cyber(self):
        with self.client.get(
            "/intelligence/v1/cyber/iocs?limit=20",
            name="/cyber/iocs",
        ) as resp:
            pass


class SNISIDBatchIngestUser(HttpUser):
    """Simulates ALPR camera batch ingestion — high throughput."""
    wait_time = between(0.1, 0.5)

    @tag("alpr_bulk")
    @task(1)
    def bulk_ingest(self):
        reads = []
        for _ in range(50):
            reads.append({
                "plate": random.choice(PLATES),
                "make": "TOYOTA",
                "model": "HILUX",
                "year": 2023,
                "color": "BLACK",
                "location": f"{18.4 + random.random():.4f},{-72.3 + random.random():.4f}",
                "timestamp": int(random.random() * 1000000000000) + 1700000000000,
                "camera_id": f"CAM-PAP-{random.randint(1,500):03d}",
            })
        with self.client.post(
            "/intelligence/v1/alpr/ingest",
            json={"reads": reads},
            name="/alpr/ingest [bulk 50]",
        ) as resp:
            pass


class SNISIDSearchIntensiveUser(HttpUser):
    """Simulates heavy search usage — analysts running queries."""
    wait_time = between(0.5, 2)

    @tag("search_heavy")
    @task(1)
    def heavy_search(self):
        for _ in range(5):
            query = random.choice(QUERIES)
            self.client.get(
                f"/intelligence/v1/search/unified?q={query}&limit=50",
                name="/search/unified [heavy]",
            )
