# Liquidity Finder

A Go-based application for finding and analyzing liquidity pools across different decentralized exchanges (DEXs).

## Features

- Gather trading pairs from various DEXs
- Detect tax tokens
- Monitor mempool transactions
- Support for multiple networks (PulseChain, Ethereum, BNB Chain, Polygon)
- High liquidity pair detection
- Real-time reserve updates

## Prerequisites

- Go 1.22 or higher
- Python 3.x (for tax token detection)
- PostgreSQL database
- Redis cache (optional)
- Docker (optional)

## Environment Variables

Copy `.env.example` to `.env` and configure the following variables:

```env
# Cache Configuration
CACHE_HOST=localhost
CACHE_PORT=6379
CACHE_PASSWORD=
CACHE_DB=0

# Database Configuration
USE_LOCAL_DB=true
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=<your-password>
POSTGRES_DB=postgres
POSTGRES_SSL_MODE=disable

# API Keys
MORALIS_API_KEY=<your-moralis-api-key>
MORALIS_BASE_URL=<moralis-base-url>

# Blockchain Configuration
BASE_DIR=${workspaceFolder}
FORGE_DIR=<your-forge-directory>

# Execution Configuration
SHOULD_EXECUTE_ARBS=true
```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/liquidityFinder.git
cd liquidityFinder
```

2. Install Go dependencies:

```bash
go mod download
```

3. Set up Python environment for tax token detection:

```bash
cd apps/python/parsebytecode
./setup.sh
```

4. Set up the database:

```bash
psql -U postgres -d postgres -f docs/db.sql
```

## Running the Application

The application supports multiple commands through the `cmd/start.go` entry point:

### Available Commands

1. **Gather Pairs**

```bash
go run cmd/start.go gatherPairs
```

2. **Path Service**

```bash
go run cmd/start.go pathService
```

3. **Update Pairs**

```bash
go run cmd/start.go updatePairs
```

4. **Update Reserves and High Liquidity**

```bash
go run cmd/start.go updateReservesAndHighLiquidity
```

5. **Detect Tax Tokens**

```bash
go run cmd/start.go detectTaxToken
```

6. **Listen to Mempool**

```bash
go run cmd/start.go listenMempool
```

### Using Docker

Build and run using Docker:

```bash
docker build -t liquidity-finder .
docker run -e POSTGRES_HOST='host.docker.internal' liquidity-finder /app/cmd/start gatherPairs
```

## Development

### Project Structure

- `/cmd` - Main application entry points
- `/internal` - Internal packages
  - `/blockchain` - Blockchain interaction logic
  - `/database` - Database models and helpers
  - `/service` - Business logic services
  - `/myConfig` - Configuration management
- `/abi` - Smart contract ABIs
- `/apps` - Additional applications
  - `/python` - Python scripts for tax token detection
- `/docs` - Documentation and database schemas

### Adding New Features

1. Add new models in `/internal/database/model`
2. Implement service logic in `/internal/service`
3. Add new commands in `cmd/start.go`
4. Update configuration in `internal/myConfig`

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Add your license here]

## Security

- All sensitive information should be stored in environment variables
- Never commit API keys or passwords
- Use secure connections for database and cache
- Rotate API keys regularly
- Monitor for suspicious transactions

## Support

[Add support information here]
