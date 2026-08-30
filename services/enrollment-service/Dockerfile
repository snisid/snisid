FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /enrollment-service ./services/enrollment-service/cmd

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /enrollment-service /usr/local/bin/enrollment-service
EXPOSE 8083
CMD ["enrollment-service"]
