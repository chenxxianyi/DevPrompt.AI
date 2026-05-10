## 快速说明 — 目标与上下文

你将辅助开发者在 DevPrompt AI 后端快速定位关键代码、修改功能并安全调用/测试 AI 相关流程。
本仓库后端位于 `server/`，基于 Go + Gin + Gorm，提供 Prompt 生成、AI Provider 管理、会员与限流等功能。

## 项目高层架构（必读引用）
- 入口与依赖注入: `server/cmd/main.go` — 加载 `config/config.yaml`、初始化 DB/Redis、注册 provider、组装 service 与 HTTP 路由（AutoMigrate 在此执行）。
- HTTP 层（handler）: `server/internal/api/*.go`（例如 `generator.go`） — 负责参数绑定、鉴权（参考 `middleware.JWTAuth`）、构建 system prompt 并调用 `GeneratorService.Generate`。
- 业务层（service）: `server/internal/service/generator_service.go` — 核心逻辑：权限/会员检查、限流、调用 provider manager、记录 `AICallLog`、保存 `GeneratedPrompt`。
- Provider 抽象: `server/internal/provider/interface.go` 与 `server/internal/provider/manager.go`（注册/切换逻辑在 `main.go`）— 新的 provider 必须实现 `AIProvider` 接口并在 `main.go` 注册。
- 仓储层（repository）: `server/internal/repository/*` — Gorm 封装，数据持久化集中处理。
- 配置/敏感信息: 优先使用环境变量覆盖 `config/config.yaml`（参见 `server/README.md` 与 `docker-compose.yml` 环境项）。

## 重要开发/运行命令（复制可用）
- 本地快速运行（需先修改 `server/config/config.yaml` 或设置 env vars）:
  - cd server; go run ./cmd/main.go
- 编译:
  - go build -o server ./cmd/main.go
- 单元测试:
  - go test ./...
- 容器化一键启动（推荐开发）：项目根目录下运行
  - docker compose up -d
  该命令会启动 MySQL、Redis、后端、前端与 Nginx；DB 初始化脚本位于 `database/init.sql`。

## 与 AI 相关的典型工作流（示例）
- 用户端请求流程：HTTP POST `/api/generator/cursor-rules` → `internal/api/generator.go` 绑定为 `GenerateCursorRulesRequest` → 构造 system prompt → 调用 `GeneratorService.Generate` → `provider.Manager.Call` 触发具体 provider（OpenAI/Claude/DeepSeek 等）→ 记录 `ai_call_logs` 并保存 `generated_prompts`。
- Provider 扩展：
  - 实现接口：`server/internal/provider/interface.go` 中的 `AIProvider`（Name, Chat）
  - 在 `server/cmd/main.go` 注册：providerMgr.Register("your-provider", provider.NewYourProvider(...))
  - 注意：API keys 仅通过环境变量注入（不要硬编码或暴露到前端）。

## 项目约定与惯例（不走寻常路的点）
- 鉴权：两套中间件 `JWTAuth(secret)`（强制）与 `OptionalAuth(secret)`（可选）；handler 里通过 `middleware.GetUserID(c)` 读取 user_id。参考：`server/internal/middleware/jwt.go`。
- 统一响应格式：所有接口使用 `internal/response/response.go` 的结构（{code,message,data}）。修改时应遵循该结构。
- 限流/计数：每日生成计数与速率限制由 Redis 管理（`rate_limit_service.go`），GeneratorService 在调用前会检查并在成功后自增计数。
- 初始化与迁移：`main.go` 在启动时执行 `db.AutoMigrate(...)` 并调用 `model.Seed(db)`，不要重复在其他位置运行迁移以避免冲突。

## 编辑/调试小贴士（对 AI 助手有用）
- 若改动 provider，请同时添加/调整 unit tests 并在 `main.go` 中条件性注册（以便 CI 在无 key 环境下不失败）。
- 调试 API 时：先用 `docker compose up -d` 启动依赖（MySQL/Redis），或本地启动 MySQL 并在 `server/config/config.yaml` 写入连接信息。
- 快速复现 AI 调用：构造 POST 到 `/api/generator/project` 或 `/api/generator/cursor-rules` 的请求（参考 `server/README.md` 中 JSON 示例），并在 headers 加上 `Authorization: Bearer <token>`。

## 关键文件索引（快速导航）
- 路由与注册：`server/cmd/main.go`
- 生成业务：`server/internal/service/generator_service.go`
- 生成 Handler：`server/internal/api/generator.go`
- Provider 接口：`server/internal/provider/interface.go`
- 中间件（JWT/CORS/Admin）：`server/internal/middleware/*.go`
- 配置文件：`server/config/config.yaml`（可被 env 覆盖）
- DB 初始化：`database/init.sql`
- Docker Compose：`docker-compose.yml`

## 安全与测试约束
- 不要将 API keys 或 JWT secrets 写入源码或 `web/dist`。所有 keys 从环境变量读取（见 `docker-compose.yml`）。
- 在 CI 环境中，mock provider 或使用无网络的回放/假数据来验证生成逻辑与数据库交互：避免真实调用外部 AI API。

——
如果你想让我把其中某一节扩展成可执行脚本或把常用 curl/postman 示例加入文件，我可以继续补充；也请指出你觉得不完整或需要更详细的区域。 
