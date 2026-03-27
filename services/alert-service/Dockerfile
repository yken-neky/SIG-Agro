# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY services/alert-service/go.mod .
COPY services/alert-service/go.sum* .
COPY api /api

# Download dependencies
RUN go mod download

# Copy source code
COPY services/alert-service .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o alert-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/alert-service .

EXPOSE 50055

CMD ["./alert-service"]
