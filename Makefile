.PHONY: up down logs k8s-apply k8s-delete k8s-status

# --- Docker Compose (local dev) ---
up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

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
