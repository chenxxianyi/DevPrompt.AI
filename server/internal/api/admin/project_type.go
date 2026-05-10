package admin

import (
	"strconv"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type ProjectTypeAdminHandler struct {
	repo *repository.ProjectTypeRepository
}

func NewProjectTypeAdminHandler(repo *repository.ProjectTypeRepository) *ProjectTypeAdminHandler {
	return &ProjectTypeAdminHandler{repo: repo}
}

// List GET /api/admin/project-types
func (h *ProjectTypeAdminHandler) List(c *gin.Context) {
	types, total, err := h.repo.FindAll()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &response.PaginatedData{
		List:     types,
		Total:    total,
		Page:     1,
		PageSize: len(types),
	})
}

// Create POST /api/admin/project-types
func (h *ProjectTypeAdminHandler) Create(c *gin.Context) {
	var pt model.ProjectType
	if err := c.ShouldBindJSON(&pt); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.repo.Create(&pt); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, pt)
}

// Update PUT /api/admin/project-types/:id
func (h *ProjectTypeAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	existing, err := h.repo.FindByID(id)
	if err != nil || existing == nil {
		response.NotFound(c, "项目类型不存在")
		return
	}

	var pt model.ProjectType
	if err := c.ShouldBindJSON(&pt); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	pt.CreatedAt = existing.CreatedAt
	pt.ID = id
	if err := h.repo.Update(&pt); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, pt)
}

// Delete DELETE /api/admin/project-types/:id
func (h *ProjectTypeAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, nil)
}
