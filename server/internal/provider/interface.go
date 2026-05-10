package provider

// ChatMessage 消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest AI 请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	MaxTokens   int           `json:"max_tokens"`
	Timeout     int           `json:"timeout"`
}

// ChatResponse AI 响应
type ChatResponse struct {
	Content          string `json:"content"`
	PromptTokens     int    `json:"prompt_tokens"`
	CompletionTokens int    `json:"completion_tokens"`
	TotalTokens      int    `json:"total_tokens"`
	Model            string `json:"model"`
	Provider         string `json:"provider"`
}

// AIProvider 统一 AI Provider 接口
type AIProvider interface {
	// Name 返回 provider 名称
	Name() string
	// Chat 发送对话请求
	Chat(req ChatRequest) (*ChatResponse, error)
}
