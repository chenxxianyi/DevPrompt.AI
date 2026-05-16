package admin

import (
	"strconv"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"

	"github.com/gin-gonic/gin"
)

type PromptRecipeAdminHandler struct {
	recipeRepo *repository.PromptRecipeRepository
}

func NewPromptRecipeAdminHandler(recipeRepo *repository.PromptRecipeRepository) *PromptRecipeAdminHandler {
	return &PromptRecipeAdminHandler{recipeRepo: recipeRepo}
}

// List GET /api/admin/prompt-recipes
func (h *PromptRecipeAdminHandler) List(c *gin.Context) {
	recipeType := c.Query("type")
	targetTool := c.Query("targetTool")
	status := c.Query("status")

	recipes, total, err := h.recipeRepo.List(recipeType, targetTool, status)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, &response.PaginatedData{
		List:     recipes,
		Total:    total,
		Page:     1,
		PageSize: len(recipes),
	})
}

// Create POST /api/admin/prompt-recipes
func (h *PromptRecipeAdminHandler) Create(c *gin.Context) {
	var recipe model.PromptRecipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if recipe.Type == "" {
		response.BadRequest(c, "类型不能为空")
		return
	}
	if recipe.Name == "" {
		response.BadRequest(c, "名称不能为空")
		return
	}
	if recipe.SystemPrompt == "" {
		response.BadRequest(c, "System Prompt 不能为空")
		return
	}

	if err := h.recipeRepo.Create(&recipe); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, recipe)
}

// Update PUT /api/admin/prompt-recipes/:id
func (h *PromptRecipeAdminHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	existing, err := h.recipeRepo.FindByID(id)
	if err != nil || existing == nil {
		response.NotFound(c, "Recipe 不存在")
		return
	}

	var recipe model.PromptRecipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	recipe.ID = id
	recipe.CreatedAt = existing.CreatedAt
	if err := h.recipeRepo.Update(&recipe); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, recipe)
}

// Delete DELETE /api/admin/prompt-recipes/:id
func (h *PromptRecipeAdminHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.recipeRepo.Delete(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// SetDefault PUT /api/admin/prompt-recipes/:id/set-default
func (h *PromptRecipeAdminHandler) SetDefault(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.recipeRepo.SetDefault(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
