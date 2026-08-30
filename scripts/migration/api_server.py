"""SNISID — REST API for migration orchestration."""

import logging
import uuid
from pathlib import Path
from typing import Dict, Any, Optional
from dataclasses import dataclass, field

from fastapi import FastAPI, BackgroundTasks, HTTPException
from pydantic import BaseModel

from config import PipelineConfig, SourceConfig, SourceType, TargetConfig
from checkpoint import CheckpointManager, CheckpointData
from source_connectors import (
    SourceConnector,
    CSVConnector,
    SQLServerConnector,
    PostgreSQLConnector,
    RESTAPIConnector,
)
from data_cleansing import DataCleansingEngine
from matching_engine import MatchingEngine
from target_loader import TargetLoader, LoaderConfig
from etl_pipeline import ETLPipeline, PipelineStats

logger = logging.getLogger(__name__)

app = FastAPI(title="SNISID Migration API", version="1.0.0")


@dataclass
class PipelineRunState:
    pipeline_id: str
    status: str = "idle"
    stats: Optional[PipelineStats] = None
    error: Optional[str] = None
    config: Optional[PipelineConfig] = None


pipeline_state: Dict[str, PipelineRunState] = {}


class RunRequest(BaseModel):
    source_type: str = "csv"
    source_path: Optional[str] = None
    api_url: Optional[str] = None
    dry_run: bool = False
    max_records: Optional[int] = None
    batch_size: int = 500


class StatusResponse(BaseModel):
    pipeline_id: str
    status: str
    stats: Optional[Dict[str, Any]] = None
    error: Optional[str] = None


class ValidateResponse(BaseModel):
    status: str
    checks: Dict[str, bool]


def _create_source_connector_from_config(source_type: str, source_path: Optional[str] = None, api_url: Optional[str] = None) -> SourceConnector:
    cfg = SourceConfig(type=SourceType(source_type))
    if source_path:
        cfg.path = source_path
    if api_url:
        cfg.api_url = api_url
    return SourceConnector.create(cfg)


def _run_pipeline_async(pipeline_id: str, source_type: str, source_path: Optional[str], api_url: Optional[str], dry_run: bool, max_records: Optional[int], batch_size: int) -> None:
    state = pipeline_state.get(pipeline_id)
    if not state:
        return

    try:
        state.status = "running"
        source = _create_source_connector_from_config(source_type, source_path, api_url)
        cleanser = DataCleansingEngine()
        matcher = MatchingEngine()
        checkpoint = CheckpointManager("./checkpoints")
        loader_cfg = LoaderConfig(api_base_url="http://localhost:8081")
        loader = TargetLoader(loader_cfg)

        pipeline = ETLPipeline(
            source=source,
            cleanser=cleanser,
            matcher=matcher,
            loader=loader,
            checkpoint_mgr=checkpoint,
            batch_size=batch_size,
        )

        stats = pipeline.run(max_records=max_records)
        state.stats = stats
        state.status = "completed"
    except Exception as e:
        logger.exception("Pipeline %s failed", pipeline_id)
        state.status = "failed"
        state.error = str(e)


@app.post("/migration/run", response_model=StatusResponse)
async def run_migration(request: RunRequest, background_tasks: BackgroundTasks) -> Dict[str, Any]:
    pipeline_id = str(uuid.uuid4())
    state = PipelineRunState(pipeline_id=pipeline_id)
    pipeline_state[pipeline_id] = state

    background_tasks.add_task(
        _run_pipeline_async,
        pipeline_id,
        request.source_type,
        request.source_path,
        request.api_url,
        request.dry_run,
        request.max_records,
        request.batch_size,
    )

    return {"pipeline_id": pipeline_id, "status": "started"}


@app.get("/migration/status", response_model=StatusResponse)
async def get_status(pipeline_id: str) -> Dict[str, Any]:
    state = pipeline_state.get(pipeline_id)
    if not state:
        raise HTTPException(status_code=404, detail=f"Pipeline {pipeline_id} not found")

    stats_dict = None
    if state.stats:
        stats_dict = {
            "total_records": state.stats.total_records,
            "processed": state.stats.processed,
            "cleaned": state.stats.cleaned,
            "matched": state.stats.matched,
            "loaded": state.stats.loaded,
            "errors": state.stats.errors,
            "duplicates": state.stats.duplicates,
            "throughput": round(state.stats.throughput, 2),
            "elapsed": round(state.stats.elapsed, 2),
            "batches": state.stats.batches,
        }

    return {"pipeline_id": pipeline_id, "status": state.status, "stats": stats_dict, "error": state.error}


@app.post("/migration/restart", response_model=StatusResponse)
async def restart_migration(pipeline_id: str) -> Dict[str, Any]:
    state = pipeline_state.get(pipeline_id)
    if not state:
        raise HTTPException(status_code=404, detail=f"Pipeline {pipeline_id} not found")

    checkpoint = CheckpointManager("./checkpoints")
    existing = checkpoint.delete(pipeline_id)

    state.status = "reset"
    state.stats = None
    state.error = None

    return {"pipeline_id": pipeline_id, "status": "reset", "checkpoint_found": existing}


@app.get("/migration/validate", response_model=ValidateResponse)
async def validate_configuration() -> Dict[str, Any]:
    checkpoint = CheckpointManager("./checkpoints")
    checks = {
        "checkpoint_writable": checkpoint.is_writable() if hasattr(checkpoint, "is_writable") else True,
    }
    status = "healthy" if all(checks.values()) else "degraded"
    return {"status": status, "checks": checks}


@app.get("/migration/pipelines")
async def list_pipelines() -> Dict[str, list]:
    return {"pipeline_ids": list(pipeline_state.keys())}
