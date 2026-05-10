package provider

import (
	"fmt"
	"sync"
	"time"

	"devprompt-ai/internal/model"
	"devprompt-ai/internal/repository"
)

// ProviderManager AI Provider 管理器，负责模型选择、失败切换、重试
type ProviderManager struct {
	providers   map[string]AIProvider
	modelRepo   *repository.AIModelRepository
	unhealthyAt map[string]time.Time // 记录不健康的 provider
	mu          sync.RWMutex
	maxRetries  int
}

func NewProviderManager(modelRepo *repository.AIModelRepository) *ProviderManager {
	return &ProviderManager{
		providers:   make(map[string]AIProvider),
		modelRepo:   modelRepo,
		unhealthyAt: make(map[string]time.Time),
		maxRetries:  2,
	}
}

// Register 注册 provider
func (m *ProviderManager) Register(name string, p AIProvider) {
	m.providers[name] = p
}

// GetProvider 获取指定名称的 provider
func (m *ProviderManager) GetProvider(name string) (AIProvider, bool) {
	p, ok := m.providers[name]
	return p, ok
}

// markUnhealthy 标记 provider 为不健康（30 秒后自动恢复）
func (m *ProviderManager) markUnhealthy(providerName string) {
	m.mu.Lock()
	m.unhealthyAt[providerName] = time.Now()
	m.mu.Unlock()
}

// isHealthy 检查 provider 是否健康
func (m *ProviderManager) isHealthy(providerName string) bool {
	m.mu.RLock()
	markedAt, exists := m.unhealthyAt[providerName]
	m.mu.RUnlock()
	if !exists {
		return true
	}
	// 30 秒后自动恢复
	return time.Since(markedAt) > 30*time.Second
}

// Call 调用 AI，支持失败自动切换和重试
func (m *ProviderManager) Call(req ChatRequest, preferredModel string) (*ChatResponse, error) {
	// 获取所有活跃模型（按优先级排序）
	models, err := m.modelRepo.FindActive()
	if err != nil {
		return nil, fmt.Errorf("获取模型列表失败: %w", err)
	}
	if len(models) == 0 {
		return nil, fmt.Errorf("没有可用的 AI 模型")
	}

	// 排序：优先使用指定模型，然后按优先级
	var candidates []model.AIModel
	for _, mdl := range models {
		if mdl.ModelName == preferredModel && mdl.Status == "active" {
			candidates = append([]model.AIModel{mdl}, candidates...)
		} else if mdl.Status == "active" {
			candidates = append(candidates, mdl)
		}
	}

	// 遍历候选模型
	for _, mdl := range candidates {
		provider, ok := m.GetProvider(mdl.Provider)
		if !ok {
			continue
		}

		// 健康检查
		if !m.isHealthy(mdl.Provider) {
			continue
		}

		req.Model = mdl.ModelName
		if req.Timeout == 0 {
			req.Timeout = mdl.TimeoutSeconds
		}

		// 重试
		for attempt := 0; attempt <= m.maxRetries; attempt++ {
			resp, err := provider.Chat(req)
			if err == nil {
				return resp, nil
			}
			if attempt < m.maxRetries {
				time.Sleep(time.Duration(attempt+1) * 500 * time.Millisecond)
			}
		}

		// 当前模型失败，标记为不健康并尝试下一个
		m.markUnhealthy(mdl.Provider)
	}

	return nil, fmt.Errorf("所有 AI Provider 调用均失败")
}
