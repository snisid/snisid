# =========================
# FRONTEND (Node)
# =========================
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# isolate npm cache (prevents EIO cascade)
ENV NPM_CONFIG_CACHE=/tmp/.npm

COPY frontend/package*.json ./

# deterministic install
RUN npm ci --no-audit --no-fund

COPY frontend/ .

RUN npm run build


# =========================
# BACKEND (Go)
# =========================
FROM golang:1.26-alpine AS backend-builder

WORKDIR /app

# isolate Go caches (CRITICAL FIX)
ENV GOCACHE=/tmp/gocache
ENV GOMODCACHE=/tmp/gomod

# Accept a build argument for the service to build (defaults to gateway)
ARG SERVICE_PATH=gateway/cmd

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./${SERVICE_PATH}


# =========================
# RUNTIME
# =========================
FROM alpine:3.20

WORKDIR /app

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /app/dist ./public

CMD ["./main"]