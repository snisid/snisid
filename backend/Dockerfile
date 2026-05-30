# ---- Build Stage ----
FROM python:3.12-slim AS builder

# Prevent Python from writing .pyc files & buffer output
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /build

# Install build dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc build-essential && \
    rm -rf /var/lib/apt/lists/*

# Install dependencies into a virtual environment
RUN python -m venv /opt/venv
ENV PATH="/opt/venv/bin:$PATH"

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . /build/app

# ---- Production Stage ----
FROM gcr.io/distroless/python3-debian12

# Set non-root user (distroless default nonroot is 65532)
USER nonroot:nonroot
WORKDIR /app

# Copy virtual environment and app code from builder
COPY --from=builder --chown=nonroot:nonroot /opt/venv /opt/venv
COPY --from=builder --chown=nonroot:nonroot /build/app /app

# Set path for python executable inside venv
ENV PATH="/opt/venv/bin:$PATH"
ENV PYTHONPATH="/app"

# Expose FastAPI port
EXPOSE 8000

# Run uvicorn server
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4", "--proxy-headers"]
