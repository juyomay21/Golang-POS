package pos

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ittechman101/go-pos/config"

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

func (handler *CashierHandler) GetAllCashier(c *fiber.Ctx) error {

	var cashiers []Cashiers = handler.repository.FindAllCashier(c)
	count := handler.repository.GetCashierCount()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": fiber.Map{
			"cashiers": cashiers,
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

		secretKey := config.Config("SECRET_KEY")
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
			"message": "ID not found",
			"error":   err,
		})
	}

	cashier, err := handler.repository.FindCashier(id)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
		})
	}

	var p struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(c.Body(), &p)
	cashier.Name = p.Name

	item, err := handler.repository.SaveCashier(cashier)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
		})
	}

	return c.JSON(item)
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
