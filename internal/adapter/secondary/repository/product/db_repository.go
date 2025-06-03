package product

import "gorm.io/gorm"

type productDbRepository struct {
	db *gorm.DB
}

func NewProductDbRepository(db *gorm.DB) *productDbRepository {
	return &productDbRepository{db: db}
}
