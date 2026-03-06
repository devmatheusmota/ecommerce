# E-commerce (Mercado Livre-style) — Distributed Architecture Study Project

A learning project to master full-cycle development: application, deploy, observability, Kubernetes, and scalability. AGENTS.md guides AI agents in the project context.

---

## Project Vision

- **Goal**: Replicate a Mercado Livre-style e-commerce as a vehicle for learning distributed architecture
- **Learning focus**: Full-cycle (dev → deploy → observability → Kubernetes → scaling)
- **Secondary goal**: Extract maximum value from AI — clear prompts, planning before coding, guided iterations

---

## Architecture Principles

- **Event-driven**: Prefer asynchronous communication between services; domain events via RabbitMQ instead of synchronous calls when it makes sense
- **Bounded contexts**: Separate domains (catalog, orders, payments, shipping, users); each context with its own API and data model
- **Resilience**: Timeout, retry, circuit breaker, dead-letter for failures
- **Observability first**: Structured logs, metrics, and traces from day one; ease debugging in distributed environments
- **Stateless**: Services with no local state to allow horizontal scaling

## Communication Patterns

- **Client → App**: HTTP only. Clients (web, mobile, external APIs) talk to the application exclusively via HTTP through Kong API Gateway. Backend services do not expose HTTP ports to the host; only Kong is exposed (e.g. port 8000), so clients cannot reach services directly.
- **Service → Service (sync)**: gRPC. When a service needs a direct, synchronous response from another, use gRPC
- **Service → Service (async)**: RabbitMQ. Domain events, eventual consistency, fire-and-forget — use the message broker

---

## Stack and Tools

- **Backend**: Go (learning focus)
- **Database**: Postgres (transactional store)
- **Message broker**: RabbitMQ (domain events, async communication)
- **Service-to-service sync**: gRPC
- **API Gateway**: Kong (client-facing HTTP; learning focus)
- **Infra**: Docker, Kubernetes (minikube/kind locally); CI/CD (GitHub Actions or similar)
- **Observability**: Prometheus + Grafana; OpenTelemetry for traces; ELK/Loki for logs (to be added)

---

## How to Work with AI on This Project

### Before coding
- **Use Planning Mode** (Shift+Tab in the agent input): plan concrete steps before writing code
- **Context in chunks**: Request one feature at a time; avoid huge prompts with many requirements
- **Cite relevant files**: `@file.ts` when you know what you need; let the agent find the rest

### During development
- **Ask for explanations**: "Why this approach?" helps solidify architecture concepts
- **Ask for alternatives**: "What are the trade-offs of X vs Y here?" reinforces technical decisions
- **Ask for documentation**: "Document this decision in AGENTS.md" — the project itself teaches future agents

### Prompt patterns
- ❌ "build an e-commerce"  
- ✅ "Implement the create order flow: REST API, validation, persistence in Postgres, order created event"
- ❌ "configure kubernetes"  
- ✅ "Add a Kubernetes Deployment and Service for the orders service; include readiness and liveness probes"

---

## Service Layer Architecture

Each microservice follows a layered structure for testability and decoupling (see [ADR-005](docs/adr/005-service-layer-architecture.md), [docs/architecture.md](docs/architecture.md)):

- **Handlers**: HTTP only — parse request, delegate to use case, map response
- **Use cases**: Business logic, validation orchestration; depend on repository **interfaces**
- **Repository**: Persistence (Postgres); implements interfaces for mockable unit tests
- **Validation**: Input validation rules; used by use cases
- **Domain**: Entities and domain errors; no outer dependencies

---

## Repository Structure (planned)

```
ecommerce/
├── services/           # Microservices
│   ├── catalog/
│   ├── orders/
│   ├── payments/
│   └── ...
├── shared/             # Libraries, contracts, schemas
│   └── proto/          # gRPC protobuf definitions
├── infra/              # Terraform, Helm, K8s manifests
├── docker/             # Dockerfiles
├── docs/               # ADRs, diagrams, runbooks
└── scripts/            # Dev/deploy utilities
```

---

## Commands and Dev Environment

- **Start local environment** (Postgres + RabbitMQ): `make up`
- **Stop**: `make down`
- **View logs**: `make logs`
- **Kubernetes apply**: `make k8s-apply` (from project root, requires cluster)
- **Kubernetes status**: `make k8s-status`
- **Tests/Lint/Build**: per-service (`make test`, `make lint`, etc. when services exist)

---

## Code Conventions

- **Language**: Code, comments, and docs in **English**
- **Variable names**: No abbreviations — use full names (e.g. `userRepository`, `registerUserUseCase` instead of `userRepo`, `registerUC`)
- **Commits**: Clear messages; prefer Conventional Commits (`feat:`, `fix:`, `docs:`)
- **Client-facing APIs**: REST with versioning (`/v1/orders`); OpenAPI/Swagger for documentation; exposed via Kong
- **Service-to-service**: gRPC with protobuf; proto files in `shared/proto/` or per-service
- **Config**: Environment variables; no secrets in code

---

## Learned User Preferences

*(Section for preferences extracted from conversations over time)*

- **Variable naming**: No abbreviated variable names (prefer `userRepository` over `userRepo`, `registerUserUseCase` over `registerUC`)


---

## Learned Workspace Facts

- Stack: Go, Postgres, RabbitMQ, gRPC, Kong (see ADR-001, ADR-004)
- Business rules: `docs/business-rules.md`
- Client → Kong → services (HTTP); Service → Service sync: gRPC; async: RabbitMQ
- Local dev: `make up` runs Postgres + RabbitMQ via Docker Compose
- K8s: `infra/k8s/` uses Kustomize; apply with `make k8s-apply`

