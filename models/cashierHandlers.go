package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"
)

type CashierHandler struct {
	repository *CashierRepository
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

func (handler *CashierHandler) GetAllCashier(c *fiber.Ctx) error {

	var cashiers []Cashiers = handler.repository.FindAllCashier(c)

	//	fmt.Println(cashiers)

	cashiersResponses := make([]CashierResponse, len(cashiers))
	count := handler.repository.GetCashierCount()

	for i, element := range cashiers {
		cashiersResponses[i].CashierId = element.CashierId
		cashiersResponses[i].Name = element.Name
	}

	//	fmt.Println(cashiersResponses)

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

func (handler *CashierHandler) GetCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	cashier, err := handler.repository.FindCashier(id)

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

func (handler *CashierHandler) Passcode(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	passcode, err := handler.repository.Passcode(id)

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

func (handler *CashierHandler) Login(c *fiber.Ctx) error {
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

	passcode, err := handler.repository.Passcode(id)

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

func (handler *CashierHandler) CreateCashier(c *fiber.Ctx) error {
	data := new(Cashiers)

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

	item, err := handler.repository.CreateCashier(*data)

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

func (handler *CashierHandler) UpdateCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Cashier Not Found",
			"error":   fiber.Map{},
		})
	}

	cashier, err := handler.repository.FindCashier(id)

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
	err = json.Unmarshal(c.Body(), &p)
	cashier.Name = p.Name
	if len(cashier.Name) == 0 {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "body ValidationError: \"name\" is required"})
	}

	item, err := handler.repository.SaveCashier(cashier)
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

func (handler *CashierHandler) DeleteCashier(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  400,
			"message": "Failed deleting cashier",
			"err":     err,
		})
	}
	RowsAffected := handler.repository.DeleteCashier(id)
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

func NewCashierHandler(repository *CashierRepository) *CashierHandler {
	return &CashierHandler{
		repository: repository,
	}
}
