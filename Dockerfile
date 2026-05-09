# Multi-stage Dockerfile for building and running the Go application

FROM golang:1.21-bullseye AS builder
WORKDIR /workspace

# Allow multi-arch builds via buildx by passing TARGETOS/TARGETARCH arguments
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG TARGETVARIANT=

# Install build tools required for cgo
RUN apt-get update && apt-get install -y build-essential ca-certificates && rm -rf /var/lib/apt/lists/*

# Download modules separately to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Enable cgo (onnxruntime uses native libs) and set target OS/ARCH from build args
ENV CGO_ENABLED=1
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

# Build the server (adjust package path if your main is elsewhere)
RUN go build -o /app/server ./cmd

FROM debian:bookworm-slim
WORKDIR /app

# Runtime deps
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy onnxruntime native libraries into standard library path (from repo's onnxruntime/lib)
COPY --from=builder /workspace/onnxruntime/lib/ /usr/lib/
RUN ldconfig || true

# Copy the built binary
COPY --from=builder /app/server /usr/local/bin/server

# Copy model and any AI assets the service expects (adjust path if your code uses a different one)
COPY --from=builder /workspace/internal/infrastructures/ai /app/internal/infrastructures/ai

# Expose an env var pointing to the ONNX model folder inside the image
ENV ONNX_MODEL_DIR=/app/internal/infrastructures/ai
ENV LD_LIBRARY_PATH=/usr/lib:$LD_LIBRARY_PATH

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
