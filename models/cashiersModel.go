package models

import "github.com/jinzhu/gorm"

type Cashiers struct {
	gorm.Model
	CashierId int64  `gorm:"Not Null" json:"cashierId"`
	Name      string `gorm:"Not Null" json:"name"`
	Passcode  string `gorm:"Not Null" json:"passcode"`
}
