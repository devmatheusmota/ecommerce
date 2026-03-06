# E-commerce (Mercado Livre-style)

A learning project for distributed architecture and full-cycle development: Go, Postgres, RabbitMQ, Docker, and Kubernetes.

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (when building services)
- kubectl + minikube or kind (for Kubernetes)

## Quick Start

### Local development (Docker Compose)

```bash
make up
```

- Postgres: `localhost:5432` (user: `ecommerce`, pass: `ecommerce_dev`, db: `ecommerce`)
- RabbitMQ: `localhost:5672` (AMQP) | Management UI: http://localhost:15672 (user: `ecommerce`, pass: `ecommerce_dev`)
- Kong (API Gateway): `localhost:8000` (proxy) | Kong Manager (UI): http://localhost:8002 | Test: `make kong-test` — *Note: Kong Manager OSS has no login; RBAC/auth require Enterprise license*
- Users (auth): `localhost:8081` (direct) | via Kong: `curl http://localhost:8000/v1/users/health`

### Kubernetes

```bash
# Ensure cluster is running (e.g. minikube start)
make k8s-apply
make k8s-status
```

## Project Structure

- [AGENTS.md](AGENTS.md) — Full context for AI agents
- [docs/architecture.md](docs/architecture.md) — Service layer architecture (handlers, use cases, repository, domain)
- [docs/adr/](docs/adr/) — Architecture decisions
- [docs/project-plan.md](docs/project-plan.md) — Features, microservices, tasks, and implementation order
