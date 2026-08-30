# =============================================================================
# SNISID - Go Microservice Runtime
# =============================================================================

FROM golang:1.26-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-w -s" \
    -o /out/service \
    ./cmd

# =============================================================================
# DISTROLESS RUNTIME
# =============================================================================

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /out/service /service

EXPOSE 8080

ENTRYPOINT ["/service"]