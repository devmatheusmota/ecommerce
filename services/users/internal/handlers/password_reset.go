package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/usecase"
	"github.com/ecommerce/services/users/internal/validation"
)

type RequestPasswordResetRequest struct {
	Email string `json:"email"`
}

func RequestPasswordReset(uc *usecase.RequestPasswordReset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		var req RequestPasswordResetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		email, err := validation.ValidateRequestPasswordResetEmail(req.Email)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		output, err := uc.Execute(usecase.RequestPasswordResetInput{Email: email})
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				respondJSON(w, http.StatusOK, map[string]string{"message": "If an account exists, a reset link has been sent."})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process request"})
			return
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"message":     "If an account exists, a reset link has been sent.",
			"reset_token": output.Token,
			"expires_at":  output.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
}

type ConfirmPasswordResetRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func ConfirmPasswordReset(uc *usecase.ConfirmPasswordReset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		var req ConfirmPasswordResetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		if err := validation.ValidateConfirmPasswordReset(req.Token, req.NewPassword); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := uc.Execute(usecase.ConfirmPasswordResetInput{Token: req.Token, NewPassword: req.NewPassword}); err != nil {
			if errors.Is(err, domain.ErrInvalidResetToken) {
				respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid or expired password reset token"})
				return
			}
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to reset password"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
