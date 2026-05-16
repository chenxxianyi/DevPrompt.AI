package api

import (
	"strconv"

	"devprompt-ai/internal/middleware"
	"devprompt-ai/internal/response"
	"devprompt-ai/internal/service"

	"github.com/gin-gonic/gin"
)

type GeneratorHandler struct {
	genService *service.GeneratorService
}

func NewGeneratorHandler(genService *service.GeneratorService) *GeneratorHandler {
	return &GeneratorHandler{genService: genService}
}

// Project POST /api/generator/project
func (h *GeneratorHandler) Project(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req service.GenerateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	qualityOpts := service.ExtractQualityOptions(
		req.QualityMode, req.OutputFormat,
		req.IncludeAcceptanceCriteria, req.IncludeRiskCheck,
		req.IncludeTestPlan, req.IncludeDeploymentNotes,
	)

	gp, err := h.genService.Generate(userID, "project", req.ProjectName, req, qualityOpts)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gp)
}

// CursorRules POST /api/generator/cursor-rules
func (h *GeneratorHandler) CursorRules(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req service.GenerateCursorRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	qualityOpts := service.ExtractQualityOptions(
		req.QualityMode, req.OutputFormat,
		req.IncludeAcceptanceCriteria, req.IncludeRiskCheck,
		req.IncludeTestPlan, req.IncludeDeploymentNotes,
	)

	gp, err := h.genService.Generate(userID, "cursor-rules",
		req.Language+" "+req.Framework+" Cursor Rules", req, qualityOpts)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gp)
}

// ClaudeCode POST /api/generator/claude-code
func (h *GeneratorHandler) ClaudeCode(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req service.GenerateClaudeCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	qualityOpts := service.ExtractQualityOptions(
		req.QualityMode, req.OutputFormat,
		req.IncludeAcceptanceCriteria, req.IncludeRiskCheck,
		req.IncludeTestPlan, req.IncludeDeploymentNotes,
	)

	gp, err := h.genService.Generate(userID, "claude-code", req.Task, req, qualityOpts)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gp)
}

// Optimize POST /api/generator/optimize
func (h *GeneratorHandler) Optimize(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	var req service.GenerateOptimizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	qualityOpts := service.ExtractQualityOptions(
		req.QualityMode, req.OutputFormat,
		req.IncludeAcceptanceCriteria, req.IncludeRiskCheck,
		req.IncludeTestPlan, req.IncludeDeploymentNotes,
	)

	gp, err := h.genService.Generate(userID, "optimize", "Prompt 优化", req, qualityOpts)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gp)
}

// History GET /api/generator/history
func (h *GeneratorHandler) History(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	genType := c.Query("type")

	result, err := h.genService.GetHistory(userID, genType, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}
