import time, sys
from datetime import datetime


class ProgressTracker:
    def __init__(self, total: int = 0, label: str = "rows"):
        self.total = total
        self.processed = 0
        self.label = label
        self.start_time = time.time()
        self.last_log = 0

    def update(self, n: int = 1):
        self.processed += n
        now = time.time()
        if now - self.last_log >= 2.0:
            self._log()
            self.last_log = now

    def _log(self):
        elapsed = time.time() - self.start_time
        rate = self.processed / max(elapsed, 0.1)
        pct = (self.processed / max(self.total, 1)) * 100 if self.total else 0
        msg = f"[{self.processed:,} {self.label}] ({rate:.0f}/s)"
        if self.total:
            eta = (self.total - self.processed) / max(rate, 0.1)
            msg += f" {pct:.1f}% ETA: {eta:.0f}s"
        print(f"\r{msg}", end="", flush=True)

    def done(self, success: bool = True):
        elapsed = time.time() - self.start_time
        print(f"\r{'✓' if success else '✗'} {self.processed:,} {self.label} in {elapsed:.1f}s{' ' * 20}")
