# DevPrompt AI 后端服务

DevPrompt AI 是一款智能 Prompt 生成与管理的 SaaS 平台后端服务。基于 Golang + Gin + Gorm + MySQL + Redis 构建。

## 技术栈

- **语言**: Golang 1.23
- **Web 框架**: Gin 1.10
- **ORM**: Gorm 1.25
- **数据库**: MySQL 8
- **缓存**: Redis 7
- **认证**: JWT (golang-jwt)
- **配置**: Viper + 环境变量
- **日志**: Zap

## 项目结构

```
server/
├── cmd/
│   └── main.go                 # 入口：依赖注入、路由注册、启动服务
├── config/
│   └── config.yaml             # 默认配置文件
├── internal/
│   ├── api/                    # HTTP Handler 层
│   │   ├── auth.go             # 注册、登录、用户信息
│   │   ├── prompt.go           # Prompt 列表、详情、点赞、收藏
│   │   ├── generator.go        # Prompt 生成接口（4 种 + 历史）
│   │   ├── membership.go       # 会员套餐查询
│   │   ├── user.go             # 用户生成统计
│   │   └── admin/              # 管理后台接口
│   │       ├── user.go         # 用户管理
│   │       ├── prompt.go       # Prompt 模板管理
│   │       ├── category.go     # 分类管理
│   │       ├── ai_model.go     # AI 模型管理
│   │       ├── call_log.go     # AI 调用日志
│   │       ├── plan.go         # 会员套餐管理
│   │       ├── order.go        # 订单管理
│   │       └── dashboard.go    # 数据概览
│   ├── middleware/
│   │   ├── jwt.go              # JWT 认证中间件（强制 + 可选）
│   │   └── admin.go            # 管理员权限中间件
│   ├── model/                  # Gorm 数据模型
│   │   ├── base.go             # 基础模型（ID、时间戳）
│   │   ├── user.go             # 用户
│   │   ├── prompt_category.go  # Prompt 分类
│   │   ├── prompt_template.go  # Prompt 模板
│   │   ├── prompt_favorite.go  # 收藏关联
│   │   ├── prompt_like.go      # 点赞关联
│   │   ├── generated_prompt.go # 生成记录
│   │   ├── ai_model.go         # AI 模型配置
│   │   ├── ai_call_log.go      # AI 调用日志
│   │   ├── membership_plan.go  # 会员套餐
│   │   └── order.go            # 订单
│   ├── repository/             # 数据访问层
│   │   ├── user_repository.go
│   │   ├── prompt_category_repository.go
│   │   ├── prompt_template_repository.go
│   │   ├── prompt_favorite_repository.go
│   │   ├── prompt_like_repository.go
│   │   ├── generated_prompt_repository.go
│   │   ├── ai_model_repository.go
│   │   ├── ai_call_log_repository.go
│   │   ├── membership_plan_repository.go
│   │   └── order_repository.go
│   ├── service/                # 业务逻辑层
│   │   ├── auth_service.go     # 用户认证
│   │   ├── prompt_service.go   # Prompt 业务
│   │   ├── generator_service.go # AI 生成编排
│   │   ├── membership_service.go # 会员套餐
│   │   └── rate_limit_service.go # Redis 限流
│   ├── provider/               # AI Provider 抽象
│   │   ├── interface.go        # 统一接口定义
│   │   ├── manager.go          # 模型选择 + 失败切换
│   │   ├── openai.go           # OpenAI / GPT-4o
│   │   ├── claude.go           # Anthropic Claude
│   │   ├── gemini.go           # Google Gemini
│   │   ├── deepseek.go         # DeepSeek
│   │   └── qwen.go             # 通义千问
│   ├── response/
│   │   └── response.go         # 统一响应格式
│   ├── config/
│   │   └── config.go           # 配置加载（YAML + 环境变量覆盖）
│   └── utils/
│       ├── jwt.go              # JWT 生成与解析
│       └── password.go         # bcrypt 密码哈希
├── Dockerfile                  # 多阶段构建
├── go.mod / go.sum             # Go 模块依赖
└── README.md                   # 本文档
```

## 快速开始

### 前置要求

- Go 1.23+
- MySQL 8.0+
- Redis 7.0+

### 开发模式

1. **克隆项目并进入目录**

```bash
cd server
```

2. **配置数据库和 Redis**

编辑 `config/config.yaml`，修改 MySQL 和 Redis 连接信息。

3. **运行**

```bash
go run ./cmd/main.go
```

服务启动在 `http://localhost:8080`。

### Docker Compose 一键启动（推荐）

在项目根目录执行：

```bash
docker compose up -d
```

此命令会启动 MySQL、Redis、后端服务、前端和 Nginx 五个容器。

## API 文档

所有接口统一返回格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

### 公开接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/register` | 用户注册 |
| POST | `/api/auth/login` | 用户登录 |
| GET | `/api/prompts` | Prompt 模板列表（支持 keyword/category/page/pageSize/sort） |
| GET | `/api/prompts/:slug` | Prompt 模板详情 |
| GET | `/api/membership/plans` | 会员套餐列表 |

### 需要认证的接口

在请求头添加 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/auth/profile` | 获取当前用户信息 |
| POST | `/api/prompts/:id/like` | 点赞 / 取消点赞 |
| POST | `/api/prompts/:id/favorite` | 收藏 / 取消收藏 |
| POST | `/api/generator/project` | 生成项目 Prompt |
| POST | `/api/generator/cursor-rules` | 生成 Cursor Rules |
| POST | `/api/generator/claude-code` | 生成 Claude Code Prompt |
| POST | `/api/generator/optimize` | 优化 Prompt |
| GET | `/api/generator/history` | 生成历史记录 |
| GET | `/api/user/generate-stats` | 生成次数统计 |

### 管理后台接口

需要 `admin` 角色 + Bearer token 访问。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/admin/dashboard` | 数据概览 |
| GET | `/api/admin/users` | 用户列表 |
| GET/POST | `/api/admin/prompts` | Prompt 模板管理 |
| PUT/DELETE | `/api/admin/prompts/:id` | 编辑/删除 Prompt |
| GET/POST | `/api/admin/categories` | 分类管理 |
| PUT/DELETE | `/api/admin/categories/:id` | 编辑/删除分类 |
| GET/POST | `/api/admin/ai-models` | AI 模型管理 |
| PUT/DELETE | `/api/admin/ai-models/:id` | 编辑/删除模型 |
| GET | `/api/admin/ai-call-logs` | AI 调用日志 |
| GET/POST | `/api/admin/membership-plans` | 会员套餐管理 |
| PUT/DELETE | `/api/admin/membership-plans/:id` | 编辑/删除套餐 |
| GET | `/api/admin/orders` | 订单管理 |

## 生成接口请求示例

### 项目 Prompt 生成

```json
POST /api/generator/project
{
  "projectName": "DevPrompt AI",
  "projectType": "SaaS",
  "techStack": ["Next.js", "Golang", "MySQL"],
  "features": ["登录注册", "Prompt 生成", "会员系统"],
  "targetAiTool": "Cursor"
}
```

### Prompt 优化

```json
POST /api/generator/optimize
{
  "rawPrompt": "帮我写一个网站",
  "targetTool": "Cursor",
  "optimizeLevel": "professional"
}
```

## AI Provider 配置

所有 AI API Key 从环境变量读取，**不会暴露给前端**。

| 环境变量 | Provider | 默认模型 |
|----------|----------|----------|
| `OPENAI_API_KEY` | OpenAI | gpt-4o |
| `CLAUDE_API_KEY` | Anthropic Claude | claude-sonnet-4-6 |
| `GEMINI_API_KEY` | Google Gemini | gemini-2.0-flash |
| `DEEPSEEK_API_KEY` | DeepSeek | deepseek-chat |
| `QWEN_API_KEY` | 通义千问 | qwen-plus |

Provider Manager 支持：
- 按优先级自动选择模型
- 调用失败后自动切换到下一个可用模型
- 自动重试（最多 2 次）
- 不健康标记（30 秒后自动恢复）

## 配置说明

配置优先级：环境变量 > `config.yaml` 文件

### 关键配置项

| 配置项 | 环境变量 | 默认值 | 说明 |
|--------|----------|--------|------|
| `app.port` | `APP_PORT` | 8080 | 服务端口 |
| `mysql` | `MYSQL_*` | localhost:3306 | MySQL 连接 |
| `redis` | `REDIS_*` | localhost:6379 | Redis 连接 |
| `jwt.secret` | `JWT_SECRET` | — | JWT 签名密钥 |
| `jwt.expire_hours` | `JWT_EXPIRE_HOURS` | 168 | Token 过期时间（小时） |
| `rate_limit.enabled` | — | true | 是否启用 Redis 限流 |

## 会员等级与每日限额

| 等级 | 每日生成次数 | 价格 |
|------|-------------|------|
| Free | 5 次 | 免费 |
| Pro | 100 次 | $29/月 |
| Team | 500 次 | $99/月 |
| Enterprise | 不限 | $299/年 |

## 开发命令

```bash
# 编译
go build -o server ./cmd/main.go

# 运行
go run ./cmd/main.go

# 测试
go test ./...

# 代码格式化
go fmt ./...
```

## 数据库

数据库初始化脚本位于 `../database/init.sql`，包含：
- 10 张数据表的 DDL
- 会员套餐种子数据
- Prompt 分类种子数据
- AI 模型种子数据
- 默认管理员账号

Docker Compose 启动时会自动执行初始化脚本。
