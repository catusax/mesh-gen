FROM alpine AS builder
LABEL org.opencontainers.image.source https://github.com/catusax/mesh-gen
RUN apk --update --no-cache add ca-certificates
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.11 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe