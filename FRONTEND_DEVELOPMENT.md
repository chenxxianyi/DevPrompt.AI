# DevPrompt AI 前端开发方案

## 1. 前端目标

DevPrompt AI 前端由两个独立应用组成：

- `web`：面向普通用户的 Prompt 生成、模板浏览、收藏、历史与会员页面。
- `admin`：面向管理员的后台管理系统，用于管理用户、Prompt、AI 模型、日志、会员套餐、订单和系统设置。

前端整体目标是提供现代、简洁、响应式的产品体验，并通过统一 API 封装与后端服务对接。

## 2. 前端技术栈

### 2.1 用户端 web

- Next.js 14
- TypeScript
- Tailwind CSS
- Shadcn UI
- Zustand
- React Hook Form
- Axios
- Markdown 渲染组件

### 2.2 管理后台 admin

- Vue3
- Vite
- TypeScript
- Element Plus
- Pinia
- Vue Router
- Axios

## 3. 前端目录结构

```text
DevPrompt AI/
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
```

## 4. 用户端 web 开发方案

### 4.1 页面规划

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

### 4.2 核心功能

- 首页展示核心生成入口、热门模板、产品能力和会员入口。
- Prompt 模板列表支持分类、搜索、排序和分页。
- Prompt 模板详情支持 Markdown 展示、一键复制、点赞和收藏。
- 项目 Prompt 生成器支持项目名称、项目类型、技术栈、功能需求、目标 AI 工具等表单字段。
- Cursor Rules 生成器支持语言、框架、代码风格和规则偏好输入。
- Claude Code Prompt 生成器支持任务、上下文和需求输入。
- Prompt 优化器支持输入原始 Prompt，选择目标工具与优化等级。
- 我的工作台展示用户生成概览、收藏概览、会员状态和快捷入口。
- 我的收藏展示已收藏 Prompt 模板。
- 生成历史展示用户历史生成记录。
- 登录注册页完成认证流程。
- 会员价格页展示套餐、每日生成次数和功能权益。

### 4.3 用户体验要求

- 页面风格现代、简洁、专业，适合开发者工具类产品。
- 支持桌面端、平板端和移动端响应式布局。
- 支持 SEO，关键页面需要配置 metadata。
- Prompt 内容和生成结果支持一键复制。
- 生成结果支持 Markdown 渲染。
- 表单使用 React Hook Form 管理，提供清晰校验反馈。
- 登录状态持久化，刷新页面后自动恢复用户信息。
- Token 失效时自动清理登录状态并跳转登录页。
- 会员限制、生成失败、网络异常需要给出友好提示。

### 4.4 用户端状态管理

建议使用 Zustand 拆分以下 store：

```text
authStore        登录状态、用户信息、token
promptStore      Prompt 搜索条件、分类、模板缓存
generatorStore   生成器表单状态、生成结果、加载状态
uiStore          全局 UI 状态，如主题、侧边栏、弹窗
```

### 4.5 用户端 API 封装

建议在 `web/lib` 下封装：

```text
web/lib/request.ts
web/lib/auth.ts
web/lib/prompts.ts
web/lib/generator.ts
web/lib/user.ts
```

统一处理：

- API Base URL
- Authorization Header
- 请求超时
- 业务错误提示
- 401 登录失效
- 响应数据解包

### 4.6 用户端接口依赖

```http
POST /api/auth/register
POST /api/auth/login
GET  /api/auth/profile

GET  /api/prompts
GET  /api/prompts/:slug
POST /api/prompts/:id/like
POST /api/prompts/:id/favorite

POST /api/generator/project
POST /api/generator/cursor-rules
POST /api/generator/claude-code
POST /api/generator/optimize
GET  /api/generator/history
```

## 5. 管理后台 admin 开发方案

### 5.1 页面规划

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

### 5.2 核心功能

- 后台登录：管理员账号登录，保存 token。
- 控制台：展示用户数、Prompt 数、生成次数、AI 调用量、订单概览等统计数据。
- 用户管理：分页列表、搜索、禁用、角色与会员信息查看。
- Prompt 分类管理：新增、编辑、删除、排序、启用禁用。
- Prompt 模板管理：新增、编辑、删除、分类筛选、状态管理。
- AI 模型管理：配置 provider、model、优先级、默认模型、状态和超时时间。
- AI 调用日志：查看 provider、model、调用状态、token、耗时和错误信息。
- 会员套餐管理：配置套餐价格、时长、每日限制和功能权益。
- 订单管理：查看订单状态、金额、支付时间和用户信息。
- 系统设置：维护基础系统配置。

### 5.3 后台体验要求

- 使用 Element Plus 构建表格、表单、弹窗、分页和确认框。
- 列表页统一支持分页、搜索和筛选。
- 新增编辑使用弹窗或抽屉。
- 删除操作必须二次确认。
- 后台路由需要登录权限控制。
- 管理员接口需要携带 token。
- Token 失效自动跳转后台登录页。
- API 调用错误统一提示。

### 5.4 后台状态管理

建议使用 Pinia 拆分以下 store：

```text
authStore       管理员登录状态、token、用户信息
appStore        侧边栏、主题、面包屑等应用状态
permissionStore 路由权限与菜单权限
```

### 5.5 后台 API 封装

建议在 `admin/src/api` 下拆分：

```text
admin/src/api/request.ts
admin/src/api/auth.ts
admin/src/api/users.ts
admin/src/api/categories.ts
admin/src/api/prompts.ts
admin/src/api/ai-models.ts
admin/src/api/ai-call-logs.ts
admin/src/api/membership-plans.ts
admin/src/api/orders.ts
admin/src/api/settings.ts
```

统一处理：

- API Base URL
- Authorization Header
- 请求超时
- 业务错误提示
- 401 登录失效
- Element Plus Message 提示
- 响应数据解包

### 5.6 后台接口依赖

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

后续扩展接口：

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

## 6. 前端开发阶段

### 第一阶段：前端项目初始化

- 初始化 `web` Next.js 项目。
- 初始化 `admin` Vue3 + Vite 项目。
- 配置 TypeScript、Tailwind CSS、Shadcn UI、Element Plus。
- 配置 Axios 请求封装。
- 配置基础路由、布局和登录状态持久化。

### 第二阶段：用户端核心页面

- 首页。
- Prompt 模板列表页。
- Prompt 模板详情页。
- 登录注册页。
- 会员价格页。

### 第三阶段：生成器页面

- 项目 Prompt 生成器。
- Cursor Rules 生成器。
- Claude Code Prompt 生成器。
- Prompt 优化器。
- Markdown 结果展示。
- 一键复制功能。

### 第四阶段：用户工作台

- 我的工作台。
- 我的收藏。
- 生成历史。
- 会员限制提示。

### 第五阶段：管理后台

- 后台登录。
- 控制台。
- 用户管理。
- Prompt 分类管理。
- Prompt 模板管理。
- AI 模型管理。
- AI 调用日志。
- 会员套餐管理。
- 订单管理。
- 系统设置。

### 第六阶段：前端优化

- 响应式适配。
- SEO 优化。
- 表单体验优化。
- 错误提示优化。
- Loading、Empty、Error 状态完善。
- Docker 构建与 Nginx 部署验证。

## 7. 前端环境变量

用户端：

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api
NEXT_PUBLIC_APP_NAME=DevPrompt AI
```

管理后台：

```env
VITE_API_BASE_URL=http://localhost:8080/api
VITE_APP_NAME=DevPrompt AI Admin
```

## 8. 前端启动命令

用户端：

```bash
cd web
npm install
npm run dev
```

管理后台：

```bash
cd admin
npm install
npm run dev
```

## 9. 前端交付标准

- 用户端和管理后台均可本地启动。
- 所有页面路由可访问。
- API 请求统一封装。
- 登录状态可持久化。
- 页面具备基础响应式能力。
- 核心表单有校验和错误提示。
- Prompt 内容支持复制。
- 生成结果支持 Markdown 展示。
- 后台列表具备分页、搜索、筛选和新增编辑删除流程。

