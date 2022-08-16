package models

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model
	CategoryId int64  `gorm:"Not Null" json:"categoryId"`
	Name       string `json:"name"`
}
