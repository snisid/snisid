"""Integration tests for ETLPipeline with mocked components."""

import unittest
from unittest.mock import MagicMock, patch
from typing import Generator, Optional

from etl_pipeline import ETLPipeline, PipelineStats
from source_connectors import SourceConnector, SourceRecord, SourceConfig, SourceType
from data_cleansing import DataCleansingEngine
from matching_engine import MatchingEngine
from target_loader import TargetLoader, LoaderConfig
from checkpoint import CheckpointManager, CheckpointData


class MockSourceConnector:
    def __init__(self, records: list, fail_on_batch: Optional[int] = None):
        self.records = records
        self.fail_on_batch = fail_on_batch
        self.batch_calls = 0
        self.stats = MagicMock()
        self.stats.total_records = 0

    def stream_batches(self, batch_size: int, offset: int = 0) -> Generator:
        start = offset
        while start < len(self.records):
            if self.fail_on_batch is not None and self.batch_calls >= self.fail_on_batch:
                raise RuntimeError("Simulated source failure")
            batch = self.records[start:start + batch_size]
            self.batch_calls += 1
            yield batch
            start += batch_size

    def extract(self, record: SourceRecord) -> dict:
        return record.data

    def test_connection(self) -> bool:
        return True

    def close(self):
        pass


class TestETLPipeline(unittest.TestCase):
    def setUp(self):
        self.records = [
            SourceRecord(data={"nom": "Dupont", "prenom": "Jean", "id": i}, source_type="test", source_name="test", offset=i)
            for i in range(10)
        ]
        self.mock_source = MockSourceConnector(self.records)
        self.cleanser = DataCleansingEngine()
        self.matcher = MatchingEngine()
        self.loader = MagicMock(spec=TargetLoader)
        self.loader.load_batch.return_value = 5
        self.loader.test_connection.return_value = True

        self.checkpoint = MagicMock(spec=CheckpointManager)
        self.checkpoint.load.return_value = None
        self.checkpoint.is_writable.return_value = True

        self.pipeline = ETLPipeline(
            source=self.mock_source,
            cleanser=self.cleanser,
            matcher=self.matcher,
            loader=self.loader,
            checkpoint_mgr=self.checkpoint,
            batch_size=5,
        )

    def test_full_pipeline_run(self):
        stats = self.pipeline.run()
        self.assertIsInstance(stats, PipelineStats)
        self.assertEqual(stats.total_records, 10)
        self.assertGreater(stats.processed, 0)
        self.assertGreater(stats.batches, 0)
        self.assertGreater(stats.elapsed, 0)
        self.loader.load_batch.assert_called()

    def test_dry_run_mode_no_loader_calls(self):
        self.loader.load_batch.return_value = 0
        stats = self.pipeline.run()
        self.assertGreater(stats.processed, 0)

    def test_pipeline_stats_tracking(self):
        stats = self.pipeline.run()
        self.assertGreaterEqual(stats.total_records, stats.processed)
        self.assertGreaterEqual(stats.processed, stats.loaded)
        self.assertIsNotNone(stats.start_time)
        self.assertIsNotNone(stats.end_time)

    def test_pipeline_throughput_calculation(self):
        stats = self.pipeline.run()
        self.assertGreater(stats.throughput, 0)
        self.assertAlmostEqual(stats.throughput, stats.processed / stats.elapsed, places=1)

    def test_max_records_limit(self):
        stats = self.pipeline.run(max_records=3)
        self.assertGreaterEqual(stats.total_records, 3)

    def test_pipeline_validate(self):
        result = self.pipeline.validate()
        self.assertIn("status", result)
        self.assertIn("checks", result)

    def test_pipeline_reset_clears_checkpoint(self):
        self.pipeline.reset()
        self.checkpoint.reset.assert_not_called()

    def test_error_handling_in_source(self):
        failing_source = MockSourceConnector(self.records, fail_on_batch=0)
        pipeline = ETLPipeline(
            source=failing_source,
            cleanser=self.cleanser,
            matcher=self.matcher,
            loader=self.loader,
            checkpoint_mgr=self.checkpoint,
            batch_size=5,
        )
        with self.assertRaises(RuntimeError):
            pipeline.run()

    @patch("etl_pipeline.logger")
    def test_pipeline_logging(self, mock_logger):
        self.pipeline.run()
        mock_logger.info.assert_any_call("Démarrage du pipeline ETL SNISID")

    def test_batch_processing(self):
        pipeline = ETLPipeline(
            source=self.mock_source,
            cleanser=self.cleanser,
            matcher=self.matcher,
            loader=self.loader,
            checkpoint_mgr=self.checkpoint,
            batch_size=3,
        )
        stats = pipeline.run()
        self.assertGreaterEqual(stats.batches, 3)

    def test_empty_source(self):
        empty_source = MockSourceConnector([])
        pipeline = ETLPipeline(
            source=empty_source,
            cleanser=self.cleanser,
            matcher=self.matcher,
            loader=self.loader,
            checkpoint_mgr=self.checkpoint,
        )
        stats = pipeline.run()
        self.assertEqual(stats.total_records, 0)
        self.assertEqual(stats.processed, 0)
