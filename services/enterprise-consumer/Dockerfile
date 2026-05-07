# Pure Go build — no CGO, no IBM MQ SDK headers required.
# Connects to IBM MQ via REST API (port 9443/HTTPS) — AMQP not available in Developer image.
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o enterprise-consumer .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1001 appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/enterprise-consumer .
EXPOSE 8080
ENTRYPOINT ["./enterprise-consumer"]
