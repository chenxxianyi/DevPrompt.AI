# DevPrompt AI Server

`server/` 是项目的 Go API 服务，负责认证、模板、生成、会员、用户统计和后台管理接口。

## 技术栈

- Go 1.23
- Gin
- Gorm
- MySQL 8
- Redis 7
- JWT

## 目录结构

```text
server/
├─ cmd/main.go                程序入口
├─ config/config.yaml         默认配置
├─ internal/api/              用户端 API
├─ internal/api/admin/        管理后台 API
├─ internal/config/           配置加载
├─ internal/middleware/       JWT / 管理员鉴权 / CORS
├─ internal/model/            数据模型与种子数据
├─ internal/provider/         AI Provider 封装
├─ internal/repository/       数据访问层
├─ internal/service/          业务逻辑
├─ internal/response/         统一响应
└─ internal/utils/            JWT / 密码工具
```

## 本地运行

### 前置条件

- Go 1.23+
- 可用的 MySQL
- 可用的 Redis

### 启动命令

```bash
cd server
go run ./cmd/main.go
```

默认监听：`http://localhost:8080`

## 配置方式

服务默认读取 `config/config.yaml`，同时支持环境变量覆盖。

常用环境变量：

| 变量 | 说明 |
| --- | --- |
| `APP_ENV` | 运行环境，常用 `development` / `production` |
| `APP_PORT` | 服务端口 |
| `MYSQL_HOST` `MYSQL_PORT` `MYSQL_USER` `MYSQL_PASSWORD` `MYSQL_DATABASE` | MySQL 连接 |
| `REDIS_HOST` `REDIS_PORT` `REDIS_PASSWORD` | Redis 连接 |
| `JWT_SECRET` | JWT 密钥 |
| `JWT_EXPIRE_HOURS` | Token 过期时间 |
| `OPENAI_API_KEY` | OpenAI Key |
| `CLAUDE_API_KEY` | Claude Key |
| `GEMINI_API_KEY` | Gemini Key |
| `DEEPSEEK_API_KEY` | DeepSeek Key |
| `QWEN_API_KEY` | Qwen Key |
| `WEB_BASE_URL` | 用户端地址，用于 CORS |
| `ADMIN_BASE_URL` | 管理后台地址，用于 CORS |

## 启动行为

服务启动后会依次执行：

1. 加载配置
2. 校验生产环境下的 `JWT_SECRET`
3. 连接 MySQL，并在数据库不存在时自动创建
4. 执行 `AutoMigrate`
5. 初始化种子数据
6. 连接 Redis
7. 注册 Provider、Service、Handler 和路由

说明：

- Redis 连接失败时，服务仍可启动，但会关闭限流能力
- 默认管理员账号只会在 `users` 表为空时创建

## 当前数据模型

核心表包括：

- `users`
- `prompt_categories`
- `prompt_templates`
- `prompt_favorites`
- `prompt_likes`
- `generated_prompts`
- `ai_models`
- `ai_call_logs`
- `membership_plans`
- `orders`
- `project_types`
- `trial_requests`

## 路由概览

### 公开接口

- 认证：注册、登录
- 模板：列表、详情、分类
- 会员方案：列表
- 项目类型：公开列表

### 登录后接口

- 当前用户信息
- 模板点赞与收藏
- 四类生成接口
- 生成历史
- 用户额度统计
- 试用申请提交

### 管理接口

- 仪表盘统计
- 用户管理
- 模板管理
- 分类管理
- AI 模型管理
- 调用日志查询
- 会员方案管理
- 订单查询
- 项目类型管理
- 试用申请管理

## 生成相关说明

当前服务支持四类生成：

- `project`
- `cursor-rules`
- `claude-code`
- `optimize`

Provider 会根据当前已配置的 API Key 进行注册。若没有任何可用 Provider，生成接口将无法正常工作，但其他非生成功能仍可使用。

## Docker

构建镜像：

```bash
cd server
docker build -t devprompt-server .
```

`Dockerfile` 使用多阶段构建，运行阶段基于 Alpine，仅包含编译后的服务与 `config/config.yaml`。

## 当前已知边界

- 后台可以维护 `ai_models`，但默认模型策略仍以当前服务实现为准，使用前建议结合代码确认
- 试用申请已有后台 API，但管理后台前端暂未提供对应页面
- `database/init.sql` 与运行时种子逻辑并不是唯一事实来源，实际以服务启动后的迁移和种子数据为准
