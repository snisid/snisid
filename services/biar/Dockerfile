FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /biar-svc ./cmd/main

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /biar-svc /usr/local/bin/biar-svc
EXPOSE 8103
ENTRYPOINT ["biar-svc"]
