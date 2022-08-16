package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ittechman101/go-pos/models"

	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
)

type Cashier struct {
	Base
}

type JwtCustomClaims struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

type CashierResponse struct {
	CashierId int64  `json:"cashierId"`
	Name      string `json:"name"`
}

func (b *Cashier) GetAllCashier(c *fiber.Ctx) error {

	var cashiers []models.Cashiers = models.FindAllCashier(c)

	cashiersResponses := make([]CashierResponse, len(cashiers))
	count := models.GetCashierCount()

	for i, element := range cashiers {
		cashiersResponses[i].CashierId = element.CashierId
		cashiersResponses[i].Name = element.Name
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"cashiers": cashiersResponses,
			"meta": fiber.Map{
				"total": count,
				"limit": c.Query("limit"),
				"skip":  c.Query("skip"),
			},
		},
	})
}

func (b *Cashier) GetCashier(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	cashier, err := models.FindCashier(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashier,
	})
}

func (b *Cashier) Passcode(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	passcode, err := models.Passcode(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"passcode": passcode,
		},
	})
}

func (b *Cashier) Login(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "param validationError: \"cashierId\" is required"})
	}

	var p struct {
		Passcode string `json:"passcode"`
	}
	err = json.Unmarshal(c.Body(), &p)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "body validationError: \"passcode\" is required"})
	}

	passcode, err := models.Passcode(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
		})
	}

	if passcode == p.Passcode {
		//		secretKey := config.Config("SECRET_KEY")
		secretKey := "goPos"
		claims := jwt.MapClaims{
			"id":     id,
			"active": true,
			"exp":    time.Now().Add(time.Hour * 6).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "JWT Token Error"})
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": fiber.Map{
				"token": tokenString,
			},
		})
	} else {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Passcode Not Match",
		})
	}
}

func (b *Cashier) CreateCashier(c *fiber.Ctx) error {
	data := new(models.Cashiers)

	if len(c.Body()) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	var p struct {
		Name string `json:"name"`
	}
	err := json.Unmarshal(c.Body(), &p)
	if err != nil || len(p.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"success": false,
			"message": "body ValidationError: \"name\" is required",
			"error": fiber.Map{
				"message": "\"name\" is required",
				"path":    "name",
				"type":    "any.required",
				"context": fiber.Map{
					"label": "name",
					"key":   "name",
				},
			},
		})
	}

	data.Name = p.Name

	item, err := models.CreateCashier(*data)
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

func (b *Cashier) UpdateCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	cashier, err := models.FindCashier(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	var p struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(c.Body(), &p)
	cashier.Name = p.Name
	if len(cashier.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	item, err := models.SaveCashier(cashier)
	item = item
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func (b *Cashier) DeleteCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting cashier",
			"err":     err,
		})
	}

	RowsAffected := models.DeleteCashier(id)
	if RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
