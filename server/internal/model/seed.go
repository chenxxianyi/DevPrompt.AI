package model

import (
	"log"

	"gorm.io/gorm"
)

// Seed initializes required lookup data.
func Seed(db *gorm.DB) error {
	if err := seedMembershipPlans(db); err != nil {
		return err
	}
	if err := seedPromptCategories(db); err != nil {
		return err
	}
	if err := seedAIModels(db); err != nil {
		return err
	}
	if err := seedAdminUser(db); err != nil {
		return err
	}
	if err := seedProjectTypes(db); err != nil {
		return err
	}
	log.Println("seed data is ready")
	return nil
}

func seedMembershipPlans(db *gorm.DB) error {
	plans := []MembershipPlan{
		{
			Name:         "Free",
			Code:         "free",
			Price:        0,
			DurationDays: 0,
			DailyLimit:   5,
			Features:     `["基础 Prompt 模板","每日 5 次生成","社区支持"]`,
			Status:       "active",
		},
		{
			Name:         "Pro",
			Code:         "pro",
			Price:        29,
			DurationDays: 30,
			DailyLimit:   100,
			Features:     `["所有 Prompt 模板","每日 100 次生成","支持 5 种 AI 模型","优先支持"]`,
			Status:       "active",
		},
		{
			Name:         "Team",
			Code:         "team",
			Price:        99,
			DurationDays: 30,
			DailyLimit:   500,
			Features:     `["所有 Pro 功能","每日 500 次生成","团队协作","API 访问"]`,
			Status:       "active",
		},
		{
			Name:         "Enterprise",
			Code:         "enterprise",
			Price:        299,
			DurationDays: 365,
			DailyLimit:   999999,
			Features:     `["所有 Team 功能","无限生成","专属支持","私有化部署","SLA 保障"]`,
			Status:       "active",
		},
	}
	for _, p := range plans {
		if err := db.Where("code = ?", p.Code).FirstOrCreate(&p).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedPromptCategories(db *gorm.DB) error {
	categories := []PromptCategory{
		{Name: "项目开发", Slug: "project-dev", Description: "项目开发相关的 Prompt 模板", Sort: 1, Status: "active"},
		{Name: "Cursor 规则", Slug: "cursor-rules", Description: "Cursor 编辑器 .cursorrules 配置模板", Sort: 2, Status: "active"},
		{Name: "Claude Code", Slug: "claude-code", Description: "Claude Code CLI 任务 Prompt 模板", Sort: 3, Status: "active"},
		{Name: "Prompt 优化", Slug: "prompt-optimize", Description: "Prompt 优化与改进模板", Sort: 4, Status: "active"},
		{Name: "代码生成", Slug: "code-generation", Description: "代码生成相关 Prompt 模板", Sort: 5, Status: "active"},
		{Name: "代码审查", Slug: "code-review", Description: "代码审查相关 Prompt 模板", Sort: 6, Status: "active"},
		{Name: "文档写作", Slug: "documentation", Description: "技术文档写作 Prompt 模板", Sort: 7, Status: "active"},
		{Name: "测试", Slug: "testing", Description: "测试相关 Prompt 模板", Sort: 8, Status: "active"},
	}
	for _, c := range categories {
		if err := db.Where("slug = ?", c.Slug).FirstOrCreate(&c).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAIModels(db *gorm.DB) error {
	var count int64
	if err := db.Model(&AIModel{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	models := []AIModel{
		{Provider: "openai", ModelName: "gpt-4o", DisplayName: "GPT-4o", IsDefault: true, Status: "active", Priority: 1, TimeoutSeconds: 60},
		{Provider: "openai", ModelName: "gpt-4o-mini", DisplayName: "GPT-4o Mini", IsDefault: false, Status: "active", Priority: 2, TimeoutSeconds: 30},
		{Provider: "claude", ModelName: "claude-sonnet-4-6", DisplayName: "Claude Sonnet 4.6", IsDefault: false, Status: "active", Priority: 3, TimeoutSeconds: 60},
		{Provider: "claude", ModelName: "claude-haiku-4-5", DisplayName: "Claude Haiku 4.5", IsDefault: false, Status: "active", Priority: 4, TimeoutSeconds: 30},
		{Provider: "deepseek", ModelName: "deepseek-chat", DisplayName: "DeepSeek Chat", IsDefault: false, Status: "active", Priority: 5, TimeoutSeconds: 60},
		{Provider: "gemini", ModelName: "gemini-2.0-flash", DisplayName: "Gemini 2.0 Flash", IsDefault: false, Status: "active", Priority: 6, TimeoutSeconds: 60},
		{Provider: "qwen", ModelName: "qwen-plus", DisplayName: "Qwen Plus", IsDefault: false, Status: "active", Priority: 7, TimeoutSeconds: 60},
	}
	for _, m := range models {
		if err := db.Where("provider = ? AND model_name = ?", m.Provider, m.ModelName).FirstOrCreate(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAdminUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	admin := User{
		Username:        "admin",
		Email:           "admin@devprompt.ai",
		PasswordHash:    "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // admin123
		Role:            "admin",
		MembershipLevel: "enterprise",
		Status:          "active",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	log.Println("default admin user created (admin@devprompt.ai / admin123)")
	return nil
}

func seedProjectTypes(db *gorm.DB) error {
	var count int64
	if err := db.Model(&ProjectType{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	types := []ProjectType{
		{Name: "SaaS", Value: "SaaS", Description: "SaaS 软件即服务项目", Sort: 1, Status: "active"},
		{Name: "电商平台", Value: "电商", Description: "电商平台类项目", Sort: 2, Status: "active"},
		{Name: "社交应用", Value: "社交", Description: "社交应用类项目", Sort: 3, Status: "active"},
		{Name: "工具类", Value: "工具", Description: "工具类项目", Sort: 4, Status: "active"},
		{Name: "企业内部系统", Value: "企业内部", Description: "企业内部系统类项目", Sort: 5, Status: "active"},
	}
	for _, t := range types {
		if err := db.Where("value = ?", t.Value).FirstOrCreate(&t).Error; err != nil {
			return err
		}
	}
	log.Println("default project types seeded")
	return nil
}
