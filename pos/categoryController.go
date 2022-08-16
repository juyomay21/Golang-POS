package pos

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type CategoryRepository struct {
	database *gorm.DB
}

func (repository *CategoryRepository) FindAll(c *fiber.Ctx) []Category {
	var categories []Category

	db := repository.database
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}
	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}
	db.Find(&categories)
	return categories
}

func (repository *CategoryRepository) GetCategoryCount() int64 {
	var categories []Category
	count := repository.database.Find(&categories).RowsAffected
	return count
}

func (repository *CategoryRepository) Find(id int) (Category, error) {
	var category Category
	err := repository.database.Where("category_id = ?", id).First(&category).Error
	if err != nil {
		err = errors.New("category not found")
	}
	return category, err
}

func (repository *CategoryRepository) Create(category Category) (Category, error) {
	var maxCategory Category

	repository.database.Raw(`
		SELECT COALESCE(MAX(category_id) + 1, 0) as category_id
		FROM categories
		`).Scan(
		&maxCategory,
	)
	category.CategoryId = maxCategory.CategoryId
	err := repository.database.Create(&category).Error
	if err != nil {
		return category, err
	}

	return category, nil
}

func (repository *CategoryRepository) Save(category Category) error {
	err := repository.database.Table("categories").Where("category_id = ?", category.CategoryId).Update(category).Error
	return err
}

func (repository *CategoryRepository) Delete(id int) int64 {
	count := repository.database.Where("category_id = ?", id).Delete(&Category{}).RowsAffected
	return count
}

func NewCategoryRepository(database *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		database: database,
	}
}
