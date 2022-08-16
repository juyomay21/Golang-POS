package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Payment struct {
	Base
}

func (b *Payment) GetAllPayment(c *fiber.Ctx) error {
	var payments []models.Payments = models.FindAllPayment(c)
	count := models.GetPaymentsCount()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"payments": payments,
			"meta": fiber.Map{
				"total": count,
				"limit": c.Query("limit"),
				"skip":  c.Query("skip"),
			},
		},
	})
}

func (b *Payment) GetPayment(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	payment, err := models.FindPayment(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    payment,
	})
}
