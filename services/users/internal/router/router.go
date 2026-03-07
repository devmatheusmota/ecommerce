// Package router configures HTTP routes for the users service.
package router

import (
	"database/sql"
	"net/http"

	"github.com/ecommerce/services/users/internal/handlers"
	authmiddleware "github.com/ecommerce/services/users/internal/middleware"
	"github.com/ecommerce/services/users/internal/openapi"
	"github.com/ecommerce/services/users/internal/repository"
	"github.com/ecommerce/services/users/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New returns a router with Postgres-backed user repository. Use for production.
func New(db *sql.DB) http.Handler {
	return NewWithRepository(repository.NewPostgresUserRepository(db))
}

// NewWithRepository returns a router with the given user repository. Use for tests with a mock repository.
func NewWithRepository(userRepository repository.UserRepository) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	registerUsecase := usecase.NewRegisterUser(userRepository)
	loginUsecase := usecase.NewLoginUser(userRepository)
	meUsecase := usecase.NewMeUser(userRepository)

	r.Get("/health", handlers.Health)
	r.Get("/docs", openapi.Handler())
	r.Post("/register", handlers.Register(registerUsecase))
	r.Post("/login", handlers.Login(loginUsecase))

	// Routes that require X-User-ID (set by Kong after JWT validation). Guarantees the header is present; use authmiddleware.UserIDFromContext in handlers.
	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.RequireUserID)
		r.Get("/me", handlers.Me(meUsecase))
	})

	return r
}
