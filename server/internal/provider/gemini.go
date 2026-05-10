package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GeminiProvider struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewGeminiProvider(apiKey, baseURL string) *GeminiProvider {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	return &GeminiProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

func (p *GeminiProvider) Chat(req ChatRequest) (*ChatResponse, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 60
	}
	p.client.Timeout = time.Duration(timeout) * time.Second

	// 构建 Gemini 格式的 contents
	var contents []map[string]interface{}
	var systemInstruction string

	for _, m := range req.Messages {
		if m.Role == "system" {
			systemInstruction = m.Content
			continue
		}
		role := "user"
		if m.Role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]interface{}{
			"role":  role,
			"parts": []map[string]string{{"text": m.Content}},
		})
	}

	body := map[string]interface{}{
		"contents": contents,
		"generationConfig": map[string]interface{}{
			"temperature": req.Temperature,
		},
	}
	if req.MaxTokens > 0 {
		body["generationConfig"].(map[string]interface{})["maxOutputTokens"] = req.MaxTokens
	}
	if systemInstruction != "" {
		body["systemInstruction"] = map[string]interface{}{
			"parts": []map[string]string{{"text": systemInstruction}},
		}
	}

	jsonBody, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, req.Model, p.apiKey)

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Gemini 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API 错误 (status=%d): %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
		UsageMetadata struct {
			PromptTokenCount     int `json:"promptTokenCount"`
			CandidatesTokenCount int `json:"candidatesTokenCount"`
			TotalTokenCount      int `json:"totalTokenCount"`
		} `json:"usageMetadata"`
		ModelVersion string `json:"modelVersion"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析 Gemini 响应失败: %w", err)
	}

	if len(result.Candidates) == 0 {
		return nil, fmt.Errorf("Gemini 返回空结果")
	}

	content := ""
	for _, part := range result.Candidates[0].Content.Parts {
		content += part.Text
	}

	return &ChatResponse{
		Content:          content,
		PromptTokens:     result.UsageMetadata.PromptTokenCount,
		CompletionTokens: result.UsageMetadata.CandidatesTokenCount,
		TotalTokens:      result.UsageMetadata.TotalTokenCount,
		Model:            result.ModelVersion,
		Provider:         "gemini",
	}, nil
}
