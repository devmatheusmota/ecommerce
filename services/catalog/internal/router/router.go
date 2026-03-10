package router

import (
	"database/sql"
	"net/http"

	"github.com/ecommerce/services/catalog/internal/handlers"
	"github.com/ecommerce/services/catalog/internal/openapi"
	"github.com/ecommerce/services/catalog/internal/repository"
	"github.com/ecommerce/services/catalog/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(database *sql.DB) http.Handler {
	return NewWithRepositories(
		repository.NewPostgresCategoryRepository(database),
		repository.NewPostgresProductRepository(database),
	)
}

func NewWithRepositories(categoryRepository repository.CategoryRepository, productRepository repository.ProductRepository) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	createCategoryUseCase := usecase.NewCreateCategory(categoryRepository)
	listCategoriesUseCase := usecase.NewListCategories(categoryRepository)
	listCategoriesTreeUseCase := usecase.NewListCategoriesTree(categoryRepository)
	getCategoryUseCase := usecase.NewGetCategory(categoryRepository)
	updateCategoryUseCase := usecase.NewUpdateCategory(categoryRepository)
	deleteCategoryUseCase := usecase.NewDeleteCategory(categoryRepository)

	createProductUseCase := usecase.NewCreateProduct(productRepository)
	listProductsUseCase := usecase.NewListProducts(productRepository)
	getProductUseCase := usecase.NewGetProduct(productRepository)
	listRelatedProductsUseCase := usecase.NewListRelatedProducts(productRepository, categoryRepository)
	updateProductUseCase := usecase.NewUpdateProduct(productRepository)
	deleteProductUseCase := usecase.NewDeleteProduct(productRepository)

	router.Get("/health", handlers.Health)
	router.Get("/docs", openapi.Handler())

	router.Route("/v1/categories", func(router chi.Router) {
		router.Post("/", handlers.CreateCategory(createCategoryUseCase))
		router.Get("/", handlers.ListCategories(listCategoriesUseCase))
		router.Get("/tree", handlers.ListCategoriesTree(listCategoriesTreeUseCase))
		router.Get("/{id}", handlers.GetCategory(getCategoryUseCase))
		router.Patch("/{id}", handlers.UpdateCategory(updateCategoryUseCase))
		router.Delete("/{id}", handlers.DeleteCategory(deleteCategoryUseCase))
	})

	router.Route("/v1/products", func(router chi.Router) {
		router.Post("/", handlers.CreateProduct(createProductUseCase))
		router.Get("/", handlers.ListProducts(listProductsUseCase))
		router.Get("/{id}/related", handlers.GetRelatedProducts(listRelatedProductsUseCase))
		router.Get("/{id}", handlers.GetProduct(getProductUseCase))
		router.Patch("/{id}", handlers.UpdateProduct(updateProductUseCase))
		router.Delete("/{id}", handlers.DeleteProduct(deleteProductUseCase))
	})

	return router
}
