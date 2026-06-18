FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/offline-sync-engine ./services/offline-sync-engine/cmd

FROM alpine:3.19
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app
COPY --from=builder /app/offline-sync-engine .
RUN mkdir -p /app/data

EXPOSE 8090
ENV SQLITE_DB_PATH=/app/data/offline_sync.db
ENV PORT=8090

ENTRYPOINT ["./offline-sync-engine"]
