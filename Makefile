.PHONY: generate generate-types generate-server clean-gen
.PHONY: test test-unit test-integration test-e2e test-all test-coverage
.PHONY: test-setup test-teardown test-db-up test-db-down
.PHONY: run

generate: clean-gen generate-types generate-server

generate-types:
	@mkdir -p internal/api/generated
	oapi-codegen \
		-package generated \
		-generate types \
		-o internal/api/generated/types.gen.go \
		spec/openapi.yml

generate-server:
	oapi-codegen \
		-package generated \
		-generate server,strict-server \
		-o internal/api/generated/server.gen.go \
		spec/openapi.yml

clean-gen:
	@rm -rf internal/api/generated/*.gen.go

# Testing

# run all tests
test-all: test-unit test-integration test-e2e

# run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v -race -count=1 ./pkg/... ./internal/domain/... ./internal/helpers/... ./internal/api/... ./internal/service/...

# run integration tests (requires database)
test-integration: test-db-up
	@echo "Waiting for test database to be ready..."
	@sleep 3
	@echo "Running integration tests..."
	TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/test_db?sslmode=disable" \
		go test -v -race -count=1 -tags=integration ./internal/storage/repos/...

# run E2E tests (requires full environment)
test-e2e: test-db-up
	@echo "Waiting for test database to be ready..."
	@sleep 3
	@echo "Running E2E tests..."
	TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/test_db?sslmode=disable" \
		go test -v -race -count=1 -tags=e2e ./test/e2e/...

# run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Quick test (unit tests only)
test: test-unit

# Setup test environment
test-setup: test-db-up
	@echo "Test environment ready"

# Teardown test environment
test-teardown: test-db-down

# Start test database
test-db-up:
	@echo "Starting test database..."
	docker-compose -f docker-compose.test.yml up -d
	@echo "Waiting for database to be healthy..."
	@timeout 30 sh -c 'until docker-compose -f docker-compose.test.yml exec -T test-db pg_isready -U postgres; do sleep 1; done' 2>/dev/null || true

# Stop test database
test-db-down:
	@echo "Stopping test database..."
	docker-compose -f docker-compose.test.yml down -v

# Clean test data
test-clean: test-db-down
	@echo "Cleaning test data..."
	rm -f coverage.out coverage.html

# Stress testing
stress:
	@echo "Running stress test..."
	go run ./cmd/stresser

# Profiling
pprof:
	@echo "Opening pprof web interface..."
	@echo "Make sure the application is running..."
	go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30

# Run the application locally (for docker, just run docker-compose up)
run:
	@echo "Running the application..."
	docker-compose up postgres migrate -d
	go run ./cmd/app -config=config.local.yaml