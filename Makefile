GOPATH_BIN := $(shell go env GOPATH)/bin
APP_MAIN := cmd/toq_server.go
BIN_DIR := bin
BINARY := $(BIN_DIR)/toq_server
COMPOSE := docker compose -f docker-compose.yml
GOLANGCI_LINT := $(GOPATH_BIN)/golangci-lint

.PHONY: all tidy build build-bin run run-race test fmt vet lint clean swagger docker-build infra-up infra-down infra-logs infra-restart

all: build

tidy:
	go mod tidy

build:
	go build ./...

build-bin:
	@echo "Building toq_server binary to ./$(BIN_DIR) directory"
	@mkdir -p $(BIN_DIR)
	go build -o $(BINARY) ./$(APP_MAIN)

run:
	go run ./$(APP_MAIN)

run-race:
	go run -race ./$(APP_MAIN)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

ci-lint: 
	golangci-lint run

lint: fmt vet ci-lint

clean:
	rm -rf $(BIN_DIR)

# Swagger documentation generation
SWAG := $(GOPATH_BIN)/swag

$(SWAG):
	go install github.com/swaggo/swag/cmd/swag@v1.16.6

swagger: $(SWAG)
	@echo "Generating Swagger docs to ./docs"
	@PATH="$(GOPATH_BIN):$$PATH" swag init -g $(APP_MAIN) -d . -o ./docs --parseDependency --parseInternal

# Docker / infrastructure helpers
docker-build:
	docker build -t toq-server:local .

infra-up:
	$(COMPOSE) up -d mysql redis prometheus grafana otel-collector jaeger swagger-ui

infra-down:
	$(COMPOSE) down

infra-logs:
	$(COMPOSE) logs -f

infra-restart:
	$(COMPOSE) restart
