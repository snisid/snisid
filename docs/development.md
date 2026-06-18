# SNISID Development Guide

## Prerequisites

- Go 1.21+ (backend services)
- Python 3.9+ (sync agent)
- Node.js 18+ (frontend)
- WiX Toolset v3.11+ (MSI builds)
- NSIS 3.0+ (installer builds)
- Windows ADK (WinPE ISO builds)

## Repository Layout

```
snisid/
├── cmd/              # Entry points (snisid-svc, snisid-api)
├── internal/         # Core libraries
├── frontend/         # Web UI source
├── tools/            # Offline CLI tools
├── build/            # Compiled output
├── deploy/           # Deployment scripts
├── installer/        # NSIS installer
├── msi/              # WiX MSI project
├── sync/             # Sync agent
├── bootable/         # WinPE ISO builder
└── docs/             # Documentation
```

## Quick Start

```powershell
# Build all components
cd cmd/snisid-svc && go build -o ../../build/bin/
cd cmd/snisid-api && go build -o ../../build/api/

# Build frontend
cd frontend && npm install && npm run build

# Build installer
cd installer && makensis setup.nsi
```

## Code Standards

- Go: follow `go fmt` and `go vet`; use `net/http` or Gin for HTTP
- Python: PEP 8; type hints required for public functions
- Frontend: ESLint + Prettier; components in `src/components/`
- All new endpoints require both unit and integration tests

## Testing

```powershell
go test ./...
python -m pytest sync/tests/
cd frontend && npm run test
```

## Building an MSI

```powershell
cd msi
.\build.ps1
```

## Building a WinPE ISO

```powershell
cd bootable
.\create-iso.ps1
```
