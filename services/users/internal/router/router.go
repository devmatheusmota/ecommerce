// Package router configures HTTP routes for the users service.
package router

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ecommerce/services/users/internal/handlers"
	authmiddleware "github.com/ecommerce/services/users/internal/middleware"
	"github.com/ecommerce/services/users/internal/openapi"
	"github.com/ecommerce/services/users/internal/repository"
	"github.com/ecommerce/services/users/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New returns a router with Postgres-backed repositories. Use for production.
func New(db *sql.DB) http.Handler {
	return NewWithRepositories(
		repository.NewPostgresUserRepository(db),
		repository.NewPostgresAddressRepository(db),
		repository.NewPostgresPasswordResetRepository(db),
	)
}

// NewWithRepositories returns a router with the given repositories. Use for tests with mocks.
func NewWithRepositories(userRepository repository.UserRepository, addressRepository repository.AddressRepository, passwordResetRepository repository.PasswordResetRepository) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	registerUsecase := usecase.NewRegisterUser(userRepository)
	loginUsecase := usecase.NewLoginUser(userRepository)
	meUsecase := usecase.NewMeUser(userRepository)
	updateProfileUsecase := usecase.NewUpdateProfile(userRepository)

	createAddressUsecase := usecase.NewCreateAddress(addressRepository)
	listAddressesUsecase := usecase.NewListAddresses(addressRepository)
	getAddressUsecase := usecase.NewGetAddress(addressRepository)
	updateAddressUsecase := usecase.NewUpdateAddress(addressRepository)
	deleteAddressUsecase := usecase.NewDeleteAddress(addressRepository)

	requestPasswordResetUsecase := usecase.NewRequestPasswordReset(userRepository, passwordResetRepository)
	confirmPasswordResetUsecase := usecase.NewConfirmPasswordReset(userRepository, passwordResetRepository)

	r.Get("/health", handlers.Health)
	r.Get("/docs", openapi.Handler())
	r.Post("/register", handlers.Register(registerUsecase))
	r.Post("/login", handlers.Login(loginUsecase))
	r.Post("/password-reset/request", handlers.RequestPasswordReset(requestPasswordResetUsecase))
	r.Post("/password-reset/confirm", handlers.ConfirmPasswordReset(confirmPasswordResetUsecase))

	fmt.Println("Router initialized")

	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.RequireUserID)
		r.Get("/me", handlers.Me(meUsecase))
		r.Patch("/me", handlers.UpdateMe(updateProfileUsecase))
		r.Route("/me/addresses", func(r chi.Router) {
			r.Post("/", handlers.CreateAddress(createAddressUsecase))
			r.Get("/", handlers.ListAddresses(listAddressesUsecase))
			r.Get("/{id}", handlers.GetAddress(getAddressUsecase))
			r.Patch("/{id}", handlers.UpdateAddress(updateAddressUsecase))
			r.Delete("/{id}", handlers.DeleteAddress(deleteAddressUsecase))
		})
	})

	return r
}

// NewWithRepository returns a router with the given user repository and mock repos. Use for tests that only need user flows.
func NewWithRepository(userRepository repository.UserRepository) http.Handler {
	return NewWithRepositories(userRepository, repository.NewMockAddressRepository(), repository.NewMockPasswordResetRepository())
}
