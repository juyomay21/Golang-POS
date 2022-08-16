package pos

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	repository *ProductRepository
	Base
}

type CategoryList struct {
	CategoryId int64  `json:"categoryId"`
	Name       string `json:"name"`
}

type ProductList struct {
	ProductId  int64        `json:"productId"`
	Sku        string       `json:"sku"`
	Name       string       `json:"name"`
	Stock      int64        `json:"stock"`
	Price      float64      `json:"price"`
	Image      string       `json:"image"`
	Category   CategoryList `json:"category"`
	DiscountId int64        `json:"discount"`
}

func (handler *ProductHandler) GetAllProduct(c *fiber.Ctx) error {

	// err := handler.Auth(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"success": false,
	// 		"message": "Authentication Failed",
	// 	})
	// }

	var products []ProductList = handler.repository.FindAllProduct(c)
	count := handler.repository.GetProductCount(c)

	if len(products) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   fiber.Map{},
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": fiber.Map{
				"products": products,
				"meta": fiber.Map{
					"total": count,
					"limit": c.Query("limit"),
					"skip":  c.Query("skip"),
				},
			},
		})
	}
}

func (handler *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	product, err := handler.repository.FindProductCategory(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    product,
	})
}

func (handler *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	data := new(Products)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	err := json.Unmarshal(c.Body(), &data)

	if err != nil || len(data.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to create product",
			"error":   err,
		})
	}

	item, err := handler.repository.CreateProduct(*data)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to create product",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    item,
	})
}

func (handler *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "ID not found",
			"error":   err,
		})
	}

	product, err := handler.repository.FindProduct(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
		})
	}

	var p struct {
		CategoryId int64   `json:"categoryId"`
		Name       string  `json:"name"`
		Image      string  `json:"image"`
		Price      float64 `json:"price"`
		Stock      int64   `json:"stock"`
	}
	err = json.Unmarshal(c.Body(), &p)
	product.CategoryId = p.CategoryId
	product.Name = p.Name
	product.Image = p.Image
	product.Price = p.Price
	product.Stock = p.Stock

	err = handler.repository.SaveProduct(product)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func (handler *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting product",
			"err":     err,
		})
	}
	RowsAffected := handler.repository.DeleteProduct(id)
	if RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func NewProductHandler(repository *ProductRepository) *ProductHandler {
	return &ProductHandler{
		repository: repository,
	}
}
