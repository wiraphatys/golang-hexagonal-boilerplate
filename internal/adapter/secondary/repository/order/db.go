package order

import "gorm.io/gorm"

type orderDbRepository struct {
	db *gorm.DB
}

func NewOrderDbRepository(db *gorm.DB) *orderDbRepository {
	return &orderDbRepository{db: db}
}
