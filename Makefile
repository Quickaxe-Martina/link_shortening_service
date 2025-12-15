MIGRATIONS_DIR=./migrations
MIGRATE_BIN=$(shell go env GOPATH)/bin/migrate
DB_URL=postgres://test:test@localhost:5432/test?sslmode=disable
GO_IMPORTS_BIN=$(shell go env GOPATH)/bin/goimports

BUILD_VERSION ?= dev
BUILD_DATE    ?= $(shell date +%Y-%m-%dT%H:%M:%S)
BUILD_COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS := -ldflags "\
-X 'main.buildVersion=$(BUILD_VERSION)' \
-X 'main.buildDate=$(BUILD_DATE)' \
-X 'main.buildCommit=$(BUILD_COMMIT)'"

start:
	@$(GO_IMPORTS_BIN) -w .
	@go run $(LDFLAGS) ./cmd/shortener/main.go \
		-d="host=127.0.0.1 port=5469 user=test password=test dbname=test sslmode=disable"

start-memory:
	@$(GO_IMPORTS_BIN) -w .
	@go run $(LDFLAGS) ./cmd/shortener/main.go

start-file:
	@$(GO_IMPORTS_BIN) -w .
	@go run $(LDFLAGS) ./cmd/shortener/main.go -f=storage.json

build:
	@$(GO_IMPORTS_BIN) -w .
	@go build $(LDFLAGS) -o ./cmd/shortener/shortener ./cmd/shortener/main.go

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=seq"; \
		exit 1; \
	fi
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE_BIN) -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1
