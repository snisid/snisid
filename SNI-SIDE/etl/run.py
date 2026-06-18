import json, logging, os, sys, argparse, time, csv
from datetime import datetime
from typing import Optional
from pathlib import Path

from pipeline import ETLPipeline
from progress import ProgressTracker

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(name)s] %(levelname)s: %(message)s",
)
logger = logging.getLogger("sniside.etl")


def parse_args():
    parser = argparse.ArgumentParser(description="SNI-SIDE ETL Migration Toolkit")
    parser.add_argument("--pipeline", "-p", help="Pipeline JSON file path")
    parser.add_argument("--source", "-s", help="Source file or connection string")
    parser.add_argument("--target", "-t", help="Target table (schema.table)")
    parser.add_argument("--mapping", "-m", help="Mapping JSON file path")
    parser.add_argument("--db-source", help="Source DB connection string (for DB adapter)")
    parser.add_argument("--db-target", help="Target DB connection string")
    parser.add_argument("--kafka-servers", default="kafka:9092", help="Kafka bootstrap servers")
    parser.add_argument("--batch-size", type=int, default=1000, help="Batch size for writes")
    parser.add_argument("--kafka-events", action="store_true", help="Emit Kafka events")
    parser.add_argument("--neo4j-update", action="store_true", help="Update Neo4j graph")
    parser.add_argument("--dry-run", action="store_true", help="Validate only, no writes")
    parser.add_argument("--report-only", action="store_true", help="Generate report only")
    parser.add_argument("--resume", help="Resume from checkpoint file")
    parser.add_argument("--workers", type=int, default=4, help="Parallel workers")
    parser.add_argument("--log-level", default="INFO", help="Logging level")
    return parser.parse_args()


def main():
    args = parse_args()
    logging.getLogger("sniside.etl").setLevel(getattr(logging, args.log_level.upper()))

    if args.pipeline:
        pipeline_path = Path(args.pipeline)
        if not pipeline_path.exists():
            logger.error(f"Pipeline file not found: {pipeline_path}")
            sys.exit(1)
        config = json.loads(pipeline_path.read_text())
    else:
        config = {
            "sources": [{"type": "csv", "path": args.source}] if args.source else [],
            "targets": [{"schema": args.target.split(".")[0], "table": args.target.split(".")[1]}] if args.target else [],
            "mapping": json.loads(Path(args.mapping).read_text()) if args.mapping else {},
            "batch_size": args.batch_size,
            "kafka_events": args.kafka_events,
            "neo4j_update": args.neo4j_update,
        }

    pipeline = ETLPipeline(config)
    tracker = ProgressTracker()

    if args.report_only:
        report = pipeline.generate_report()
        print(json.dumps(report, indent=2))
        return

    try:
        start = time.monotonic()
        summary = pipeline.run(
            dry_run=args.dry_run,
            db_source=args.db_source,
            db_target=args.db_target,
            kafka_servers=args.kafka_servers,
            workers=args.workers,
            resume=args.resume,
        )
        elapsed = time.monotonic() - start
        summary["elapsed_seconds"] = round(elapsed, 2)
        summary["rows_per_second"] = round(summary.get("rows_processed", 0) / max(elapsed, 0.1), 1)

        print(f"\n=== Migration Summary ===")
        print(f"Status:        {'SUCCESS' if summary.get('success') else 'FAILED'}")
        print(f"Rows processed: {summary.get('rows_processed', 0):,}")
        print(f"Rows inserted:  {summary.get('rows_inserted', 0):,}")
        print(f"Rows errors:    {summary.get('rows_errors', 0):,}")
        print(f"Batches:        {summary.get('batches_completed', 0):,}")
        print(f"Elapsed:        {summary['elapsed_seconds']:.1f}s")
        print(f"Throughput:     {summary['rows_per_second']:.0f} rows/s")
        print(f"Checkpoint:     {summary.get('checkpoint', 'N/A')}")

        if summary.get("errors"):
            print(f"\nErrors ({len(summary['errors'])}):")
            for err in summary["errors"][:10]:
                print(f"  - {err}")

        if summary.get("warnings"):
            print(f"\nWarnings ({len(summary['warnings'])}):")
            for w in summary["warnings"][:10]:
                print(f"  - {w}")

        if not summary.get("success"):
            sys.exit(1)

    except Exception as e:
        logger.error(f"Pipeline failed: {e}", exc_info=True)
        sys.exit(1)


if __name__ == "__main__":
    main()
