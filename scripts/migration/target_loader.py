"""SNISID — Chargeur cible pour l'API SNISID.

Client HTTP avec batch processing, rate limiting, retry avec backoff,
et validation des réponses pour le chargement des données dans
le registre SNISID.
"""

import time
import logging
from typing import Dict, Any, List, Optional
from dataclasses import dataclass
from urllib.parse import urljoin

import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

logger = logging.getLogger(__name__)


@dataclass
class LoaderConfig:
    api_base_url: str = "http://localhost:8081"
    api_key: str = ""
    batch_size: int = 100
    max_retries: int = 3
    rate_limit_per_second: int = 50
    timeout_seconds: int = 30
    verify_ssl: bool = True


class TargetLoader:
    """Chargeur vers l'API REST SNISID avec résilience."""

    ENDPOINTS = {
        "identity": "/api/v1/identities",
        "biometric": "/api/v1/biometrics",
        "document": "/api/v1/documents",
        "enrollment": "/api/v1/enrollments",
    }

    def __init__(self, config: LoaderConfig):
        self.config = config
        self._session = self._create_session()
        self._last_request_time = 0.0
        self._stats = {"loaded": 0, "failed": 0, "batches": 0}

    def _create_session(self) -> requests.Session:
        """Crée une session HTTP avec retry et pooling."""
        session = requests.Session()

        retry_strategy = Retry(
            total=self.config.max_retries,
            backoff_factor=0.5,
            status_forcelist=[429, 500, 502, 503, 504],
            allowed_methods=["POST", "PUT", "GET"],
        )

        adapter = HTTPAdapter(
            pool_connections=10,
            pool_maxsize=20,
            max_retries=retry_strategy,
        )

        session.mount("http://", adapter)
        session.mount("https://", adapter)

        session.headers.update({
            "Content-Type": "application/json",
            "Accept": "application/json",
            "User-Agent": "SNISID-Migration/1.0",
        })

        if self.config.api_key:
            session.headers["X-API-Key"] = self.config.api_key

        return session

    def _rate_limit(self) -> None:
        """Respecte le rate limiting."""
        elapsed = time.time() - self._last_request_time
        min_interval = 1.0 / self.config.rate_limit_per_second
        if elapsed < min_interval:
            time.sleep(min_interval - elapsed)
        self._last_request_time = time.time()

    def load_batch(self, records: List[Dict[str, Any]]) -> int:
        """Charge un lot d'enregistrements vers l'API SNISID."""
        if not records:
            return 0

        self._rate_limit()
        endpoint = self._select_endpoint(records[0])
        url = urljoin(self.config.api_base_url, endpoint)

        try:
            response = self._session.post(
                url,
                json={"records": records, "batch": True},
                timeout=self.config.timeout_seconds,
                verify=self.config.verify_ssl,
            )
            response.raise_for_status()
            result = response.json()
            loaded = result.get("loaded", len(records))
            self._stats["loaded"] += loaded
            self._stats["batches"] += 1
            logger.debug("Lot chargé: %d/%d records", loaded, len(records))
            return loaded

        except requests.exceptions.RequestException as e:
            self._stats["failed"] += len(records)
            logger.error("Échec chargement lot: %s", str(e))
            raise

    def load_single(self, record: Dict[str, Any]) -> Optional[str]:
        """Charge un seul enregistrement et retourne son ID."""
        self._rate_limit()
        endpoint = self._select_endpoint(record)
        url = urljoin(self.config.api_base_url, endpoint)

        try:
            response = self._session.post(
                url,
                json=record,
                timeout=self.config.timeout_seconds,
                verify=self.config.verify_ssl,
            )
            response.raise_for_status()
            result = response.json()
            self._stats["loaded"] += 1
            return result.get("id")
        except requests.exceptions.RequestException as e:
            self._stats["failed"] += 1
            logger.error("Échec chargement unique: %s", str(e))
            return None

    def test_connection(self) -> bool:
        """Teste la connexion à l'API SNISID."""
        try:
            response = self._session.get(
                urljoin(self.config.api_base_url, "/health"),
                timeout=10,
            )
            return response.status_code == 200
        except requests.exceptions.RequestException:
            return False

    def _select_endpoint(self, record: Dict[str, Any]) -> str:
        """Sélectionne le bon endpoint selon le type d'enregistrement."""
        record_type = record.get("type", "identity")
        return self.ENDPOINTS.get(record_type, self.ENDPOINTS["identity"])

    @property
    def stats(self) -> Dict[str, Any]:
        return dict(self._stats)
