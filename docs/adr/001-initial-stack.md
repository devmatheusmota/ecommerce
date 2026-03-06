# ADR-001: Initial Technology Stack

## Status

Accepted

## Context

We need a foundation to build a Mercado Livre-style e-commerce as a learning vehicle for distributed architecture. The stack should support event-driven communication, transactional data, and full-cycle development (dev, deploy, observability, Kubernetes, scaling).

## Decision

- **Backend language**: Go  
- **Database**: PostgreSQL  
- **Message broker**: RabbitMQ  

## Rationale

### Go
- Learning goal: the team wants to become proficient in Go
- Strong fit for cloud-native: small binaries, fast startup, built-in concurrency
- Good ecosystem for microservices (standard library, chi/gin, pgx, amqp)
- Widely used in infra tools (Kubernetes, Docker, Terraform) — useful for full-cycle understanding

### PostgreSQL
- Mature, robust, ACID-compliant relational database
- Excellent for transactional domains (orders, inventory, payments)
- JSON support for flexible schemas when needed
- Well-supported in Go (pgx, gorm)

### RabbitMQ
- Proven AMQP broker; easier to run locally than Kafka for learning
- Management UI (port 15672) helps inspect queues and messages
- Good fit for event-driven architecture
- Can be swapped for Kafka later if throughput requirements justify it

## Consequences

- All services will be written in Go
- Database migrations will target PostgreSQL
- Domain events will use RabbitMQ exchanges/queues
- Redis can be added later for caching without changing the core stack
