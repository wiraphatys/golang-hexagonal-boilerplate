package order

import "gorm.io/gorm"

// Represents the product schema in database
type ProductModel struct {
	gorm.Model
	Name        string
	Description string
	Price       float32
	Stock       uint
}

func (p *ProductModel) TableName() string {
	return "Product"
}
