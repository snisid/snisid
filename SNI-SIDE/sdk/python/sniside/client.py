import os
from typing import Optional
from .services.ncid import NCIDService
from .services.biometrics import BiometricsService
from .services.alpr import ALPRService
from .services.financial import FinancialService
from .services.cyber import CyberService
from .services.watchlist import WatchlistService
from .services.graphrag import GraphRAGService
from .services.alerts import AlertsService
from .services.search import SearchService
from .services.ai import AIService
from .http import HttpClient


class SNISIDClient:
    def __init__(
        self,
        api_key: Optional[str] = None,
        token: Optional[str] = None,
        base_url: str = "https://api.sniside.ht",
        cert: Optional[tuple[str, str]] = None,
        timeout: int = 30,
    ):
        self.http = HttpClient(
            api_key=api_key or os.getenv("SNISID_API_KEY"),
            token=token,
            base_url=base_url,
            cert=cert,
            timeout=timeout,
        )

        self.ncid = NCIDService(self.http)
        self.biometrics = BiometricsService(self.http)
        self.alpr = ALPRService(self.http)
        self.financial = FinancialService(self.http)
        self.cyber = CyberService(self.http)
        self.watchlist = WatchlistService(self.http)
        self.graphrag = GraphRAGService(self.http)
        self.alerts = AlertsService(self.http)
        self.search = SearchService(self.http)
        self.ai = AIService(self.http)

    def search(self, query: str, **kwargs) -> dict:
        return self.search.unified(query, **kwargs)


class AsyncSNISIDClient:
    def __init__(
        self,
        api_key: Optional[str] = None,
        token: Optional[str] = None,
        base_url: str = "https://api.sniside.ht",
        cert: Optional[tuple[str, str]] = None,
        timeout: int = 30,
    ):
        from .http import AsyncHttpClient

        self.http = AsyncHttpClient(
            api_key=api_key or os.getenv("SNISID_API_KEY"),
            token=token,
            base_url=base_url,
            cert=cert,
            timeout=timeout,
        )
        self.ncid = NCIDService(self.http)
        self.biometrics = BiometricsService(self.http)
        self.alpr = ALPRService(self.http)
        self.financial = FinancialService(self.http)
        self.cyber = CyberService(self.http)
        self.watchlist = WatchlistService(self.http)
        self.graphrag = GraphRAGService(self.http)
        self.alerts = AlertsService(self.http)
        self.search = SearchService(self.http)
        self.ai = AIService(self.http)

    async def search(self, query: str, **kwargs) -> dict:
        return await self.search.unified(query, **kwargs)

    async def __aenter__(self):
        return self

    async def __aexit__(self, *args):
        pass
