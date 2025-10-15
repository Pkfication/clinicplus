# ClinicPlus API

ClinicPlus is a Go-based backend API service for managing clinic operations and patient data.

## Features

- RESTful API endpoints for clinic management
- JWT-based authentication
- Database migrations using Goose
- Scheduled tasks using cron
- Environment-based configuration
- **Observability stack**: OpenTelemetry tracing, Prometheus metrics, pprof profiling
- **Monitoring**: Jaeger for distributed tracing, Grafana for visualization

## Prerequisites

- Go 1.23.4 or higher
- PostgreSQL database
- Make (for using Makefile commands)

## Project Structure

```
.
├── cmd/            # Application entry points
├── internal/       # Private application code
├── migrations/     # Database migrations
├── pkg/           # Public library code
├── go.mod         # Go module definition
├── go.sum         # Go module checksums
├── Makefile       # Build automation
└── render.yaml    # Deployment configuration
```

## Getting Started

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd clinicplus
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up environment variables:
   Create a `.env` file in the root directory with the following variables:

   ```
   # Database Configuration
   DATABASE_URL=postgres://username:password@localhost:5432/clinicplus?sslmode=disable

   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-here

   # Server Configuration
   PORT=8080

   # Observability Configuration
   SERVICE_NAME=clinicplus-api
   SERVICE_VERSION=1.0.0
   OTLP_ENDPOINT=http://localhost:4318/v1/traces
   ```

4. Run database migrations:

   ```bash
   make migrate
   ```

5. Build and run the application:
   ```bash
   make run
   ```

## Available Make Commands

- `make build` - Build the application
- `make run` - Build and run the application
- `make migrate` - Run database migrations
- `make test` - Run tests
- `make clean` - Clean up build artifacts
- `make deps` - Download and tidy dependencies
- `make observability-up` - Start observability stack (Jaeger, Prometheus, Grafana)
- `make observability-down` - Stop observability stack
- `make observability-logs` - View observability stack logs

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router and dispatcher
- [jinzhu/gorm](https://github.com/jinzhu/gorm) - ORM library
- [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) - JWT implementation
- [robfig/cron](https://github.com/robfig/cron) - Cron job scheduling
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variable management

## Observability

The application includes comprehensive observability features:

### Tracing

- **OpenTelemetry** integration for distributed tracing
- Traces are sent to Jaeger (default: http://localhost:16686)
- All HTTP requests and database operations are automatically traced

### Metrics

- **Prometheus** metrics exposed at `/metrics`
- HTTP request metrics (count, duration, status codes)
- Database operation metrics (queries, duration, errors)
- Cron job execution metrics
- Runtime metrics (Go GC, goroutines, memory)

### Profiling

- **pprof** endpoints available at `/debug/pprof/`
- CPU, memory, goroutine, and mutex profiling
- Access via: `go tool pprof http://localhost:8080/debug/pprof/profile`

### Monitoring Stack

Start the complete observability stack:

```bash
make observability-up
```

This starts:

- **Jaeger** (http://localhost:16686) - Distributed tracing
- **Prometheus** (http://localhost:9090) - Metrics collection
- **Grafana** (http://localhost:3000) - Visualization (admin/admin)
- **Loki** (http://localhost:3100) - Log aggregation

### Environment Variables

Configure observability endpoints:

```bash
SERVICE_NAME=clinicplus-api
SERVICE_VERSION=1.0.0
OTLP_ENDPOINT=http://localhost:4318/v1/traces
```

## API Documentation

API documentation is available at `/docs` when running the server locally.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
