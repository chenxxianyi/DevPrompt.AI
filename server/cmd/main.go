package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"devprompt-ai/internal/api"
	"devprompt-ai/internal/api/admin"
	"devprompt-ai/internal/config"
	"devprompt-ai/internal/middleware"
	"devprompt-ai/internal/model"
	"devprompt-ai/internal/provider"
	"devprompt-ai/internal/repository"
	"devprompt-ai/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 2. 初始化日志
	zapLogger, _ := zap.NewDevelopment()
	defer zapLogger.Sync()

	// 3. 连接 MySQL（先确保数据库存在）
	dsnNoDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port,
		cfg.MySQL.Charset)
	tmpDB, err := gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接 MySQL 失败: %v", err)
	}
	tmpDB.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.MySQL.Database + " DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.MySQL.User, cfg.MySQL.Password, cfg.MySQL.Host, cfg.MySQL.Port,
		cfg.MySQL.Database, cfg.MySQL.Charset)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库 %s 失败: %v", cfg.MySQL.Database, err)
	}
	log.Printf("MySQL 连接成功 (数据库: %s)", cfg.MySQL.Database)

	// Auto-migrate：自动创建/更新表结构 + 种子数据
	if err := db.AutoMigrate(
		&model.User{},
		&model.PromptCategory{},
		&model.PromptTemplate{},
		&model.PromptFavorite{},
		&model.PromptLike{},
		&model.GeneratedPrompt{},
		&model.AIModel{},
		&model.AICallLog{},
		&model.MembershipPlan{},
		&model.Order{},
		&model.ProjectType{},
	); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	log.Println("数据库表结构已就绪")
	if err := model.Seed(db); err != nil {
		log.Fatalf("初始化种子数据失败: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 4. 连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("警告: 连接 Redis 失败: %v (限流功能将不可用)", err)
		cfg.RateLimit.Enabled = false
	} else {
		log.Println("Redis 连接成功")
	}

	// 5. 初始化依赖
	// Repositories
	userRepo := repository.NewUserRepository(db)
	promptTemplateRepo := repository.NewPromptTemplateRepository(db)
	promptCategoryRepo := repository.NewPromptCategoryRepository(db)
	promptFavoriteRepo := repository.NewPromptFavoriteRepository(db)
	promptLikeRepo := repository.NewPromptLikeRepository(db)
	genRepo := repository.NewGeneratedPromptRepository(db)
	aiModelRepo := repository.NewAIModelRepository(db)
	callLogRepo := repository.NewAICallLogRepository(db)
	planRepo := repository.NewMembershipPlanRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	projectTypeRepo := repository.NewProjectTypeRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	promptService := service.NewPromptService(promptTemplateRepo, promptCategoryRepo, promptFavoriteRepo, promptLikeRepo)
	rateLimitSvc := service.NewRateLimitService(rdb, cfg.RateLimit.Enabled)
	membershipSvc := service.NewMembershipService(planRepo)

	// AI Provider Manager
	providerMgr := provider.NewProviderManager(aiModelRepo)

	// 注册所有 provider
	if cfg.AI.OpenAI.APIKey != "" {
		providerMgr.Register("openai", provider.NewOpenAIProvider(cfg.AI.OpenAI.APIKey, cfg.AI.OpenAI.BaseURL))
	}
	if cfg.AI.Claude.APIKey != "" {
		providerMgr.Register("claude", provider.NewClaudeProvider(cfg.AI.Claude.APIKey, cfg.AI.Claude.BaseURL))
	}
	if cfg.AI.Gemini.APIKey != "" {
		providerMgr.Register("gemini", provider.NewGeminiProvider(cfg.AI.Gemini.APIKey, cfg.AI.Gemini.BaseURL))
	}
	if cfg.AI.DeepSeek.APIKey != "" {
		providerMgr.Register("deepseek", provider.NewDeepSeekProvider(cfg.AI.DeepSeek.APIKey, cfg.AI.DeepSeek.BaseURL))
	}
	if cfg.AI.Qwen.APIKey != "" {
		providerMgr.Register("qwen", provider.NewQwenProvider(cfg.AI.Qwen.APIKey, cfg.AI.Qwen.BaseURL))
	}

	defaultModel := cfg.AI.DeepSeek.DefaultModel
	genService := service.NewGeneratorService(
		providerMgr, genRepo, callLogRepo, userRepo, planRepo,
		rateLimitSvc, membershipSvc, defaultModel,
	)

	// Handlers
	projectTypeHandler := api.NewProjectTypeHandler(projectTypeRepo)
	authHandler := api.NewAuthHandler(authService)
	promptHandler := api.NewPromptHandler(promptService)
	generatorHandler := api.NewGeneratorHandler(genService)
	membershipHandler := api.NewMembershipHandler(membershipSvc)
	userHandler := api.NewUserHandler(genService)

	// Admin handlers
	userAdminHandler := admin.NewUserAdminHandler(userRepo)
	promptAdminHandler := admin.NewPromptAdminHandler(promptService)
	aiModelAdminHandler := admin.NewAIModelAdminHandler(aiModelRepo)
	categoryAdminHandler := admin.NewCategoryAdminHandler(promptCategoryRepo)
	callLogAdminHandler := admin.NewCallLogAdminHandler(callLogRepo)
	planAdminHandler := admin.NewPlanAdminHandler(planRepo)
	orderAdminHandler := admin.NewOrderAdminHandler(orderRepo)
	projectTypeAdminHandler := admin.NewProjectTypeAdminHandler(projectTypeRepo)
	dashboardHandler := admin.NewDashboardHandler(db)

	// 6. 创建 Gin 路由
	r := gin.Default()

	// CORS
	r.Use(middleware.Cors(cfg.App.WebBaseURL, cfg.App.AdminBaseURL))

	// ============ 公开路由 ============
	apiGroup := r.Group("/api")
	{
		// Auth
		apiGroup.POST("/auth/register", authHandler.Register)
		apiGroup.POST("/auth/login", authHandler.Login)

		// Prompts (公开)
		apiGroup.GET("/prompts", middleware.OptionalAuth(cfg.JWT.Secret), promptHandler.List)
		apiGroup.GET("/prompts/:slug", middleware.OptionalAuth(cfg.JWT.Secret), promptHandler.GetBySlug)

		// Categories (公开)
		apiGroup.GET("/categories", promptHandler.ListCategories)

		// Membership (公开)
		apiGroup.GET("/membership/plans", membershipHandler.ListPlans)
		// Project Types (公开)
		apiGroup.GET("/project-types", projectTypeHandler.ListPublic)
	}

	// ============ 需要认证的路由 ============
	authGroup := apiGroup.Group("", middleware.JWTAuth(cfg.JWT.Secret))
	{
		// Auth
		authGroup.GET("/auth/profile", authHandler.Profile)

		// Prompts 互动
		authGroup.POST("/prompts/:id/like", promptHandler.ToggleLike)
		authGroup.POST("/prompts/:id/favorite", promptHandler.ToggleFavorite)

		// Generator
		authGroup.POST("/generator/project", generatorHandler.Project)
		authGroup.POST("/generator/cursor-rules", generatorHandler.CursorRules)
		authGroup.POST("/generator/claude-code", generatorHandler.ClaudeCode)
		authGroup.POST("/generator/optimize", generatorHandler.Optimize)
		authGroup.GET("/generator/history", generatorHandler.History)

		// User
		authGroup.GET("/user/generate-stats", userHandler.GenerateStats)
	}

	// ============ 管理后台路由 ============
	adminGroup := apiGroup.Group("/admin", middleware.JWTAuth(cfg.JWT.Secret), middleware.AdminAuth())
	{
		// Dashboard
		adminGroup.GET("/dashboard", dashboardHandler.Overview)

		// Users
		adminGroup.GET("/users", userAdminHandler.List)

		// Prompts CRUD
		adminGroup.GET("/prompts", promptAdminHandler.List)
		adminGroup.POST("/prompts", promptAdminHandler.Create)
		adminGroup.PUT("/prompts/:id", promptAdminHandler.Update)
		adminGroup.DELETE("/prompts/:id", promptAdminHandler.Delete)

		// AI Models CRUD
		adminGroup.GET("/ai-models", aiModelAdminHandler.List)
		adminGroup.POST("/ai-models", aiModelAdminHandler.Create)
		adminGroup.PUT("/ai-models/:id", aiModelAdminHandler.Update)
		adminGroup.DELETE("/ai-models/:id", aiModelAdminHandler.Delete)

		// Categories CRUD
		adminGroup.GET("/categories", categoryAdminHandler.List)
		adminGroup.POST("/categories", categoryAdminHandler.Create)
		adminGroup.PUT("/categories/:id", categoryAdminHandler.Update)
		adminGroup.DELETE("/categories/:id", categoryAdminHandler.Delete)

		// AI Call Logs
		adminGroup.GET("/ai-call-logs", callLogAdminHandler.List)

		// Membership Plans CRUD
		adminGroup.GET("/membership-plans", planAdminHandler.List)
		adminGroup.POST("/membership-plans", planAdminHandler.Create)
		adminGroup.PUT("/membership-plans/:id", planAdminHandler.Update)
		adminGroup.DELETE("/membership-plans/:id", planAdminHandler.Delete)

		// Orders
		adminGroup.GET("/orders", orderAdminHandler.List)
		// Project Types CRUD
		adminGroup.GET("/project-types", projectTypeAdminHandler.List)
		adminGroup.POST("/project-types", projectTypeAdminHandler.Create)
			adminGroup.PUT("/project-types/:id", projectTypeAdminHandler.Update)
			adminGroup.DELETE("/project-types/:id", projectTypeAdminHandler.Delete)
	}

	// 7. 启动服务
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	log.Printf("DevPrompt AI 后端服务启动于 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
