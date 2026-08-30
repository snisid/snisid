FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gang-svc ./cmd/main

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /gang-svc /usr/local/bin/gang-svc
EXPOSE 8095
ENTRYPOINT ["gang-svc"]
