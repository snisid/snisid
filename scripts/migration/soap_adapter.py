"""SNISID — SOAP adapter for legacy system integration.

Zeep-based SOAP client that fetches citizen, document, and biometric
data from legacy systems and normalizes responses to SourceRecord format.
"""

import logging
import time
from datetime import date
from typing import Any, Dict, Generator, List, Optional, Tuple

try:
    from zeep import Client, Settings
    from zeep.transports import Transport
    from requests import Session
    from requests.adapters import HTTPAdapter
    from urllib3.util.retry import Retry
except ImportError:
    raise ImportError("zeep is required. Install with: pip install zeep")

from source_connectors import SourceRecord

logger = logging.getLogger(__name__)


class SOAPAdapterError(Exception):
    pass


class SOAPAdapter:
    def __init__(
        self,
        wsdl_url: str,
        service_name: Optional[str] = None,
        port_name: Optional[str] = None,
        verify_ssl: bool = True,
        max_retries: int = 3,
        timeout: int = 30,
    ):
        self.wsdl_url = wsdl_url
        self.service_name = service_name
        self.port_name = port_name
        self.verify_ssl = verify_ssl
        self.max_retries = max_retries
        self.timeout = timeout
        self._client: Optional[Client] = None

    def _create_session(self) -> Session:
        session = Session()
        session.verify = self.verify_ssl

        retry = Retry(
            total=self.max_retries,
            backoff_factor=0.5,
            status_forcelist=[429, 500, 502, 503, 504],
            allowed_methods=["POST"],
        )
        adapter = HTTPAdapter(max_retries=retry)
        session.mount("http://", adapter)
        session.mount("https://", adapter)

        return session

    def _get_client(self) -> Client:
        if self._client is not None:
            return self._client

        session = self._create_session()
        transport = Transport(session=session, timeout=self.timeout)
        settings = Settings(strict=False, xml_huge_tree=True)

        try:
            self._client = Client(
                self.wsdl_url,
                transport=transport,
                settings=settings,
                service_name=self.service_name,
                port_name=self.port_name,
            )
            logger.info("SOAP client initialized: %s", self.wsdl_url)
        except Exception as e:
            raise SOAPAdapterError(f"Failed to initialize SOAP client: {e}")

        return self._client

    def _call_with_retry(self, operation: str, *args, **kwargs) -> Any:
        last_error = None
        client = self._get_client()

        service = client.service
        if self.service_name:
            service = getattr(client.bind(self.service_name, self.port_name), self.service_name)

        for attempt in range(1, self.max_retries + 1):
            try:
                result = getattr(service, operation)(*args, **kwargs)
                return result
            except Exception as e:
                last_error = e
                logger.warning("SOAP %s attempt %d/%d failed: %s", operation, attempt, self.max_retries, e)
                if attempt < self.max_retries:
                    time.sleep(0.5 * (2 ** (attempt - 1)))

        raise SOAPAdapterError(f"SOAP {operation} failed after {self.max_retries} retries: {last_error}")

    def _to_source_record(self, raw: Dict[str, Any], record_type: str, offset: int = 0) -> SourceRecord:
        return SourceRecord(
            data=dict(raw),
            source_type="soap",
            source_name=self.wsdl_url,
            offset=offset,
        )

    def fetch_citizens(self, from_date: date, to_date: date) -> List[SourceRecord]:
        result = self._call_with_retry("GetCitizens", fromDate=from_date.isoformat(), toDate=to_date.isoformat())
        records = self._normalize_list_result(result, "citizen")
        return [self._to_source_record(r, "citizen", i) for i, r in enumerate(records)]

    def fetch_documents(self, national_id: str) -> List[SourceRecord]:
        result = self._call_with_retry("GetDocuments", nationalId=national_id)
        records = self._normalize_list_result(result, "document")
        return [self._to_source_record(r, "document", i) for i, r in enumerate(records)]

    def fetch_biometrics(self, national_id: str) -> List[SourceRecord]:
        result = self._call_with_retry("GetBiometrics", nationalId=national_id)
        records = self._normalize_list_result(result, "biometric")
        return [self._to_source_record(r, "biometric", i) for i, r in enumerate(records)]

    @staticmethod
    def _normalize_list_result(result: Any, default_key: str) -> List[Dict[str, Any]]:
        if result is None:
            return []

        if isinstance(result, list):
            return [r if isinstance(r, dict) else {"value": str(r)} for r in result]

        if hasattr(result, "__iter__"):
            try:
                return [dict(r) if not isinstance(r, dict) else r for r in result]
            except (TypeError, ValueError):
                pass

        if isinstance(result, dict):
            for key in ("records", "items", "data", "result", default_key + "s", default_key):
                val = result.get(key)
                if isinstance(val, list):
                    return [v if isinstance(v, dict) else {"value": str(v)} for v in val]
            return [result]

        return [{"value": str(result)}]

    def close(self) -> None:
        if self._client is not None:
            self._client.transport.session.close()
            self._client = None
