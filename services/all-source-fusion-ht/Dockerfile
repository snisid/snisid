FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /all-source-fusion-ht ./cmd

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /all-source-fusion-ht /usr/local/bin/all-source-fusion-ht
EXPOSE 8309
ENTRYPOINT ["all-source-fusion-ht"]
