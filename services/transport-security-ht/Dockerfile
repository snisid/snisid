FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /transport-security-ht ./cmd

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /transport-security-ht /usr/local/bin/transport-security-ht
EXPOSE 8307
ENTRYPOINT ["transport-security-ht"]
