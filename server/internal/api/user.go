package api

import (
	"devprompt-ai/internal/middleware"
	"devprompt-ai/internal/response"
	"devprompt-ai/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	genService *service.GeneratorService
}

func NewUserHandler(genService *service.GeneratorService) *UserHandler {
	return &UserHandler{genService: genService}
}

// GenerateStats GET /api/user/generate-stats
func (h *UserHandler) GenerateStats(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	stats, err := h.genService.GetUserStats(userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, stats)
}
