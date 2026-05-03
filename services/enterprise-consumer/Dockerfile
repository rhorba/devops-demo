# Stage 1 — extract IBM MQ client libraries from the official image
FROM icr.io/ibm-messaging/mq:latest AS mq-libs

# Stage 2 — build the Go binary with CGO against MQ headers
FROM golang:1.22-bullseye AS builder
COPY --from=mq-libs /opt/mqm/inc  /opt/mqm/inc
COPY --from=mq-libs /opt/mqm/lib64 /opt/mqm/lib64
ENV CGO_CFLAGS="-I/opt/mqm/inc"
ENV CGO_LDFLAGS="-L/opt/mqm/lib64 -Wl,-rpath,/opt/mqm/lib64"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o enterprise-consumer .

# Stage 3 — minimal runtime with only the MQ shared libraries
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
        libstdc++6 libgcc-s1 ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=mq-libs /opt/mqm/lib64 /opt/mqm/lib64
ENV LD_LIBRARY_PATH=/opt/mqm/lib64
RUN useradd -u 1001 -m -s /sbin/nologin appuser
USER appuser
WORKDIR /app
COPY --from=builder /app/enterprise-consumer .
EXPOSE 8080
ENTRYPOINT ["./enterprise-consumer"]
