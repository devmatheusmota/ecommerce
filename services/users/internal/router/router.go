// Package router configures HTTP routes for the users service.
package router

import (
	"database/sql"
	"net/http"

	"github.com/ecommerce/services/users/internal/handlers"
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

	r.Get("/health", handlers.Health)
	r.Get("/docs", openapi.Handler())
	r.Post("/register", handlers.Register(registerUsecase))
	r.Post("/login", handlers.Login(loginUsecase))

	return r
}
