package admin

import (
	"strconv"

	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type CallLogAdminHandler struct {
	callLogRepo *repository.AICallLogRepository
}

func NewCallLogAdminHandler(callLogRepo *repository.AICallLogRepository) *CallLogAdminHandler {
	return &CallLogAdminHandler{callLogRepo: callLogRepo}
}

// List GET /api/admin/ai-call-logs
func (h *CallLogAdminHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	logs, total, err := h.callLogRepo.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &response.PaginatedData{
		List:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
