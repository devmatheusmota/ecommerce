package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/usecase"
	"github.com/ecommerce/services/users/internal/validation"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	CPF      string `json:"cpf"`
}

type RegisterResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func Register(uc *usecase.RegisterUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}

		input := &validation.RegisterInput{
			Email: req.Email, Password: req.Password,
			Name: req.Name, Phone: req.Phone, CPF: req.CPF,
		}
		if err := validation.ValidateRegisterInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		user, err := uc.Execute(usecase.RegisterUserInput{
			Email: input.Email, Password: input.Password,
			Name: input.Name, Phone: input.Phone, CPF: input.CPF,
		})
		if err != nil {
			respondRegisterError(w, err)
			return
		}

		respondJSON(w, http.StatusCreated, RegisterResponse{
			ID: user.ID, Email: user.Email, Name: user.Name,
		})
	}
}

func respondRegisterError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrDuplicateEmail):
		respondJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
	}
}
