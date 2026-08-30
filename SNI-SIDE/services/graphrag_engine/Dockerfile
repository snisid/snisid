FROM python:3.12-slim

WORKDIR /app

RUN pip install --no-cache-dir \
    fastapi uvicorn[standard] aiokafka neo4j prometheus-client httpx pydantic && \
    rm -rf /root/.cache/pip

COPY *.py .
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    rm -rf /root/.cache/pip

RUN useradd -r -u 1001 sniside && chown -R sniside:sniside /app
USER sniside

HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=3 \
    CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8080/health')" || exit 1

ENV PYTHONUNBUFFERED=1
ENV NEO4J_URI=bolt://neo4j:7687
ENV LLM_ENDPOINT=http://ollama:11434
ENV LLM_MODEL=mistral:7b-instruct-v0.3-q4_K_M
ENV METRICS_PORT=9103
ENV KAFKA_BOOTSTRAP_SERVERS=kafka:9092
ENV KAFKA_GROUP_ID=sniside-graphrag

EXPOSE 8080

CMD ["python", "-m", "uvicorn", "api:app", "--host", "0.0.0.0", "--port", "8080"]
