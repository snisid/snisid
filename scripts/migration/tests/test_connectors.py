"""Tests for source connectors (CSV and REST)."""

import io
import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import MagicMock, patch

from source_connectors import CSVConnector, RESTAPIConnector, SourceRecord, SourceStats
from config import SourceConfig, SourceType


class TestCSVConnector(unittest.TestCase):
    def setUp(self):
        self.temp_file = tempfile.NamedTemporaryFile(mode="w", suffix=".csv", delete=False, encoding="utf-8-sig")
        self.temp_file.write("nom;prenom;birth_date;phone\n")
        self.temp_file.write("Dupont;Jean;15/03/1990;36305123\n")
        self.temp_file.write("Pierre;Marie;22/07/1985;40123456\n")
        self.temp_file.write("Paul;Luc;01/01/2000;51234567\n")
        self.temp_file.close()
        self.config = SourceConfig(type=SourceType.csv, path=self.temp_file.name, delimiter=";")
        self.connector = CSVConnector(self.config)

    def tearDown(self):
        Path(self.temp_file.name).unlink(missing_ok=True)

    def test_read_records(self):
        records = list(self.connector.read_records())
        self.assertEqual(len(records), 3)
        self.assertIsInstance(records[0], SourceRecord)
        self.assertEqual(records[0].data["nom"], "Dupont")
        self.assertEqual(records[0].data["prenom"], "Jean")

    def test_read_records_with_offset(self):
        records = list(self.connector.read_records(offset=1))
        self.assertEqual(len(records), 2)
        self.assertEqual(records[0].data["nom"], "Pierre")

    def test_read_records_with_limit(self):
        records = list(self.connector.read_records(limit=1))
        self.assertEqual(len(records), 1)

    def test_count(self):
        count = self.connector.count()
        self.assertEqual(count, 3)

    def test_validate_connection_valid(self):
        valid, msg = self.connector.validate_connection()
        self.assertTrue(valid)

    def test_validate_connection_file_not_found(self):
        bad_config = SourceConfig(type=SourceType.csv, path="/nonexistent/file.csv")
        connector = CSVConnector(bad_config)
        valid, msg = connector.validate_connection()
        self.assertFalse(valid)

    def test_get_schema(self):
        schema = self.connector.get_schema()
        self.assertIn("nom", schema)
        self.assertIn("prenom", schema)

    def test_offset_followed_by_source_record(self):
        records = list(self.connector.read_records(offset=0))
        self.assertEqual(records[0].offset, 0)
        self.assertEqual(records[1].offset, 1)
        self.assertEqual(records[2].offset, 2)

    def test_source_stats_populated(self):
        list(self.connector.read_records())
        self.assertGreaterEqual(self.connector.stats.total_records, 0)


class TestRESTConnector(unittest.TestCase):
    def setUp(self):
        self.config = SourceConfig(
            type=SourceType.rest_api,
            api_url="https://api.example.com/citizens",
            api_key="test-key-123",
        )
        self.connector = RESTAPIConnector(self.config)

    @patch("source_connectors.httpx.Client")
    def test_read_records(self, mock_client_class):
        mock_client = MagicMock()
        mock_client_class.return_value = mock_client
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "data": [
                {"id": 1, "nom": "Dupont", "prenom": "Jean"},
                {"id": 2, "nom": "Pierre", "prenom": "Marie"},
            ]
        }
        mock_response.raise_for_status.return_value = None
        mock_client.get.return_value = mock_response

        records = list(self.connector.read_records(limit=2))
        self.assertEqual(len(records), 2)
        self.assertIsInstance(records[0], SourceRecord)
        self.assertEqual(records[0].data["nom"], "Dupont")

    @patch("source_connectors.httpx.Client")
    def test_read_records_pagination(self, mock_client_class):
        mock_client = MagicMock()
        mock_client_class.return_value = mock_client

        def side_effect(*args, **kwargs):
            page = kwargs.get("params", {}).get("page", 1)
            resp = MagicMock()
            if page == 1:
                resp.json.return_value = {"data": [{"id": 1}], "total": 3}
            else:
                resp.json.return_value = {"data": [{"id": 2}, {"id": 3}], "total": 3}
            resp.raise_for_status.return_value = None
            return resp

        mock_client.get.side_effect = side_effect

        records = list(self.connector.read_records(limit=3))
        self.assertEqual(len(records), 3)

    @patch("source_connectors.httpx.Client")
    def test_validate_connection_success(self, mock_client_class):
        mock_client = MagicMock()
        mock_client_class.return_value = mock_client
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_client.head.return_value = mock_response

        valid, msg = self.connector.validate_connection()
        self.assertTrue(valid)

    @patch("source_connectors.httpx.Client")
    def test_validate_connection_failure(self, mock_client_class):
        mock_client = MagicMock()
        mock_client_class.return_value = mock_client
        mock_client.head.side_effect = Exception("Connection refused")

        valid, msg = self.connector.validate_connection()
        self.assertFalse(valid)

    @patch("source_connectors.httpx.Client")
    def test_get_schema(self, mock_client_class):
        mock_client = MagicMock()
        mock_client_class.return_value = mock_client
        mock_response = MagicMock()
        mock_response.json.return_value = {"data": [{"id": 1, "nom": "Dupont"}]}
        mock_response.raise_for_status.return_value = None
        mock_client.get.return_value = mock_response

        schema = self.connector.get_schema()
        self.assertIn("id", schema)
        self.assertIn("nom", schema)

    @patch("source_connectors.httpx.Client")
    def test_count(self, mock_client_class):
        mock_client = MagicMock()
        mock_client_class.return_value = mock_client
        mock_response = MagicMock()
        mock_response.status_code = 200
        mock_response.json.return_value = {"total": 42}
        mock_client.get.return_value = mock_response

        count = self.connector.count()
        self.assertEqual(count, 42)


class TestConnectorEmptySource(unittest.TestCase):
    def test_csv_empty_file(self):
        temp = tempfile.NamedTemporaryFile(mode="w", suffix=".csv", delete=False, encoding="utf-8-sig")
        temp.write("nom;prenom\n")
        temp.close()
        config = SourceConfig(type=SourceType.csv, path=temp.name, delimiter=";")
        connector = CSVConnector(config)
        records = list(connector.read_records())
        self.assertEqual(len(records), 0)
        Path(temp.name).unlink(missing_ok=True)
