package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	gorm.Model
	CategoryId int64  `gorm:"Not Null" json:"categoryId"`
	Name       string `json:"name"`
}

type CategoryList struct {
	CategoryId int64  `json:"categoryId"`
	Name       string `json:"name"`
}

func FindAllCategory(c *fiber.Ctx) []Categories {
	var categories []Categories

	db := GetDB()

	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}

	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}

	db.Find(&categories)

	return categories
}

func GetCategoryCount() int64 {
	var categories []Categories

	count := GetDB().Find(&categories).RowsAffected

	return count
}

func FindCategory(id int) (Categories, error) {
	var category Categories

	err := GetDB().Where("category_id = ?", id).First(&category).Error
	if err != nil {
		err = errors.New("category not found")
	}

	return category, err
}

func CreateCategory(category Categories) (Categories, error) {
	var maxCategory Categories

	GetDB().Raw(`
		SELECT COALESCE(MAX(category_id) + 1, 1) as category_id
		FROM categories
		`).Scan(
		&maxCategory,
	)

	category.CategoryId = maxCategory.CategoryId

	err := GetDB().Create(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func SaveCategory(category Categories) error {
	err := GetDB().Table("categories").Where("category_id = ?", category.CategoryId).Update(category).Error

	return err
}

func DeleteCategory(id int) int64 {
	count := GetDB().Where("category_id = ?", id).Delete(&Categories{}).RowsAffected

	return count
}
