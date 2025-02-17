package models

import "gorm.io/gorm"

//import "time"

type Receipt struct {
	gorm.Model   `json:"-"`
	Retailer     string `json:"retailer"`
	PurchaseTime string `json:"purchaseTime"`
	PurchaseDate string `json:"purchaseDate"`
	Items        []Item `json:"items"`
	ReceiptTotal string `json:"total"`
}
