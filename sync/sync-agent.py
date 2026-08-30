#!/usr/bin/env python3
"""SNISID Sync Agent - Offline-to-Online Synchronization"""

import argparse
import hashlib
import json
import logging
import os
import shutil
import sys
import time
from pathlib import Path
from datetime import datetime, timezone

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
log = logging.getLogger("sync-agent")


class SyncState:
    def __init__(self, state_path: Path):
        self.state_path = state_path
        self.data = self._load()

    def _load(self) -> dict:
        if self.state_path.exists():
            with open(self.state_path) as f:
                return json.load(f)
        return {"files": {}, "last_sync": None, "conflicts": []}

    def save(self):
        self.state_path.parent.mkdir(parents=True, exist_ok=True)
        with open(self.state_path, "w") as f:
            json.dump(self.data, f, indent=2)

    def get_file_hash(self, path: Path) -> str:
        h = hashlib.sha256()
        with open(path, "rb") as f:
            for chunk in iter(lambda: f.read(65536), b""):
                h.update(chunk)
        return h.hexdigest()

    def register_file(self, path: Path):
        rel = str(path.relative_to(path.anchor) if path.is_absolute() else path)
        self.data["files"][rel] = {
            "hash": self.get_file_hash(path),
            "mtime": datetime.fromtimestamp(path.stat().st_mtime, tz=timezone.utc).isoformat(),
            "synced": datetime.now(timezone.utc).isoformat(),
        }
        self.data["last_sync"] = datetime.now(timezone.utc).isoformat()
        self.save()

    def has_conflict(self, path: Path) -> bool:
        rel = str(path.relative_to(path.anchor) if path.is_absolute() else path)
        if rel not in self.data["files"]:
            return False
        prev = self.data["files"][rel]
        new_hash = self.get_file_hash(path)
        return prev["hash"] != new_hash

    def resolve_conflict(self, path: Path, strategy: str = "last-write-wins"):
        rel = str(path.relative_to(path.anchor) if path.is_absolute() else path)
        self.data["conflicts"].append({
            "file": rel,
            "resolved_at": datetime.now(timezone.utc).isoformat(),
            "strategy": strategy,
        })
        self.register_file(path)


class SyncAgent:
    def __init__(self, watch_dir: Path, archive_dir: Path, state_path: Path, etl_script: str = None):
        self.watch_dir = watch_dir
        self.archive_dir = archive_dir
        self.state = SyncState(state_path)
        self.etl_script = etl_script

    def etl_pipeline(self, path: Path):
        if self.etl_script:
            log.info("Running ETL on %s", path)
            os.system(f"{self.etl_script} {path}")
        else:
            log.info("ETL step skipped (no script configured) for %s", path)

    def process_file(self, path: Path):
        log.info("Processing %s", path)
        if self.state.has_conflict(path):
            log.warning("Conflict detected for %s — applying last-write-wins", path)
            self.state.resolve_conflict(path, "last-write-wins")

        self.etl_pipeline(path)
        self.state.register_file(path)

        dest = self.archive_dir / path.name
        dest.parent.mkdir(parents=True, exist_ok=True)
        shutil.move(str(path), str(dest))
        log.info("Archived to %s", dest)

    def watch(self, interval: float = 5.0):
        log.info("Watching %s every %.1f seconds", self.watch_dir, interval)
        while True:
            for item in sorted(self.watch_dir.iterdir()):
                if item.is_file() and not item.name.startswith("."):
                    self.process_file(item)
            time.sleep(interval)

    def sync_once(self):
        for item in sorted(self.watch_dir.iterdir()):
            if item.is_file() and not item.name.startswith("."):
                self.process_file(item)


def main():
    parser = argparse.ArgumentParser(description="SNISID Sync Agent")
    parser.add_argument("--watch", default="./sync/inbox", help="Directory to watch for new data files")
    parser.add_argument("--archive", default="./sync/archive", help="Directory to move processed files")
    parser.add_argument("--state", default="./sync/state.json", help="Sync state file")
    parser.add_argument("--etl", help="ETL script or command to run on each file")
    parser.add_argument("--once", action="store_true", help="Run once then exit")
    parser.add_argument("--interval", type=float, default=5.0, help="Watch interval in seconds")
    args = parser.parse_args()

    agent = SyncAgent(
        watch_dir=Path(args.watch),
        archive_dir=Path(args.archive),
        state_path=Path(args.state),
        etl_script=args.etl,
    )

    if args.once:
        agent.sync_once()
    else:
        agent.watch(interval=args.interval)


if __name__ == "__main__":
    main()
