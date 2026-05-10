package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DeepSeekProvider DeepSeek 兼容 OpenAI 接口格式
type DeepSeekProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewDeepSeekProvider(apiKey, baseURL string) *DeepSeekProvider {
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}
	return &DeepSeekProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *DeepSeekProvider) Name() string {
	return "deepseek"
}

func (p *DeepSeekProvider) Chat(req ChatRequest) (*ChatResponse, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 60
	}
	p.client.Timeout = time.Duration(timeout) * time.Second

	body := map[string]interface{}{
		"model":       req.Model,
		"messages":    req.Messages,
		"temperature": req.Temperature,
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}

	jsonBody, _ := json.Marshal(body)
	httpReq, err := http.NewRequest("POST", p.baseURL+"/chat/completions", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("DeepSeek 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepSeek API 错误 (status=%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Choices []struct {
			Message ChatMessage `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析 DeepSeek 响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("DeepSeek 返回空结果")
	}

	return &ChatResponse{
		Content:          result.Choices[0].Message.Content,
		PromptTokens:     result.Usage.PromptTokens,
		CompletionTokens: result.Usage.CompletionTokens,
		TotalTokens:      result.Usage.TotalTokens,
		Model:            result.Model,
		Provider:         "deepseek",
	}, nil
}
