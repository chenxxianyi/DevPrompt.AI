package service

import (
	"encoding/json"
	"fmt"
	"time"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/provider"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/response"
)

// GeneratorService Prompt 生成服务
type GeneratorService struct {
	providerMgr     *provider.ProviderManager
	genRepo         *repository.GeneratedPromptRepository
	callLogRepo     *repository.AICallLogRepository
	userRepo        *repository.UserRepository
	planRepo        *repository.MembershipPlanRepository
	rateLimitSvc    *RateLimitService
	membershipSvc   *MembershipService
	defaultModel    string
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
		defaultModel:  defaultModel,
	}
}

// GenerateProjectRequest 项目 Prompt 生成参数
type GenerateProjectRequest struct {
	ProjectName  string   `json:"projectName"`
	ProjectType  string   `json:"projectType"`
	TechStack    []string `json:"techStack"`
	Features     []string `json:"features"`
	TargetAiTool string   `json:"targetAiTool"`
}

// GenerateCursorRulesRequest Cursor Rules 生成参数
type GenerateCursorRulesRequest struct {
	Language  string   `json:"language"`
	Framework string   `json:"framework"`
	CodeStyle string   `json:"codeStyle"`
	Rules     []string `json:"rules"`
}

// GenerateClaudeCodeRequest Claude Code 生成参数
type GenerateClaudeCodeRequest struct {
	Task         string   `json:"task"`
	Context      string   `json:"context"`
	Requirements []string `json:"requirements"`
}

// GenerateOptimizeRequest Prompt 优化参数
type GenerateOptimizeRequest struct {
	RawPrompt     string `json:"rawPrompt"`
	TargetTool    string `json:"targetTool"`
	OptimizeLevel string `json:"optimizeLevel"`
}

// Generate 执行 AI 生成
func (s *GeneratorService) Generate(userID uint64, genType, title string, input interface{}, systemPrompt string) (*model.GeneratedPrompt, error) {
	startTime := time.Now()

	// 1. 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return nil, fmt.Errorf("用户不存在")
	}

	// 2. 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("账号已被禁用")
	}

	// 3. 检查会员是否有效
	if !CheckMembership(user) {
		return nil, fmt.Errorf("会员已过期，请续费后使用")
	}

	// 4. 获取每日限制
	dailyLimit := s.membershipSvc.GetDailyLimit(user.MembershipLevel)

	// 5. 检查每日生成次数
	count, limit, allowed, err := s.rateLimitSvc.CheckDailyLimit(userID, dailyLimit)
	if err != nil {
		return nil, fmt.Errorf("限流检查失败: %w", err)
	}
	if !allowed {
		return nil, fmt.Errorf("今日生成次数已达上限 (%d/%d)，请明天再试", count, limit)
	}

	// 6. 接口级限流
	ok, err := s.rateLimitSvc.CheckRateLimit(userID, genType, 30)
	if err != nil || !ok {
		return nil, fmt.Errorf("请求过于频繁，请稍后重试")
	}

	// 7. 构建消息
	inputJSON, _ := json.Marshal(input)
	messages := []provider.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: string(inputJSON)},
	}

	// 8. 调用 AI
	resp, err := s.providerMgr.Call(provider.ChatRequest{
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
	}, s.defaultModel)

	latency := int(time.Since(startTime).Milliseconds())

	// 9. 记录调用日志
	callLog := &model.AICallLog{
		UserID:      userID,
		Provider:    "",
		Model:       "",
		RequestType: genType,
		Status:      "success",
		LatencyMs:   latency,
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
	_ = s.callLogRepo.Create(callLog)

	// 10. 增加每日生成次数
	_ = s.rateLimitSvc.IncrementDailyCount(userID)

	// 11. 保存生成记录
	gp := &model.GeneratedPrompt{
		UserID:   userID,
		Type:     genType,
		Title:    title,
		Input:    string(inputJSON),
		Output:   resp.Content,
		Model:    resp.Model,
		Provider: resp.Provider,
		Tokens:   resp.TotalTokens,
	}

	if err := s.genRepo.Create(gp); err != nil {
		return nil, fmt.Errorf("保存生成记录失败: %w", err)
	}

	return gp, nil
}

// GetHistory 获取用户生成历史
func (s *GeneratorService) GetHistory(userID uint64, page, pageSize int) (*response.PaginatedData, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	prompts, total, err := s.genRepo.FindByUserID(userID, page, pageSize)
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

// GetUserStats 获取用户生成统计
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
