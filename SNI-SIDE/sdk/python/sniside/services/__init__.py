class BaseService:
    def __init__(self, http):
        self.http = http


class NCIDService(BaseService):
    def search_wanted(self, risk_level=None, status=None, agency=None, limit=50, offset=0):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/ncid/wanted-persons", params=params)

    def get_wanted(self, niu: str):
        return self.http.get(f"/intelligence/v1/ncid/wanted-persons/{niu}")

    def get_warrants(self, niu: str):
        return self.http.get(f"/intelligence/v1/ncid/wanted-persons/{niu}/warrants")

    def get_aliases(self, niu: str):
        return self.http.get(f"/intelligence/v1/ncid/wanted-persons/{niu}/aliases")

    def create_wanted(self, data: dict):
        return self.http.post("/intelligence/v1/ncid/wanted-persons", data)

    def create_warrant(self, niu: str, data: dict):
        return self.http.post(f"/intelligence/v1/ncid/wanted-persons/{niu}/warrants", data)

    def search_cases(self, case_type=None, status=None, agency=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/ncid/cases", params=params)

    def get_case(self, case_id: str):
        return self.http.get(f"/intelligence/v1/ncid/cases/{case_id}")

    def search_gangs(self, name=None, territory=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/ncid/gangs", params=params)

    def get_interpol_notices(self, notice_type=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/ncid/interpol-notices", params=params)


class BiometricsService(BaseService):
    def verify(self, face_image: str = None, fingerprint: str = None, iris: str = None, niu: str = None):
        data = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.post("/intelligence/v1/biometrics/verify", data)

    def identify(self, face_image: str = None, fingerprint: str = None):
        data = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.post("/intelligence/v1/biometrics/identify", data)

    def search_face(self, image: str, limit=20):
        return self.http.post("/intelligence/v1/biometrics/search/face", {"image": image, "limit": limit})

    def search_fingerprint(self, fingerprint: str, limit=20):
        return self.http.post("/intelligence/v1/biometrics/search/fingerprint", {"fingerprint": fingerprint, "limit": limit})

    def enroll(self, niu: str, biometric_type: str, data: str):
        return self.http.post("/intelligence/v1/biometrics/enroll", {"niu": niu, "biometric_type": biometric_type, "data": data})


class ALPRService(BaseService):
    def search(self, plate=None, camera_id=None, since=None, until=None, limit=100, offset=0):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/alpr/reads", params=params)

    def ingest(self, read: dict):
        return self.http.post("/intelligence/v1/alpr/ingest", read)

    def ingest_bulk(self, reads: list):
        return self.http.post("/intelligence/v1/alpr/ingest", {"reads": reads})

    def route_analysis(self, plate: str, since: str, until: str):
        return self.http.post("/intelligence/v1/alpr/route-analysis", {"plate": plate, "since": since, "until": until})

    def heatmap(self, start: str, end: str, location=None):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/alpr/heatmap", params=params)


class FinancialService(BaseService):
    def search_suspicious(self, amount_min=None, amount_max=None, currency=None, bank=None, status=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/financial/suspicious", params=params)

    def get_transaction(self, tx_id: str):
        return self.http.get(f"/intelligence/v1/financial/suspicious/{tx_id}")

    def search_pep(self, name=None, country=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/financial/pep", params=params)

    def report_suspicious(self, data: dict):
        return self.http.post("/intelligence/v1/financial/report", data)


class CyberService(BaseService):
    def search_iocs(self, ioc_type=None, value=None, threat_type=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/cyber/iocs", params=params)

    def get_ioc(self, ioc_id: str):
        return self.http.get(f"/intelligence/v1/cyber/iocs/{ioc_id}")

    def search_incidents(self, severity=None, status=None, limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/cyber/incidents", params=params)

    def get_incident(self, incident_id: str):
        return self.http.get(f"/intelligence/v1/cyber/incidents/{incident_id}")

    def submit_ioc(self, data: dict):
        return self.http.post("/intelligence/v1/cyber/iocs", data)

    def search_campaigns(self, limit=50):
        return self.http.get("/intelligence/v1/cyber/campaigns", params={"limit": limit})


class WatchlistService(BaseService):
    def search(self, category=None, risk_level=None, status="ACTIVE", limit=50):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/watchlist/entries", params=params)

    def get_entry(self, entry_id: str):
        return self.http.get(f"/intelligence/v1/watchlist/entries/{entry_id}")

    def add_entry(self, data: dict):
        return self.http.post("/intelligence/v1/watchlist/entries", data)

    def matches(self, entity_id: str):
        return self.http.get(f"/intelligence/v1/watchlist/matches/{entity_id}")


class GraphRAGService(BaseService):
    def generate_report(self, entity_id: str, report_type: str = "ENTITY_PROFILE",
                        entity_type: str = "Citizen", entity_label: str = None,
                        depth: int = 2, entity_id2: str = None):
        data = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.post("/intelligence/v1/ai/report", data)

    def get_report(self, entity_id: str, report_type: str = "ENTITY_PROFILE"):
        return self.http.get(f"/intelligence/v1/ai/report/{report_type}/{entity_id}")

    def cross_search(self, query: str, max_results: int = 20, include_graph: bool = True):
        return self.http.post("/intelligence/v1/ai/cross-search", {
            "query": query, "max_results": max_results, "include_graph": include_graph,
        })


class AlertsService(BaseService):
    def list(self, severity=None, alert_type=None, source=None, since=None, limit=50, offset=0):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/alerts", params=params)

    def get(self, alert_id: str):
        return self.http.get(f"/intelligence/v1/alerts/{alert_id}")

    def acknowledge(self, alert_id: str):
        return self.http.post(f"/intelligence/v1/alerts/{alert_id}/acknowledge")

    def resolve(self, alert_id: str, resolution: str = ""):
        return self.http.post(f"/intelligence/v1/alerts/{alert_id}/resolve", {"resolution": resolution})


class SearchService(BaseService):
    def unified(self, query: str, include_graph: bool = True, limit: int = 50):
        return self.http.get("/intelligence/v1/search/unified", params={
            "q": query, "include_graph": str(include_graph).lower(), "limit": limit,
        })

    def graph(self, entity_type: str, entity_id: str, depth: int = 2):
        return self.http.get(f"/intelligence/v1/search/graph/{entity_type}/{entity_id}", params={"depth": depth})

    def federated(self, query: str, sources: str = None, limit: int = 50):
        params = {"q": query, "limit": limit}
        if sources:
            params["sources"] = sources
        return self.http.get("/intelligence/v1/search/federated", params=params)


class AIService(BaseService):
    def score_entity(self, entity_id: str):
        return self.http.get(f"/intelligence/v1/ai/score/{entity_id}")

    def predictive_alerts(self, zone_id: str = None, days: int = 7):
        params = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.get("/intelligence/v1/ai/predictive", params=params)

    def deepfake_detect(self, image: str, video: str = None):
        data = {k: v for k, v in locals().items() if v is not None and k != 'self'}
        return self.http.post("/intelligence/v1/ai/deepfake-detect", data)

    def behavioral_analyze(self, user_id: str, days: int = 30):
        return self.http.post("/intelligence/v1/ai/behavioral", {"user_id": user_id, "days": days})
