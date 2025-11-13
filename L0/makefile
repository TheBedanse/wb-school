GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=app
BINARY_UNIX=$(BINARY_NAME)_unix

CMD_DIR=cmd
APP_DIR=$(CMD_DIR)/app
MIGRATE_DIR=$(CMD_DIR)/migrate
GENERATOR_DIR=$(CMD_DIR)/generator
INTERNAL_DIR=internal

all: test build

build: build-app build-generator

build-app:
	@echo "Building application..."
	$(GOBUILD) -o bin/app $(APP_DIR)/main.go

build-generator:
	@echo "Building order generator..."
	$(GOBUILD) -o bin/generator $(GENERATOR_DIR)/main.go

build-migrate:
	@echo "Building migration tool..."
	$(GOBUILD) -o bin/migrate $(MIGRATE_DIR)/main.go

test:
	@echo "Running tests..."
	$(GOTEST) ./$(INTERNAL_DIR)/models
	$(GOTEST) ./$(INTERNAL_DIR)/cache
	$(GOTEST) ./$(INTERNAL_DIR)/service
	$(GOTEST) ./$(INTERNAL_DIR)/handler

test-verbose:
	@echo "Running verbose tests..."
	$(GOTEST) -v ./$(INTERNAL_DIR)/models
	$(GOTEST) -v ./$(INTERNAL_DIR)/cache
	$(GOTEST) -v ./$(INTERNAL_DIR)/service
	$(GOTEST) -v ./$(INTERNAL_DIR)/handler

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -cover ./$(INTERNAL_DIR)/models
	$(GOTEST) -cover ./$(INTERNAL_DIR)/cache
	$(GOTEST) -cover ./$(INTERNAL_DIR)/service
	$(GOTEST) -cover ./$(INTERNAL_DIR)/handler

docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting services..."
	docker-compose up

docker-down:
	@echo "Stopping services..."
	docker-compose down

docker-restart:
	@echo "Restarting services..."
	docker-compose restart

docker-logs:
	@echo "Showing logs..."
	docker-compose logs -f


migrate-up: build-migrate
	@echo "Applying migrations..."
	./bin/migrate up

migrate-down: build-migrate
	@echo "Rolling back migrations..."
	./bin/migrate down

migrate-status: build-migrate
	@echo "Migration status:"
	./bin/migrate status

run-app: build-app
	@echo "Running application..."
	./bin/app

generate-mocks:
	@echo "Generating mocks..."
	$(GOCMD) generate ./...

.PHONY: all build build-app build-generator build-migrate test test-verbose test-coverage \
        docker-build docker-up docker-up-detached docker-down docker-restart docker-logs docker-clean \
        migrate-up migrate-down migrate-status run-app run-generator dev clean generate-mocks help