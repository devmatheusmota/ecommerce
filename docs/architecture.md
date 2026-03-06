# Service Layer Architecture

Clean, layered architecture for microservices to support testability, maintainability, and separation of concerns.

---

## Layers

```
┌─────────────────────────────────────────────────────────┐
│  Handlers (HTTP)                                        │
│  - Parse request, delegate to use case, map response     │
│  - No business logic, no direct DB access                │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  Use Cases (Application)                                 │
│  - Orchestrate validation, domain logic, persistence     │
│  - Depends on repository interfaces (injectable)         │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  Repository (Infrastructure)                              │
│  - Implements persistence interfaces                     │
│  - SQL, external APIs, etc.                              │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  Domain                                                  │
│  - Entities, domain errors                              │
│  - No dependencies on outer layers                      │
└─────────────────────────────────────────────────────────┘

  Validation: separate package, used by use cases
```

---

## Per-Layer Responsibilities

| Layer | Responsibility | Depends on |
|-------|----------------|------------|
| **Handlers** | HTTP parsing, status codes, JSON encode/decode | Use cases |
| **Use Cases** | Business flow, validation, orchestration | Domain, Repository (interface), Validation |
| **Repository** | Persistence (Postgres, etc.) | Domain |
| **Validation** | Input validation rules | Domain (errors) |
| **Domain** | Entities, domain errors | Nothing |

---

## Dependency Rule

- Inner layers do **not** depend on outer layers
- Dependencies point **inward**: Handlers → Use Cases → Repositories
- Use **interfaces** for repositories so use cases can be tested with mocks

---

## Project Structure (Users Service)

```
services/users/
├── main.go
├── internal/
│   ├── domain/           # Entities and domain errors
│   │   ├── user.go
│   │   └── errors.go
│   ├── validation/       # Input validation (DTO → domain rules)
│   │   └── register.go
│   ├── repository/       # Persistence
│   │   ├── user_repository.go      # Interface
│   │   └── postgres_user_repository.go
│   ├── usecase/          # Application logic
│   │   ├── register_user.go
│   │   └── register_user_test.go
│   └── handlers/         # HTTP layer
│       └── register.go
└── ...
```

---

## Testing

- **Use cases**: Mock `UserRepository`; test validation, duplicate email, success
- **Handlers**: Use a fake use case or integration tests with real use case + mock repo
- **Repositories**: Integration tests with test database

Example (usecase):

```go
type mockUserRepo struct {
    createFunc func(*domain.User) (*domain.User, error)
}

func (m *mockUserRepo) Create(user *domain.User) (*domain.User, error) {
    return m.createFunc(user)
}

uc := NewRegisterUser(&mockUserRepo{createFunc: ...})
user, err := uc.Execute(input)
```

---

## References

- ADR-005: Service Layer Architecture
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) (Uncle Bob)
