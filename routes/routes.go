package routes

import (
	"chi-sqlx/config"
	"chi-sqlx/database/repository"
	"chi-sqlx/handler"
	"chi-sqlx/service"

	"github.com/jmoiron/sqlx"
)

func RegisterRoutes(db *sqlx.DB) {
	port := config.Env("APP_PORT", "8080")

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductController(productService)

	handler.ProductHandler(productHandler)
	handler.Start(":" + port)
}
