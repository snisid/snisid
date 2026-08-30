FROM python:3.12-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    libpq-dev gcc && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && rm -rf /root/.cache/pip

COPY run.py pipeline.py progress.py .
COPY adapters/ adapters/
COPY transformers/ transformers/
COPY validators/ validators/
COPY writers/ writers/
COPY pipelines/ pipelines/
COPY mappings/ mappings/
COPY README.md .

RUN useradd -r -u 1001 sniside && chown -R sniside:sniside /app
USER sniside

ENTRYPOINT ["python", "run.py"]
CMD ["--help"]
