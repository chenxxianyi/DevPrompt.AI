package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构体
type Config struct {
	App       AppConfig       `yaml:"app"`
	MySQL     MySQLConfig     `yaml:"mysql"`
	Redis     RedisConfig     `yaml:"redis"`
	JWT       JWTConfig       `yaml:"jwt"`
	AI        AIConfig        `yaml:"ai"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type AppConfig struct {
	Env          string `yaml:"env"`
	Port         int    `yaml:"port"`
	WebBaseURL   string `yaml:"web_base_url"`
	AdminBaseURL string `yaml:"admin_base_url"`
}

type MySQLConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Database     string `yaml:"database"`
	Charset      string `yaml:"charset"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type AIConfig struct {
	OpenAI  AIProviderConfig `yaml:"openai"`
	Claude  AIProviderConfig `yaml:"claude"`
	Gemini  AIProviderConfig `yaml:"gemini"`
	DeepSeek AIProviderConfig `yaml:"deepseek"`
	Qwen    AIProviderConfig `yaml:"qwen"`
}

type AIProviderConfig struct {
	APIKey       string `yaml:"api_key"`
	BaseURL      string `yaml:"base_url"`
	DefaultModel string `yaml:"default_model"`
}

type RateLimitConfig struct {
	Enabled            bool `yaml:"enabled"`
	RequestsPerMinute  int  `yaml:"requests_per_minute"`
}

// Load 从 yaml 文件加载配置，然后用环境变量覆盖
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// 环境变量覆盖
	cfg.overrideFromEnv()

	return cfg, nil
}

// overrideFromEnv 用环境变量覆盖 yaml 配置
func (c *Config) overrideFromEnv() {
	if v := os.Getenv("APP_ENV"); v != "" {
		c.App.Env = v
	}
	if v := os.Getenv("APP_PORT"); v != "" {
		c.App.Port = parseIntEnv(v)
	}
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		c.MySQL.Host = v
	}
	if v := os.Getenv("MYSQL_PORT"); v != "" {
		c.MySQL.Port = parseIntEnv(v)
	}
	if v := os.Getenv("MYSQL_USER"); v != "" {
		c.MySQL.User = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		c.MySQL.Password = v
	}
	if v := os.Getenv("MYSQL_DATABASE"); v != "" {
		c.MySQL.Database = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		c.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		c.Redis.Port = parseIntEnv(v)
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		c.Redis.Password = v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.JWT.Secret = v
	}
	if v := os.Getenv("JWT_EXPIRE_HOURS"); v != "" {
		c.JWT.ExpireHours = parseIntEnv(v)
	}
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		c.AI.OpenAI.APIKey = v
	}
	if v := os.Getenv("CLAUDE_API_KEY"); v != "" {
		c.AI.Claude.APIKey = v
	}
	if v := os.Getenv("GEMINI_API_KEY"); v != "" {
		c.AI.Gemini.APIKey = v
	}
	if v := os.Getenv("DEEPSEEK_API_KEY"); v != "" {
		c.AI.DeepSeek.APIKey = v
	}
	if v := os.Getenv("QWEN_API_KEY"); v != "" {
		c.AI.Qwen.APIKey = v
	}
	if v := os.Getenv("WEB_BASE_URL"); v != "" {
		c.App.WebBaseURL = v
	}
	if v := os.Getenv("ADMIN_BASE_URL"); v != "" {
		c.App.AdminBaseURL = v
	}
}

func parseIntEnv(s string) int {
	var v int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			v = v*10 + int(c-'0')
		}
	}
	return v
}
