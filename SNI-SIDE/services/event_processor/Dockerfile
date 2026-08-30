FROM python:3.12-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    rm -rf /root/.cache/pip

COPY *.py .

RUN useradd -r -u 1001 sniside && chown -R sniside:sniside /app
USER sniside

HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:9102/')" || exit 1

ENV PYTHONUNBUFFERED=1
ENV KAFKA_BOOTSTRAP_SERVERS=kafka:9092
ENV KAFKA_GROUP_ID=sniside-event-processor
ENV METRICS_PORT=9102
ENV NEO4J_URI=bolt://neo4j:7687
ENV REDIS_HOST=redis
ENV CLICKHOUSE_HOST=clickhouse

CMD ["python", "main.py"]
