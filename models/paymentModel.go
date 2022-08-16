package models

import "github.com/jinzhu/gorm"

type Payment struct {
	gorm.Model
	PaymentId int64  `gorm:"Not Null" json:"paymentId"`
	Name      string `gorm:"Not Null" json:"name"`
	Type      string `gorm:"Not Null" json:"type"`
	Logo      string `gorm:"Not Null" json:"logo"`
	Card      string `gorm:"Not Null" json:"card"`
}
