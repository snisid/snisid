FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /identity-mesh ./services/identity-mesh

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /identity-mesh /usr/local/bin/identity-mesh
EXPOSE 8084
CMD ["identity-mesh"]
