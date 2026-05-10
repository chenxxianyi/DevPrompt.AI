package admin

import (
	"strconv"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/response"
	"devprompt-ai/internal/service"

	"github.com/gin-gonic/gin"
)

type PromptAdminHandler struct {
	promptService *service.PromptService
}

func NewPromptAdminHandler(promptService *service.PromptService) *PromptAdminHandler {
	return &PromptAdminHandler{promptService: promptService}
}

// List GET /api/admin/prompts
func (h *PromptAdminHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	result, err := h.promptService.AdminList(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// Create POST /api/admin/prompts
func (h *PromptAdminHandler) Create(c *gin.Context) {
	var t model.PromptTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.promptService.AdminCreate(&t); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, t)
}

// Update PUT /api/admin/prompts/:id
func (h *PromptAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var t model.PromptTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.promptService.AdminUpdate(id, &t); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, t)
}

// Delete DELETE /api/admin/prompts/:id
func (h *PromptAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.promptService.AdminDelete(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, nil)
}
