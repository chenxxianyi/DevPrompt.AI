# DevPrompt AI Admin

`admin/` 是项目的管理后台前端，服务于运营和内容维护。

## 技术栈

- Vue 3
- Vite
- TypeScript
- Element Plus
- Pinia
- Vue Router
- Axios

## 本地开发

```bash
cd admin
npm install
npm run dev
```

默认地址：`http://localhost:5173`

构建命令：

```bash
npm run build
```

## 环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `VITE_API_BASE_URL` | `/api` | API 基础地址 |
| `VITE_APP_NAME` | `DevPrompt AI Admin` | 应用名称 |

开发模式下，Vite 会将 `/api` 代理到 `http://localhost:8080`。

## 登录说明

后台依赖服务端管理员账号。

本地默认账号来自后端种子数据：

- 邮箱：`admin@devprompt.ai`
- 密码：`admin123`

仅建议用于本地开发。

## 页面范围

当前已有页面：

- `/login`
- `/dashboard`
- `/users`
- `/categories`
- `/prompts`
- `/ai-models`
- `/ai-call-logs`
- `/membership-plans`
- `/orders`
- `/project-types`

## 已接入 API

对应的 API 模块位于：

- `api/auth.ts`
- `api/dashboard.ts`
- `api/users.ts`
- `api/categories.ts`
- `api/prompts.ts`
- `api/aiModels.ts`
- `api/callLogs.ts`
- `api/plans.ts`
- `api/orders.ts`
- `api/projectTypes.ts`

请求层会自动：

- 从 `localStorage` 读取 `admin_token`
- 为请求附加 `Authorization: Bearer <token>`
- 在 `401` / `403` 时清理登录状态并跳回登录页

## 当前能力边界

- 后端已经提供试用申请管理 API，但当前后台还没有对应页面
- Docker Compose 目前不会自动部署 `admin/`
- 页面以 CRUD 和基础统计为主，运营分析能力仍比较轻
