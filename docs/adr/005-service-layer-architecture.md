# ADR-005: Service Layer Architecture

## Status

Accepted

## Context

Services need a clear structure that supports unit testing, maintainability, and separation of concerns. Avoiding tight coupling between HTTP handling, validation, business logic, and persistence makes changes safer and tests faster.

## Decision

Adopt a layered architecture for each service:

1. **Handlers (HTTP layer)**: Thin; parse request, call use case, map response. No business logic or direct DB access.
2. **Use Cases**: Orchestrate validation, domain logic, and persistence. Depend on repository **interfaces** (injectable).
3. **Repository**: Implements persistence; concrete implementation (e.g. Postgres). Abstracts DB details.
4. **Validation**: Separate package for input validation; used by use cases.
5. **Domain**: Entities and domain errors. No dependencies on outer layers.

## Rationale

- **Testability**: Use cases can be tested with mocked repositories; no DB required for unit tests.
- **Decoupling**: Validation, persistence, and HTTP are independent; changes in one layer don't ripple across.
- **Clarity**: Each layer has a single responsibility; new developers can follow the flow easily.
- **Flexibility**: Swap implementations (e.g. different DB) by providing a new repository implementation.

## Structure (per service)

```
internal/
├── domain/       # Entities, domain errors
├── validation/   # Input validation
├── repository/   # Interface + Postgres impl
├── usecase/      # Application logic
└── handlers/     # HTTP
```

## Consequences

- New features follow this pattern; initial setup is more verbose but pays off in tests and refactors.
- Documentation: see `docs/architecture.md`.
- Other services (catalog, orders, etc.) should adopt the same structure.
