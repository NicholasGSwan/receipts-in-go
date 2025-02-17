package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model       `json:"-"`
	ReceiptId        uint
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
