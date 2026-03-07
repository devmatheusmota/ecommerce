package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/middleware"
	"github.com/ecommerce/services/users/internal/usecase"
	"github.com/ecommerce/services/users/internal/validation"
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

type UpdateMeRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	CPF   string `json:"cpf"`
}

func UpdateMe(updateProfileUseCase *usecase.UpdateProfile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req UpdateMeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}

		input := &validation.UpdateProfileInput{Name: req.Name, Phone: req.Phone, CPF: req.CPF}
		if err := validation.ValidateUpdateProfileInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		userID := middleware.UserIDFromContext(r.Context())
		output, err := updateProfileUseCase.Execute(usecase.UpdateProfileInput{
			UserID: userID,
			Name:   input.Name,
			Phone:  input.Phone,
			CPF:    input.CPF,
		})
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
