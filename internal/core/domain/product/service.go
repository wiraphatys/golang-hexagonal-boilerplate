package product

type productService struct {
	orderDbRepo ProductDbRepository
}

func NewProductService(orderDbRepo ProductDbRepository) *productService {
	return &productService{orderDbRepo: orderDbRepo}
}
