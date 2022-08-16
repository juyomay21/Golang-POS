package pos

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

func Register(router fiber.Router, database *gorm.DB) {
	database.AutoMigrate(&Cashiers{})
	database.AutoMigrate(&Category{})
	database.AutoMigrate(&Products{})

	cashierRepository := NewCashierRepository(database)
	cashierHandler := NewCashierHandler(cashierRepository)

	router.Get("/cashiers", cashierHandler.GetAllCashier)
	router.Get("/cashiers/:id", cashierHandler.GetCashier)
	router.Get("/cashiers/:id/passcode", cashierHandler.Passcode)
	router.Post("/cashiers/:id/login", cashierHandler.Login)
	router.Put("/cashiers/:id", cashierHandler.UpdateCashier)
	router.Post("/cashiers", cashierHandler.CreateCashier)
	router.Delete("/cashiers/:id", cashierHandler.DeleteCashier)

	categoryRepository := NewCategoryRepository(database)
	categoryHandler := NewCategoryHandler(categoryRepository)

	router.Get("/categories", categoryHandler.GetAll)
	router.Get("/categories/:id", categoryHandler.Get)
	router.Post("/categories", categoryHandler.Create)
	router.Delete("/categories/:id", categoryHandler.Delete)
	router.Put("/categories/:id", categoryHandler.Update)

	productRepository := NewProductRepository(database)
	productHandler := NewProductHandler(productRepository)

	router.Post("/products", productHandler.CreateProduct)
	router.Get("/products", productHandler.GetAllProduct)
	router.Get("/products/:id", productHandler.GetProduct)
	router.Put("/products/:id", productHandler.UpdateProduct)
	router.Delete("/products/:id", productHandler.DeleteProduct)
}
