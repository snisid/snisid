FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /afis-svc ./cmd/server

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /afis-svc /usr/local/bin/afis-svc

EXPOSE 8091
ENTRYPOINT ["afis-svc"]
