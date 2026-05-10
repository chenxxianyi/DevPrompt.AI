package admin

import (
	"strconv"

	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type UserAdminHandler struct {
	userRepo *repository.UserRepository
}

func NewUserAdminHandler(userRepo *repository.UserRepository) *UserAdminHandler {
	return &UserAdminHandler{userRepo: userRepo}
}

// List GET /api/admin/users
func (h *UserAdminHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	users, total, err := h.userRepo.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &response.PaginatedData{
		List:     users,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
