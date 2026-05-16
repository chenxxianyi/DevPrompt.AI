package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/provider"
	"devprompt-ai/internal/recipe"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"
)

const (
	defaultGenerateMaxTokens = 4096
	expertGenerateMaxTokens  = 8192
)

// GeneratorService handles prompt generation requests.
type GeneratorService struct {
	providerMgr   *provider.ProviderManager
	genRepo       *repository.GeneratedPromptRepository
	callLogRepo   *repository.AICallLogRepository
	userRepo      *repository.UserRepository
	planRepo      *repository.MembershipPlanRepository
	rateLimitSvc  *RateLimitService
	membershipSvc *MembershipService
	recipeEngine  *recipe.Engine
	defaultModel  string
}

func NewGeneratorService(
	providerMgr *provider.ProviderManager,
	genRepo *repository.GeneratedPromptRepository,
	callLogRepo *repository.AICallLogRepository,
	userRepo *repository.UserRepository,
	planRepo *repository.MembershipPlanRepository,
	rateLimitSvc *RateLimitService,
	membershipSvc *MembershipService,
	defaultModel string,
) *GeneratorService {
	return &GeneratorService{
		providerMgr:   providerMgr,
		genRepo:       genRepo,
		callLogRepo:   callLogRepo,
		userRepo:      userRepo,
		planRepo:      planRepo,
		rateLimitSvc:  rateLimitSvc,
		membershipSvc: membershipSvc,
		recipeEngine:  recipe.NewEngine(),
		defaultModel:  defaultModel,
	}
}

type GenerateProjectRequest struct {
	ProjectName  string   `json:"projectName"`
	ProjectType  string   `json:"projectType"`
	TechStack    []string `json:"techStack"`
	Features     []string `json:"features"`
	TargetAiTool string   `json:"targetAiTool"`

	QualityMode               *string `json:"qualityMode,omitempty"`
	OutputFormat              *string `json:"outputFormat,omitempty"`
	IncludeAcceptanceCriteria *bool   `json:"includeAcceptanceCriteria,omitempty"`
	IncludeRiskCheck          *bool   `json:"includeRiskCheck,omitempty"`
	IncludeTestPlan           *bool   `json:"includeTestPlan,omitempty"`
	IncludeDeploymentNotes    *bool   `json:"includeDeploymentNotes,omitempty"`
}

type GenerateCursorRulesRequest struct {
	Language  string   `json:"language"`
	Framework string   `json:"framework"`
	CodeStyle string   `json:"codeStyle"`
	Rules     []string `json:"rules"`

	QualityMode               *string `json:"qualityMode,omitempty"`
	OutputFormat              *string `json:"outputFormat,omitempty"`
	IncludeAcceptanceCriteria *bool   `json:"includeAcceptanceCriteria,omitempty"`
	IncludeRiskCheck          *bool   `json:"includeRiskCheck,omitempty"`
	IncludeTestPlan           *bool   `json:"includeTestPlan,omitempty"`
	IncludeDeploymentNotes    *bool   `json:"includeDeploymentNotes,omitempty"`
}

type GenerateClaudeCodeRequest struct {
	Task         string   `json:"task"`
	Context      string   `json:"context"`
	Requirements []string `json:"requirements"`

	QualityMode               *string `json:"qualityMode,omitempty"`
	OutputFormat              *string `json:"outputFormat,omitempty"`
	IncludeAcceptanceCriteria *bool   `json:"includeAcceptanceCriteria,omitempty"`
	IncludeRiskCheck          *bool   `json:"includeRiskCheck,omitempty"`
	IncludeTestPlan           *bool   `json:"includeTestPlan,omitempty"`
	IncludeDeploymentNotes    *bool   `json:"includeDeploymentNotes,omitempty"`
}

type GenerateOptimizeRequest struct {
	RawPrompt     string `json:"rawPrompt"`
	TargetTool    string `json:"targetTool"`
	OptimizeLevel string `json:"optimizeLevel"`

	QualityMode               *string `json:"qualityMode,omitempty"`
	OutputFormat              *string `json:"outputFormat,omitempty"`
	IncludeAcceptanceCriteria *bool   `json:"includeAcceptanceCriteria,omitempty"`
	IncludeRiskCheck          *bool   `json:"includeRiskCheck,omitempty"`
	IncludeTestPlan           *bool   `json:"includeTestPlan,omitempty"`
	IncludeDeploymentNotes    *bool   `json:"includeDeploymentNotes,omitempty"`
}

func (s *GeneratorService) Generate(userID uint64, genType, title string, input interface{}, qualityOpts recipe.QualityOptions) (*model.GeneratedPrompt, error) {
	startTime := time.Now()

	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if user.Status != "active" {
		return nil, fmt.Errorf("账号已被禁用")
	}
	if user.Role != "admin" && !CheckMembership(user) {
		return nil, fmt.Errorf("会员已过期，请续费后使用")
	}

	qualityOpts, err = normalizeQualityOptions(user, input, qualityOpts)
	if err != nil {
		return nil, err
	}

	recipeInput, err := buildRecipeInput(genType, input)
	if err != nil {
		return nil, fmt.Errorf("请求参数不合法: %w", err)
	}

	if user.Role != "admin" {
		dailyLimit := s.membershipSvc.GetDailyLimit(user.MembershipLevel)
		count, limit, allowed, err := s.rateLimitSvc.CheckDailyLimit(userID, dailyLimit)
		if err != nil {
			return nil, fmt.Errorf("限流检查失败: %w", err)
		}
		if !allowed {
			return nil, fmt.Errorf("今日生成次数已达上限 (%d/%d)，请明天再试", count, limit)
		}

		ok, err := s.rateLimitSvc.CheckRateLimit(userID, genType, 30)
		if err != nil || !ok {
			return nil, fmt.Errorf("请求过于频繁，请稍后重试")
		}
	}

	built, err := s.recipeEngine.Build(recipe.BuildContext{
		Type:    recipe.GeneratorType(genType),
		Input:   recipeInput,
		Quality: qualityOpts,
	})
	if err != nil {
		return nil, fmt.Errorf("构建 Prompt 失败: %w", err)
	}

	resp, err := s.providerMgr.Call(provider.ChatRequest{
		Messages: []provider.ChatMessage{
			{Role: "system", Content: built.SystemPrompt},
			{Role: "user", Content: built.UserPrompt},
		},
		Temperature: 0.7,
		MaxTokens:   maxTokensForGenerate(input, qualityOpts),
	}, s.defaultModel)

	callLog := &model.AICallLog{
		UserID:      userID,
		Provider:    "",
		Model:       "",
		RequestType: genType,
		Status:      "success",
		LatencyMs:   int(time.Since(startTime).Milliseconds()),
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	if err != nil {
		callLog.Status = "failed"
		callLog.ErrorMessage = err.Error()
		_ = s.callLogRepo.Create(callLog)
		return nil, fmt.Errorf("AI 调用失败: %w", err)
	}

	callLog.Provider = resp.Provider
	callLog.Model = resp.Model
	callLog.PromptTokens = resp.PromptTokens
	callLog.CompletionTokens = resp.CompletionTokens
	callLog.TotalTokens = resp.TotalTokens

	// The upstream call already consumed quota and cost at this point.
	_ = s.rateLimitSvc.IncrementDailyCount(userID)

	output := resp.Content
	if qualityOpts.OutputFormat == "json" {
		output, err = normalizeJSONOutput(resp.Content)
		if err != nil {
			callLog.Status = "failed"
			callLog.ErrorMessage = "AI returned invalid JSON: " + err.Error()
			_ = s.callLogRepo.Create(callLog)
			return nil, fmt.Errorf("AI 返回的 JSON 格式无效，请重试")
		}
	}

	_ = s.callLogRepo.Create(callLog)

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败: %w", err)
	}

	gp := &model.GeneratedPrompt{
		UserID:        userID,
		Type:          genType,
		Title:         title,
		Input:         string(inputJSON),
		Output:        output,
		Model:         resp.Model,
		Provider:      resp.Provider,
		Tokens:        resp.TotalTokens,
		RecipeID:      built.RecipeID,
		RecipeVersion: built.RecipeVersion,
	}
	if err := s.genRepo.Create(gp); err != nil {
		return nil, fmt.Errorf("保存生成记录失败: %w", err)
	}

	return gp, nil
}

func ExtractQualityOptions(qm *string, of *string, ac *bool, rc *bool, tp *bool, dn *bool) recipe.QualityOptions {
	opts := recipe.DefaultQualityOptions()
	if qm != nil {
		opts.QualityMode = *qm
	}
	if of != nil {
		opts.OutputFormat = *of
	}
	if ac != nil {
		opts.IncludeAcceptanceCriteria = *ac
	}
	if rc != nil {
		opts.IncludeRiskCheck = *rc
	}
	if tp != nil {
		opts.IncludeTestPlan = *tp
	}
	if dn != nil {
		opts.IncludeDeploymentNotes = *dn
	}
	return opts
}

func buildRecipeInput(genType string, input interface{}) (interface{}, error) {
	switch genType {
	case string(recipe.GeneratorProject):
		switch req := input.(type) {
		case GenerateProjectRequest:
			return recipe.ProjectInput{
				ProjectName:  req.ProjectName,
				ProjectType:  req.ProjectType,
				TechStack:    req.TechStack,
				Features:     req.Features,
				TargetAiTool: req.TargetAiTool,
			}, nil
		case *GenerateProjectRequest:
			if req == nil {
				return nil, fmt.Errorf("project request is nil")
			}
			return recipe.ProjectInput{
				ProjectName:  req.ProjectName,
				ProjectType:  req.ProjectType,
				TechStack:    req.TechStack,
				Features:     req.Features,
				TargetAiTool: req.TargetAiTool,
			}, nil
		case recipe.ProjectInput, *recipe.ProjectInput:
			return input, nil
		}

	case string(recipe.GeneratorCursorRules):
		switch req := input.(type) {
		case GenerateCursorRulesRequest:
			return recipe.CursorRulesInput{
				Language:  req.Language,
				Framework: req.Framework,
				CodeStyle: req.CodeStyle,
				Rules:     req.Rules,
			}, nil
		case *GenerateCursorRulesRequest:
			if req == nil {
				return nil, fmt.Errorf("cursor rules request is nil")
			}
			return recipe.CursorRulesInput{
				Language:  req.Language,
				Framework: req.Framework,
				CodeStyle: req.CodeStyle,
				Rules:     req.Rules,
			}, nil
		case recipe.CursorRulesInput, *recipe.CursorRulesInput:
			return input, nil
		}

	case string(recipe.GeneratorClaudeCode):
		switch req := input.(type) {
		case GenerateClaudeCodeRequest:
			return recipe.ClaudeCodeInput{
				Task:         req.Task,
				Context:      req.Context,
				Requirements: req.Requirements,
			}, nil
		case *GenerateClaudeCodeRequest:
			if req == nil {
				return nil, fmt.Errorf("claude code request is nil")
			}
			return recipe.ClaudeCodeInput{
				Task:         req.Task,
				Context:      req.Context,
				Requirements: req.Requirements,
			}, nil
		case recipe.ClaudeCodeInput, *recipe.ClaudeCodeInput:
			return input, nil
		}

	case string(recipe.GeneratorOptimize):
		level, err := requestedOptimizeLevel(input)
		if err != nil {
			return nil, err
		}

		switch req := input.(type) {
		case GenerateOptimizeRequest:
			return recipe.OptimizeInput{
				RawPrompt:     req.RawPrompt,
				TargetTool:    req.TargetTool,
				OptimizeLevel: level,
			}, nil
		case *GenerateOptimizeRequest:
			if req == nil {
				return nil, fmt.Errorf("optimize request is nil")
			}
			return recipe.OptimizeInput{
				RawPrompt:     req.RawPrompt,
				TargetTool:    req.TargetTool,
				OptimizeLevel: level,
			}, nil
		case recipe.OptimizeInput:
			req.OptimizeLevel = level
			return req, nil
		case *recipe.OptimizeInput:
			if req == nil {
				return nil, fmt.Errorf("optimize input is nil")
			}
			req.OptimizeLevel = level
			return req, nil
		}
	}

	return nil, fmt.Errorf("unsupported %s input type %T", genType, input)
}

func normalizeQualityOptions(user *model.User, input interface{}, opts recipe.QualityOptions) (recipe.QualityOptions, error) {
	defaults := recipe.DefaultQualityOptions()

	opts.QualityMode = strings.ToLower(strings.TrimSpace(opts.QualityMode))
	if opts.QualityMode == "" {
		opts.QualityMode = defaults.QualityMode
	}
	if !recipe.IsValidQualityMode(opts.QualityMode) {
		return opts, fmt.Errorf("不支持的 qualityMode: %s", opts.QualityMode)
	}

	opts.OutputFormat = strings.ToLower(strings.TrimSpace(opts.OutputFormat))
	if opts.OutputFormat == "" {
		opts.OutputFormat = defaults.OutputFormat
	}
	if !recipe.IsValidOutputFormat(opts.OutputFormat) {
		return opts, fmt.Errorf("不支持的 outputFormat: %s", opts.OutputFormat)
	}

	optimizeLevel, err := requestedOptimizeLevel(input)
	if err != nil {
		return opts, err
	}

	if user.Role == "admin" || user.MembershipLevel != "free" {
		return opts, nil
	}

	if opts.QualityMode == "expert" {
		return opts, fmt.Errorf("免费用户不支持 expert 质量模式")
	}
	if opts.OutputFormat != "markdown" {
		return opts, fmt.Errorf("免费用户仅支持 markdown 输出格式")
	}
	if opts.IncludeAcceptanceCriteria || opts.IncludeRiskCheck || opts.IncludeTestPlan || opts.IncludeDeploymentNotes {
		return opts, fmt.Errorf("当前会员等级不支持附加质量选项")
	}
	if optimizeLevel == "expert" {
		return opts, fmt.Errorf("免费用户不支持专家级 Prompt 优化")
	}

	return opts, nil
}

func requestedOptimizeLevel(input interface{}) (string, error) {
	switch req := input.(type) {
	case GenerateOptimizeRequest:
		return normalizeOptimizeLevelValue(req.OptimizeLevel)
	case *GenerateOptimizeRequest:
		if req == nil {
			return "", fmt.Errorf("optimize request is nil")
		}
		return normalizeOptimizeLevelValue(req.OptimizeLevel)
	case recipe.OptimizeInput:
		return normalizeOptimizeLevelValue(req.OptimizeLevel)
	case *recipe.OptimizeInput:
		if req == nil {
			return "", fmt.Errorf("optimize input is nil")
		}
		return normalizeOptimizeLevelValue(req.OptimizeLevel)
	default:
		return "", nil
	}
}

func normalizeOptimizeLevelValue(level string) (string, error) {
	level = strings.ToLower(strings.TrimSpace(level))
	if !recipe.IsValidOptimizeLevel(level) {
		return "", fmt.Errorf("不支持的 optimizeLevel: %s", level)
	}
	return level, nil
}

func maxTokensForGenerate(input interface{}, qualityOpts recipe.QualityOptions) int {
	if qualityOpts.QualityMode == "expert" {
		return expertGenerateMaxTokens
	}

	optimizeLevel, err := requestedOptimizeLevel(input)
	if err == nil && optimizeLevel == "expert" {
		return expertGenerateMaxTokens
	}

	return defaultGenerateMaxTokens
}

func normalizeJSONOutput(content string) (string, error) {
	content = unwrapCodeFence(strings.TrimSpace(content))

	var payload interface{}
	decoder := json.NewDecoder(strings.NewReader(content))
	decoder.UseNumber()

	if err := decoder.Decode(&payload); err != nil {
		return "", err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return "", fmt.Errorf("JSON contains trailing content")
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func unwrapCodeFence(content string) string {
	if !strings.HasPrefix(content, "```") {
		return content
	}

	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return content
	}
	if strings.TrimSpace(lines[len(lines)-1]) != "```" {
		return content
	}

	return strings.Join(lines[1:len(lines)-1], "\n")
}

func (s *GeneratorService) GetHistory(userID uint64, genType string, page, pageSize int) (*response.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	prompts, total, err := s.genRepo.FindByUserID(userID, genType, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &response.PaginatedData{
		List:     prompts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *GeneratorService) GetUserStats(userID uint64) (map[string]interface{}, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	dailyLimit := s.membershipSvc.GetDailyLimit(user.MembershipLevel)
	dailyCount := s.rateLimitSvc.GetDailyCount(userID)

	return map[string]interface{}{
		"membershipLevel":     user.MembershipLevel,
		"dailyLimit":          dailyLimit,
		"dailyUsed":           dailyCount,
		"dailyRemaining":      max(0, dailyLimit-dailyCount),
		"membershipExpiredAt": user.MembershipExpiredAt,
	}, nil
}
