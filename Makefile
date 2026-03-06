.PHONY: up down dev dev-down restart-dev kong-reset logs build build-users k8s-apply k8s-delete k8s-status migrate-users-up migrate-users-down kong-test

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

kong-test:
	@code=$$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8000/health 2>/dev/null); \
	if [ "$$code" = "200" ]; then echo "Kong OK — HTTP $$code"; else echo "Kong FAIL — gateway unreachable (run 'make up'?)"; fi

# --- Build (official/release: version from git) ---
# Builds the users service image with VERSION from git (same as dev). Image tagged as ecommerce-users:$(VERSION).
# Use for release: make build, then push the image to your registry.
build: build-users

build-users:
	docker build --build-arg VERSION=$(VERSION) -t ecommerce-users:$(VERSION) -f docker/services/users/Dockerfile .
	@echo "Built ecommerce-users:$(VERSION)"

# --- Kubernetes (minikube/kind) ---
k8s-apply:
	kubectl apply -k infra/k8s/

k8s-delete:
	kubectl delete -k infra/k8s/

k8s-status:
	kubectl get all -n ecommerce
