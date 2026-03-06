# ADR-003: Local Development vs Kubernetes

## Status

Accepted

## Context

We need to run Postgres and RabbitMQ for local development and also have them available in Kubernetes for production-like environments (minikube, kind, or real clusters).

## Decision

- **Local dev**: Docker Compose (`docker-compose.yml`) runs Postgres and RabbitMQ
- **Kubernetes**: Manifests in `infra/k8s/` deploy the same stack via Kustomize
- Use the same credentials and port expectations where possible for consistency

## Rationale

- Docker Compose is the fastest path for local iteration: `make up` and code
- Kubernetes manifests serve the learning goal: understanding Deployments, StatefulSets, Services, probes
- Postgres uses a StatefulSet for stable storage; RabbitMQ uses a Deployment (single replica for now)
- Health checks (liveness/readiness) are included to practice production patterns

## Consequences

- Developers can choose: `make up` (Compose) or `kubectl apply -k infra/k8s/` (K8s)
- Services will read `POSTGRES_URL` and `RABBITMQ_URL` from env; same format for both environments
- In production, consider managed Postgres/RabbitMQ instead of self-hosted in K8s
