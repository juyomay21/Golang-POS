package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ittechman101/go-pos/models"
)

type Order struct {
	Base
}

func (b *Order) GetAllOrder(c *fiber.Ctx) error {

	orderList := models.FindAllOrders(c)

	count := models.GetOrderCount(c)
	if len(orderList) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "order Not Found",
			"error":   fiber.Map{},
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": fiber.Map{
				"orders": orderList,
				"meta": fiber.Map{
					"total": count,
					"limit": c.Query("limit"),
					"skip":  c.Query("skip"),
				},
			},
		})
	}
}

func (b *Order) GetSubTotal(c *fiber.Ctx) error {

	data := new([]models.ReqDetailSubtotalOrder)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	err := json.Unmarshal(c.Body(), &data)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found",
		})
	}
	resOrderSubTotal := models.FindSubTotal(*data)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    resOrderSubTotal,
	})
}

func (b *Order) GetOrder(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var data = new([]models.ReqDetailSubtotalOrder)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}
	err := json.Unmarshal(c.Body(), &data)
	resOrderDetail, err := models.FindOrder(id, *data)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    resOrderDetail,
	})
}

func (b *Order) GetOrderStatus(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var data = new([]models.ReqDetailSubtotalOrder)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}
	err := json.Unmarshal(c.Body(), &data)
	resOrderDetail, err := models.FindOrder(id, *data)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found",
		})
	}
	var isDownload struct {
		IsDownload bool `json:"isDownload"`
	}

	isDownload.IsDownload = false
	if resOrderDetail.DetailOrder.OrderId != 0 {
		isDownload.IsDownload = true
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    isDownload,
	})
}

func (b *Order) CreateOrder(c *fiber.Ctx) error {
	data := new(models.ReqCreateOrder)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	err := json.Unmarshal(c.Body(), &data)
	if err != nil || len(data.Products) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to create order",
			"error":   err,
		})
	}

	item, err := models.CreateOrder(*data)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed to create order",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    item,
	})
}
