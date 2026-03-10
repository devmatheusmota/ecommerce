package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecommerce/services/users/internal/domain"
	"github.com/ecommerce/services/users/internal/middleware"
	"github.com/ecommerce/services/users/internal/usecase"
	"github.com/ecommerce/services/users/internal/validation"
	"github.com/go-chi/chi/v5"
)

type AddressResponse struct {
	ID                string `json:"id"`
	UserID            string `json:"user_id"`
	Street            string `json:"street"`
	Number            string `json:"number"`
	Complement        string `json:"complement"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	State             string `json:"state"`
	ZipCode           string `json:"zip_code"`
	Type              string `json:"type"`
	IsDefaultBilling  bool   `json:"is_default_billing"`
	IsDefaultShipping bool   `json:"is_default_shipping"`
}

func addressToResponse(a *domain.Address) AddressResponse {
	return AddressResponse{
		ID: a.ID, UserID: a.UserID, Street: a.Street, Number: a.Number,
		Complement: a.Complement, Neighborhood: a.Neighborhood, City: a.City, State: a.State, ZipCode: a.ZipCode,
		Type: a.Type, IsDefaultBilling: a.IsDefaultBilling, IsDefaultShipping: a.IsDefaultShipping,
	}
}

type CreateAddressRequest struct {
	Street            string `json:"street"`
	Number            string `json:"number"`
	Complement        string `json:"complement"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	State             string `json:"state"`
	ZipCode           string `json:"zip_code"`
	Type              string `json:"type"`
	IsDefaultBilling  bool   `json:"is_default_billing"`
	IsDefaultShipping bool   `json:"is_default_shipping"`
}

func CreateAddress(uc *usecase.CreateAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		var req CreateAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input := &validation.AddressInput{
			Street: req.Street, Number: req.Number, Complement: req.Complement,
			Neighborhood: req.Neighborhood, City: req.City, State: req.State, ZipCode: req.ZipCode,
			Type: req.Type, IsDefaultBilling: req.IsDefaultBilling, IsDefaultShipping: req.IsDefaultShipping,
		}
		if err := validation.ValidateAddressInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		userID := middleware.UserIDFromContext(r.Context())
		address, err := uc.Execute(usecase.CreateAddressInput{
			UserID: userID, Street: input.Street, Number: input.Number, Complement: input.Complement,
			Neighborhood: input.Neighborhood, City: input.City, State: input.State, ZipCode: input.ZipCode,
			Type: input.Type, IsDefaultBilling: input.IsDefaultBilling, IsDefaultShipping: input.IsDefaultShipping,
		})
		if err != nil {
			respondAddressError(w, err)
			return
		}
		respondJSON(w, http.StatusCreated, addressToResponse(address))
	}
}

func ListAddresses(uc *usecase.ListAddresses) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		userID := middleware.UserIDFromContext(r.Context())
		list, err := uc.Execute(userID)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list addresses"})
			return
		}
		out := make([]AddressResponse, len(list))
		for i, a := range list {
			out[i] = addressToResponse(a)
		}
		respondJSON(w, http.StatusOK, out)
	}
}

func GetAddress(uc *usecase.GetAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		addressID := chi.URLParam(r, "id")
		userID := middleware.UserIDFromContext(r.Context())
		address, err := uc.Execute(addressID, userID)
		if err != nil {
			respondAddressError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, addressToResponse(address))
	}
}

type UpdateAddressRequest struct {
	Street            string `json:"street"`
	Number            string `json:"number"`
	Complement        string `json:"complement"`
	Neighborhood      string `json:"neighborhood"`
	City              string `json:"city"`
	State             string `json:"state"`
	ZipCode           string `json:"zip_code"`
	Type              string `json:"type"`
	IsDefaultBilling  bool   `json:"is_default_billing"`
	IsDefaultShipping bool   `json:"is_default_shipping"`
}

func UpdateAddress(uc *usecase.UpdateAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		addressID := chi.URLParam(r, "id")
		var req UpdateAddressRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input := &validation.AddressInput{
			Street: req.Street, Number: req.Number, Complement: req.Complement,
			Neighborhood: req.Neighborhood, City: req.City, State: req.State, ZipCode: req.ZipCode,
			Type: req.Type, IsDefaultBilling: req.IsDefaultBilling, IsDefaultShipping: req.IsDefaultShipping,
		}
		if err := validation.ValidateAddressInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		userID := middleware.UserIDFromContext(r.Context())
		address, err := uc.Execute(usecase.UpdateAddressInput{
			AddressID: addressID, UserID: userID,
			Street: input.Street, Number: input.Number, Complement: input.Complement,
			Neighborhood: input.Neighborhood, City: input.City, State: input.State, ZipCode: input.ZipCode,
			Type: input.Type, IsDefaultBilling: input.IsDefaultBilling, IsDefaultShipping: input.IsDefaultShipping,
		})
		if err != nil {
			respondAddressError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, addressToResponse(address))
	}
}

func DeleteAddress(uc *usecase.DeleteAddress) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		addressID := chi.URLParam(r, "id")
		userID := middleware.UserIDFromContext(r.Context())
		if err := uc.Execute(addressID, userID); err != nil {
			respondAddressError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func respondAddressError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrAddressNotFound):
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "address not found"})
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process address"})
	}
}
