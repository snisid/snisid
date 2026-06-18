"""
Tests for ETL Pipeline components
"""
import pytest
import json
import tempfile
from pathlib import Path


class TestTransformers:
    def test_generate_niu(self):
        from etl.transformers.id_generators import generate_niu
        niu = generate_niu("test-seed")
        assert len(niu) == 10
        assert niu.isalnum()

    def test_generate_niu_unique(self):
        from etl.transformers.id_generators import generate_niu
        niu1 = generate_niu("seed1")
        niu2 = generate_niu("seed2")
        assert niu1 != niu2

    def test_normalize_phone(self):
        from etl.transformers.id_generators import normalize_phone
        assert normalize_phone("12345678") == "+50912345678"
        assert normalize_phone("50912345678") == "+50912345678"
        assert normalize_phone("+50912345678") == "+50912345678"

    def test_normalize_plate(self):
        from etl.transformers.id_generators import normalize_plate
        assert normalize_plate("ab 123 cd") == "AB-123-CD"
        assert normalize_plate("AA-123-BB") == "AA-123-BB"
        assert normalize_plate("  abc  ") == "ABC"

    def test_normalize_name(self):
        from etl.transformers.id_generators import normalize_name
        assert normalize_name("JÉAN DÚPONT") == "JEAN DUPONT"
        assert normalize_name("  marie-franÇoise  ") == "MARIE-FRANCOISE"


class TestValidators:
    def test_nui_format_valid(self):
        from etl.validators.field_validators import validate
        errors = validate({"niu": "HT12345678"}, [{"field": "niu", "rule": "niu_format"}])
        assert len(errors) == 0

    def test_nui_format_invalid(self):
        from etl.validators.field_validators import validate
        errors = validate({"niu": "123"}, [{"field": "niu", "rule": "niu_format"}])
        assert len(errors) == 1

    def test_required_field_missing(self):
        from etl.validators.field_validators import validate
        errors = validate({"name": ""}, [{"field": "name", "rule": "required"}])
        assert len(errors) == 1

    def test_required_field_present(self):
        from etl.validators.field_validators import validate
        errors = validate({"name": "John"}, [{"field": "name", "rule": "required"}])
        assert len(errors) == 0

    def test_email_valid(self):
        from etl.validators.field_validators import validate
        errors = validate({"email": "test@sniside.ht"}, [{"field": "email", "rule": "email"}])
        assert len(errors) == 0

    def test_email_invalid(self):
        from etl.validators.field_validators import validate
        errors = validate({"email": "not-an-email"}, [{"field": "email", "rule": "email"}])
        assert len(errors) == 1

    def test_plate_valid(self):
        from etl.validators.field_validators import validate
        errors = validate({"plate": "AB-123-CD"}, [{"field": "plate", "rule": "plate"}])
        assert len(errors) == 0

    def test_not_future(self):
        from etl.validators.field_validators import validate
        errors = validate({"date": "2025-01-01"}, [{"field": "date", "rule": "not_future"}])
        assert len(errors) == 0  # past date

    def test_phone_valid(self):
        from etl.validators.field_validators import validate
        errors = validate({"phone": "+50912345678"}, [{"field": "phone", "rule": "phone"}])
        assert len(errors) == 0

    def test_phone_invalid(self):
        from etl.validators.field_validators import validate
        errors = validate({"phone": "12"}, [{"field": "phone", "rule": "phone"}])
        assert len(errors) == 1


class TestCsvAdapter:
    def test_csv_read(self):
        from etl.adapters import CsvAdapter
        import tempfile, csv
        with tempfile.NamedTemporaryFile(mode="w", suffix=".csv", delete=False, encoding="utf-8") as f:
            writer = csv.writer(f)
            writer.writerow(["niu", "name", "risk"])
            writer.writerow(["HT00000001", "Jean", "HIGH"])
            writer.writerow(["HT00000002", "Marie", "LOW"])
            tmp = f.name
        try:
            adapter = CsvAdapter(tmp)
            rows = adapter.read()
            assert len(rows) == 2
            assert rows[0]["niu"] == "HT00000001"
            assert rows[1]["name"] == "Marie"
        finally:
            import os
            os.unlink(tmp)


class TestPipeline:
    def test_basic_mapping(self):
        from etl.pipeline import ETLPipeline
        pipeline = ETLPipeline({
            "mapping": {
                "field_mappings": [
                    {"source": "NOM", "target": "last_name", "type": "string", "transformer": "upper"},
                    {"source": "AGE", "target": "age", "type": "int"},
                    {"source": "ACTIF", "target": "status", "type": "enum",
                     "mapping": {"OUI": "ACTIVE", "NON": "INACTIVE"}},
                ]
            },
            "batch_size": 1000,
        })
        record = pipeline._apply_mapping({"NOM": "dupont", "AGE": "30", "ACTIF": "OUI"}, pipeline.config["mapping"]["field_mappings"])
        assert record["last_name"] == "DUPONT"
        assert record["age"] == 30
        assert record["status"] == "ACTIVE"
