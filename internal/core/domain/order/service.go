package order

type orderService struct {
	orderDbRepo OrderDbRepository
}

func NewOrderService(orderDbRepo OrderDbRepository) *orderService {
	return &orderService{orderDbRepo: orderDbRepo}
}
