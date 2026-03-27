# SIG-Agro Architecture Documentation

## System Overview

SIG-Agro is a distributed microservices-based Geographic Information System for agricultural management. It follows Clean Architecture principles and uses modern technologies for scalability and maintainability.

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         Frontend (Next.js)                      │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                    HAProxy API Gateway                          │
│  • gRPC routing (TCP mode)                                      │
│  • JWT validation (prepared)                                    │
│  • SSL/TLS termination                                          │
└─────────────────────────────────────────────────────────────────┘
         ↓           ↓           ↓           ↓           ↓
    ┌─────────┬──────────┬─────────┬──────────┬──────────┐
    │  User   │ Producer │ Parcel  │Production│ Alert    │
    │ Service │ Service  │ Service │ Service  │ Service  │
    └─────────┴──────────┴─────────┴──────────┴──────────┘
         ↓           ↓           ↓
    ┌────────────────────────────────────────────┐
    │      PostgreSQL with PostGIS               │
    │  • Users DB          • Producers DB         │
    │  • Parcels DB        • Production DB        │
    │  • Alerts DB         • Notifications DB     │
    │  • Reports DB                              │
    └────────────────────────────────────────────┘
         ↓
    ┌────────────────────────────────────────────┐
    │           RabbitMQ Message Broker          │
    │  • Async event publishing                  │
    │  • Service-to-service communication        │
    └────────────────────────────────────────────┘
```

## Clean Architecture Layers

Each microservice is structured in 4 layers:

### 1. Domain Layer (`internal/domain/`)

**Responsibility:** Business entities and models

**Files:**
- `*.go` - Domain models (User, Producer, Parcel, etc.)

**Example:**
```go
type User struct {
    ID        int64
    Email     string
    FullName  string
    Roles     []string
    CreatedAt int64
}
```

### 2. Repository Layer (`internal/repository/`)

**Responsibility:** Data access abstraction

**Files:**
- `repository.go` - Database operations
- Implements Create, Read, Update, Delete operations
- Uses type-safe queries from SQLC

**Characteristics:**
- Database connection pooling
- Error handling
- Transaction support (where needed)

### 3. Use Case Layer (`internal/usecase/`)

**Responsibility:** Business logic and orchestration

**Files:**
- `usecase.go` - Business logic
- Coordinates domain models and repositories
- Implements business rules

### 4. Handler/Adapter Layer (`internal/handler/`)

**Responsibility:** gRPC API endpoints

**Files:**
- `handler.go` - gRPC service implementation
- Implements Protocol Buffer service definitions
- Converts gRPC messages to domain models

## Microservices

### 1. User Service (Port 50051)

**Responsibilities:**
- User registration and authentication
- JWT token management
- Role-based access control

**Key Operations:**
- `Register`: Create new user
- `Login`: Authenticate and get JWT
- `ValidateToken`: Verify token validity
- `GetUser`: Retrieve user information
- `ListUsers`: List all users with pagination

**Database:** `sig_agro_users`

**Tables:**
```sql
- users (id, email, password_hash, full_name, phone, created_at)
- user_roles (id, user_id, role, created_at)
- tokens (id, user_id, token_hash, expires_at, created_at)
```

### 2. Producer Service (Port 50052)

**Responsibilities:**
- Manage agricultural producers/farmers
- Producer profile management

**Key Operations:**
- `CreateProducer`: Register new producer
- `GetProducer`: Retrieve producer by ID
- `ListProducers`: List user's producers
- `UpdateProducer`: Update producer info
- `DeleteProducer`: Remove producer

**Database:** `sig_agro_producers`

**Tables:**
```sql
- producers (id, user_id, name, document_id, phone, email, address, created_at)
```

### 3. Parcel Service (Port 50053)

**Responsibilities:**
- Geospatial data management using PostGIS
- Parcel boundary storage and querying
- Spatial queries

**Key Operations:**
- `CreateParcel`: Create parcel with WKT geometry
- `GetParcel`: Retrieve parcel details
- `ListParcels`: List producer's parcels
- `UpdateParcel`: Update parcel information
- `QueryByGeometry`: Spatial query (intersection, contains, etc.)

**Database:** `sig_agro_parcels`

**Key Features:**
- PostGIS GEOMETRY(POLYGON, 4326) column
- Spatial indexing with GIST
- WKT input/output support

**Example WKT Geometries:**
```
POLYGON((-3.5 40.5, -3.4 40.5, -3.4 40.6, -3.5 40.6, -3.5 40.5))
LINESTRING(-3.5 40.5, -3.4 40.6)
POINT(-3.45 40.55)
```

### 4. Production Service (Port 50054)

**Responsibilities:**
- Record agricultural activities
- Activity tracking and logging

**Key Operations:**
- `RecordActivity`: Log production activity
- `GetActivity`: Retrieve activity details
- `ListActivities`: List parcel activities

**Database:** `sig_agro_production`

### 5. Alert Service (Port 50055)

**Responsibilities:**
- Climate and sanitary alert management
- Alert creation and evaluation

**Key Operations:**
- `CreateAlert`: Create new alert
- `GetAlert`: Retrieve alert details
- `ListAlerts`: List alerts with filtering
- `EvaluateAlerts`: Evaluate alert conditions

**Database:** `sig_agro_alerts`

### 6. Notification Service (Port 50056)

**Responsibilities:**
- Send notifications (push, SMS, email)
- Notification tracking

**Key Operations:**
- `SendNotification`: Send notification
- `GetNotification`: Retrieve notification
- `ListNotifications`: List user's notifications
- `MarkAsRead`: Mark notification as read

**Database:** `sig_agro_notifications`

### 7. Report Service (Port 50057)

**Responsibilities:**
- Generate reports and dashboards
- Data aggregation and analysis

**Key Operations:**
- `GenerateReport`: Create new report
- `GetReport`: Retrieve report
- `ListReports`: List producer's reports

**Database:** `sig_agro_reports`

## Database Design

### PostgreSQL Setup

- **Version:** 16+
- **Extension:** PostGIS 3.4+
- **Topology:** One database per service (micro-database pattern)

### Database Separation Strategy

Each service has its own database to:
- Ensure loose coupling
- Allow independent scaling
- Simplify backups and recovery
- Enable different schemas per service need

**Connection String Format:**
```
postgres://user:password@postgres:5432/sig_agro_[service_name]?sslmode=disable
```

### Key Tables

**Users (User Service):**
- Contains authentication data
- Password stored as SHA256 hash (improve with bcrypt in production)
- JWT tokens tracked separately

**Producers (Producer Service):**
- Links to User Service by user_id
- Represents agricultural entities

**Parcels (Parcel Service):**
- PostGIS GEOMETRY column for coordinates
- GIST index for spatial queries
- WKT format for input/output

**Activities, Alerts, Notifications, Reports:**
- Temporal data with timestamps
- JSON metadata for flexible attributes
- Indexes on common queries

## Communication Patterns

### Synchronous (gRPC)

All service-to-service communication currently uses gRPC:

```
Client A → gRPC Call → Server B
         (TCP port 50051-50057)
```

**Advantages:**
- Type-safe with Protocol Buffers
- Efficient binary protocol
- Streaming support
- Built-in error handling

### Asynchronous (RabbitMQ - Prepared)

Future event-driven architecture:

```
Service A → Emit Event → RabbitMQ → Service B
         (amqp://rabbitmq:5672)
```

**Prepared Queues:**
- `parcel.created` - When new parcel is registered
- `activity.recorded` - When activity is logged
- `alert.triggered` - When alert condition is met
- `notification.sent` - When notification is dispatched

## API Gateway (HAProxy)

### Configuration

**File:** `infrastructure/haproxy/haproxy.cfg`

**Ports:**
- 80 (HTTP)
- 443 (HTTPS with SSL)
- 8404 (Stats)

### Features

1. **TCP Load Balancing**
   - Routes gRPC traffic based on service path
   - Example: `/user.UserService/*` → User Service

2. **SSL/TLS Termination**
   - Self-signed certificate (development)
   - Listens on port 443

3. **Health Checks**
   - Inter-service checks every 10 seconds
   - Automatic failover

4. **Statistics Dashboard**
   - Available at `http://localhost:8404/stats`
   - Real-time connection and request metrics

### Routing Rules

```
/user.UserService/* → 50051
/producer.ProducerService/* → 50052
/parcel.ParcelService/* → 50053
/production.ProductionService/* → 50054
/alert.AlertService/* → 50055
/notification.NotificationService/* → 50056
/report.ReportService/* → 50057
```

### JWT Validation (Prepared)

**Current Status:** ACL rules prepared, implementation pending

**Planned Implementation:**
- Extract JWT token from request header
- Validate signature with public key (RS256)
- Inject `X-User-Id` and `X-User-Roles` headers
- Forward authenticated requests to backends

## Security Architecture

### Authentication Flow

1. Client calls `User Service/Login` with credentials
2. User Service validates and generates JWT
3. JWT returned to client (valid for 1 hour)
4. Client includes JWT in subsequent requests
5. HAProxy validates token (future)
6. Services can verify token with `User Service/ValidateToken`

### Password Security

**Current:** SHA256 hashing
**TODO:** Implement bcrypt or Argon2

### Token Security

**Type:** JWT (RS256 prepared)
**Payload:**
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "roles": ["user", "producer"],
  "exp": 1234567890
}
```

### Access Control

**Role-based (RBAC):**
- `user` - Basic user access
- `producer` - Producer data access
- `admin` - Full system access

## Deployment Architecture

### Docker Compose Structure

**Services:**
1. PostgreSQL - Data persistence
2. RabbitMQ - Message broker
3. 7 Microservices - Business logic
4. HAProxy - API Gateway

**Networks:**
- `sig-agro-network` - Bridge network for service communication
- Isolated from host except HAProxy ports

**Volumes:**
- `postgres_data` - Database persistence
- `rabbitmq_data` - Queue persistence

### Container Health Checks

Each service includes:
- TCP probes for database connectivity
- Startup delays to ensure dependencies ready
- Restart policies (unless stopped)

**Health Check Command:**
```bash
docker-compose ps  # Check service status
docker-compose logs [service] # View logs
```

## Development Workflow

### Local Development Setup

**Option 1: Full Docker Compose**
```bash
docker-compose up -d
# Services accessible on localhost:50051-50057
```

**Option 2: Database Only**
```bash
docker-compose -f docker-compose.dev.yml up -d
# Run services locally on different ports
```

### Code Generation

1. **Protocol Buffers**
   ```bash
   make proto
   # Generates .pb.go and *_grpc.pb.go files
   ```

2. **SQLC**
   ```bash
   make sqlc
   # Generates type-safe database functions
   ```

3. **Build**
   ```bash
   make build
   # Compiles all services
   ```

## Monitoring and Observability

### Current Logging

- All services log to stdout
- Docker Compose aggregates logs
- View with: `docker-compose logs`

### Metrics (Future)

**Planned Integration:**
- Prometheus - Metrics collection
- Grafana - Visualization
- Jaeger - Distributed tracing

### Health Endpoints (Future)

```
GET /health
GET /ready
GET /metrics
```

## Performance Considerations

### Database Optimization

1. **Indexes:**
   - Email on users table
   - Producer foreign keys
   - Parcel spatial indexes (GIST)
   - Timestamp indexes for range queries

2. **Connection Pooling:**
   - Configured in database drivers
   - Visible in service health checks

3. **Batch Operations:**
   - ListUsers/ListProducers with pagination
   - Offset-limit pattern

### gRPC Optimization

1. **Protocol Buffers:**
   - Binary serialization (efficient)
   - Protobuf v3 syntax

2. **Connection Multiplexing:**
   - gRPC uses HTTP/2
   - Multiple streams over single connection

3. **Message Compression:**
   - Supported natively
   - Enabled per-call or per-channel

## Scaling Strategies

### Horizontal Scaling

**Current State:** Ready for Docker/Kubernetes orchestration

**Future Preparation:**
```yaml
# Example Kubernetes-style scaling
replicas: 3
resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "512Mi"
```

### Database Scaling

**Read Replicas:**
- PostgreSQL streaming replication
- Connection pooling with PgBouncer

**Sharding:**
- Partition by user_id or producer_id
- Future enhancement

## Testing Strategy

### Unit Tests

**Location:** `*_test.go` files
**Coverage:** Domain logic, repository queries

### Integration Tests

**Location:** `tests/` directory
**Coverage:** Service handlers, database operations

### End-to-End Tests

**Tool:** grpcurl
**Coverage:** Full API workflows

## Deployment Checklist

- [ ] Generate proper SSL certificates
- [ ] Change default passwords (PostgreSQL, RabbitMQ)
- [ ] Configure JWT secret (random, long key)
- [ ] Enable bcrypt for passwords
- [ ] Set resource limits
- [ ] Configure backup strategy
- [ ] Set up monitoring
- [ ] Configure auto-scaling policies
- [ ] Enable request logging
- [ ] Set up alerting

## File Structure Summary

```
.
├── api/
│   └── proto/              (Protocol Buffer definitions)
├── services/
│   └── [service-name]/
│       ├── cmd/server/     (Entry point)
│       ├── internal/       (Private code)
│       ├── db/             (Migrations + SQLC)
│       ├── Dockerfile      (Multi-stage build)
│       └── go.mod
├── infrastructure/
│   ├── haproxy/           (API Gateway config)
│   └── postgres/          (Database init script)
├── docker-compose.yml     (Production-like)
├── docker-compose.dev.yml (Development)
├── Makefile               (Automation)
├── setup.sh/setup.bat     (Quick setup)
└── README.md
```

## References

- [Protocol Buffers Documentation](https://developers.google.com/protocol-buffers)
- [gRPC Documentation](https://grpc.io)
- [PostGIS Documentation](https://postgis.net)
- [PostgreSQL Documentation](https://www.postgresql.org/docs)
- [SQLC Documentation](https://sqlc.dev)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [gRPC Best Practices](https://grpc.io/docs/guides/performance-best-practices)
