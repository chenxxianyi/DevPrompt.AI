# DevPrompt AI

面向开发者的 AI Prompt 工作台，当前仓库包含用户端、管理后台、Go API 服务、数据库初始化脚本和 Docker 部署配置。

## 当前实现

- 用户端 `web/`：浏览 Prompt 模板、查看详情、登录注册、四类 Prompt 生成、生成历史、个人工作台、会员方案页
- 管理后台 `admin/`：用户、分类、模板、AI 模型、调用日志、会员方案、订单、项目类型管理
- 服务端 `server/`：认证、模板、生成、会员、用户统计、试用申请、后台管理 API
- 基础设施：MySQL、Redis、Nginx、`docker-compose.yml`

## 技术栈

| 模块 | 技术 |
| --- | --- |
| `web` | Vue 3 + Vite + TypeScript + Tailwind CSS + Pinia + Vue Router |
| `admin` | Vue 3 + Vite + TypeScript + Element Plus + Pinia + Vue Router |
| `server` | Go 1.23 + Gin + Gorm + MySQL + Redis |
| 部署 | Docker Compose + Nginx |

## 仓库结构

```text
DevPrompt AI/
├─ web/                 用户端前端
├─ admin/               管理后台前端
├─ server/              Go API 服务
├─ database/            数据库初始化脚本
├─ deploy/nginx/        Nginx 配置
├─ docker-compose.yml   一体化部署编排
└─ .env.example         环境变量示例
```

## 本地启动

### 1. 准备环境

- Node.js 18+
- npm 9+
- Go 1.23+
- Docker Desktop 或本地 MySQL 8 / Redis 7

### 2. 复制环境变量

将根目录的 `.env.example` 复制为 `.env`，至少确认以下值正确：

- `JWT_SECRET`
- `MYSQL_*`
- `REDIS_*`
- `WEB_BASE_URL`
- `ADMIN_BASE_URL`
- 任一可用的 AI Provider Key，例如 `DEEPSEEK_API_KEY`

说明：

- 仅浏览模板、查看页面时可以不配 AI Key
- 调用生成接口时必须至少配置一个可用 Provider

### 3. 启动 MySQL 和 Redis

```bash
docker compose up -d mysql redis
```

如果你本地已经有 MySQL 和 Redis，也可以直接复用，只要把配置改到 `server/config/config.yaml` 或环境变量里即可。

### 4. 启动后端

```bash
cd server
go run ./cmd/main.go
```

默认地址：`http://localhost:8080`

后端启动时会自动：

- 创建数据库（若不存在）
- 执行 `AutoMigrate`
- 初始化基础种子数据
- 创建默认管理员账号（仅在用户表为空时）

默认管理员账号：

- 邮箱：`admin@devprompt.ai`
- 密码：`admin123`

仅建议本地开发使用，首次登录后应立即修改。

### 5. 启动用户端

```bash
cd web
npm install
npm run dev
```

默认地址：`http://localhost:3000`

### 6. 启动管理后台

```bash
cd admin
npm install
npm run dev
```

默认地址：`http://localhost:5173`

## Docker Compose

根目录 `docker-compose.yml` 当前会启动：

- `mysql`
- `redis`
- `server`
- `web`

说明：

- 对外入口是 `http://localhost:3000`
- `web` 容器内使用 Nginx 反向代理 `/api` 到 `server:8080`
- 管理后台 `admin/` 目前不在 Compose 中，需要本地单独运行或自行补充部署

完整启动：

```bash
docker compose up -d --build
```

## 主要页面

### 用户端

- `/`：首页与快速生成入口
- `/prompts`：模板列表
- `/prompts/:slug`：模板详情
- `/generator`：生成工作台
- `/generator/history`：生成历史
- `/dashboard`：个人工作台
- `/pricing`：会员方案
- `/login`：登录注册

### 管理后台

- `/dashboard`
- `/users`
- `/categories`
- `/prompts`
- `/ai-models`
- `/ai-call-logs`
- `/membership-plans`
- `/orders`
- `/project-types`

## API 概览

### 公开接口

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/prompts`
- `GET /api/prompts/:slug`
- `GET /api/categories`
- `GET /api/membership/plans`
- `GET /api/project-types`

### 登录后接口

- `GET /api/auth/profile`
- `POST /api/prompts/:id/like`
- `POST /api/prompts/:id/favorite`
- `POST /api/generator/project`
- `POST /api/generator/cursor-rules`
- `POST /api/generator/claude-code`
- `POST /api/generator/optimize`
- `GET /api/generator/history`
- `GET /api/user/generate-stats`
- `POST /api/trial-requests`

### 管理接口

- `GET /api/admin/dashboard`
- `GET /api/admin/users`
- `GET/POST/PUT/DELETE /api/admin/prompts`
- `GET/POST/PUT/DELETE /api/admin/categories`
- `GET/POST/PUT/DELETE /api/admin/ai-models`
- `GET /api/admin/ai-call-logs`
- `GET/POST/PUT/DELETE /api/admin/membership-plans`
- `GET /api/admin/orders`
- `GET/POST/PUT/DELETE /api/admin/project-types`
- `GET /api/admin/trial-requests`
- `PUT /api/admin/trial-requests/:id`

## 需要知道的现状

- `database/init.sql` 仍可作为初始化参考，但运行时以后端的 `AutoMigrate + Seed` 为准
- 用户端已经接真实 API，不再是纯 mock Demo
- 管理后台已有试用申请相关后端接口，但当前前端还没有对应管理页面
- 根目录存在若干规划类文档，下面的链接已按当前仓库状态更新

## 文档索引

- [server/README.md](server/README.md)
- [web/README.md](web/README.md)
- [admin/README.md](admin/README.md)
- [NEXT_OPTIMIZATION_PLAN.md](NEXT_OPTIMIZATION_PLAN.md)
- [PRODUCT_OPTIMIZATION_DEVELOPMENT.md](PRODUCT_OPTIMIZATION_DEVELOPMENT.md)
- [PR_FIX_GUIDE.md](PR_FIX_GUIDE.md)
