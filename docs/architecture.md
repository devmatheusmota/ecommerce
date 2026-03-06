# Service Layer Architecture

Clean, layered architecture for microservices to support testability, maintainability, and separation of concerns.

---

## Layers

```
┌─────────────────────────────────────────────────────────┐
│  Handlers (HTTP)                                        │
│  - Parse request, validate input (required, format)      │
│  - Delegate to use case, map response                   │
│  - No business logic, no direct DB access                │
└────────────────────────┬────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────┐
│  Use Cases (Application)                                 │
│  - Domain / business logic and persistence only           │
│  - Depends on repository interfaces (injectable)          │
│  - No input validation (handler does that)               │
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

  Validation: input checks (required, format) in handlers; shared rules in validation package
```

---

## Per-Layer Responsibilities

| Layer | Responsibility | Depends on |
|-------|----------------|------------|
| **Handlers** | HTTP parsing, **input validation** (required, format), status codes, JSON | Use cases, Validation |
| **Use Cases** | **Domain / business logic only** (e.g. hash password, check credentials, create user) | Domain, Repository (interface) |
| **Repository** | Persistence (Postgres, etc.) | Domain |
| **Validation** | Input validation helpers (required, format); shared rules in `common.go` | Domain (errors) |
| **Domain** | Entities, domain errors | Nothing |

**Where validation lives:** Input checks (e.g. “email present”, “email format”, “password required”) run in **handlers**; handlers call the `validation` package and return 400 before invoking the use case. Use cases contain **domain rules only** (e.g. duplicate email, invalid credentials, hash and persist). Shared rules (e.g. email format) live in `validation/common.go` and are reused by handlers. See `services/users/internal/validation/doc.go`.

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
│   │   ├── doc.go        # Package policy: reuse shared rules from common.go
│   │   ├── common.go     # Shared rules (e.g. ValidateEmail) — reuse, do not duplicate
│   │   ├── register.go
│   │   ├── login.go
│   │   └── cpf.go
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

- **Use cases**: Mock `UserRepository`; test domain rules (duplicate email, success); no input validation
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

## Code Style

- **Variable names**: No abbreviations — use full names (e.g. `userRepository`, `registerUserUseCase` instead of `userRepo`, `registerUC`)

---

## References

- ADR-005: Service Layer Architecture
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) (Uncle Bob)
