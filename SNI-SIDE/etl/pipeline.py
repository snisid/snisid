import asyncio, csv, json, logging, os, uuid, re, io, hashlib
from datetime import datetime
from pathlib import Path
from typing import Optional
from concurrent.futures import ThreadPoolExecutor

logger = logging.getLogger("sniside.etl.pipeline")


class ETLPipeline:
    def __init__(self, config: dict):
        self.config = config
        self.stats = {
            "rows_processed": 0, "rows_inserted": 0, "rows_errors": 0,
            "batches_completed": 0, "batch_errors": 0, "errors": [], "warnings": [],
            "checkpoint": None, "success": False,
        }
        self.adapter = None
        self.writer = None

    def generate_report(self) -> dict:
        sources = self.config.get("sources", [])
        targets = self.config.get("targets", [])
        return {
            "pipeline": self.config.get("name", "unnamed"),
            "sources": len(sources),
            "targets": len(targets),
            "mappings": len(self.config.get("mapping", {}).get("field_mappings", [])),
            "estimated_rows": self.config.get("estimated_rows", 0),
            "batch_size": self.config.get("batch_size", 1000),
            "validations": len(self.config.get("validations", [])),
            "post_process": self.config.get("post_process", []),
        }

    def run(self, dry_run=False, db_source=None, db_target=None, kafka_servers=None,
            workers=4, resume=None) -> dict:
        logger.info(f"Starting pipeline: {self.config.get('name', 'unnamed')}")
        logger.info(f"Dry run: {dry_run}")

        source_type = self.config.get("sources", [{}])[0].get("type", "csv") if self.config.get("sources") else "csv"
        source_path = self.config.get("sources", [{}])[0].get("path", "") if self.config.get("sources") else ""

        target = self.config.get("targets", [{}])[0] if self.config.get("targets") else {}
        target_schema = target.get("schema", "snisid_ncid")
        target_table = target.get("table", "citizens")

        mapping = self.config.get("mapping", {})
        field_mappings = mapping.get("field_mappings", [])
        batch_size = self.config.get("batch_size", 1000)
        validations = mapping.get("validations", [])
        post_process = mapping.get("post_process", [])

        self.adapter = self._get_adapter(source_type, source_path, db_source)
        rows = self.adapter.read()
        logger.info(f"Read {len(rows)} rows from source")

        transformed = []
        progress = {"skipped": 0, "errors": 0}
        resume_from = int(resume) if resume else 0

        for i, row in enumerate(rows):
            if i < resume_from:
                continue
            try:
                record = self._apply_mapping(row, field_mappings)
                errors = self._validate(record, validations)
                if errors:
                    progress["errors"] += 1
                    self.stats["rows_errors"] += 1
                    for e in errors:
                        self.stats["errors"].append(f"Row {i}: {e}")
                    continue
                transformed.append(record)
                self.stats["rows_processed"] += 1
            except Exception as e:
                progress["errors"] += 1
                self.stats["rows_errors"] += 1
                self.stats["errors"].append(f"Row {i}: {e}")

            if len(transformed) >= batch_size:
                if not dry_run:
                    self._write_batch(transformed, target_schema, target_table, db_target)
                self.stats["batches_completed"] += 1
                self.stats["rows_inserted"] += len(transformed)
                self.stats["checkpoint"] = str(i + 1)
                logger.info(f"Batch {self.stats['batches_completed']}: {len(transformed)} rows inserted (row {i+1})")
                transformed = []

        if transformed:
            if not dry_run:
                self._write_batch(transformed, target_schema, target_table, db_target)
            self.stats["batches_completed"] += 1
            self.stats["rows_inserted"] += len(transformed)
            logger.info(f"Final batch: {len(transformed)} rows")

        if not dry_run:
            for step in post_process:
                self._run_post_process(step, db_target, kafka_servers)

        self.stats["success"] = progress["errors"] == 0
        logger.info(f"Pipeline complete: {self.stats['rows_inserted']} inserted, "
                    f"{self.stats['rows_errors']} errors, "
                    f"{self.stats['batches_completed']} batches")
        return self.stats

    def _get_adapter(self, source_type: str, source_path: str, db_source: str):
        if source_type == "csv":
            from adapters.csv_adapter import CsvAdapter
            return CsvAdapter(source_path)
        elif source_type == "excel":
            from adapters.excel_adapter import ExcelAdapter
            return ExcelAdapter(source_path)
        elif source_type == "json":
            from adapters.json_adapter import JsonAdapter
            return JsonAdapter(source_path)
        elif source_type in ("postgres", "mysql", "mssql"):
            from adapters.db_adapter import DbAdapter
            return DbAdapter(source_type, db_source)
        elif source_type == "pdf":
            from adapters.pdf_adapter import PdfAdapter
            return PdfAdapter(source_path)
        else:
            raise ValueError(f"Unsupported source type: {source_type}")

    def _apply_mapping(self, row: dict, mappings: list) -> dict:
        record = {}
        for mapping in mappings:
            source_key = mapping.get("source")
            target_key = mapping.get("target")
            field_type = mapping.get("type", "string")
            default = mapping.get("default")
            transformer = mapping.get("transformer")
            enum_map = mapping.get("mapping", {})

            raw_value = row.get(source_key, default)
            if raw_value is None:
                if mapping.get("required"):
                    raise ValueError(f"Required field '{source_key}' is missing")
                record[target_key] = None
                continue

            value = raw_value
            if transformer:
                value = self._apply_transformer(value, transformer, mapping)
            if field_type == "enum" and enum_map:
                value = enum_map.get(str(value).strip(), value)
            if field_type in ("int", "integer"):
                try:
                    value = int(value)
                except (ValueError, TypeError):
                    value = 0
            elif field_type == "float":
                try:
                    value = float(value)
                except (ValueError, TypeError):
                    value = 0.0
            elif field_type == "date":
                fmt = mapping.get("format", "%Y-%m-%d")
                try:
                    value = datetime.strptime(str(value)[:25].strip(), fmt).isoformat()
                except ValueError:
                    value = None

            record[target_key] = value
        return record

    def _apply_transformer(self, value, transformer: str, mapping: dict) -> str:
        t = transformer.lower()
        if t == "upper":
            return str(value).upper()
        elif t == "lower":
            return str(value).lower()
        elif t == "trim":
            return str(value).strip()
        elif t == "capitalize":
            return str(value).strip().capitalize()
        elif t == "niu_generator":
            from transformers.id_generators import generate_niu
            return generate_niu(value)
        elif t == "phone_normalize":
            digits = re.sub(r"\D", "", str(value))
            if len(digits) == 8:
                return f"+509{digits}"
            if len(digits) == 10 and digits.startswith("509"):
                return f"+{digits}"
            return digits
        elif t == "plate_normalize":
            return str(value).upper().replace(" ", "-")[:15]
        elif t == "hash_pii":
            return hashlib.sha256(str(value).encode()).hexdigest()[:16]
        elif t == "file_copy":
            from transformers.file_handlers import copy_to_minio
            return copy_to_minio(str(value), mapping.get("copy_to"))
        elif t == "base64_decode":
            from transformers.file_handlers import base64_to_file
            return base64_to_file(value, mapping.get("target"))
        elif t == "concat":
            return str(value)
        elif "regex_extract" in t:
            pattern = mapping.get("pattern", r"(.*)")
            m = re.search(pattern, str(value))
            return m.group(1) if m else ""
        elif t == "geo_encode":
            from transformers.lookup_transformers import geo_encode
            return geo_encode(value)
        else:
            return str(value)

    def _validate(self, record: dict, validations: list) -> list:
        errors = []
        for v in validations:
            field = v.get("field")
            rule = v.get("rule")
            value = record.get(field)
            if rule == "required" and (value is None or str(value).strip() == ""):
                errors.append(f"'{field}' is required")
            elif rule == "niu_format" and value:
                if not re.match(r"^[A-Z0-9]{10}$", str(value)):
                    errors.append(f"'{field}': invalid NIU format: {value}")
            elif rule == "email" and value:
                if not re.match(r"^[^@\s]+@[^@\s]+\.[^@\s]+$", str(value)):
                    errors.append(f"'{field}': invalid email: {value}")
            elif rule == "phone" and value:
                digits = re.sub(r"\D", "", str(value))
                if len(digits) < 8:
                    errors.append(f"'{field}': invalid phone: {value}")
            elif rule == "plate" and value:
                if not re.match(r"^[A-Z0-9-]{2,15}$", str(value).upper()):
                    errors.append(f"'{field}': invalid plate: {value}")
            elif rule == "not_future" and value:
                try:
                    from datetime import datetime
                    dt = datetime.fromisoformat(str(value))
                    if dt > datetime.utcnow():
                        errors.append(f"'{field}': date is in the future")
                except ValueError:
                    pass
            elif rule == "min_length" and value:
                min_len = v.get("min", 1)
                if len(str(value)) < min_len:
                    errors.append(f"'{field}': too short ({len(str(value))} < {min_len})")
        return errors

    def _write_batch(self, records: list, schema: str, table: str, db_target: str):
        from writers.pg_writer import PostgresWriter
        writer = PostgresWriter(db_target) if db_target else PostgresWriter()
        writer.write_batch(records, schema, table)

    def _run_post_process(self, step: dict, db_target: str, kafka_servers: str):
        action = step.get("action")
        if action == "emit_kafka":
            from writers.kafka_writer import KafkaWriter
            w = KafkaWriter(kafka_servers)
            w.emit_topic(step.get("topic"), self.stats)
        elif action == "update_neo4j":
            from writers.neo4j_writer import Neo4jWriter
            w = Neo4jWriter()
            w.update_graph(step.get("label"), self.stats)
