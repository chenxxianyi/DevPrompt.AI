# DevPrompt AI 开发文档

## 1. 项目概述

DevPrompt AI 是一个面向程序员的 AI 编程 Prompt 生成、优化、管理与复用平台。

平台目标是帮助开发者快速生成适用于以下工具的高质量 Prompt：

- Cursor
- Claude Code
- GPT / OpenAI
- Gemini
- DeepSeek
- 通义千问 / Qwen

核心能力包括：

- Prompt 模板库
- 项目开发 Prompt 生成器
- Cursor Rules 生成器
- Claude Code Prompt 生成器
- Prompt 优化器
- 用户工作台
- Prompt 收藏
- Prompt 生成历史
- 会员系统
- 管理后台
- AI 模型管理
- AI 调用日志

## 2. 技术架构

整体采用前后端分离架构：

```text
DevPrompt AI
├── web       用户端 Next.js
├── admin     管理后台 Vue3
├── server    Golang 后端 API
├── deploy    Nginx / Docker 配置
└── docs      项目文档
```

### 2.1 用户端技术栈

- Next.js 14
- TypeScript
- Tailwind CSS
- Shadcn UI
- Zustand
- React Hook Form
- Axios
- Markdown 渲染

### 2.2 管理后台技术栈

- Vue3
- Vite
- TypeScript
- Element Plus
- Pinia
- Vue Router
- Axios

### 2.3 后端技术栈

- Golang
- Gin
- Gorm
- JWT
- MySQL 8
- Redis
- Zap / Logrus 日志
- Viper 配置管理

### 2.4 部署技术栈

- Docker
- Docker Compose
- Nginx
- MySQL 8
- Redis

## 3. 项目目录设计

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
├── web/
│   ├── app/
│   ├── components/
│   ├── lib/
│   ├── store/
│   ├── types/
│   ├── public/
│   ├── package.json
│   └── Dockerfile
│
├── admin/
│   ├── src/
│   │   ├── api/
│   │   ├── router/
│   │   ├── store/
│   │   ├── views/
│   │   ├── components/
│   │   └── utils/
│   ├── package.json
│   └── Dockerfile
│
├── deploy/
│   └── nginx/
│       └── nginx.conf
│
├── database/
│   └── init.sql
│
├── docs/
│   ├── api.md
│   ├── database.md
│   └── deployment.md
│
├── docker-compose.yml
├── .env.example
└── README.md
```

## 4. 后端模块设计

### 4.1 后端分层

```text
api          HTTP Handler，处理请求和响应
service      业务逻辑层
repository   数据访问层
model        Gorm 数据模型
provider     AI Provider 抽象与实现
middleware   JWT、CORS、限流、权限
response     统一响应结构
utils        JWT、密码加密、配置、日志等工具
```

### 4.2 核心后端模块

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

## 5. 数据库设计

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

### 5.1 users

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

### 5.2 prompt_categories

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

### 5.3 prompt_templates

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

### 5.4 generated_prompts

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

### 5.5 prompt_favorites

```text
id
user_id
prompt_template_id
created_at
```

### 5.6 ai_models

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

### 5.7 ai_call_logs

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

### 5.8 membership_plans

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

### 5.9 orders

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

## 6. API 设计

### 6.1 统一响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

错误示例：

```json
{
  "code": 400,
  "message": "invalid request",
  "data": null
}
```

## 7. 认证接口

```http
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/profile
```

### 7.1 注册

```json
{
  "username": "dev",
  "email": "dev@example.com",
  "password": "123456"
}
```

### 7.2 登录

```json
{
  "email": "dev@example.com",
  "password": "123456"
}
```

## 8. Prompt 接口

```http
GET  /api/prompts
GET  /api/prompts/:slug
POST /api/prompts/:id/like
POST /api/prompts/:id/favorite
```

支持查询参数：

```text
keyword
category
page
page_size
sort
```

## 9. 生成接口

```http
POST /api/generator/project
POST /api/generator/cursor-rules
POST /api/generator/claude-code
POST /api/generator/optimize
GET  /api/generator/history
```

### 9.1 项目 Prompt 生成参数

```json
{
  "project_name": "DevPrompt AI",
  "project_type": "SaaS",
  "tech_stack": ["Next.js", "Golang", "MySQL"],
  "features": ["登录注册", "Prompt 生成", "会员系统"],
  "target_ai_tool": "Cursor"
}
```

### 9.2 Cursor Rules 生成参数

```json
{
  "language": "TypeScript",
  "framework": "Next.js",
  "code_style": "clean architecture",
  "rules": ["使用组件化", "禁止硬编码 API Key"]
}
```

### 9.3 Claude Code Prompt 生成参数

```json
{
  "task": "开发用户登录注册模块",
  "context": "Golang Gin + JWT",
  "requirements": ["分层架构", "统一错误处理", "密码加密"]
}
```

### 9.4 Prompt 优化参数

```json
{
  "raw_prompt": "帮我写一个网站",
  "target_tool": "Cursor",
  "optimize_level": "professional"
}
```

## 10. 管理后台接口

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

后续可扩展：

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

## 11. AI Provider 设计

统一接口：

```go
type AIProvider interface {
    Chat(req ChatRequest) (*ChatResponse, error)
}
```

### 11.1 ChatRequest

```go
type ChatRequest struct {
    Model       string
    Messages    []ChatMessage
    Temperature float64
    MaxTokens   int
    Timeout     int
}
```

### 11.2 ChatResponse

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

### 11.3 支持 Provider

- OpenAI
- Claude
- Gemini
- DeepSeek
- Qwen

### 11.4 Provider 能力要求

- 失败自动切换模型
- 超时控制
- 错误重试
- Token 统计
- 调用日志记录
- API Key 只存后端环境变量
- 前端不暴露任何模型密钥

## 12. 用户端页面设计

### 12.1 页面列表

```text
/                        首页
/prompts                 Prompt 模板列表
/prompts/[slug]          Prompt 模板详情
/generator/project       项目 Prompt 生成器
/generator/cursor-rules  Cursor Rules 生成器
/generator/claude-code   Claude Code Prompt 生成器
/generator/optimize      Prompt 优化器
/dashboard               我的工作台
/favorites               我的收藏
/history                 生成历史
/login                   登录
/register                注册
/pricing                 会员价格
```

### 12.2 用户端核心体验

- 首页直接展示核心生成入口
- Prompt 模板支持分类、搜索、排序
- 模板详情支持一键复制
- 生成结果支持 Markdown 展示
- 登录状态持久化
- 会员限制友好提示
- 移动端适配

## 13. 管理后台页面设计

### 13.1 页面列表

```text
/login             后台登录
/dashboard         控制台
/users             用户管理
/categories        Prompt 分类管理
/prompts           Prompt 模板管理
/ai-models         AI 模型管理
/ai-call-logs      AI 调用日志
/membership-plans  会员套餐管理
/orders            订单管理
/settings          系统设置
```

### 13.2 后台核心体验

- Element Plus 表格分页
- 搜索筛选
- 新增编辑弹窗
- 删除确认
- 权限控制
- API 统一封装
- Token 失效自动跳转登录

## 14. 权限与会员设计

### 14.1 用户角色

```text
user   普通用户
admin  管理员
```

### 14.2 会员等级

```text
free
pro
team
enterprise
```

### 14.3 每日生成限制示例

```text
free        5 次 / 天
pro         100 次 / 天
team        500 次 / 天
enterprise  不限
```

后端生成接口需要在调用 AI 前检查：

- 用户是否登录
- 用户是否被禁用
- 会员是否过期
- 当日生成次数是否超限
- Redis 限流是否通过

## 15. Docker 部署设计

容器组成：

```text
devprompt-server
devprompt-web
devprompt-admin
devprompt-mysql
devprompt-redis
devprompt-nginx
```

### 15.1 端口规划

```text
web:     3000
admin:   5173
server:  8080
mysql:   3306
redis:   6379
nginx:   80 / 443
```

### 15.2 访问路径建议

```text
https://example.com          用户端
https://example.com/admin    管理后台
https://example.com/api      后端 API
```

## 16. 环境变量设计

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

## 17. 开发阶段规划

### 第一阶段：项目初始化与认证

目标：

- 初始化项目结构
- 配置 Docker Compose
- 配置 MySQL、Redis
- 实现后端基础框架
- 实现用户登录注册

交付物：

- Monorepo 目录
- 后端 Gin 项目
- Gorm 连接 MySQL
- Redis 连接
- JWT 登录注册
- Docker Compose 基础环境

### 第二阶段：Prompt 模板系统

目标：

- Prompt 分类管理
- Prompt 模板管理
- Prompt 列表和详情
- 收藏和点赞

交付物：

- 分类 Model / API
- 模板 Model / API
- 收藏 API
- 点赞 API
- 用户端模板列表页
- 用户端模板详情页

### 第三阶段：AI 生成能力

目标：

- AI Provider 抽象
- OpenAI / DeepSeek 接入
- Cursor Rules 生成
- Claude Code Prompt 生成
- Prompt 优化

交付物：

- Provider 接口
- Provider 实现
- 自动重试
- 模型切换
- Token 统计
- 四个生成接口
- 用户端生成器页面

### 第四阶段：用户工作台与会员

目标：

- 用户工作台
- 生成历史
- 会员限制
- AI 调用日志

交付物：

- 工作台页面
- 收藏页面
- 历史页面
- 会员价格页
- 生成次数限制
- AI 调用日志记录

### 第五阶段：管理后台

目标：

- 后台基础框架
- 用户管理
- Prompt 管理
- AI 模型管理
- 系统设置

交付物：

- Vue3 后台
- Element Plus 布局
- 用户管理表格
- Prompt 管理表格
- AI 模型管理
- 调用日志页面
- 会员套餐管理
- 订单管理

### 第六阶段：优化与交付

目标：

- UI 优化
- 错误处理完善
- Docker 部署完善
- README 文档
- API 文档

交付物：

- 完整 README
- API 文档
- SQL 初始化文件
- Docker Compose
- Nginx 配置
- 本地一键启动说明

## 18. 启动命令规划

最终本地启动方式：

```bash
docker compose up -d
```

开发模式：

```bash
cd server
go run ./cmd/main.go
```

```bash
cd web
npm install
npm run dev
```

```bash
cd admin
npm install
npm run dev
```

## 19. 第一阶段开发清单

下一步真正开始开发时，建议从以下文件开始：

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

第一阶段完成后，应具备：

- MySQL / Redis 可启动
- 后端服务可启动
- 用户可注册
- 用户可登录
- 可通过 JWT 获取 profile
- 项目目录结构完整

## 20. 后续扩展建议

后续可以扩展：

- Prompt 版本管理
- 团队空间
- Prompt 分享页
- Prompt 评分系统
- Prompt 变量表单化
- Prompt A/B 测试
- AI 调用成本统计
- 多模型响应对比
- Stripe / Lemon Squeezy 支付
- Webhook 订单同步
- MCP Prompt Server
- 浏览器插件
- VS Code 插件
- Cursor Rules 一键导出

