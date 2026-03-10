package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/usecase"
	"github.com/ecommerce/services/catalog/internal/validation"
	"github.com/go-chi/chi/v5"
)

type CategoryResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ParentID  string `json:"parent_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func categoryToResponse(c *domain.Category) CategoryResponse {
	resp := CategoryResponse{
		ID: c.ID, Name: c.Name, Slug: c.Slug,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
	if c.ParentID != "" {
		resp.ParentID = c.ParentID
	}
	return resp
}

type CreateCategoryRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID string `json:"parent_id"`
}

func CreateCategory(createCategoryUseCase *usecase.CreateCategory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		var request CreateCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input := &validation.CategoryInput{
			Name: request.Name, Slug: request.Slug, ParentID: request.ParentID,
		}
		if err := validation.ValidateCategoryInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		category, err := createCategoryUseCase.Execute(usecase.CreateCategoryInput{
			Name: input.Name, Slug: input.Slug, ParentID: input.ParentID,
		})
		if err != nil {
			respondCategoryError(w, err)
			return
		}
		respondJSON(w, http.StatusCreated, categoryToResponse(category))
	}
}

func ListCategoriesTree(listCategoriesTreeUseCase *usecase.ListCategoriesTree) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		tree, err := listCategoriesTreeUseCase.Execute()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list categories tree"})
			return
		}
		respondJSON(w, http.StatusOK, tree)
	}
}

func ListCategories(listCategoriesUseCase *usecase.ListCategories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		parentID := r.URL.Query().Get("parent_id")
		categories, err := listCategoriesUseCase.Execute(parentID)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list categories"})
			return
		}
		response := make([]CategoryResponse, len(categories))
		for i, category := range categories {
			response[i] = categoryToResponse(category)
		}
		respondJSON(w, http.StatusOK, response)
	}
}

func GetCategory(getCategoryUseCase *usecase.GetCategory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		categoryID := chi.URLParam(r, "id")
		category, err := getCategoryUseCase.Execute(categoryID)
		if err != nil {
			respondCategoryError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, categoryToResponse(category))
	}
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID string `json:"parent_id"`
}

func UpdateCategory(updateCategoryUseCase *usecase.UpdateCategory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		categoryID := chi.URLParam(r, "id")
		var request UpdateCategoryRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input := &validation.CategoryInput{
			Name: request.Name, Slug: request.Slug, ParentID: request.ParentID,
		}
		if err := validation.ValidateCategoryInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		category, err := updateCategoryUseCase.Execute(usecase.UpdateCategoryInput{
			ID: categoryID, Name: input.Name, Slug: input.Slug, ParentID: input.ParentID,
		})
		if err != nil {
			respondCategoryError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, categoryToResponse(category))
	}
}

func DeleteCategory(deleteCategoryUseCase *usecase.DeleteCategory) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		categoryID := chi.URLParam(r, "id")
		if err := deleteCategoryUseCase.Execute(categoryID); err != nil {
			respondCategoryError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func respondCategoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrCategoryNotFound):
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "category not found"})
	case errors.Is(err, domain.ErrDuplicateSlug):
		respondJSON(w, http.StatusConflict, map[string]string{"error": "category slug already exists"})
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process category"})
	}
}
