.PHONY: help proto build docker-build docker-up docker-down clean-up migrate test grpcurl

help:
	@echo "SIG-Agro Makefile Commands"
	@echo "=========================="
	@echo "proto              - Generate Go code from .proto files"
	@echo "sqlc               - Generate SQLC code for all services"
	@echo "build              - Build all microservices"
	@echo "docker-build       - Build Docker images"
	@echo "docker-up          - Start all services with docker-compose"
	@echo "docker-down        - Stop all services"
	@echo "docker-logs        - View logs from all services"
	@echo "logs-user          - View logs from user service"
	@echo "logs-producer      - View logs from producer service"
	@echo "logs-parcel        - View logs from parcel service"
	@echo "migrate-up         - Run database migrations"
	@echo "migrate-down       - Rollback database migrations"
	@echo "grpcurl-user       - Test user service with grpcurl"
	@echo "grpcurl-producer   - Test producer service with grpcurl"
	@echo "grpcurl-parcel     - Test parcel service with grpcurl"
	@echo "clean-up           - Remove all generated files"
	@echo "test               - Run unit tests"

# Generate Protocol Buffer code
proto:
	@echo "Generating Protocol Buffer code..."
	protoc --go_out=. --go_opt=paths=source_relative \
	        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	        api/proto/*.proto
	@echo "Protocol Buffer code generated successfully"

# Generate SQLC code
sqlc:
	@echo "Generating SQLC code..."
	cd services/user-service && sqlc generate && cd ../..
	cd services/producer-service && sqlc generate && cd ../..
	cd services/parcel-service && sqlc generate && cd ../..
	cd services/production-service && sqlc generate && cd ../..
	cd services/alert-service && sqlc generate && cd ../..
	cd services/notification-service && sqlc generate && cd ../..
	cd services/report-service && sqlc generate && cd ../..
	@echo "SQLC code generated successfully"

# Build all microservices
build:
	@echo "Building user-service..."
	cd services/user-service && go build -v -o bin/user-service ./cmd/server && cd ../..
	@echo "Building producer-service..."
	cd services/producer-service && go build -v -o bin/producer-service ./cmd/server && cd ../..
	@echo "Building parcel-service..."
	cd services/parcel-service && go build -v -o bin/parcel-service ./cmd/server && cd ../..
	@echo "Building production-service..."
	cd services/production-service && go build -v -o bin/production-service ./cmd/server && cd ../..
	@echo "Building alert-service..."
	cd services/alert-service && go build -v -o bin/alert-service ./cmd/server && cd ../..
	@echo "Building notification-service..."
	cd services/notification-service && go build -v -o bin/notification-service ./cmd/server && cd ../..
	@echo "Building report-service..."
	cd services/report-service && go build -v -o bin/report-service ./cmd/server && cd ../..
	@echo "All services built successfully"

# Build Docker images
docker-build:
	@echo "Building Docker images..."
	docker-compose -f docker-compose.yml build
	@echo "Docker images built successfully"

# Start all services with docker-compose
docker-up:
	@echo "Starting all services with docker-compose..."
	docker-compose -f docker-compose.yml up -d
	@echo "Services started. Waiting for database to be ready..."
	sleep 10
	@echo "Running database migrations..."
	docker-compose -f docker-compose.yml exec -T postgres psql -U user -d postgres -f /docker-entrypoint-initdb.d/01-init.sql || true
	@echo "All services are running. Check with: make docker-logs"

# Stop all services
docker-down:
	@echo "Stopping all services..."
	docker-compose -f docker-compose.yml down
	@echo "Services stopped"

# View logs from all services
docker-logs:
	@echo "Viewing logs from all services..."
	docker-compose -f docker-compose.yml logs -f

# View logs from specific service
logs-user:
	docker-compose -f docker-compose.yml logs -f user-service

logs-producer:
	docker-compose -f docker-compose.yml logs -f producer-service

logs-parcel:
	docker-compose -f docker-compose.yml logs -f parcel-service

# Health check
health:
	@echo "Checking service health..."
	docker-compose -f docker-compose.yml ps
	@echo "Checking database connection..."
	docker-compose -f docker-compose.yml exec -T postgres pg_isready -U user

# Test endpoints with grpcurl
grpcurl-user:
	@echo "Testing User Service..."
	grpcurl -plaintext localhost:50051 user.UserService/GetUser <<EOF
	{"user_id": 1}
	EOF

grpcurl-producer:
	@echo "Testing Producer Service..."
	grpcurl -plaintext localhost:50052 producer.ProducerService/ListProducers <<EOF
	{"user_id": 1, "limit": 10, "offset": 0}
	EOF

grpcurl-parcel:
	@echo "Testing Parcel Service..."
	grpcurl -plaintext localhost:50053 parcel.ParcelService/ListParcels <<EOF
	{"producer_id": 1, "limit": 10, "offset": 0}
	EOF

# Generate self-signed certificate for HAProxy
generate-cert:
	@echo "Generating self-signed certificate..."
	mkdir -p infrastructure/haproxy/certs
	openssl req -x509 -newkey rsa:2048 -keyout infrastructure/haproxy/certs/selfsigned.key -out infrastructure/haproxy/certs/selfsigned.crt -days 365 -nodes -subj "/C=ES/ST=State/L=City/O=Organization/CN=localhost"
	cat infrastructure/haproxy/certs/selfsigned.crt infrastructure/haproxy/certs/selfsigned.key > infrastructure/haproxy/certs/selfsigned.pem
	@echo "Certificate generated successfully"

# Clean up generated files
clean-up:
	@echo "Cleaning up generated files..."
	find services -name "sqlc" -type d -exec rm -rf {} + 2>/dev/null || true
	find services -name "bin" -type d -exec rm -rf {} + 2>/dev/null || true
	find api/proto -name "*.pb.go" -delete
	find api/proto -name "*_grpc.pb.go" -delete
	@echo "Cleanup complete"

# Run tests
test:
	@echo "Running tests..."
	cd services/user-service && go test -v ./... && cd ../..
	cd services/producer-service && go test -v ./... && cd ../..
	cd services/parcel-service && go test -v ./... && cd ../..
	cd services/production-service && go test -v ./... && cd ../..
	cd services/alert-service && go test -v ./... && cd ../..
	cd services/notification-service && go test -v ./... && cd ../..
	cd services/report-service && go test -v ./... && cd ../..
	@echo "Tests completed"

# Full setup from scratch
setup: clean-up proto sqlc build docker-build generate-cert docker-up health
	@echo "Setup complete! Services are running."
	@echo "HAProxy API Gateway: https://localhost:443"
	@echo "RabbitMQ Management: http://localhost:15672"
	@echo "PostgreSQL: localhost:5432"

# Development mode (live rebuild with hot reload)
dev:
	@echo "Starting development mode..."
	docker-compose -f docker-compose.yml up --build

# Generate gRPC documentation
proto-docs:
	@echo "Generating protocol buffer documentation..."
	protoc --doc_out=./docs/proto --doc_opt=markdown,proto.md \
	        api/proto/*.proto
	@echo "Documentation generated in docs/proto/"
