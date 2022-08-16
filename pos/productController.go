package pos

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
)

type ProductRepository struct {
	database *gorm.DB
}

func (repository *ProductRepository) FindAllProduct(c *fiber.Ctx) []ProductList {
	var products []ProductList
	var product ProductList

	db := repository.database.Table("products").Select(
		"products.product_id, products.sku, products.name, products.stock, products.price, products.image, products.category_id, categories.name as category_name")
	db = db.Where("products.deleted_at is NULL")
	if len(c.Query("limit")) > 0 {
		db = db.Limit(c.Query("limit"))
	}
	if len(c.Query("skip")) > 0 {
		db = db.Offset(c.Query("skip"))
	}
	if len(c.Query("categoryId")) > 0 {
		db = db.Where("products.category_id = ?", c.Query("categoryId"))
	}
	if len(c.Query("q")) > 0 {
		db = db.Where("products.name LIKE ?", `%`+c.Query("q")+`%`)
	}
	rows, _ := db.Joins("left join categories on products.category_id=categories.category_id").Rows()
	defer rows.Close()

	for rows.Next() {
		rows.Scan(
			&product.ProductId,
			&product.Sku,
			&product.Name,
			&product.Stock,
			&product.Price,
			&product.Image,
			&product.Category.CategoryId,
			&product.Category.Name,
		)
		products = append(products, product)
	}
	return products
}

func (repository *ProductRepository) GetProductCount(c *fiber.Ctx) int64 {
	var products []Products
	db := repository.database
	if len(c.Query("categoryId")) > 0 {
		db = db.Where("category_id = ?", c.Query("categoryId"))
	}
	if len(c.Query("q")) > 0 {
		db = db.Where("name LIKE ?", `%`+c.Query("q")+`%`)
	}
	count := db.Find(&products).RowsAffected
	return count
}

func (repository *ProductRepository) FindProduct(id int) (Products, error) {
	var product Products
	err := repository.database.Where("product_id = ?", id).First(&product).Error
	if err != nil {
		err = errors.New("product not found")
	}
	return product, err
}

func (repository *ProductRepository) FindProductCategory(id int) (ProductList, error) {
	var product ProductList

	fmt.Println("---------- Test --------")

	db := repository.database.Table("products").Select(
		"products.product_id, products.sku, products.name, products.stock, products.price, products.image, products.category_id, categories.name as category_name")
	db = db.Where("products.deleted_at is NULL AND products.product_id = ?", id)
	rows, err := db.Joins("left join categories on products.category_id=categories.category_id").Rows()

	defer rows.Close()

	if !rows.Next() {
		return product, errors.New("Product Not Found")
	}

	rows.Scan(
		&product.ProductId,
		&product.Sku,
		&product.Name,
		&product.Stock,
		&product.Price,
		&product.Image,
		&product.Category.CategoryId,
		&product.Category.Name,
	)

	return product, err
}

func (repository *ProductRepository) CreateProduct(product Products) (Products, error) {
	// Get Max productId
	var maxProduct Products
	var category Category

	count := repository.database.Table("categories").Where("category_id = ?", product.CategoryId).Find(&category).RowsAffected
	if count == 0 {
		return product, errors.New("category not found")
	}

	repository.database.Raw(`
		SELECT COALESCE(MAX(product_id) + 1, 1) as product_id
		FROM products
		`).Scan(
		&maxProduct,
	)

	product.Sku = fmt.Sprintf("ID%03d", maxProduct.ProductId)
	product.ProductId = maxProduct.ProductId
	err := repository.database.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (repository *ProductRepository) SaveProduct(product Products) error {
	var category Category

	count := repository.database.Table("categories").Where("category_id = ?", product.CategoryId).Find(&category).RowsAffected
	if count == 0 {
		return errors.New("category not found")
	}

	err := repository.database.Table("products").Where("product_id = ?", product.ProductId).Update(product).Error
	return err
}

func (repository *ProductRepository) DeleteProduct(id int) int64 {
	count := repository.database.Where("product_id = ?", id).Delete(&Products{}).RowsAffected
	return count
}

func NewProductRepository(database *gorm.DB) *ProductRepository {
	return &ProductRepository{
		database: database,
	}
}
