package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Category struct {
	Base
}

func (b *Category) GetAllCategory(c *fiber.Ctx) error {
	var categories []models.Categories = models.FindAllCategory(c)
	count := models.GetCategoryCount()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"categories": categories,
			"meta": fiber.Map{
				"total": count,
				"limit": c.Query("limit"),
				"skip":  c.Query("skip"),
			},
		},
	})
}

func (b *Category) GetCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	category, err := models.FindCategory(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    category,
	})
}

func (b *Category) CreateCategory(c *fiber.Ctx) error {
	data := new(models.Categories)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	var p struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(c.Body(), &p)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	data.Name = p.Name
	if len(data.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	item, err := models.CreateCategory(*data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed creating item",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    item,
	})
}

func (b *Category) UpdateCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "ID not found",
			"error":   err,
		})
	}

	category, err := models.FindCategory(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
		})
	}

	var p struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(c.Body(), &p)
	category.Name = p.Name
	if len(category.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	save_err := models.SaveCategory(category)
	if save_err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error updating category",
			"error":   err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func (b *Category) DeleteCategory(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting category",
			"err":     err,
		})
	}

	RowsAffected := models.DeleteCategory(id)
	if RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
