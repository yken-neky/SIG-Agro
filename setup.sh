#!/bin/bash
set -e

echo "╔════════════════════════════════════════════════════════╗"
echo "║      SIG-Agro: Sistema de Información Geográfica       ║"
echo "║                   Setup Script                         ║"
echo "╚════════════════════════════════════════════════════════╝"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check prerequisites
echo -e "\n${BLUE}Checking prerequisites...${NC}"

if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Docker is not installed. Please install Docker first.${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Docker Compose is not installed. Please install Docker Compose first.${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker and Docker Compose are installed${NC}"

# Generate self-signed certificate
echo -e "\n${BLUE}Generating self-signed certificate for HAProxy...${NC}"
mkdir -p infrastructure/haproxy/certs

if [ -f infrastructure/haproxy/certs/selfsigned.pem ]; then
    echo -e "${YELLOW}Certificate already exists, skipping generation${NC}"
else
    openssl req -x509 -newkey rsa:2048 \
        -keyout infrastructure/haproxy/certs/selfsigned.key \
        -out infrastructure/haproxy/certs/selfsigned.crt \
        -days 365 -nodes \
        -subj "/C=ES/ST=Madrid/L=Madrid/O=SIG-Agro/CN=localhost"

    cat infrastructure/haproxy/certs/selfsigned.crt \
        infrastructure/haproxy/certs/selfsigned.key > \
        infrastructure/haproxy/certs/selfsigned.pem

    echo -e "${GREEN}✓ Self-signed certificate generated${NC}"
fi

# Build Docker images
echo -e "\n${BLUE}Building Docker images...${NC}"
docker-compose build

echo -e "${GREEN}✓ Docker images built successfully${NC}"

# Start services
echo -e "\n${BLUE}Starting all services...${NC}"
docker-compose up -d

echo -e "${GREEN}✓ Services started${NC}"

# Wait for services to be ready
echo -e "\n${BLUE}Waiting for services to be ready...${NC}"
sleep 15

# Check service health
echo -e "\n${BLUE}Checking service health...${NC}"
docker-compose ps

# Print access information
echo -e "\n${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║         SIG-Agro is now running! 🎉                    ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"

echo -e "\n${BLUE}Service Endpoints:${NC}"
echo -e "  ${YELLOW}User Service:${NC}         grpc://localhost:50051"
echo -e "  ${YELLOW}Producer Service:${NC}     grpc://localhost:50052"
echo -e "  ${YELLOW}Parcel Service:${NC}       grpc://localhost:50053"
echo -e "  ${YELLOW}Production Service:${NC}   grpc://localhost:50054"
echo -e "  ${YELLOW}Alert Service:${NC}        grpc://localhost:50055"
echo -e "  ${YELLOW}Notification Service:${NC} grpc://localhost:50056"
echo -e "  ${YELLOW}Report Service:${NC}       grpc://localhost:50057"

echo -e "\n${BLUE}API Gateway:${NC}"
echo -e "  ${YELLOW}HTTP:${NC}  http://localhost:80"
echo -e "  ${YELLOW}GRPC:${NC}  grpc://localhost:443 (insecure: -insecure flag)"

echo -e "\n${BLUE}Management Interfaces:${NC}"
echo -e "  ${YELLOW}RabbitMQ:${NC}     http://localhost:15672 (user/password)"
echo -e "  ${YELLOW}HAProxy Stats:${NC} http://localhost:8404/stats"
echo -e "  ${YELLOW}PostgreSQL:${NC}    localhost:5432 (user/password)"

echo -e "\n${BLUE}Quick Start Commands:${NC}"
echo -e "  ${YELLOW}View logs:${NC}           make docker-logs"
echo -e "  ${YELLOW}Stop services:${NC}       docker-compose down"
echo -e "  ${YELLOW}Test User Service:${NC}   make grpcurl-user"
echo -e "  ${YELLOW}See all commands:${NC}    make help"

echo -e "\n${BLUE}Example gRPC call (Register User):${NC}"
echo -e "  grpcurl -plaintext \\"
echo -e "    -d '{\"email\":\"test@example.com\",\"password\":\"test123\",\"full_name\":\"Test\",\"phone\":\"555-1234\"}' \\"
echo -e "    localhost:50051 \\"
echo -e "    user.UserService/Register"

echo -e "\n${GREEN}Happy coding! 🚀${NC}\n"
