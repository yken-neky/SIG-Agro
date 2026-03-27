# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY services/production-service/go.mod .
COPY services/production-service/go.sum* .
COPY api /api

# Download dependencies
RUN go mod download

# Copy source code
COPY services/production-service .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o production-service ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/production-service .

EXPOSE 50054

CMD ["./production-service"]
