package admin

import (
	"strconv"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type CategoryAdminHandler struct {
	categoryRepo *repository.PromptCategoryRepository
}

func NewCategoryAdminHandler(categoryRepo *repository.PromptCategoryRepository) *CategoryAdminHandler {
	return &CategoryAdminHandler{categoryRepo: categoryRepo}
}

// List GET /api/admin/categories
func (h *CategoryAdminHandler) List(c *gin.Context) {
	categories, total, err := h.categoryRepo.ListAll()
	if err != nil {
		response.InternalError(c, err.Error())
return
		}

	response.Success(c, &response.PaginatedData{
		List:     categories,
		Total:    total,
		Page:     1,
		PageSize: len(categories),
	})
}

// Create POST /api/admin/categories
func (h *CategoryAdminHandler) Create(c *gin.Context) {
	var cat model.PromptCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
return
		}

	if err := h.categoryRepo.Create(&cat); err != nil {
		response.InternalError(c, err.Error())
return
		}

	response.Success(c, cat)
}

// Update PUT /api/admin/categories/:id
func (h *CategoryAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
return
		}

	var cat model.PromptCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
return
		}

		existing, err := h.categoryRepo.FindByID(id)
		if err != nil || existing == nil {
			response.NotFound(c, "分类不存在")
			return
		}

	cat.CreatedAt = existing.CreatedAt
	cat.ID = id
	if err := h.categoryRepo.Update(&cat); err != nil {
		response.InternalError(c, err.Error())
return
		}

	response.Success(c, cat)
}

// Delete DELETE /api/admin/categories/:id
func (h *CategoryAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
return
		}

	if err := h.categoryRepo.Delete(id); err != nil {
			response.NotFound(c, err.Error())
return
		}

	response.Success(c, nil)
}
