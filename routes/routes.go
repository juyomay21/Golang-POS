package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/controllers"
	"github.com/jinzhu/gorm"
)

func Register(router fiber.Router, database *gorm.DB) {

	cashier := new(controllers.Cashier)
	router.Get("/cashiers", cashier.GetAllCashier)
	router.Get("/cashiers/:id", cashier.GetCashier)
	router.Get("/cashiers/:id/passcode", cashier.Passcode)
	router.Post("/cashiers/:id/login", cashier.Login)
	router.Put("/cashiers/:id", cashier.UpdateCashier)
	router.Post("/cashiers", cashier.CreateCashier)
	router.Delete("/cashiers/:id", cashier.DeleteCashier)

	category := new(controllers.Category)
	router.Get("/categories", category.GetAllCategory)
	router.Get("/categories/:id", category.GetCategory)
	router.Post("/categories", category.CreateCategory)
	router.Delete("/categories/:id", category.DeleteCategory)
	router.Put("/categories/:id", category.UpdateCategory)

	product := new(controllers.Product)
	router.Post("/products", product.CreateProduct)
	router.Get("/products", product.GetAllProduct)
	router.Get("/products/:id", product.GetProduct)
	router.Put("/products/:id", product.UpdateProduct)
	router.Delete("/products/:id", product.DeleteProduct)

	payment := new(controllers.Payment)
	router.Get("/payments", payment.GetAllPayment)
	router.Get("/payments/:id", payment.GetPayment)

	order := new(controllers.Order)
	router.Post("/orders", order.CreateOrder)
	router.Post("/subtotal", order.GetSubTotal)
	router.Get("/orders", order.GetAllOrder)
	router.Get("/orders/:id", order.GetOrder)
	router.Get("/orders/:id/download", order.GetOrder)
	router.Get("/orders/:id/check-download", order.GetOrderStatus)

}
