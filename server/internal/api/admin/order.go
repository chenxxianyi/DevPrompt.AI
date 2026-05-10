package admin

import (
	"strconv"

	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type OrderAdminHandler struct {
	orderRepo *repository.OrderRepository
}

func NewOrderAdminHandler(orderRepo *repository.OrderRepository) *OrderAdminHandler {
	return &OrderAdminHandler{orderRepo: orderRepo}
}

// List GET /api/admin/orders
func (h *OrderAdminHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	orders, total, err := h.orderRepo.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &response.PaginatedData{
		List:     orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
