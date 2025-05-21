# ClinicPlus API

ClinicPlus is a Go-based backend API service for managing clinic operations and patient data.

## Features

- RESTful API endpoints for clinic management
- JWT-based authentication
- Database migrations using Goose
- Scheduled tasks using cron
- Environment-based configuration

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
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=clinicplus
   JWT_SECRET=your_jwt_secret
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

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router and dispatcher
- [jinzhu/gorm](https://github.com/jinzhu/gorm) - ORM library
- [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) - JWT implementation
- [robfig/cron](https://github.com/robfig/cron) - Cron job scheduling
- [joho/godotenv](https://github.com/joho/godotenv) - Environment variable management

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
