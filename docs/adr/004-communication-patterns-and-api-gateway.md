# ADR-004: Communication Patterns and API Gateway

## Status

Accepted

## Context

We need clear rules for how clients talk to the system and how services talk to each other. We also want to learn Kong as an API Gateway.

## Decision

- **Client → Application**: HTTP only. All client traffic (web, mobile, external integrations) goes through Kong API Gateway. Kong handles routing, rate limiting, auth, and request/response handling.
- **Service → Service (synchronous)**: gRPC. When a service needs a direct, immediate response from another (e.g. catalog service checking stock for the orders service), use gRPC with protobuf.
- **Service → Service (asynchronous)**: RabbitMQ. For domain events, eventual consistency, and fire-and-forget, use the message broker.
- **API Gateway**: Kong. Learning focus; provides routing, plugins, and a single entry point for client HTTP traffic.

## Rationale

### Client via HTTP + Kong
- Clients (browsers, mobile apps) expect HTTP/REST; Kong is built for that
- Kong centralizes cross-cutting concerns: auth, rate limiting, logging, CORS
- One public entry point simplifies security and deployment
- Kong has a strong plugin ecosystem and is widely used in production

### gRPC for sync service-to-service
- Binary protocol is more efficient than JSON over HTTP for internal calls
- Strong typing via protobuf reduces integration bugs
- Streaming support if needed later
- Go has excellent gRPC support (google.golang.org/grpc)

### RabbitMQ for async
- Already chosen in ADR-001; decoupling and resilience for events
- Services don’t need to be up at the same time; queues buffer messages

## Consequences

- Services expose two interfaces: HTTP (for Kong to route client requests) and gRPC (for other services)
- Proto definitions will live in `shared/proto/` or similar for reuse
- Kong must be deployed and configured before client-facing features; can be added to Docker Compose and K8s
- When designing a flow, choose: client need → HTTP via Kong; sync service call → gRPC; async/event → RabbitMQ
- **Auth**: Kong validates the JWT and forwards the subject (user ID) to upstream via a header (e.g. `X-User-ID`). Services behind Kong trust this header and do not validate JWT; they focus on business logic only.
