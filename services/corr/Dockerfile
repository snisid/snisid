FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /corr-svc ./cmd/main

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /corr-svc /usr/local/bin/corr-svc
EXPOSE 8130
ENTRYPOINT ["corr-svc"]
