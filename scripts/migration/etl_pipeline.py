"""SNISID — ETL Pipeline pour migration des données nationales.

Pipeline complet avec checkpoint/restart pour migration résiliente
des données d'identité, biométriques et documentaires des sources
existantes vers le registre SNISID.
"""

import logging
import time
from typing import Optional, Dict, Any
from dataclasses import dataclass, field

from checkpoint import CheckpointManager, Checkpoint
from source_connectors import (
    BaseSourceConnector,
    CSVConnector,
    SQLServerConnector,
    PostgreSQLConnector,
    RESTConnector,
)
from data_cleansing import DataCleanser
from matching_engine import MatchingEngine
from target_loader import TargetLoader

logger = logging.getLogger(__name__)


@dataclass
class PipelineStats:
    """Statistiques d'exécution du pipeline."""
    total_records: int = 0
    processed: int = 0
    cleaned: int = 0
    matched: int = 0
    loaded: int = 0
    errors: int = 0
    duplicates: int = 0
    start_time: float = 0.0
    end_time: float = 0.0
    batches: int = 0

    @property
    def elapsed(self) -> float:
        return (self.end_time or time.time()) - self.start_time

    @property
    def throughput(self) -> float:
        elapsed = self.elapsed
        return self.processed / elapsed if elapsed > 0 else 0


class ETLPipeline:
    """Pipeline ETL orchestré avec capacité de reprise."""

    def __init__(
        self,
        source: BaseSourceConnector,
        cleanser: DataCleanser,
        matcher: MatchingEngine,
        loader: TargetLoader,
        checkpoint_mgr: CheckpointManager,
        batch_size: int = 500,
    ):
        self.source = source
        self.cleanser = cleanser
        self.matcher = matcher
        self.loader = loader
        self.checkpoint = checkpoint_mgr
        self.batch_size = batch_size
        self.stats = PipelineStats()

    def run(self, max_records: Optional[int] = None) -> PipelineStats:
        """Exécute le pipeline ETL complet."""
        self.stats.start_time = time.time()
        last_checkpoint = self.checkpoint.load()
        processed_count = 0

        logger.info("Démarrage du pipeline ETL SNISID")
        logger.info("Source: %s", self.source.__class__.__name__)
        logger.info("Checkpoint trouvé: %s", last_checkpoint is not None)

        try:
            for batch_idx, batch in enumerate(self.source.stream_batches(
                self.batch_size,
                offset=last_checkpoint.offset if last_checkpoint else 0
            )):
                self.stats.total_records += len(batch)
                self.stats.batches += 1

                # Extract phase
                raw_records = [self.source.extract(record) for record in batch]

                # Cleanse phase
                cleaned = []
                for record in raw_records:
                    result = self.cleanser.cleanse(record)
                    if result.is_valid:
                        cleaned.append(result.data)
                    else:
                        logger.warning("Record rejeté: %s", result.errors)
                        self.stats.errors += 1
                self.stats.cleaned += len(cleaned)

                # Match phase (déduplication)
                matched = []
                for record in cleaned:
                    match_result = self.matcher.match(record)
                    if match_result.is_duplicate:
                        self.stats.duplicates += 1
                        if match_result.merged:
                            matched.append(match_result.merged)
                    else:
                        matched.append(record)
                self.stats.matched += len(matched)

                # Load phase
                loaded = self.loader.load_batch(matched)
                self.stats.loaded += loaded
                self.stats.processed += len(matched)

                # Checkpoint
                self.checkpoint.save(Checkpoint(
                    offset=last_checkpoint.offset + sum(
                        len(b) for b in [batch]
                    ) if last_checkpoint else sum(len(b) for b in [batch]),
                    batch=batch_idx + 1,
                    records_processed=self.stats.processed,
                ))

                processed_count += len(batch)
                if max_records and processed_count >= max_records:
                    logger.info("Limite max_records atteinte: %d", max_records)
                    break

                if (batch_idx + 1) % 10 == 0:
                    logger.info(
                        "Batch %d: %d records | Speed: %.0f rec/s",
                        batch_idx + 1,
                        self.stats.processed,
                        self.stats.throughput,
                    )

        except Exception as e:
            logger.error("Pipeline interrompu: %s", str(e))
            self.stats.errors += 1
            raise
        finally:
            self.stats.end_time = time.time()
            self.summarize()

        return self.stats

    def summarize(self) -> None:
        """Affiche le résumé du pipeline."""
        logger.info(
            "Pipeline terminé : %d rec en %.1fs (%.0f rec/s)",
            self.stats.processed,
            self.stats.elapsed,
            self.stats.throughput,
        )
        logger.info(
            "Total: %d | Nettoyés: %d | Matchés: %d | Chargés: %d | Erreurs: %d | Doublons: %d",
            self.stats.total_records,
            self.stats.cleaned,
            self.stats.matched,
            self.stats.loaded,
            self.stats.errors,
            self.stats.duplicates,
        )

    def validate(self) -> Dict[str, Any]:
        """Valide la configuration du pipeline."""
        checks = {
            "source_connected": self.source.test_connection(),
            "cleanser_ready": True,
            "matcher_ready": self.matcher.is_ready(),
            "loader_connected": self.loader.test_connection(),
            "checkpoint_writable": self.checkpoint.is_writable(),
        }
        return {
            "status": "healthy" if all(checks.values()) else "degraded",
            "checks": checks,
        }

    def reset(self) -> None:
        """Réinitialise le pipeline (efface le checkpoint)."""
        self.checkpoint.reset()
        self.stats = PipelineStats()
        logger.warning("Pipeline réinitialisé — tous les checkpoints ont été effacés")
