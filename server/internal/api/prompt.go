package api

import (
	"strconv"

	"devprompt-ai/internal/middleware"
	"devprompt-ai/internal/response"
	"devprompt-ai/internal/service"

	"github.com/gin-gonic/gin"
)

type PromptHandler struct {
	promptService *service.PromptService
}

func NewPromptHandler(promptService *service.PromptService) *PromptHandler {
	return &PromptHandler{promptService: promptService}
}

// List GET /api/prompts
func (h *PromptHandler) List(c *gin.Context) {
	keyword := c.Query("keyword")
	category := c.Query("category")
	sort := c.Query("sort")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var userID uint64
	if uid, ok := middleware.GetUserID(c); ok {
		userID = uid
	}

	result, err := h.promptService.List(keyword, category, sort, page, pageSize, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetBySlug GET /api/prompts/:slug
func (h *PromptHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var userID uint64
	if uid, ok := middleware.GetUserID(c); ok {
		userID = uid
	}

	template, err := h.promptService.GetBySlug(slug, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, template)
}

// ToggleLike POST /api/prompts/:id/like
func (h *PromptHandler) ToggleLike(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	promptID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	liked, err := h.promptService.ToggleLike(userID, promptID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"liked": liked})
}

// ListCategories GET /api/categories
func (h *PromptHandler) ListCategories(c *gin.Context) {
	categories, err := h.promptService.ListCategories()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, &response.PaginatedData{
		List:     categories,
		Total:    int64(len(categories)),
		Page:     1,
		PageSize: len(categories),
	})
}

// ToggleFavorite POST /api/prompts/:id/favorite
func (h *PromptHandler) ToggleFavorite(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	promptID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	favorited, err := h.promptService.ToggleFavorite(userID, promptID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"favorited": favorited})
}
