package pos

import "github.com/jinzhu/gorm"

type Products struct {
	gorm.Model
	ProductId  int64   `gorm:"Not Null" json:"productid"`
	Name       string  `json:"name"`
	Stock      int64   `json:"stock"`
	Price      float64 `json:"price"`
	Image      string  `json:"image"`
	Sku        string  `json:"SKU"`
	CategoryId int64   `json:"categoryId"`
	DiscountId int64   `json:"discount"`
}
