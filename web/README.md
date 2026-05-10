# DevPrompt AI Web

`web/` 是面向普通用户的前端应用，负责模板浏览、Prompt 生成、历史记录、个人工作台和会员页。

## 技术栈

- Vue 3
- Vite
- TypeScript
- Tailwind CSS
- Pinia
- Vue Router
- Axios
- `marked`
- `highlight.js`

## 本地开发

```bash
cd web
npm install
npm run dev
```

默认地址：`http://localhost:3000`

构建命令：

```bash
npm run build
```

## 环境变量

当前使用的变量：

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `VITE_API_BASE_URL` | `http://localhost:8080/api` 或 `/api` | API 基础地址 |
| `VITE_APP_NAME` | `DevPrompt AI` | 应用名称 |

说明：

- 开发环境下 `vite.config.ts` 已将 `/api` 代理到 `http://localhost:8080`
- 仓库中的 `web/.env` 默认直接指向 `http://localhost:8080/api`

## 页面与路由

| 路由 | 说明 |
| --- | --- |
| `/` | 首页，包含快速输入和四类生成入口 |
| `/prompts` | 模板列表 |
| `/prompts/:slug` | 模板详情 |
| `/generator` | 生成工作台 |
| `/generator/history` | 生成历史 |
| `/dashboard` | 个人工作台 |
| `/login` | 登录注册 |
| `/pricing` | 会员方案页 |

## 当前功能

### 首页

- 快速输入任务或原始 Prompt
- 选择目标 AI 工具
- 跳转到对应生成场景
- 展示热门模板

### 模板库

- 模板列表
- 分类筛选
- 关键词搜索
- 排序
- 模板详情
- 登录后点赞与收藏

### 生成工作台

当前支持四个 Tab：

- 项目 Prompt
- Cursor Rules
- Claude Code
- Prompt 优化

生成成功后支持：

- 复制结果
- 再次生成
- 项目 Prompt / Claude Code 导出 Markdown
- Prompt 优化前后对比

### 个人工作台

- 今日额度查看
- 当前会员等级
- 最近生成记录
- 从历史记录再次生成

### 会员页

- 展示会员方案
- 登录用户可提交试用申请
- `free` 方案跳登录，其余方案调用 `/api/trial-requests`

## 与后端的接口关系

主要 API 模块：

- `api/auth.ts`
- `api/prompts.ts`
- `api/categories.ts`
- `api/generator.ts`
- `api/membership.ts`
- `api/user.ts`
- `api/trial.ts`

认证方式：

- `request.ts` 会自动从 `localStorage` 读取 `token`
- 遇到 `401` 会清理本地登录态并跳转 `/login`

## 当前实现细节

- `projectTypes` 来自后端接口 `/api/project-types`
- `techStacks` 和 `aiTools` 仍来自 `src/mock/data.ts`
- 历史页和工作台已经接真实后端接口，不再是纯前端假数据流程

## 适合继续补强的地方

- 模板适用工具、示例输入输出等元数据还不够完整
- 会员页已支持提交试用申请，但前端未提供联系人和备注输入
- 常用技术栈与常用 AI 工具仍是前端静态配置
