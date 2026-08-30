FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /radiation-safety-svc ./cmd

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /radiation-safety-svc /usr/local/bin/radiation-safety-svc
EXPOSE 8308
ENTRYPOINT ["radiation-safety-svc"]
