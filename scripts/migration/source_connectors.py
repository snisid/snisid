from __future__ import annotations

import abc
import csv
import io
import json
import time
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, Dict, Generator, List, Optional, Tuple

import httpx
import pandas as pd

from config import SourceConfig, SourceType


@dataclass
class SourceRecord:
    data: Dict[str, Any]
    source_type: str
    source_name: str
    raw_bytes: Optional[bytes] = None
    offset: int = 0


@dataclass
class SourceStats:
    total_records: int = 0
    total_bytes: int = 0
    read_time_seconds: float = 0.0
    parse_errors: int = 0
    skipped: int = 0
    columns: List[str] = field(default_factory=list)


class SourceConnector(abc.ABC):
    def __init__(self, config: SourceConfig):
        self.config = config
        self.stats = SourceStats()
        self._start_time: Optional[float] = None

    @abc.abstractmethod
    def read_records(self, offset: int = 0, limit: Optional[int] = None) -> Generator[SourceRecord, None, None]:
        pass

    @abc.abstractmethod
    def count(self) -> int:
        pass

    @abc.abstractmethod
    def validate_connection(self) -> Tuple[bool, str]:
        pass

    @abc.abstractmethod
    def get_schema(self) -> Dict[str, str]:
        pass

    def open(self) -> None:
        self._start_time = time.time()

    def close(self) -> None:
        if self._start_time is not None:
            self.stats.read_time_seconds = time.time() - self._start_time

    @staticmethod
    def create(config: SourceConfig) -> SourceConnector:
        mapping = {
            SourceType.csv: CSVConnector,
            SourceType.sqlserver: SQLServerConnector,
            SourceType.postgresql: PostgreSQLConnector,
            SourceType.rest_api: RESTAPIConnector,
        }
        cls = mapping.get(config.type)
        if cls is None:
            raise ValueError(f"Unsupported source type: {config.type}")
        return cls(config)


class CSVConnector(SourceConnector):
    def __init__(self, config: SourceConfig):
        super().__init__(config)
        self._file_path = Path(config.path) if config.path else None
        self._file_handle: Optional[io.TextIOWrapper] = None
        self._reader: Optional[csv.DictReader] = None
        self._total_lines: Optional[int] = None

    def validate_connection(self) -> Tuple[bool, str]:
        if self._file_path is None:
            return False, "No file path specified"
        if not self._file_path.exists():
            return False, f"File not found: {self._file_path}"
        if not self._file_path.is_file():
            return False, f"Path is not a file: {self._file_path}"
        if self._file_path.stat().st_size == 0:
            return False, f"File is empty: {self._file_path}"
        try:
            with open(self._file_path, "r", encoding=self.config.encoding) as f:
                reader = csv.DictReader(f, delimiter=self.config.delimiter)
                if not reader.fieldnames:
                    return False, "No column headers found in CSV"
            return True, "OK"
        except Exception as e:
            return False, str(e)

    def get_schema(self) -> Dict[str, str]:
        if self._reader is None:
            self._open_file()
        fieldnames = self._reader.fieldnames or []
        schema = {}
        for fn in fieldnames:
            col_sample = self._infer_column_type(fn)
            schema[fn] = col_sample
        self._close_file()
        return schema

    def _infer_column_type(self, col_name: str) -> str:
        try:
            self._open_file()
            for i, row in enumerate(self._reader):
                if i >= 100:
                    break
                val = row.get(col_name, "")
                if val and val.strip():
                    from dateutil.parser import parse as parse_date
                    try:
                        int(val)
                        return "int"
                    except ValueError:
                        pass
                    try:
                        float(val)
                        return "float"
                    except ValueError:
                        pass
                    try:
                        parse_date(val)
                        return "date"
                    except (ValueError, OverflowError):
                        pass
                    return "str"
            return "str"
        except Exception:
            return "str"
        finally:
            self._close_file()

    def count(self) -> int:
        if self._total_lines is not None:
            return self._total_lines
        try:
            self._open_file()
            count = 0
            for _ in self._reader:
                count += 1
            self._total_lines = count
            return count
        finally:
            self._close_file()

    def _open_file(self) -> None:
        if self._file_handle is None:
            self._file_handle = open(self._file_path, "r", encoding=self.config.encoding)
            self._reader = csv.DictReader(self._file_handle, delimiter=self.config.delimiter)
            self.stats.columns = self._reader.fieldnames or []

    def _close_file(self) -> None:
        if self._file_handle is not None:
            self._file_handle.close()
            self._file_handle = None
            self._reader = None

    def read_records(self, offset: int = 0, limit: Optional[int] = None) -> Generator[SourceRecord, None, None]:
        try:
            self._open_file()
            current = 0
            emitted = 0

            for row_num, row in enumerate(self._reader):
                if row_num < offset:
                    continue

                if limit is not None and emitted >= limit:
                    break

                cleaned = {}
                for k, v in row.items():
                    if k is None:
                        continue
                    cleaned[k.strip()] = v.strip() if v and isinstance(v, str) else v

                yield SourceRecord(
                    data=cleaned,
                    source_type="csv",
                    source_name=self._file_path.name if self._file_path else "csv",
                    offset=row_num,
                )
                emitted += 1
                current = row_num

            self.stats.total_records = current - offset + 1 if current >= offset else 0
        finally:
            self._close_file()


class SQLServerConnector(SourceConnector):
    def __init__(self, config: SourceConfig):
        super().__init__(config)
        self._connection = None
        self._cursor = None

    def validate_connection(self) -> Tuple[bool, str]:
        return True, "OK (pymssql not loaded; install pymssql for actual connection)"

    def get_schema(self) -> Dict[str, str]:
        return {"_placeholder": "str"}

    def count(self) -> int:
        return 0

    def read_records(self, offset: int = 0, limit: Optional[int] = None) -> Generator[SourceRecord, None, None]:
        if limit is None:
            limit = self.config.fetch_size
        query = self.config.query or f"SELECT * FROM {self.config.table}"
        paginated_query = f"{query} ORDER BY (SELECT NULL) OFFSET {offset} ROWS FETCH NEXT {limit} ROWS ONLY"
        try:
            import pymssql

            conn = pymssql.connect(
                server=self.config.host,
                port=self.config.port or 1433,
                user=self.config.username,
                password=self.config.password,
                database=self.config.database,
                timeout=self.config.connect_timeout,
            )
            cursor = conn.cursor(as_dict=True)
            cursor.execute(paginated_query)
            for i, row in enumerate(cursor):
                yield SourceRecord(
                    data=dict(row),
                    source_type="sqlserver",
                    source_name=f"{self.config.host}/{self.config.database}",
                    offset=offset + i,
                )
            cursor.close()
            conn.close()
        except ImportError:
            raise ImportError("pymssql is required to connect to SQL Server. Install with: pip install pymssql")
        except Exception as e:
            raise RuntimeError(f"SQL Server connection failed: {e}")

    def close(self) -> None:
        super().close()


class PostgreSQLConnector(SourceConnector):
    def __init__(self, config: SourceConfig):
        super().__init__(config)
        self._connection = None

    def validate_connection(self) -> Tuple[bool, str]:
        return True, "OK (psycopg2 not loaded; install psycopg2-binary for actual connection)"

    def get_schema(self) -> Dict[str, str]:
        return {"_placeholder": "str"}

    def count(self) -> int:
        return 0

    def read_records(self, offset: int = 0, limit: Optional[int] = None) -> Generator[SourceRecord, None, None]:
        if limit is None:
            limit = self.config.fetch_size
        query = self.config.query or f"SELECT * FROM {self.config.table}"
        paginated_query = f"{query} LIMIT {limit} OFFSET {offset}"
        try:
            import psycopg2
            import psycopg2.extras

            conn = psycopg2.connect(
                host=self.config.host,
                port=self.config.port or 5432,
                user=self.config.username,
                password=self.config.password,
                dbname=self.config.database,
                connect_timeout=self.config.connect_timeout,
            )
            cursor = conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor)
            cursor.execute(paginated_query)
            for i, row in enumerate(cursor):
                yield SourceRecord(
                    data=dict(row),
                    source_type="postgresql",
                    source_name=f"{self.config.host}/{self.config.database}",
                    offset=offset + i,
                )
            cursor.close()
            conn.close()
        except ImportError:
            raise ImportError("psycopg2 is required to connect to PostgreSQL. Install with: pip install psycopg2-binary")
        except Exception as e:
            raise RuntimeError(f"PostgreSQL connection failed: {e}")

    def close(self) -> None:
        super().close()


class RESTAPIConnector(SourceConnector):
    def __init__(self, config: SourceConfig):
        super().__init__(config)
        self._client: Optional[httpx.Client] = None
        self._total_count: Optional[int] = None

    def validate_connection(self) -> Tuple[bool, str]:
        try:
            if self._client is None:
                self._client = httpx.Client(timeout=self.config.connect_timeout)
            resp = self._client.head(self.config.api_url, headers=self._headers())
            if resp.status_code < 500:
                return True, "OK"
            return False, f"API returned status {resp.status_code}"
        except Exception as e:
            return False, str(e)

    def _headers(self) -> Dict[str, str]:
        headers = {"Accept": "application/json", "User-Agent": "SNISID-Migration/1.0"}
        if self.config.api_key:
            headers["Authorization"] = f"Bearer {self.config.api_key}"
        return headers

    def get_schema(self) -> Dict[str, str]:
        records = list(self.read_records(offset=0, limit=1))
        if records:
            return {k: type(v).__name__ for k, v in records[0].data.items()}
        return {}

    def count(self) -> int:
        if self._total_count is not None:
            return self._total_count
        if self._client is None:
            self._client = httpx.Client(timeout=self.config.connect_timeout)
        try:
            resp = self._client.get(
                self.config.api_url,
                headers=self._headers(),
                params={"count": True},
            )
            if resp.status_code == 200:
                data = resp.json()
                self._total_count = int(data.get("total", data.get("count", len(data))))
                return self._total_count
        except Exception:
            pass
        return 0

    def read_records(self, offset: int = 0, limit: Optional[int] = None) -> Generator[SourceRecord, None, None]:
        if self._client is None:
            self._client = httpx.Client(timeout=self.config.connect_timeout)

        page_size = limit or self.config.fetch_size
        page = (offset // page_size) + 1
        local_offset = offset % page_size
        emitted = 0

        while True:
            if limit is not None and emitted >= limit:
                break

            params = {
                "page": page,
                "per_page": page_size,
            }

            try:
                resp = self._client.get(
                    self.config.api_url,
                    headers=self._headers(),
                    params=params,
                )
                resp.raise_for_status()
                data = resp.json()
            except httpx.HTTPStatusError as e:
                raise RuntimeError(f"REST API error: {e.response.status_code} {e.response.text[:200]}")
            except httpx.RequestError as e:
                raise RuntimeError(f"REST API request failed: {e}")

            records = data if isinstance(data, list) else data.get("data", data.get("results", [data]))
            if not records:
                break

            for i, record in enumerate(records):
                if i < local_offset:
                    continue
                if limit is not None and emitted >= limit:
                    break
                yield SourceRecord(
                    data=record if isinstance(record, dict) else {"value": record},
                    source_type="rest_api",
                    source_name=self.config.api_url or "rest_api",
                    offset=offset + emitted,
                )
                emitted += 1

            if len(records) < page_size:
                break
            local_offset = 0
            page += 1

    def close(self) -> None:
        if self._client is not None:
            self._client.close()
            self._client = None
        super().close()
