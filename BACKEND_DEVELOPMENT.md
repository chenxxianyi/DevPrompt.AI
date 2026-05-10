# DevPrompt AI 后端开发方案

## 1. 后端目标

DevPrompt AI 后端负责提供用户认证、Prompt 模板、Prompt 生成、AI Provider 调用、会员限制、调用日志和管理后台接口。

后端需要保证：

- 分层清晰。
- 接口响应统一。
- 错误处理统一。
- API Key 不暴露给前端。
- 支持 MySQL 数据持久化。
- 支持 Redis 限流与生成次数控制。
- 支持多 AI Provider 和失败自动切换。
- 支持 Docker Compose 本地一键启动。

## 2. 后端技术栈

- Golang
- Gin
- Gorm
- JWT
- MySQL 8
- Redis
- Viper 配置管理
- Zap / Logrus 日志
- Docker
- Docker Compose
- Nginx

## 3. 后端目录结构

```text
DevPrompt AI/
├── server/
│   ├── cmd/
│   │   └── main.go
│   ├── config/
│   │   └── config.yaml
│   ├── internal/
│   │   ├── api/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── provider/
│   │   ├── response/
│   │   └── utils/
│   ├── go.mod
│   └── Dockerfile
│
├── database/
│   └── init.sql
│
├── deploy/
│   └── nginx/
│       └── nginx.conf
│
├── docker-compose.yml
├── .env.example
└── README.md
```

## 4. 后端分层设计

```text
api          HTTP Handler，处理请求参数、认证上下文和响应
service      业务逻辑层，处理核心规则、权限、调用编排
repository   数据访问层，封装 Gorm 查询
model        Gorm 数据模型
provider     AI Provider 抽象与具体实现
middleware   JWT、CORS、限流、管理员权限
response     统一响应结构
utils        JWT、密码加密、配置、日志等工具
```

## 5. 核心模块设计

| 模块 | 说明 |
| --- | --- |
| Auth | 注册、登录、获取用户信息 |
| User | 用户管理、角色、会员信息 |
| PromptCategory | Prompt 分类管理 |
| PromptTemplate | Prompt 模板管理 |
| Favorite | Prompt 收藏 |
| Like | Prompt 点赞 |
| Generator | 四类 Prompt 生成接口 |
| AIProvider | OpenAI、Claude、Gemini、DeepSeek、Qwen |
| Membership | 会员套餐与权限限制 |
| AICallLog | AI 调用日志 |
| Admin | 后台管理接口 |
| RateLimit | Redis 限流 |

## 6. 数据库设计

核心表：

```text
users
prompt_categories
prompt_templates
generated_prompts
prompt_favorites
ai_models
ai_call_logs
membership_plans
orders
```

### 6.1 users

```text
id
username
email
password_hash
avatar
role
membership_level
membership_expired_at
daily_generate_count
last_generate_date
status
created_at
updated_at
deleted_at
```

### 6.2 prompt_categories

```text
id
name
slug
description
sort
status
created_at
updated_at
deleted_at
```

### 6.3 prompt_templates

```text
id
category_id
title
slug
description
content
tags
use_count
like_count
favorite_count
status
created_at
updated_at
deleted_at
```

### 6.4 generated_prompts

```text
id
user_id
type
title
input
output
model
provider
tokens
created_at
updated_at
deleted_at
```

### 6.5 prompt_favorites

```text
id
user_id
prompt_template_id
created_at
```

### 6.6 ai_models

```text
id
provider
model_name
display_name
api_base_url
is_default
status
priority
timeout_seconds
created_at
updated_at
deleted_at
```

### 6.7 ai_call_logs

```text
id
user_id
provider
model
request_type
prompt_tokens
completion_tokens
total_tokens
status
error_message
latency_ms
created_at
```

### 6.8 membership_plans

```text
id
name
code
price
duration_days
daily_limit
features
status
created_at
updated_at
deleted_at
```

### 6.9 orders

```text
id
user_id
plan_id
order_no
amount
status
paid_at
created_at
updated_at
deleted_at
```

## 7. 统一响应与错误处理

所有接口返回格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

错误响应：

```json
{
  "code": 400,
  "message": "invalid request",
  "data": null
}
```

后端需要统一处理：

- 参数绑定错误。
- 参数校验错误。
- 认证失败。
- 权限不足。
- 资源不存在。
- 数据库错误。
- AI Provider 调用失败。
- Redis 限流失败。
- 会员额度不足。

## 8. 认证与权限设计

### 8.1 认证接口

```http
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/profile
```

注册参数：

```json
{
  "username": "dev",
  "email": "dev@example.com",
  "password": "123456"
}
```

登录参数：

```json
{
  "email": "dev@example.com",
  "password": "123456"
}
```

### 8.2 用户角色

```text
user   普通用户
admin  管理员
```

### 8.3 会员等级

```text
free
pro
team
enterprise
```

### 8.4 每日生成限制示例

```text
free        5 次 / 天
pro         100 次 / 天
team        500 次 / 天
enterprise  不限
```

生成接口调用 AI 前必须检查：

- 用户是否登录。
- 用户是否被禁用。
- 会员是否过期。
- 当日生成次数是否超限。
- Redis 限流是否通过。

## 9. Prompt API 设计

```http
GET  /api/prompts
GET  /api/prompts/:slug
POST /api/prompts/:id/like
POST /api/prompts/:id/favorite
```

列表查询参数：

```text
keyword
category
page
page_size
sort
```

## 10. 生成 API 设计

```http
POST /api/generator/project
POST /api/generator/cursor-rules
POST /api/generator/claude-code
POST /api/generator/optimize
GET  /api/generator/history
```

### 10.1 项目 Prompt 生成参数

```json
{
  "project_name": "DevPrompt AI",
  "project_type": "SaaS",
  "tech_stack": ["Next.js", "Golang", "MySQL"],
  "features": ["登录注册", "Prompt 生成", "会员系统"],
  "target_ai_tool": "Cursor"
}
```

### 10.2 Cursor Rules 生成参数

```json
{
  "language": "TypeScript",
  "framework": "Next.js",
  "code_style": "clean architecture",
  "rules": ["使用组件化", "禁止硬编码 API Key"]
}
```

### 10.3 Claude Code Prompt 生成参数

```json
{
  "task": "开发用户登录注册模块",
  "context": "Golang Gin + JWT",
  "requirements": ["分层架构", "统一错误处理", "密码加密"]
}
```

### 10.4 Prompt 优化参数

```json
{
  "raw_prompt": "帮我写一个网站",
  "target_tool": "Cursor",
  "optimize_level": "professional"
}
```

## 11. 管理后台 API 设计

```http
GET    /api/admin/users

GET    /api/admin/prompts
POST   /api/admin/prompts
PUT    /api/admin/prompts/:id
DELETE /api/admin/prompts/:id

GET    /api/admin/ai-models
POST   /api/admin/ai-models
PUT    /api/admin/ai-models/:id
DELETE /api/admin/ai-models/:id
```

后续扩展：

```http
GET    /api/admin/categories
POST   /api/admin/categories
PUT    /api/admin/categories/:id
DELETE /api/admin/categories/:id

GET    /api/admin/ai-call-logs
GET    /api/admin/membership-plans
GET    /api/admin/orders
GET    /api/admin/settings
```

## 12. AI Provider 设计

统一接口：

```go
type AIProvider interface {
    Chat(req ChatRequest) (*ChatResponse, error)
}
```

请求结构：

```go
type ChatRequest struct {
    Model       string
    Messages    []ChatMessage
    Temperature float64
    MaxTokens   int
    Timeout     int
}
```

响应结构：

```go
type ChatResponse struct {
    Content          string
    PromptTokens     int
    CompletionTokens int
    TotalTokens      int
    Model            string
    Provider         string
}
```

支持 Provider：

- OpenAI
- Claude
- Gemini
- DeepSeek
- Qwen

Provider 能力要求：

- 支持模型失败自动切换。
- 支持超时控制。
- 支持错误重试。
- 支持 Token 统计。
- 支持调用日志。
- API Key 必须来自后端环境变量。
- 不允许在前端暴露 API Key。

## 13. Redis 设计

Redis 用途：

- 接口限流。
- 用户每日生成次数缓存。
- 登录态辅助缓存。
- 热门 Prompt 缓存。

建议 Key 设计：

```text
rate_limit:{user_id}:{route}
daily_generate:{user_id}:{yyyy-mm-dd}
prompt_hot:list
```

## 14. Docker 与部署设计

容器组成：

```text
devprompt-server
devprompt-web
devprompt-admin
devprompt-mysql
devprompt-redis
devprompt-nginx
```

端口规划：

```text
web:     3000
admin:   5173
server:  8080
mysql:   3306
redis:   6379
nginx:   80 / 443
```

访问路径建议：

```text
https://example.com          用户端
https://example.com/admin    管理后台
https://example.com/api      后端 API
```

## 15. 环境变量

```env
APP_ENV=development
APP_PORT=8080

MYSQL_HOST=mysql
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=devprompt_password
MYSQL_DATABASE=devprompt_ai

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=change_me
JWT_EXPIRE_HOURS=168

OPENAI_API_KEY=
CLAUDE_API_KEY=
GEMINI_API_KEY=
DEEPSEEK_API_KEY=
QWEN_API_KEY=

WEB_BASE_URL=http://localhost:3000
ADMIN_BASE_URL=http://localhost:5173
```

## 16. 后端开发阶段

### 第一阶段：项目初始化与认证

- 初始化 `server` 目录。
- 配置 `go.mod`。
- 配置 Gin 基础服务。
- 配置 Viper 读取配置。
- 配置 Gorm 连接 MySQL。
- 配置 Redis 客户端。
- 配置 Docker Compose。
- 实现用户注册、登录和 profile。
- 实现 JWT 中间件。

第一阶段核心文件：

```text
docker-compose.yml
.env.example

server/go.mod
server/cmd/main.go
server/config/config.yaml
server/internal/model/user.go
server/internal/api/auth.go
server/internal/service/auth_service.go
server/internal/repository/user_repository.go
server/internal/middleware/jwt.go
server/internal/response/response.go
server/internal/utils/jwt.go
server/internal/utils/password.go

database/init.sql
README.md
```

### 第二阶段：Prompt 模板系统

- 实现 Prompt 分类 Model、Repository、Service、API。
- 实现 Prompt 模板 Model、Repository、Service、API。
- 实现 Prompt 列表和详情。
- 实现收藏。
- 实现点赞。

### 第三阶段：AI 生成能力

- 定义 AI Provider 抽象接口。
- 接入 OpenAI。
- 接入 DeepSeek。
- 扩展 Claude、Gemini、Qwen。
- 实现失败自动切换。
- 实现超时控制。
- 实现错误重试。
- 实现 Token 统计。
- 实现四个生成接口。

### 第四阶段：用户工作台与会员

- 实现生成历史。
- 实现会员限制。
- 实现每日生成次数限制。
- 实现 AI 调用日志。
- 实现 Redis 限流。

### 第五阶段：管理后台接口

- 实现用户管理接口。
- 实现 Prompt 分类管理接口。
- 实现 Prompt 模板管理接口。
- 实现 AI 模型管理接口。
- 实现 AI 调用日志接口。
- 实现会员套餐管理接口。
- 实现订单管理接口。
- 实现系统设置接口。

### 第六阶段：后端优化与交付

- 完善统一错误处理。
- 完善参数校验。
- 完善日志。
- 完善 Dockerfile。
- 完善 Docker Compose。
- 完善 Nginx 配置。
- 完成 SQL 初始化文件。
- 完成 README 与 API 文档。

## 17. 后端启动命令

开发模式：

```bash
cd server
go run ./cmd/main.go
```

完整环境：

```bash
docker compose up -d
```

## 18. 后端交付标准

- 后端服务可启动。
- MySQL / Redis 可连接。
- Gorm Model 与迁移逻辑可用。
- 注册、登录、profile 接口可用。
- Prompt 分类和模板接口可用。
- 四个生成接口可用。
- AI 调用日志可记录。
- 会员生成次数限制可生效。
- 管理员接口具备权限控制。
- 所有接口使用统一响应格式。
- API Key 全部来自环境变量。
- Docker Compose 可一键启动。

