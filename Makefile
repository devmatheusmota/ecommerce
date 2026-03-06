.PHONY: up down logs k8s-apply k8s-delete k8s-status migrate-users-up migrate-users-down

# --- Docker Compose (local dev) ---
up:
	docker compose up -d

down:
	docker compose down

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

# --- Kubernetes (minikube/kind) ---
k8s-apply:
	kubectl apply -k infra/k8s/

k8s-delete:
	kubectl delete -k infra/k8s/

k8s-status:
	kubectl get all -n ecommerce
