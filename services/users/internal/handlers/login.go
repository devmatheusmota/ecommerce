package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/usecase"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

func Login(uc *usecase.LoginUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}

		output, err := uc.Execute(usecase.LoginUserInput{
			Email: req.Email, Password: req.Password,
		})
		if err != nil {
			respondLoginError(w, err)
			return
		}

		respondJSON(w, http.StatusOK, LoginResponse{
			Token: output.Token, ExpireAt: output.ExpireAt,
		})
	}
}

func respondLoginError(w http.ResponseWriter, err error) {
	var valErr domain.ErrValidation
	switch {
	case errors.As(err, &valErr):
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to login user"})
	}
}
