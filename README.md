# Prodyo Backend API

An application to manager/control productivity in the software development environment

## Features

- **Project Management**: Create, read, update, and delete projects with team members
- **Iteration Tracking**: Manage development iterations with tasks and metrics
- **Productivity Metrics**: Track speed, rework, and instability indicators
- **Quality Tracking**: Monitor bugs and improvements per task
- **Action Planning**: Create causes and actions based on productivity analysis
- **Automatic Migrations**: Database migrations run automatically on startup
- **Docker Support**: Easy deployment with Docker Compose

## Documentation

- [Architecture & Database Schema](ARCHITECTURE.md) - Domain model and entity relationships
- [API Documentation](docs/) - Swagger/OpenAPI documentation

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Docker & Docker Compose (optional)

### Running with Docker

```bash
docker-compose up
```

The API will be available at `http://localhost:8080`

### Running Locally

1. Set up PostgreSQL database
2. Configure environment variables (or use defaults)
3. Run the application:

```bash
go run cmd/api/main.go
```

## Project Structure

```
prodyo-backend/
├── cmd/
│   ├── api/                    # Main API entry point
│   ├── fix-migration/          # Migration utilities
│   └── internal/
│       ├── config/             # Configuration
│       ├── handlers/           # HTTP handlers
│       ├── migrations/         # Database migrations
│       ├── models/             # Domain models
│       ├── repositories/       # Data access layer
│       ├── services/           # Business services
│       └── usecases/           # Use case implementations
├── docs/                       # API documentation
├── ARCHITECTURE.md             # Architecture documentation
├── docker-compose.yml
└── Dockerfile
```

## Testing

Use the provided test script to test the API:

```bash
./test_api.sh
```

Make sure you have `jq` installed for JSON formatting:
```bash
# On Ubuntu/Debian
sudo apt-get install jq

# On macOS
brew install jq

# On Windows (with Chocolatey)
choco install jq
```

## License

MIT
