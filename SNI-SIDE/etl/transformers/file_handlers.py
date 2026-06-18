import io, base64, logging
from pathlib import Path

logger = logging.getLogger("sniside.etl.transformers")

MINIO_ENDPOINT = "minio:9000"
MINIO_ACCESS_KEY = "sniside"
MINIO_SECRET_KEY = "sniside-minio-secret"


def copy_to_minio(source_path: str, target_prefix: str = None) -> str:
    try:
        from minio import Minio
    except ImportError:
        logger.warning("minio not available, returning source path")
        return source_path

    client = Minio(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, secure=False)
    bucket = "sniside-media"
    if not client.bucket_exists(bucket):
        client.make_bucket(bucket)

    src = Path(source_path)
    if not src.exists():
        logger.warning(f"File not found for MinIO copy: {source_path}")
        return source_path

    target = f"{target_prefix or 'etl'}/{src.name}"
    client.fput_object(bucket, target, str(src))
    logger.info(f"Copied {source_path} → minio://{bucket}/{target}")
    return f"minio://{bucket}/{target}"


def base64_to_file(data: str, target_field: str = "file") -> str:
    bucket = "sniside-media"
    try:
        from minio import Minio
        client = Minio(MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, secure=False)
        if not client.bucket_exists(bucket):
            client.make_bucket(bucket)
        blob = base64.b64decode(data)
        path = f"etl/base64/{uuid.uuid4().hex}"
        client.put_object(bucket, path, io.BytesIO(blob), len(blob))
        return f"minio://{bucket}/{path}"
    except Exception as e:
        logger.warning(f"base64 decode failed: {e}")
        return data[:100] + "..."
