package handlers

import (
	"errors"
	"net/http"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/middleware"
	"github.com/ecommerce/services/users/internal/usecase"
)

type MeResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	CPF   string `json:"cpf"`
}

func Me(uc *usecase.MeUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		userID := middleware.UserIDFromContext(r.Context())

		output, err := uc.Execute(usecase.MeUserInput{UserID: userID})
		if err != nil {
			respondMeError(w, err)
			return
		}

		respondJSON(w, http.StatusOK, MeResponse{
			ID:    output.ID,
			Email: output.Email,
			Name:  output.Name,
			Phone: output.Phone,
			CPF:   output.CPF,
		})

	}
}

func respondMeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
	}
}
