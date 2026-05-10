package api

import (
	"encoding/json"
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

	inputBytes, _ := json.Marshal(req)
	systemPrompt := buildProjectPrompt(req.TargetAiTool)

	gp, err := h.genService.Generate(userID, "project", req.ProjectName, json.RawMessage(inputBytes), systemPrompt)
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

	inputBytes, _ := json.Marshal(req)
	systemPrompt := buildCursorRulesPrompt()

	gp, err := h.genService.Generate(userID, "cursor-rules",
		req.Language+" "+req.Framework+" Cursor Rules",
		json.RawMessage(inputBytes), systemPrompt)
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

	inputBytes, _ := json.Marshal(req)
	systemPrompt := buildClaudeCodePrompt()

	gp, err := h.genService.Generate(userID, "claude-code", req.Task,
		json.RawMessage(inputBytes), systemPrompt)
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

	inputBytes, _ := json.Marshal(req)
	systemPrompt := buildOptimizePrompt(req.OptimizeLevel)

	gp, err := h.genService.Generate(userID, "optimize", "Prompt 优化",
		json.RawMessage(inputBytes), systemPrompt)
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

// --- System Prompt Builders ---

func buildProjectPrompt(tool string) string {
	return "你是一个资深的软件架构师和技术顾问。\n" +
		"请根据用户提供的项目信息，生成一份详细的" + tool + "项目开发 Prompt。\n" +
		"输出应包含：项目概述、技术架构、目录结构、核心模块设计、开发步骤、注意事项。\n" +
		"用 Markdown 格式输出，确保内容专业、结构化、可执行。"
}

func buildCursorRulesPrompt() string {
	return "你是一个 Cursor IDE 配置专家。\n" +
		"请根据用户提供的技术栈和编码规范，生成一份完整的 .cursorrules 配置文件。\n" +
		"输出应包含：AI 行为指令、代码风格规则、命名规范、文件组织规则、最佳实践建议。\n" +
		"用清晰的段落结构输出，每条规则用 - 开头。"
}

func buildClaudeCodePrompt() string {
	return "你是一个 Claude Code CLI 工具的使用专家。\n" +
		"请根据用户描述的任务和上下文，生成一份详细的 Claude Code 任务 Prompt。\n" +
		"输出应包含：任务目标、上下文背景、具体要求、验收标准、注意事项。\n" +
		"确保 Prompt 清晰、完整，可以直接复制到 Claude Code 中使用。"
}

func buildOptimizePrompt(level string) string {
	levelDesc := "基础优化"
	if level == "professional" {
		levelDesc = "专业级优化"
	} else if level == "expert" {
		levelDesc = "专家级优化"
	}

	return "你是一个 Prompt Engineering 专家。\n" +
		"请对用户提供的原始 Prompt 进行" + levelDesc + "。\n" +
		"优化方向：明确角色定位、结构化指令、添加约束条件、指定输出格式、补充上下文。\n" +
		"先输出优化后的 Prompt，然后用 --- 分隔，再输出优化说明。"
}
