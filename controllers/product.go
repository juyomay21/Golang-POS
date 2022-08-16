package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Product struct {
	Base
}

func (b *Product) GetAllProduct(c *fiber.Ctx) error {

	// err := b.Auth(c)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"success": false,
	// 		"message": "Authentication Failed",
	// 	})
	// }

	var products []models.ProductList = models.FindAllProduct(c)

	count := models.GetProductCount(c)
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

func (b *Product) GetProduct(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	product, err := models.FindProductCategory(id)
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

func (b *Product) CreateProduct(c *fiber.Ctx) error {
	data := new(models.Products)

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

	item, err := models.CreateProduct(*data)
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

func (b *Product) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "ID not found",
			"error":   err,
		})
	}

	product, err := models.FindProduct(id)
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
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   fiber.Map{},
		})
	}

	product.CategoryId = p.CategoryId
	product.Name = p.Name
	product.Image = p.Image
	product.Price = p.Price
	product.Stock = p.Stock

	err = models.SaveProduct(product)
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

func (b *Product) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting product",
			"err":     err,
		})
	}

	RowsAffected := models.DeleteProduct(id)
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
