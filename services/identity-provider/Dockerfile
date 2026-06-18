FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /out/service ./services/identity-provider/cmd

FROM gcr.io/distroless/static-debian12
COPY --from=builder /out/service /service
EXPOSE 8096
ENTRYPOINT ["/service"]
