# Redis Monitoring System

A Go-based monitoring system for Redis instances, designed to work with Magento environments. This tool provides real-time metrics collection, performance monitoring, and alerting capabilities for Redis servers.

## Features

- Real-time Redis metrics collection
    - Memory usage statistics
    - Connection monitoring
    - Cache hit/miss rates
- HTTP API for metrics retrieval
- Concurrent monitoring of multiple Redis instances
- Configurable alerting system

## Prerequisites

- Go 1.21 or higher
- Redis server (tested with Redis 6.x, 7.x)

## Installation

```bash
git clone https://github.com/yourusername/redis-monitor
cd redis-monitor
go mod tidy
```

## Project Structure

```
redis-monitor/
├── cmd/
│   └── monitor/
│       └── main.go         # Application entry point
├── internal/
│   ├── collector/          # Redis metrics collection
│   ├── models/            # Data structures
│   └── server/            # HTTP server implementation
├── config/                # Configuration files
└── tests/                # Test files
```

## Getting Started

1. Configure Redis connection in `config/config.yaml`:
```yaml
redis:
  host: localhost
  port: 6379
```

2. Run the monitoring system:
```bash
go run cmd/monitor/main.go
```

3. Access metrics via HTTP endpoints:
```
GET /metrics/memory       # Memory statistics
GET /metrics/cache       # Cache hit/miss rates
GET /metrics/connections # Connection information
```

## Development

To start development:

1. Fork the repository
2. Create a feature branch
3. Submit a pull request

## Testing

Run tests with:
```bash
go test ./...
```

## License

MIT License - see LICENSE file for details