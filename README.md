# Kratos Layout

A Go microservice project template based on [Kratos](https://go-kratos.dev/) framework with DDD/Onion architecture.

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make

### Setup

```bash
# Install dependencies and tools
make init

# Start development environment (MySQL, Redis, Nacos, Apollo)
make dev-up

# Build and run
make build
./bin/server -conf ./configs
```

## Development

### Common Commands

```bash
make init        # Install protoc plugins and tools
make api         # Generate API code (pb.go, http, grpc, openapi)
make config      # Generate internal config proto
make generate    # Run go generate and tidy
make all         # Generate all (api + config + generate)
make build       # Build binary
make test        # Run tests
make lint        # Run linter
make coverage    # Run tests with coverage check
make dev-up      # Start dev environment
make dev-down    # Stop dev environment
```

### Wire (Dependency Injection)

```bash
cd cmd/server && wire
```

### Run Single Test

```bash
go test -v ./pkg/log/... -run TestInitJSONLogger
```

## Architecture

```
API (Proto) → Server (HTTP/gRPC) → Service → Biz → Data
```

### Project Structure

```
├── api/                 # Protocol Buffer definitions
├── cmd/server/          # Application entry point
├── configs/             # Runtime configuration
├── internal/
│   ├── biz/            # Business logic layer
│   ├── conf/           # Config proto definitions
│   ├── data/           # Data access layer
│   ├── server/         # HTTP/gRPC server setup
│   └── service/        # Service layer (implements proto)
├── pkg/                 # Public packages
└── scripts/             # Development scripts
```

## Development Environment

Docker Compose provides MySQL, Redis, Nacos, and Apollo for local development.

See [docs/dev-environment.md](docs/dev-environment.md) for details.

## Docker

```bash
# Build
docker build -t <image-name> .

# Run
docker run --rm -p 8000:8000 -p 9000:9000 -v /path/to/configs:/data/conf <image-name>
```

## References

- [Kratos Documentation](https://go-kratos.dev/)
- [Kratos GitHub](https://github.com/go-kratos/kratos)
