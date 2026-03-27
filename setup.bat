@echo off
setlocal enabledelayedexpansion

echo.
echo ========================================================
echo      SIG-Agro: Sistema de Informacion Geografica
echo                   Setup Script
echo ========================================================
echo.

REM Check prerequisites
echo Checking prerequisites...

docker.exe --version >nul 2>&1
if errorlevel 1 (
    echo Docker is not installed. Please install Docker first.
    exit /b 1
)

docker-compose.exe --version >nul 2>&1
if errorlevel 1 (
    echo Docker Compose is not installed. Please install Docker Compose first.
    exit /b 1
)

echo [OK] Docker and Docker Compose are installed
echo.

REM Generate self-signed certificate
echo Generating self-signed certificate for HAProxy...
if not exist "infrastructure\haproxy\certs" mkdir infrastructure\haproxy\certs

if exist "infrastructure\haproxy\certs\selfsigned.pem" (
    echo Certificate already exists, skipping generation
) else (
    REM Using a different approach for Windows without OpenSSL
    echo Creating certificate directories...
    echo.
    echo NOTE: For production, generate a proper certificate using OpenSSL or other tools.
    echo For development, you can use: https://github.com/FiloSottile/mkcert
    echo.
    
    REM Create a dummy certificate for demonstration
    REM In production, use: openssl or mkcert
    echo [WARNING] Skipping certificate generation. 
    echo Please generate it manually using:
    echo   openssl req -x509 -newkey rsa:2048 -keyout infrastructure/haproxy/certs/selfsigned.key -out infrastructure/haproxy/certs/selfsigned.crt -days 365 -nodes -subj "/C=ES/ST=Madrid/L=Madrid/O=SIG-Agro/CN=localhost"
    echo Or use mkcert: https://github.com/FiloSottile/mkcert
)

echo.
echo Building Docker images...
docker-compose build

if errorlevel 1 (
    echo Failed to build Docker images.
    exit /b 1
)

echo [OK] Docker images built successfully
echo.

echo Starting all services...
docker-compose up -d

if errorlevel 1 (
    echo Failed to start services.
    exit /b 1
)

echo [OK] Services started
echo.

REM Wait for services
echo Waiting for services to be ready (15 seconds)...
timeout /t 15 /nobreak

echo.
echo Checking service health...
docker-compose ps

echo.
echo ========================================================
echo         SIG-Agro is now running! ^^!
echo ========================================================
echo.

echo Service Endpoints:
echo   User Service:         grpc://localhost:50051
echo   Producer Service:     grpc://localhost:50052
echo   Parcel Service:       grpc://localhost:50053
echo   Production Service:   grpc://localhost:50054
echo   Alert Service:        grpc://localhost:50055
echo   Notification Service: grpc://localhost:50056
echo   Report Service:       grpc://localhost:50057
echo.

echo API Gateway:
echo   HTTP: http://localhost:80
echo   GRPC: grpc://localhost:443
echo.

echo Management Interfaces:
echo   RabbitMQ:     http://localhost:15672 (user/password)
echo   HAProxy Stats: http://localhost:8404/stats
echo   PostgreSQL:    localhost:5432 (user/password)
echo.

echo Quick Start Commands:
echo   View logs:         docker-compose logs -f
echo   Stop services:     docker-compose down
echo   See services:      docker-compose ps
echo.

echo Example gRPC call (from PowerShell):
echo   $body = @{email="test@example.com"; password="test123"; full_name="Test"; phone="555-1234"} ^| ConvertTo-Json
echo   grpcurl -plaintext -d $body localhost:50051 user.UserService/Register
echo.

echo Happy coding! ^^_^/
echo.
pause
