package models

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	repository *PaymentRepository
	Base
}

func (handler *PaymentHandler) GetAll(c *fiber.Ctx) error {

	var payments []Payment = handler.repository.FindAll(c)
	count := handler.repository.GetPaymentsCount()
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

func (handler *PaymentHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	payment, err := handler.repository.Find(id)

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

func NewPaymentHandler(repository *PaymentRepository) *PaymentHandler {
	return &PaymentHandler{
		repository: repository,
	}
}
