# E-commerce — Mercado Livre-style

Learning project for **distributed architecture** and full-cycle development: from application and API Gateway to deploy and observability.

**Stack:** Go · Postgres · RabbitMQ · Kong · Docker · Kubernetes

---

## What this project is

A study project that mimics a Mercado Livre-style e-commerce to practice:

- Microservices with bounded contexts (users, catalog, orders, payments)
- Event-driven communication (RabbitMQ) and synchronous calls (gRPC)
- API Gateway (Kong): routing, JWT, rate limiting, path handling
- Full cycle: local dev → Docker Compose → Kubernetes

Currently implemented: **Users service** (register, login with JWT, GET /me) behind Kong with per-route plugins.

---

## Prerequisites

- **Docker** and **Docker Compose**
- **Go 1.26+** (for building or running services locally)
- **kubectl** + **minikube** or **kind** (only for Kubernetes)

---

## Quick start

### Run everything locally (Docker Compose)

```bash
make up
```

| Service   | URL / Port | Notes |
|----------|------------|--------|
| Kong (proxy) | http://localhost:8000 | Client-facing API |
| Kong Manager | http://localhost:8002 | Routes, plugins (OSS: no login) |
| Postgres | localhost:5432 | User: `ecommerce`, Pass: `ecommerce_dev` |
| RabbitMQ | localhost:5672 (AMQP) | Management: http://localhost:15672 |
| Users API | via Kong: `http://localhost:8000/v1/users/...` | e.g. `/health`, `/login`, `/register`, `/me` |

**Smoke test:** `make kong-test`

### Dev mode (hot reload)

```bash
make dev
```

Uses Air for the users service; Kong and infra are the same.

### Kubernetes

```bash
make k8s-apply
make k8s-status
```

Requires a running cluster (e.g. `minikube start`).

---

## Project structure

```
ecommerce/
├── services/users/     # Auth: register, login (JWT), GET /me
├── infra/              # Kong config, K8s, Postgres init
│   └── kong/           # Routes, JWT, request-transformer (see README there)
├── docker/             # Dockerfiles
├── docs/               # Architecture, ADRs, project plan
└── scripts/            # Utilities
```

- **[AGENTS.md](AGENTS.md)** — Context and rules for AI agents
- **[docs/architecture.md](docs/architecture.md)** — Service layer (handlers, use cases, repository)
- **[docs/adr/](docs/adr/)** — Architecture decision records
- **[infra/kong/README.md](infra/kong/README.md)** — Kong routes, JWT, path handling

---

## Main commands

| Command | Description |
|--------|-------------|
| `make up` | Start Postgres, RabbitMQ, Kong, users |
| `make down` | Stop all |
| `make dev` | Like `up` but with hot reload for users |
| `make kong-reset` | Wipe Kong DB and re-import config (dev only) |
| `make kong-test` | Check Kong proxy (HTTP 200 on /health) |
| `make kong-jwt-test` | Full JWT flow via Kong: register → login → GET /me (requires jq) |
| `make k8s-apply` | Apply K8s manifests |
| `make build` | Build users service image |

---

## Versioning

Releases are tagged in git (`v0.1.0`, `v1.0.0`, `v1.1.0`, …). The users service exposes version in responses and via build-time ldflags (see [docs/git-tags.md](docs/git-tags.md)).
