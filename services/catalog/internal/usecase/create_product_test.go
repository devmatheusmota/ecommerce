package usecase

import (
	"testing"

	"github.com/ecommerce/services/catalog/internal/domain"
	"github.com/ecommerce/services/catalog/internal/repository"
)

func TestCreateProduct_Execute_Success(t *testing.T) {
	mockRepository := &repository.MockProductRepository{
		CreateFunc: func(product *domain.Product) (*domain.Product, error) {
			product.ID = "prod-456"
			return product, nil
		},
	}
	createProductUseCase := NewCreateProduct(mockRepository)

	product, err := createProductUseCase.Execute(CreateProductInput{
		SellerID:    "seller-1",
		CategoryID:  "cat-1",
		Title:       "Smartphone",
		Description: "A great phone",
		Price:       "999.90",
		Images:      []string{"https://example.com/img1.jpg"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != "prod-456" {
		t.Errorf("expected ID prod-456, got %s", product.ID)
	}
	if product.Title != "Smartphone" {
		t.Errorf("expected Title Smartphone, got %s", product.Title)
	}
	if len(product.Images) != 1 || product.Images[0] != "https://example.com/img1.jpg" {
		t.Errorf("expected one image, got %v", product.Images)
	}
}

func TestCreateProduct_Execute_NilImages(t *testing.T) {
	var capturedProduct *domain.Product
	mockRepository := &repository.MockProductRepository{
		CreateFunc: func(product *domain.Product) (*domain.Product, error) {
			capturedProduct = product
			product.ID = "prod-789"
			return product, nil
		},
	}
	createProductUseCase := NewCreateProduct(mockRepository)

	product, err := createProductUseCase.Execute(CreateProductInput{
		SellerID:    "seller-1",
		CategoryID:  "cat-1",
		Title:       "Laptop",
		Description: "A great laptop",
		Price:       "1999.00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedProduct.Images == nil {
		t.Error("expected non-nil Images (empty slice)")
	}
	if len(capturedProduct.Images) != 0 {
		t.Errorf("expected empty Images, got %v", capturedProduct.Images)
	}
	if product.ID != "prod-789" {
		t.Errorf("expected ID prod-789, got %s", product.ID)
	}
}
