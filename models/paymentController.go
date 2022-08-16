package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type PaymentRepository struct {
	database *gorm.DB
}

func (repository *PaymentRepository) FindAll(c *fiber.Ctx) []Payment {
	var payments []Payment

	db := repository.database
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}
	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}
	db.Find(&payments)
	return payments
}

func (repository *PaymentRepository) GetPaymentsCount() int64 {
	var payments []Payment
	count := repository.database.Find(&payments).RowsAffected
	return count
}

func (repository *PaymentRepository) Find(id int) (Payment, error) {
	var payment Payment
	err := repository.database.Where("category_id = ?", id).First(&payment).Error
	if err != nil {
		err = errors.New("payment not found")
	}
	return payment, err
}

func NewPaymentRepository(database *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		database: database,
	}
}
