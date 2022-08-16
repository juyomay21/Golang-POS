package models

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	repository *CategoryRepository
	Base
}

func (handler *CategoryHandler) GetAll(c *fiber.Ctx) error {

	var categories []Category = handler.repository.FindAll(c)
	count := handler.repository.GetCategoryCount()
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

func (handler *CategoryHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	category, err := handler.repository.Find(id)

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

func (handler *CategoryHandler) Create(c *fiber.Ctx) error {
	data := new(Category)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	var p struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(c.Body(), &p)
	data.Name = p.Name
	if len(data.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	item, err := handler.repository.Create(*data)

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

func (handler *CategoryHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "ID not found",
			"error":   err,
		})
	}

	category, err := handler.repository.Find(id)

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

	save_err := handler.repository.Save(category)

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

func (handler *CategoryHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting category",
			"err":     err,
		})
	}
	RowsAffected := handler.repository.Delete(id)
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

func NewCategoryHandler(repository *CategoryRepository) *CategoryHandler {
	return &CategoryHandler{
		repository: repository,
	}
}
