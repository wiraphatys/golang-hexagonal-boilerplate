package http

import (
	core "basedir/internal/core/domain/order"

	"github.com/gofiber/fiber/v2"
)

type OrderHttpHandler struct {
	orderSvc core.OrderService
}

func NewOrderHttpHandler(orderSvc core.OrderService) *OrderHttpHandler {
	return &OrderHttpHandler{orderSvc: orderSvc}
}

func (h *OrderHttpHandler) GetOrderByID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"msg": "get order by id successful."})
}
