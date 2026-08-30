# syntax=docker/dockerfile:1.6
# =============================================================================
# SNISID - Dockerfile multi-stage - VENDOR OFFLINE
# Backend Go  : go.mod + go.work à la racine, code source dans ./backend/
# Frontend    : ./frontend/
# Build 100% offline, pas besoin de proxy.golang.org
#
# Prérequis sur la machine hôte (1 fois) :
#   go work vendor          ← workspace project: use go work vendor, NOT go mod vendor
#   # Enlever "vendor/" de .dockerignore si présent
# =============================================================================

# ---------- 1. Backend builder ----------
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION}-alpine AS backend-builder

WORKDIR /app

# ca-certificates only — git not needed in vendor mode
RUN apk add --no-cache ca-certificates

# Vendor mode: no network calls, no checksum DB needed
ENV GONOSUMDB=* \
    GOFLAGS=-mod=vendor \
    CGO_ENABLED=0 \
    GOOS=linux

# Copy module definition files — go.work required for workspace vendor
COPY go.mod go.sum ./
COPY go.work go.work.sum ./

# Copy vendor directory (created by: go work vendor)
COPY vendor ./vendor

# Copy all local modules referenced in go.work
# Adjust this list to match the `use` directives in your go.work file
COPY backend ./backend
COPY pkg ./pkg
COPY internal ./internal

ARG BACKEND_PATH=./backend
ARG BINARY_NAME=server
RUN go build -ldflags="-s -w" -o /app/${BINARY_NAME} ${BACKEND_PATH}

# ---------- 2. Frontend builder ----------
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package*.json ./

# npm ci uses npm cache — correct mount target
RUN --mount=type=cache,target=/root/.npm npm ci --no-audit --no-fund

COPY frontend/ ./

# npm run build uses node_modules, not npm cache — no mount here
RUN npm run build

# ---------- 3. Image finale ----------
FROM alpine:3.20 AS runtime

WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -g 10001 -S app && \
    adduser -u 10001 -S app -G app

ARG BINARY_NAME=server
COPY --from=backend-builder /app/${BINARY_NAME} /app/server
COPY --from=frontend-builder /app/frontend/dist /app/web

RUN chown -R app:app /app
USER app

ENV PORT=8080 \
    WEB_ROOT=/app/web \
    GIN_MODE=release

EXPOSE 8080
ENTRYPOINT ["/app/server"]
