.PHONY: up down dev dev-down restart-dev kong-reset logs build build-users build-catalog k8s-apply k8s-delete k8s-status migrate-users-up migrate-users-down migrate-catalog-up migrate-catalog-down kong-test kong-jwt-test test-users test-catalog cover-users cover-users-html cover-catalog

# Version for dev: from latest git tag (e.g. 1.0.1 or 1.0.1-2-gabc123). Exported so docker-compose.dev.yml can use ${VERSION}.
VERSION := $(shell git describe --tags --always 2>/dev/null | sed 's/^v//' || echo "dev")
export VERSION

# --- Docker Compose (local dev) ---
up:
	docker compose up -d --build --force-recreate --remove-orphans

down:
	docker compose down

# Dev mode: users service with hot reload (Air) — code changes reflect without rebuild
dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

dev-down:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down

restart-dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build

# Create catalog database (run if Postgres was initialized before catalog was added)
catalog-db-create:
	docker exec ecommerce-postgres psql -U ecommerce -d postgres -c "CREATE DATABASE catalog;" 2>/dev/null || true
	@echo "Catalog database ready."

# Remove Kong's DB volume so kong-seed can re-import cleanly. Use when kong-seed fails with "UNIQUE violation on key".
kong-reset:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down
	docker volume rm ecommerce_kong_data 2>/dev/null || true
	@echo "Kong volume removed. Run 'make dev' or 'make restart-dev' to start fresh."

logs:
	docker compose logs -f

# --- Migrations (golang-migrate) ---
# Migrations run automatically when the users service starts (make up).
# For manual control, install migrate CLI and use (PG_HOST=localhost when outside Docker):
#   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
USERS_DB_URL ?= postgres://ecommerce:ecommerce_dev@localhost:5432/users?sslmode=disable

migrate-users-up:
	@cd services/users && migrate -path internal/database/migrations -database "$(USERS_DB_URL)" up

migrate-users-down:
	@cd services/users && migrate -path internal/database/migrations -database "$(USERS_DB_URL)" down 1

CATALOG_DB_URL ?= postgres://ecommerce:ecommerce_dev@localhost:5432/catalog?sslmode=disable

migrate-catalog-up:
	@cd services/catalog && migrate -path internal/database/migrations -database "$(CATALOG_DB_URL)" up

migrate-catalog-down:
	@cd services/catalog && migrate -path internal/database/migrations -database "$(CATALOG_DB_URL)" down 1

kong-test:
	@code=$$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8000/health 2>/dev/null); \
	if [ "$$code" = "200" ]; then echo "Kong OK — HTTP $$code"; else echo "Kong FAIL — gateway unreachable (run 'make up'?)"; fi

# Full JWT flow via Kong: register → login → GET /me. Requires curl and jq. Run after 'make up'.
kong-jwt-test:
	@base="http://localhost:8000/v1/users"; \
	email="kong-jwt-test-$$(date +%s)@example.com"; \
	echo "Registering $$email..."; \
	reg=$$(curl -s -X POST "$$base/register" -H "Content-Type: application/json" \
	  -d '{"email":"'$$email'","password":"testpass123","name":"JWT Test","phone":"+5511999999999","cpf":"529.982.247-25"}'); \
	if echo "$$reg" | grep -q '"error"'; then echo "Register FAIL: $$reg"; exit 1; fi; \
	echo "Logging in..."; \
	login=$$(curl -s -X POST "$$base/login" -H "Content-Type: application/json" -d '{"email":"'$$email'","password":"testpass123"}'); \
	token=$$(echo "$$login" | jq -r '.data.token'); \
	if [ -z "$$token" ] || [ "$$token" = "null" ]; then echo "Login FAIL: $$login"; exit 1; fi; \
	code=$$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $$token" "$$base/me"); \
	if [ "$$code" = "200" ]; then echo "Kong JWT OK — register, login, GET /me = $$code"; else echo "Kong JWT FAIL — GET /me = $$code (expected 200)"; exit 1; fi

# --- Tests & coverage (users service) ---
# Run tests only.
test-users:
	cd services/users && go test ./... -count=1

# Run tests and print coverage per function. Writes coverage to services/users/coverage.out.
# Requires a single Go version (go and go tool must match); if you see "version does not match", fix PATH or go.mod.
cover-users:
	cd services/users && go test ./... -coverprofile=coverage.out -count=1 && go tool cover -func=coverage.out

# Run tests and open HTML coverage report in the browser (coverage.out must exist; run make cover-users first, or this runs tests).
cover-users-html: cover-users
	cd services/users && go tool cover -html=coverage.out -o coverage.html
	@echo "Open services/users/coverage.html in your browser."

test-catalog:
	cd services/catalog && go test ./... -count=1

cover-catalog:
	cd services/catalog && go test ./... -coverprofile=coverage.out -count=1 && go tool cover -func=coverage.out

# --- Build (official/release: version from git) ---
# Builds service images with VERSION from git (same as dev).
# Use for release: make build, then push the images to your registry.
build: build-users build-catalog

build-users:
	docker build --build-arg VERSION=$(VERSION) -t ecommerce-users:$(VERSION) -f docker/services/users/Dockerfile .
	@echo "Built ecommerce-users:$(VERSION)"

build-catalog:
	docker build --build-arg VERSION=$(VERSION) -t ecommerce-catalog:$(VERSION) -f docker/services/catalog/Dockerfile .
	@echo "Built ecommerce-catalog:$(VERSION)"

# --- Kubernetes (minikube/kind) ---
k8s-apply:
	kubectl apply -k infra/k8s/

k8s-delete:
	kubectl delete -k infra/k8s/

k8s-status:
	kubectl get all -n ecommerce
