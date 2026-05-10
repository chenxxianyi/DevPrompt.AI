package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClaudeProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClaudeProvider(apiKey, baseURL string) *ClaudeProvider {
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}
	return &ClaudeProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *ClaudeProvider) Name() string {
	return "claude"
}

func (p *ClaudeProvider) Chat(req ChatRequest) (*ChatResponse, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 60
	}
	p.client.Timeout = time.Duration(timeout) * time.Second

	// 转换消息格式：system 消息需要提取到顶层
	var systemMsg string
	var messages []ChatMessage
	for _, m := range req.Messages {
		if m.Role == "system" {
			systemMsg = m.Content
		} else {
			messages = append(messages, m)
		}
	}

	body := map[string]interface{}{
		"model":       req.Model,
		"messages":    messages,
		"temperature": req.Temperature,
		"max_tokens":  4096,
	}
	if req.MaxTokens > 0 {
		body["max_tokens"] = req.MaxTokens
	}
	if systemMsg != "" {
		body["system"] = systemMsg
	}

	jsonBody, _ := json.Marshal(body)
	httpReq, err := http.NewRequest("POST", p.baseURL+"/messages", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Claude 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Claude API 错误 (status=%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析 Claude 响应失败: %w", err)
	}

	content := ""
	for _, c := range result.Content {
		if c.Type == "text" {
			content += c.Text
		}
	}

	return &ChatResponse{
		Content:          content,
		PromptTokens:     result.Usage.InputTokens,
		CompletionTokens: result.Usage.OutputTokens,
		TotalTokens:      result.Usage.InputTokens + result.Usage.OutputTokens,
		Model:            result.Model,
		Provider:         "claude",
	}, nil
}
