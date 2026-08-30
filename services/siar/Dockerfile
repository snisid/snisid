FROM golang:1.26-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/siar-ht ./cmd/main.go

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /build/siar-ht /usr/local/bin/siar-ht

EXPOSE 8102
ENTRYPOINT ["/usr/local/bin/siar-ht"]
