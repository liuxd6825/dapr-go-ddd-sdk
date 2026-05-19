# AGENTS.md

## Build & Run

```bash
# Build binary for current OS/arch
make build

# Build for Linux amd64 with version
make build GOOS=linux GOARCH=amd64 REL_VERSION=v1.15-251027

# Output: dist/<OS>_<ARCH>/master
```

## Project Structure

- `app/` - Application domains
  - `domain/master/` - Core business (accounts, records, transactions)
  - `domain/graph/` - Graph database operations (nodes, relations)
  - `domain/tag/` - Tag management
  - `domain/rag/` - Retrieval-Augmented Generation
  - `domain/import/` - Excel file import
  - `domain/analysis/` - Suspicious transaction analysis
  - `domain/sys/` - System services (portal, code, bank, currency)
  - `cmd/` - Binary entry points (master, worker, ai-mcp)
- `pkg/` - Shared packages
  - `ddd/` - DDD core (Aggregate, Command, Event, Entity, Event Sourcing, Repository)
  - `core/dapr/` - Dapr client wrapper (HTTP/gRPC, Pub/Sub, Actor, Event Store)
  - `core/restapp/` - REST application framework
  - `db/` - Database abstractions (MongoDB, SQL, RSQL)
  - `micro/` - Microservice utilities (Feign, State)
  - `lowcode/` - Low-code engine

## Key Conventions

- **DDD Layer Order**: Application → Domain → Infrastructure
- **Event Sourcing**: Use `pkg/ddd/event_store.go` for event sourcing with snapshots
- **Multi-tenancy**: All entities have `TenantId` field; use `appctx.NewTenantContext()`
- **Service Communication**: Use `pkg/micro/micro_feign` for inter-service HTTP calls
- **RestSQL**: Custom query language in `pkg/db/rsql/` for complex queries

## Local Dependencies

This project uses `replace` directives in `go.mod` for local development:

```
github.com/dapr/dapr => ../dapr
github.com/dapr/go-sdk => ../dapr-go-sdk
github.com/dapr/components-contrib => ../dapr-components-contrib
github.com/liuxd6825/jsonschema/v6 => ../../jsonschema
github.com/liuxd6825/k6server => ../../k6server
gorm.io/gorm => ../../gorm
```

When building outside the expected directory structure, these replacements may cause build failures.

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/ddd/...
```

## Entry Points

- `app/cmd/master/main.go` - Main application server (registers all domain APIs)
- `app/cmd/worker/main.go` - Background worker process
- `app/cmd/ai-mcp/main.go` - AI MCP server

## Framework Quirks

- Uses Iris web framework (`github.com/kataras/iris/v12`)
- Actor model via Dapr (`pkg/core/dapr/actor/`)
- Event registry pattern in `pkg/ddd/event_registry.go`
- Chain-style DAO pattern in `pkg/db/dao/`