FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /sigint-ht ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /sigint-ht /sigint-ht
EXPOSE 8301
ENTRYPOINT ["/sigint-ht"]
