# ADR-002: Repository Structure

## Status

Accepted

## Context

We need a clear layout for a polyglot-ready monorepo that will host multiple Go services, shared code, infrastructure definitions, and documentation.

## Decision

```
ecommerce/
├── services/       # Microservices (catalog, orders, payments, etc.)
├── shared/         # Shared libraries, contracts, schemas
├── infra/          # Terraform, Helm, K8s manifests
├── docker/         # Dockerfiles (can live in services/ or here)
├── docs/           # ADRs, diagrams, runbooks
└── scripts/        # Dev/deploy utilities
```

## Rationale

- **services/**: One folder per bounded context; each is independently deployable
- **shared/**: Avoid duplication of event schemas, DB helpers, common types
- **infra/**: All infra-as-code in one place; Kustomize for K8s
- **docker/**: Centralized Dockerfiles or colocated in services — flexible
- **docs/**: ADRs follow the format `docs/adr/XXX-title.md`; runbooks for ops
- **scripts/**: Makefile, helper scripts for local dev and CI

## Consequences

- New services are added under `services/<name>/`
- Shared Go packages live in `shared/` and are imported via module path
- K8s resources are applied via `kubectl apply -k infra/k8s/`
