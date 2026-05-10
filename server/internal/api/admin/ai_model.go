package admin

import (
	"strconv"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type AIModelAdminHandler struct {
	modelRepo *repository.AIModelRepository
}

func NewAIModelAdminHandler(modelRepo *repository.AIModelRepository) *AIModelAdminHandler {
	return &AIModelAdminHandler{modelRepo: modelRepo}
}

// List GET /api/admin/ai-models
func (h *AIModelAdminHandler) List(c *gin.Context) {
	models, total, err := h.modelRepo.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &response.PaginatedData{
		List:     models,
		Total:    total,
		Page:     1,
		PageSize: len(models),
	})
}

// Create POST /api/admin/ai-models
func (h *AIModelAdminHandler) Create(c *gin.Context) {
	var mdl model.AIModel
	if err := c.ShouldBindJSON(&mdl); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.modelRepo.Create(&mdl); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, mdl)
}

// Update PUT /api/admin/ai-models/:id
func (h *AIModelAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	existing, err := h.modelRepo.FindByID(id)
	if err != nil || existing == nil {
		response.NotFound(c, "模型不存在")
		return
	}

	var mdl model.AIModel
	if err := c.ShouldBindJSON(&mdl); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	mdl.ID = id
	if err := h.modelRepo.Update(&mdl); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, mdl)
}

// Delete DELETE /api/admin/ai-models/:id
func (h *AIModelAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.modelRepo.Delete(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, nil)
}
