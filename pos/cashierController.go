package pos

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type CashierRepository struct {
	database *gorm.DB
}

func (repository *CashierRepository) FindAllCashier(c *fiber.Ctx) []Cashiers {
	var cashiers []Cashiers

	db := repository.database
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}
	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}
	db.Find(&cashiers)
	return cashiers
}

func (repository *CashierRepository) GetCashierCount() int64 {
	var cashiers []Cashiers
	count := repository.database.Find(&cashiers).RowsAffected
	return count
}

func (repository *CashierRepository) FindCashier(id int) (Cashiers, error) {
	var cashier Cashiers
	err := repository.database.Where("cashier_id = ?", id).First(&cashier).Error
	if err != nil {
		err = errors.New("cashier not found")
	}
	return cashier, err
}

func (repository *CashierRepository) Passcode(id int) (string, error) {
	var cashier Cashiers
	err := repository.database.Where("cashier_id = ?", id).First(&cashier).Error
	if err != nil {
		err = errors.New("cashier not found")
	}
	return cashier.Passcode, err
}

func (repository *CashierRepository) CreateCashier(cashier Cashiers) (Cashiers, error) {
	// Get Max cashierId
	var maxCashier Cashiers

	repository.database.Raw(`
		SELECT COALESCE(MAX(cashier_id) + 1, 0) as cashier_id
		FROM cashiers
		`).Scan(
		&maxCashier,
	)
	cashier.Passcode = strconv.Itoa(rand.Intn(899999) + 100000)
	cashier.CashierId = maxCashier.CashierId
	err := repository.database.Create(&cashier).Error
	if err != nil {
		return cashier, err
	}

	return cashier, nil
}

func (repository *CashierRepository) SaveCashier(cashier Cashiers) (Cashiers, error) {
	err := repository.database.Table("cashiers").Where("cashier_id = ?", cashier.CashierId).Update(cashier).Error
	return cashier, err
}

func (repository *CashierRepository) DeleteCashier(id int) int64 {
	count := repository.database.Where("cashier_id = ?", id).Delete(&Cashiers{}).RowsAffected
	return count
}

func NewCashierRepository(database *gorm.DB) *CashierRepository {
	return &CashierRepository{
		database: database,
	}
}
