# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- `go build -o compartment .` - Build the binary
- `go mod tidy` - Clean up dependencies after adding/removing packages
- `go fmt ./...` - Format all Go code (run after updates)
- `./compartment [flags] <command> <service> [version]` - Run the built binary

### Development Environment
- Go version: 1.24.3 (managed via mise.toml)
- No tests currently exist in the project

### Documentation Website
- Website source: `docs/` directory (GitHub Pages)
- Live site: https://lcmen.github.io/compartment
- Auto-deployed via GitHub Actions on push to main branch
- Built with Evil Martians LaunchKit template

## Architecture Overview

### Core Design Pattern
Compartment uses a factory pattern for service creation with Docker API abstraction:

1. **CLI Entry** (`main.go` → `cmd/compartment.go`): Parses flags and routes commands
2. **Service Factory** (`pkg/service/service.go:31`): Creates service instances based on type
3. **Container Wrapper** (`pkg/container/container.go`): Abstracts Docker client operations with state management
4. **Service Implementations**: Individual service configurations (postgres, redis, devdns)

### Key Architectural Decisions

**Service State Management**: Container states (Running/Stopped/Removed/Error) are tracked at `pkg/container/container.go:15` with corresponding error handling.

**Data Persistence**: Each service gets isolated data directories under `$XDG_DATA_HOME/compartment/{service}/{name}` for volume mounts.

**Configuration Pattern**: Environment-based config initialization at `pkg/configuration/config.go:16` using XDG Base Directory specification.

**Docker API Usage**: Uses specific Docker API imports (`github.com/docker/docker/api/types/container`) rather than general types package, following project's rules.

### Service Implementation Pattern
New services follow this structure:
- Factory function `NewXXXService()` in `pkg/service/xxx.go`
- Default environment variables as package variables
- Volume preparation with automatic directory creation
- Port mapping configuration (if needed)

### Command Flow
1. `parseArgs()` extracts flags (`-n` for name, `-e` for env) and arguments
2. `getServiceForCommand()` creates appropriate service instance or returns nil for check command
3. Service methods (`Start()`, `Stop()`, `Status()`) delegate to container operations
4. Container operations manage Docker API calls with state transitions

### Logging and Debugging
- Uses structured logging via `log/slog`
- Debug mode enabled with `DEBUG=1` environment variable
- Container operations include detailed debug logging for troubleshooting

## Development Guidelines
- Prefer built-in packages over 3rd party libraries unless there is a good reason
- Always run `go mod tidy` after adding or removing dependencies to keep them clean
- Follow TDD - implement tests before writing implementation (though no tests currently exist)
- Use Go's built-in `testing` package for all tests
- Run formatter after updating code to ensure it compiles
- Run tests after updating code to ensure they still pass
- Import from specific Docker API packages (e.g., `docker/docker/api/types/container`) rather than general `docker/docker/api/types`
