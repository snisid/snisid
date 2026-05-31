# SNISID — Master Makefile
.PHONY: all build test lint run clean docker-up docker-down help

APP_NAME := snisid-platform
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "*/node_modules/*")

all: lint test build

# --- Build ---
build: ## Compile tous les binaires Go
	go build -o bin/gateway ./cmd/gateway
	go build -o bin/identity-api ./services/identity-api/cmd
	go build -o bin/fraud-engine ./services/fraud-engine/cmd
	go build -o bin/verification-api ./services/verification-api/cmd
	go build -o bin/audit-service ./services/audit-service/cmd
	go build -o bin/nexus-server ./services/nexus/cmd/nexus-server
	go build -o bin/ws-gateway ./services/ws-gateway
	@echo "Build complete: binaries in bin/"

build-all: ## Compile tous les services Go
	go build ./internal/...
	go build ./services/...

# --- Tests ---
test: ## Exécute tous les tests Go
	go test ./internal/... -v -count=1 -cover

test-short: ## Tests rapides (sans intégration)
	go test ./internal/... -v -count=1 -short -cover

test-race: ## Tests avec race detector
	go test ./internal/... -v -count=1 -race -coverprofile=coverage.out

test-integration: ## Tests d'intégration
	go test ./internal/... -v -count=1 -tags=integration ./...

coverage: test-race ## Affiche le coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# --- Lint ---
lint: ## Analyse statique du code
	staticcheck ./...
	golangci-lint run ./... --timeout=5m 2>/dev/null || true

vet: ## go vet
	go vet ./...

# --- Exécution ---
run-gateway: ## Lance le gateway
	go run ./cmd/gateway

run-identity-api: ## Lance l'API d'identité
	go run ./services/identity-api/cmd

run-fraud-engine: ## Lance le moteur de fraude
	go run ./services/fraud-engine/cmd

# --- Docker ---
docker-up: ## Démarre les services avec Docker Compose
	docker-compose -f docker-compose.yml up -d

docker-down: ## Arrête les services Docker
	docker-compose -f docker-compose.yml down

docker-logs: ## Logs des services Docker
	docker-compose -f docker-compose.yml logs -f

docker-build: ## Construit les images Docker
	docker-compose -f docker-compose.yml build

# --- Base de données ---
db-migrate: ## Applique les migrations
	go run ./scripts/migrations/main.go

db-rollback: ## Annule la dernière migration
	go run ./scripts/migrations/main.go --down

# --- Nettoyage ---
clean: ## Nettoie les artefacts de build
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache

# --- Dépendances ---
deps: ## Télécharge les dépendances Go
	go mod download
	go mod tidy

deps-update: ## Met à jour les dépendances
	go get -u ./...
	go mod tidy

# --- Utilitaires ---
help: ## Affiche cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
