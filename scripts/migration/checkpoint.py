from __future__ import annotations

import json
import os
import threading
from datetime import datetime, timezone
from pathlib import Path
from typing import Any, Dict, List, Optional


class CheckpointData:
    def __init__(
        self,
        pipeline_id: str,
        source_name: str,
        batch_index: int,
        records_processed: int,
        records_succeeded: int,
        records_failed: int,
        last_processed_offset: int,
        failed_records: Optional[List[Dict[str, Any]]] = None,
        partial_batch: bool = False,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        self.pipeline_id = pipeline_id
        self.source_name = source_name
        self.batch_index = batch_index
        self.records_processed = records_processed
        self.records_succeeded = records_succeeded
        self.records_failed = records_failed
        self.last_processed_offset = last_processed_offset
        self.failed_records = failed_records or []
        self.partial_batch = partial_batch
        self.timestamp = datetime.now(timezone.utc).isoformat()
        self.metadata = metadata or {}

    def to_dict(self) -> Dict[str, Any]:
        return {
            "pipeline_id": self.pipeline_id,
            "source_name": self.source_name,
            "batch_index": self.batch_index,
            "records_processed": self.records_processed,
            "records_succeeded": self.records_succeeded,
            "records_failed": self.records_failed,
            "last_processed_offset": self.last_processed_offset,
            "failed_records": self.failed_records,
            "partial_batch": self.partial_batch,
            "timestamp": self.timestamp,
            "metadata": self.metadata,
        }

    @classmethod
    def from_dict(cls, data: Dict[str, Any]) -> CheckpointData:
        return cls(
            pipeline_id=data["pipeline_id"],
            source_name=data["source_name"],
            batch_index=data["batch_index"],
            records_processed=data["records_processed"],
            records_succeeded=data["records_succeeded"],
            records_failed=data["records_failed"],
            last_processed_offset=data["last_processed_offset"],
            failed_records=data.get("failed_records", []),
            partial_batch=data.get("partial_batch", False),
            metadata=data.get("metadata", {}),
        )


class CheckpointManager:
    def __init__(self, checkpoint_dir: str):
        self.checkpoint_dir = Path(checkpoint_dir)
        self.checkpoint_dir.mkdir(parents=True, exist_ok=True)
        self._lock = threading.Lock()
        self._session_checkpoints: Dict[str, CheckpointData] = {}

    def _checkpoint_path(self, pipeline_id: str) -> Path:
        return self.checkpoint_dir / f"checkpoint_{pipeline_id}.json"

    def _latest_path(self, source_name: str) -> Path:
        sanitized = source_name.replace(" ", "_").replace("/", "_")
        return self.checkpoint_dir / f"latest_{sanitized}.json"

    def save(self, checkpoint: CheckpointData) -> None:
        with self._lock:
            path = self._checkpoint_path(checkpoint.pipeline_id)
            tmp_path = path.with_suffix(".tmp")
            with open(tmp_path, "w", encoding="utf-8") as f:
                json.dump(checkpoint.to_dict(), f, indent=2, ensure_ascii=False)
            tmp_path.replace(path)

            latest_path = self._latest_path(checkpoint.source_name)
            with open(latest_path, "w", encoding="utf-8") as f:
                json.dump(checkpoint.to_dict(), f, indent=2, ensure_ascii=False)

            self._session_checkpoints[checkpoint.pipeline_id] = checkpoint

    def load(self, pipeline_id: str) -> Optional[CheckpointData]:
        path = self._checkpoint_path(pipeline_id)
        if not path.exists():
            return None
        try:
            with open(path, "r", encoding="utf-8") as f:
                data = json.load(f)
            cp = CheckpointData.from_dict(data)
            self._session_checkpoints[pipeline_id] = cp
            return cp
        except (json.JSONDecodeError, KeyError, ValueError):
            return None

    def load_latest(self, source_name: str) -> Optional[CheckpointData]:
        latest_path = self._latest_path(source_name)
        if not latest_path.exists():
            return None
        try:
            with open(latest_path, "r", encoding="utf-8") as f:
                data = json.load(f)
            cp = CheckpointData.from_dict(data)
            self._session_checkpoints[cp.pipeline_id] = cp
            return cp
        except (json.JSONDecodeError, KeyError, ValueError):
            return None

    def update(
        self,
        pipeline_id: str,
        batch_index: int,
        records_processed: int,
        records_succeeded: int,
        records_failed: int,
        last_processed_offset: int,
        failed_records: Optional[List[Dict[str, Any]]] = None,
        partial_batch: bool = False,
    ) -> CheckpointData:
        existing = self._session_checkpoints.get(pipeline_id)
        if existing:
            existing.batch_index = batch_index
            existing.records_processed = records_processed
            existing.records_succeeded = records_succeeded
            existing.records_failed = records_failed
            existing.last_processed_offset = last_processed_offset
            existing.failed_records = (existing.failed_records or []) + (failed_records or [])
            existing.partial_batch = partial_batch
            existing.timestamp = datetime.now(timezone.utc).isoformat()
            self.save(existing)
            return existing

        cp = CheckpointData(
            pipeline_id=pipeline_id,
            source_name="unknown",
            batch_index=batch_index,
            records_processed=records_processed,
            records_succeeded=records_succeeded,
            records_failed=records_failed,
            last_processed_offset=last_processed_offset,
            failed_records=failed_records,
            partial_batch=partial_batch,
        )
        self.save(cp)
        return cp

    def delete(self, pipeline_id: str) -> bool:
        path = self._checkpoint_path(pipeline_id)
        if path.exists():
            path.unlink()
            self._session_checkpoints.pop(pipeline_id, None)
            return True
        return False

    def list_checkpoints(self) -> List[str]:
        return [p.stem.replace("checkpoint_", "") for p in self.checkpoint_dir.glob("checkpoint_*.json")]

    def get_progress(self, pipeline_id: str) -> Dict[str, Any]:
        cp = self.load(pipeline_id)
        if cp is None:
            return {"pipeline_id": pipeline_id, "status": "not_found"}
        return {
            "pipeline_id": cp.pipeline_id,
            "source_name": cp.source_name,
            "batch_index": cp.batch_index,
            "records_processed": cp.records_processed,
            "records_succeeded": cp.records_succeeded,
            "records_failed": cp.records_failed,
            "last_processed_offset": cp.last_processed_offset,
            "timestamp": cp.timestamp,
            "status": "in_progress" if not cp.partial_batch else "partial_failure",
        }

    def cleanup_old_checkpoints(self, keep_last: int = 5) -> int:
        checkpoints = sorted(
            self.checkpoint_dir.glob("checkpoint_*.json"),
            key=os.path.getmtime,
        )
        to_remove = len(checkpoints) - keep_last
        removed = 0
        if to_remove > 0:
            for cp in checkpoints[:to_remove]:
                cp.unlink()
                removed += 1
        return removed

    def close(self) -> None:
        self._session_checkpoints.clear()
