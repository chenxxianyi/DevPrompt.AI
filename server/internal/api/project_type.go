package api

import (
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type ProjectTypeHandler struct {
	repo *repository.ProjectTypeRepository
}

func NewProjectTypeHandler(repo *repository.ProjectTypeRepository) *ProjectTypeHandler {
	return &ProjectTypeHandler{repo: repo}
}

// ListPublic GET /api/project-types — 返回启用的项目类型列表
func (h *ProjectTypeHandler) ListPublic(c *gin.Context) {
	types, err := h.repo.FindActive()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, types)
}
