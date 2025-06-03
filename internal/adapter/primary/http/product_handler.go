package http

import (
	core "basedir/internal/core/domain/product"

	"github.com/gofiber/fiber/v2"
)

type ProductHttpHandler struct {
	orderSvc core.ProductService
}

func NewProductHttpHandler(orderSvc core.ProductService) *ProductHttpHandler {
	return &ProductHttpHandler{orderSvc: orderSvc}
}

func (h *ProductHttpHandler) GetProductByID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"msg": "get product by id successful."})
}
