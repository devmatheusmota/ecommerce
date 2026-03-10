package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/usecase"
	"github.com/ecommerce/services/catalog/internal/validation"
	"github.com/go-chi/chi/v5"
)

type ProductResponse struct {
	ID          string   `json:"id"`
	SellerID    string   `json:"seller_id"`
	CategoryID  string   `json:"category_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	Images      []string `json:"images"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func productToResponse(p *domain.Product) ProductResponse {
	images := p.Images
	if images == nil {
		images = []string{}
	}
	return ProductResponse{
		ID: p.ID, SellerID: p.SellerID, CategoryID: p.CategoryID,
		Title: p.Title, Description: p.Description, Price: p.Price,
		Images: images, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

type CreateProductRequest struct {
	SellerID    string   `json:"seller_id"`
	CategoryID  string   `json:"category_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	Images      []string `json:"images"`
}

func CreateProduct(createProductUseCase *usecase.CreateProduct) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		var request CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input := &validation.ProductInput{
			SellerID: request.SellerID, CategoryID: request.CategoryID,
			Title: request.Title, Description: request.Description,
			Price: request.Price, Images: request.Images,
		}
		if err := validation.ValidateProductInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		product, err := createProductUseCase.Execute(usecase.CreateProductInput{
			SellerID: input.SellerID, CategoryID: input.CategoryID,
			Title: input.Title, Description: input.Description,
			Price: input.Price, Images: input.Images,
		})
		if err != nil {
			respondProductError(w, err)
			return
		}
		respondJSON(w, http.StatusCreated, productToResponse(product))
	}
}

func ListProducts(listProductsUseCase *usecase.ListProducts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		filter := usecase.ListProductsFilter{
			SellerID:   r.URL.Query().Get("seller_id"),
			CategoryID: r.URL.Query().Get("category_id"),
			Limit:      limit,
			Offset:     offset,
		}

		result, err := listProductsUseCase.Execute(filter)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list products"})
			return
		}

		products := make([]ProductResponse, len(result.Products))
		for i, product := range result.Products {
			products[i] = productToResponse(product)
		}
		respondJSON(w, http.StatusOK, map[string]any{
			"products": products,
			"total":    result.Total,
		})
	}
}

func GetProduct(getProductUseCase *usecase.GetProduct) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		productID := chi.URLParam(r, "id")
		product, err := getProductUseCase.Execute(productID)
		if err != nil {
			respondProductError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, productToResponse(product))
	}
}

type UpdateProductRequest struct {
	SellerID    string   `json:"seller_id"`
	CategoryID  string   `json:"category_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       string   `json:"price"`
	Images      []string `json:"images"`
}

func UpdateProduct(updateProductUseCase *usecase.UpdateProduct) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		productID := chi.URLParam(r, "id")
		var request UpdateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		input := &validation.ProductInput{
			SellerID: request.SellerID, CategoryID: request.CategoryID,
			Title: request.Title, Description: request.Description,
			Price: request.Price, Images: request.Images,
		}
		if err := validation.ValidateProductInput(input); err != nil {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		product, err := updateProductUseCase.Execute(usecase.UpdateProductInput{
			ID: productID, SellerID: input.SellerID, CategoryID: input.CategoryID,
			Title: input.Title, Description: input.Description,
			Price: input.Price, Images: input.Images,
		})
		if err != nil {
			respondProductError(w, err)
			return
		}
		respondJSON(w, http.StatusOK, productToResponse(product))
	}
}

func DeleteProduct(deleteProductUseCase *usecase.DeleteProduct) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		productID := chi.URLParam(r, "id")
		if err := deleteProductUseCase.Execute(productID); err != nil {
			respondProductError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func respondProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
	default:
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to process product"})
	}
}
