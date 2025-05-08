# syntax=docker/dockerfile:1

#####################################
# 1. Build stage: compile the binary
#####################################
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service ./cmd/orders

#####################################
# 2. Final stage: minimal runtime
#####################################
FROM alpine:3.18
WORKDIR /app

# SSL certs for TLS if needed
RUN apk --no-cache add ca-certificates

# Copy the statically-linked binary
COPY --from=builder /app/order-service .

EXPOSE 8080
ENTRYPOINT ["/app/order-service"]
