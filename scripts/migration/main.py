"""SNISID — CLI de migration de données.

Point d'entrée principal pour l'outil de migration.
Supporte les commandes: run, validate, status, restart.
"""

import logging
import sys
from pathlib import Path
from typing import Optional

import click

from config import MigrationConfig
from checkpoint import CheckpointManager
from source_connectors import (
    CSVConnector,
    SQLServerConnector,
    PostgreSQLConnector,
    RESTConnector,
)
from data_cleansing import DataCleanser
from matching_engine import MatchingEngine
from target_loader import TargetLoader, LoaderConfig
from etl_pipeline import ETLPipeline

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    handlers=[
        logging.StreamHandler(),
        logging.FileHandler("migration.log"),
    ],
)
logger = logging.getLogger("snisid.migration")


@click.group()
@click.option("--config", "-c", default="config.yaml", help="Fichier de configuration")
@click.option("--verbose", "-v", is_flag=True, help="Logging verbose")
@click.pass_context
def cli(ctx: click.Context, config: str, verbose: bool) -> None:
    """SNISID — Outil de Migration de Données Nationales."""
    if verbose:
        logging.getLogger().setLevel(logging.DEBUG)

    cfg = MigrationConfig.from_yaml(config)
    ctx.ensure_object(dict)
    ctx.obj["config"] = cfg


@cli.command()
@click.option("--source", help="Type de source (csv, sqlserver, postgresql, rest)")
@click.option("--dry-run", is_flag=True, help="Simule sans écrire")
@click.option("--max-records", type=int, help="Limite de records à traiter")
@click.pass_context
def run(ctx: click.Context, source: Optional[str], dry_run: bool, max_records: Optional[int]) -> None:
    """Exécute le pipeline ETL de migration."""
    cfg: MigrationConfig = ctx.obj["config"]

    if source:
        cfg.source.type = source

    connector = _create_source_connector(cfg)
    cleanser = DataCleanser()
    matcher = MatchingEngine()
    checkpoint = CheckpointManager(Path(cfg.checkpoint.path))

    if dry_run:
        logger.info("=== MODE DRY RUN — Aucune donnée ne sera écrite ===")

    loader_cfg = LoaderConfig(
        api_base_url=cfg.target.api_url,
        api_key=cfg.target.api_key,
        batch_size=cfg.target.batch_size,
        rate_limit_per_second=cfg.target.rate_limit,
    )
    loader = TargetLoader(loader_cfg)

    pipeline = ETLPipeline(
        source=connector,
        cleanser=cleanser,
        matcher=matcher,
        loader=loader,
        checkpoint_mgr=checkpoint,
        batch_size=cfg.pipeline.batch_size,
    )

    stats = pipeline.run(max_records=max_records)
    click.echo(f"\n✅ Migration terminée: {stats.processed} enregistrements traités")
    click.echo(f"   Débit: {stats.throughput:.0f} rec/s")
    click.echo(f"   Erreurs: {stats.errors}")
    click.echo(f"   Doublons détectés: {stats.duplicates}")


@cli.command()
@click.pass_context
def validate(ctx: click.Context) -> None:
    """Valide la configuration et les connexions."""
    cfg: MigrationConfig = ctx.obj["config"]
    connector = _create_source_connector(cfg)
    cleanser = DataCleanser()
    matcher = MatchingEngine()
    checkpoint = CheckpointManager(Path(cfg.checkpoint.path))

    loader_cfg = LoaderConfig(
        api_base_url=cfg.target.api_url,
        api_key=cfg.target.api_key,
    )
    loader = TargetLoader(loader_cfg)

    pipeline = ETLPipeline(
        source=connector,
        cleanser=cleanser,
        matcher=matcher,
        loader=loader,
        checkpoint_mgr=checkpoint,
    )

    result = pipeline.validate()
    if result["status"] == "healthy":
        click.echo("✅ Configuration valide — toutes les connexions sont OK")
    else:
        click.echo("⚠️  Configuration dégradée:")
        for check, status in result["checks"].items():
            status_str = "✅" if status else "❌"
            click.echo(f"   {status_str} {check}")


@cli.command()
@click.pass_context
def status(ctx: click.Context) -> None:
    """Affiche l'état du dernier checkpoint."""
    cfg: MigrationConfig = ctx.obj["config"]
    checkpoint = CheckpointManager(Path(cfg.checkpoint.path))
    cp = checkpoint.load()

    if cp:
        click.echo(f"📊 État de la migration:")
        click.echo(f"   Offset: {cp.offset}")
        click.echo(f"   Batch: {cp.batch}")
        click.echo(f"   Records traités: {cp.records_processed}")
    else:
        click.echo("📊 Aucun checkpoint trouvé — migration non démarrée")


@cli.command()
@click.confirmation_option(prompt="⚠️  Réinitialiser la progression ?")
@click.pass_context
def restart(ctx: click.Context) -> None:
    """Réinitialise la progression (efface le checkpoint)."""
    cfg: MigrationConfig = ctx.obj["config"]
    connector = _create_source_connector(cfg)
    cleanser = DataCleanser()
    matcher = MatchingEngine()
    loader_cfg = LoaderConfig(api_base_url=cfg.target.api_url)
    loader = TargetLoader(loader_cfg)
    checkpoint = CheckpointManager(Path(cfg.checkpoint.path))

    pipeline = ETLPipeline(
        source=connector,
        cleanser=cleanser,
        matcher=matcher,
        loader=loader,
        checkpoint_mgr=checkpoint,
    )

    pipeline.reset()
    click.echo("✅ Migration réinitialisée — tous les checkpoints effacés")


def _create_source_connector(cfg: MigrationConfig):
    """Crée le connecteur source approprié."""
    source_type = cfg.source.type.lower()

    if source_type == "csv":
        return CSVConnector(
            filepath=cfg.source.path,
            delimiter=cfg.source.csv_delimiter,
            encoding=cfg.source.encoding,
        )
    elif source_type == "sqlserver":
        return SQLServerConnector(
            connection_string=cfg.source.connection_string,
            query=cfg.source.query,
        )
    elif source_type == "postgresql":
        return PostgreSQLConnector(
            host=cfg.source.host,
            port=cfg.source.port,
            database=cfg.source.database,
            user=cfg.source.user,
            password=cfg.source.password,
            query=cfg.source.query,
        )
    elif source_type == "rest":
        return RESTConnector(
            base_url=cfg.source.api_url,
            endpoint=cfg.source.api_endpoint,
            api_key=cfg.source.api_key,
        )
    else:
        raise ValueError(f"Type de source inconnu: {source_type}")


if __name__ == "__main__":
    cli()
