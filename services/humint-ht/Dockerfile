FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /humint-ht ./cmd/main.go

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /humint-ht /humint-ht
EXPOSE 8302
ENTRYPOINT ["/humint-ht"]
