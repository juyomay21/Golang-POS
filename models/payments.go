package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Payments struct {
	gorm.Model
	PaymentId int64  `gorm:"Not Null" json:"paymentId"`
	Name      string `gorm:"Not Null" json:"name"`
	Type      string `gorm:"Not Null" json:"type"`
	Logo      string `gorm:"Not Null" json:"logo"`
	Card      string `gorm:"Not Null" json:"card"`
}

func FindAllPayment(c *fiber.Ctx) []Payments {
	var payments []Payments

	db := GetDB()
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}

	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}

	db.Find(&payments)

	return payments
}

func GetPaymentsCount() int64 {
	var payments []Payments

	count := GetDB().Find(&payments).RowsAffected

	return count
}

func FindPayment(id int) (Payments, error) {
	var payment Payments

	err := GetDB().Where("payment_id = ?", id).First(&payment).Error
	if err != nil {
		err = errors.New("payment not found")
	}

	return payment, err
}
