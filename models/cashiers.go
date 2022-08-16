package models

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Cashiers struct {
	gorm.Model
	CashierId int64  `gorm:"Not Null" json:"cashierId"`
	Name      string `gorm:"Not Null" json:"name"`
	Passcode  string `gorm:"Not Null" json:"passcode"`
}

func FindAllCashier(c *fiber.Ctx) []Cashiers {
	var cashiers []Cashiers

	db := GetDB()
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}
	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}
	db.Find(&cashiers)

	// rows, err := repository.database.Raw(`SELECT cashiers.cashier_id, cashiers.name FROM cashiers`).Rows()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	rows.Scan(
	// 		&cashier.Name,
	// 		&cashier.CashierId,
	// 	)
	// 	cashiers = append(cashiers, cashier)
	// }
	return cashiers
}

func GetCashierCount() int64 {
	var cashiers []Cashiers
	count := GetDB().Find(&cashiers).RowsAffected

	return count
}

func FindCashier(id int) (Cashiers, error) {
	var cashier Cashiers

	err := GetDB().Where("cashier_id = ?", id).First(&cashier).Error
	if err != nil {
		err = errors.New("cashier not found")
	}

	return cashier, err
}

func Passcode(id int) (string, error) {
	var cashier Cashiers

	err := GetDB().Where("cashier_id = ?", id).First(&cashier).Error
	if err != nil {
		err = errors.New("cashier not found")
	}

	return cashier.Passcode, err
}

func CreateCashier(cashier Cashiers) (Cashiers, error) {
	// Get Max cashierId
	var maxCashier Cashiers

	GetDB().Raw(`
		SELECT COALESCE(MAX(cashier_id) + 1, 1) as cashier_id
		FROM cashiers
		`).Scan(
		&maxCashier,
	)

	cashier.Passcode = strconv.Itoa(rand.Intn(899999) + 100000)
	cashier.CashierId = maxCashier.CashierId

	err := GetDB().Create(&cashier).Error
	if err != nil {
		return cashier, err
	}

	return cashier, nil
}

func SaveCashier(cashier Cashiers) (Cashiers, error) {
	err := GetDB().Table("cashiers").Where("cashier_id = ?", cashier.CashierId).Update(cashier).Error

	return cashier, err
}

func DeleteCashier(id int) int64 {
	count := GetDB().Where("cashier_id = ?", id).Delete(&Cashiers{}).RowsAffected

	return count
}
