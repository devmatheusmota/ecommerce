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

func New(db *sql.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userRepository := repository.NewPostgresUserRepository(db)
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
