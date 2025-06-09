package product

import "time"

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"desc"`
	Price       float32 `json:"price"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}
