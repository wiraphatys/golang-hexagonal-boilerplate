package order

import "gorm.io/gorm"

type orderCacheRepository struct {
	db *gorm.DB
}

func NewOrderCacheRepository(db *gorm.DB) *orderCacheRepository {
	return &orderCacheRepository{db: db}
}
