"""Tests for CheckpointManager."""

import json
import tempfile
import unittest
from pathlib import Path

from checkpoint import CheckpointManager, CheckpointData


class TestCheckpointManager(unittest.TestCase):
    def setUp(self):
        self.temp_dir = tempfile.TemporaryDirectory()
        self.checkpoint_dir = self.temp_dir.name
        self.manager = CheckpointManager(self.checkpoint_dir)
        self.pipeline_id = "test-pipeline-1"
        self.checkpoint = CheckpointData(
            pipeline_id=self.pipeline_id,
            source_name="test-csv",
            batch_index=5,
            records_processed=1000,
            records_succeeded=995,
            records_failed=5,
            last_processed_offset=5000,
            failed_records=[{"row": 42, "error": "Invalid date"}],
            partial_batch=False,
            metadata={"source_type": "csv"},
        )

    def tearDown(self):
        self.temp_dir.cleanup()

    def test_save_and_load_checkpoint(self):
        self.manager.save(self.checkpoint)
        loaded = self.manager.load(self.pipeline_id)
        self.assertIsNotNone(loaded)
        self.assertEqual(loaded.pipeline_id, self.pipeline_id)
        self.assertEqual(loaded.batch_index, 5)
        self.assertEqual(loaded.records_processed, 1000)
        self.assertEqual(loaded.last_processed_offset, 5000)

    def test_load_nonexistent_checkpoint(self):
        loaded = self.manager.load("nonexistent-pipeline")
        self.assertIsNone(loaded)

    def test_checkpoint_reset_via_delete(self):
        self.manager.save(self.checkpoint)
        loaded = self.manager.load(self.pipeline_id)
        self.assertIsNotNone(loaded)

        deleted = self.manager.delete(self.pipeline_id)
        self.assertTrue(deleted)

        loaded_after = self.manager.load(self.pipeline_id)
        self.assertIsNone(loaded_after)

    def test_delete_nonexistent(self):
        deleted = self.manager.delete("nonexistent")
        self.assertFalse(deleted)

    def test_resume_from_offset(self):
        self.manager.save(self.checkpoint)
        loaded = self.manager.load(self.pipeline_id)
        self.assertEqual(loaded.last_processed_offset, 5000)

        new_offset = loaded.last_processed_offset + 100
        updated = self.manager.update(
            pipeline_id=self.pipeline_id,
            batch_index=6,
            records_processed=1100,
            records_succeeded=1095,
            records_failed=5,
            last_processed_offset=new_offset,
        )
        self.assertEqual(updated.last_processed_offset, new_offset)
        self.assertEqual(updated.batch_index, 6)

    def test_checkpoint_file_created(self):
        self.manager.save(self.checkpoint)
        expected_path = Path(self.checkpoint_dir) / f"checkpoint_{self.pipeline_id}.json"
        self.assertTrue(expected_path.exists())

        with open(expected_path, "r", encoding="utf-8") as f:
            data = json.load(f)
        self.assertEqual(data["pipeline_id"], self.pipeline_id)
        self.assertEqual(data["records_processed"], 1000)

    def test_load_latest_checkpoint(self):
        self.manager.save(self.checkpoint)
        loaded = self.manager.load_latest("test-csv")
        self.assertIsNotNone(loaded)
        self.assertEqual(loaded.source_name, "test-csv")

    def test_load_latest_nonexistent(self):
        loaded = self.manager.load_latest("nonexistent-source")
        self.assertIsNone(loaded)

    def test_list_checkpoints(self):
        self.manager.save(self.checkpoint)
        cp2 = CheckpointData(
            pipeline_id="test-pipeline-2",
            source_name="test-csv",
            batch_index=1,
            records_processed=100,
            records_succeeded=100,
            records_failed=0,
            last_processed_offset=500,
        )
        self.manager.save(cp2)

        checkpoints = self.manager.list_checkpoints()
        self.assertIn(self.pipeline_id, checkpoints)
        self.assertIn("test-pipeline-2", checkpoints)

    def test_get_progress(self):
        self.manager.save(self.checkpoint)
        progress = self.manager.get_progress(self.pipeline_id)
        self.assertEqual(progress["status"], "in_progress")
        self.assertEqual(progress["batch_index"], 5)

    def test_get_progress_not_found(self):
        progress = self.manager.get_progress("nonexistent")
        self.assertEqual(progress["status"], "not_found")

    def test_cleanup_old_checkpoints(self):
        for i in range(10):
            cp = CheckpointData(
                pipeline_id=f"pipeline-{i}",
                source_name="test",
                batch_index=i,
                records_processed=i * 100,
                records_succeeded=i * 100,
                records_failed=0,
                last_processed_offset=i * 500,
            )
            self.manager.save(cp)

        removed = self.manager.cleanup_old_checkpoints(keep_last=3)
        self.assertGreater(removed, 0)

        remaining = self.manager.list_checkpoints()
        self.assertLessEqual(len(remaining), 3)

    def test_close_clears_session(self):
        self.manager.save(self.checkpoint)
        self.manager.close()
        progress = self.manager.get_progress(self.pipeline_id)
        self.assertEqual(progress["pipeline_id"], self.pipeline_id)

    def test_checkpoint_has_timestamp(self):
        self.manager.save(self.checkpoint)
        loaded = self.manager.load(self.pipeline_id)
        self.assertIsNotNone(loaded.timestamp)
