# DevPrompt AI - 前端

为开发者打造的 AI Prompt 生成平台，支持 Cursor、Claude Code、GPT、Gemini、DeepSeek、Qwen 等主流 AI 工具的 Prompt 生成、模板浏览与优化。

## 技术栈

| 层级 | 技术 |
|---|---|
| 框架 | Vue 3.4 (Composition API + `<script setup>`) |
| 语言 | TypeScript 5.4 |
| 构建工具 | Vite 5.3 |
| 样式 | Tailwind CSS 3.4 + CSS 变量 (暗色玻璃拟态主题) |
| 状态管理 | Pinia 2.1 |
| 路由 | Vue Router 4.3 (HTML5 History 模式) |
| HTTP 客户端 | Axios 1.7 |
| Markdown 渲染 | marked 12.0 |
| 代码高亮 | highlight.js 11.9 |
| 字体 | Noto Sans SC (正文) / JetBrains Mono (代码) |

## 项目结构

```
web/
├── index.html                  # 入口 HTML
├── package.json                # 项目依赖
├── vite.config.ts              # Vite 配置 (别名 + API 代理)
├── tsconfig.json               # TypeScript 严格模式配置
├── tailwind.config.js          # 自定义暗色调色板
├── postcss.config.js           # PostCSS 配置
├── public/
│   └── favicon.svg
└── src/
    ├── main.ts                 # 应用入口 (Pinia + Router)
    ├── App.vue                 # 根组件 (NavBar + router-view + Footer + Toast)
    ├── assets/
    │   └── main.css            # 全局样式、CSS 变量、玻璃拟态组件
    ├── types/
    │   └── index.ts            # TypeScript 类型定义
    ├── mock/
    │   └── data.ts             # 模拟数据 (分类、模板、套餐、技术栈、AI 工具)
    ├── api/
    │   ├── request.ts          # Axios 实例与拦截器
    │   ├── auth.ts             # 认证 API (登录、注册、用户信息)
    │   ├── prompts.ts          # 模板 API (列表、详情、点赞、收藏)
    │   └── generator.ts       # 生成器 API (项目、Cursor Rules、Claude Code、优化、历史)
    ├── store/
    │   ├── auth.ts             # 认证状态 (token、用户、localStorage 持久化)
    │   ├── prompt.ts           # 模板状态 (筛选、搜索、排序、点赞/收藏)
    │   ├── generator.ts        # 生成器表单状态
    │   └── ui.ts               # UI 状态 (Toast 通知)
    ├── router/
    │   └── index.ts            # 路由配置 (6 条路由，懒加载)
    ├── components/
    │   ├── NavBar.vue          # 顶部导航栏
    │   ├── AppFooter.vue       # 页脚
    │   ├── ToastContainer.vue  # Toast 通知容器
    │   └── PromptCard.vue      # 模板卡片组件
    └── views/
        ├── HomeView.vue        # 首页
        ├── PromptsView.vue     # 模板库
        ├── PromptDetailView.vue # 模板详情
        ├── GeneratorView.vue   # Prompt 生成器
        ├── LoginView.vue       # 登录/注册
        └── PricingView.vue     # 会员套餐
```

## 页面与路由

| 路径 | 页面 | 说明 |
|---|---|---|
| `/` | 首页 | Hero 区域、生成器入口卡片、热门模板、核心能力展示、Pro 会员引导 |
| `/prompts` | 模板库 | 分类侧边栏、关键词搜索、排序 (热门/最新/最多点赞)、模板网格 |
| `/prompts/:slug` | 模板详情 | 模板内容展示、点赞/收藏/一键复制、相关模板推荐 |
| `/generator` | 生成器 | 项目名称、类型、技术栈多选、功能描述、目标 AI 工具选择 |
| `/login` | 登录/注册 | 邮箱密码登录注册、社交登录 (GitHub/Google/微信) |
| `/pricing` | 会员套餐 | Free / Pro / Team / Enterprise 四档方案 |

## 核心功能

- **Prompt 生成** — 支持项目 Prompt、Cursor Rules、Claude Code Prompt、Prompt 优化四种生成模式
- **模板库** — 7 大分类 (前端/后端/DevOps/数据库/AI-ML/移动端)，支持搜索与排序
- **模板详情** — 完整内容展示、点赞/收藏互动、一键复制到剪贴板、相关模板推荐
- **会员体系** — Free (5次/天)、Pro ($9.9/月, 100次/天)、Team ($29.9/月, 500次/天)、Enterprise (无限, 定制)
- **Toast 通知** — 全局通知系统，2.5 秒自动消失，滑入/滑出动画

## API 接口

所有接口通过 `/api` 前缀代理至后端 `http://localhost:8080`：

| 模块 | 接口 | 方法 |
|---|---|---|
| 认证 | `/auth/login` | POST |
| 认证 | `/auth/register` | POST |
| 认证 | `/auth/profile` | GET |
| 模板 | `/prompts` | GET |
| 模板 | `/prompts/:slug` | GET |
| 模板 | `/prompts/:id/like` | POST |
| 模板 | `/prompts/:id/favorite` | POST |
| 生成器 | `/generator/project` | POST |
| 生成器 | `/generator/cursor-rules` | POST |
| 生成器 | `/generator/claude-code` | POST |
| 生成器 | `/generator/optimize` | POST |
| 生成器 | `/generator/history` | GET |

请求拦截器自动附加 `Bearer` Token，响应拦截器处理业务错误码与 401 跳转。

## 开发指南

### 环境要求

- Node.js >= 18
- npm >= 9

### 安装与启动

```bash
# 安装依赖
npm install

# 启动开发服务器 (端口 3000)
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview
```

### 环境变量

| 变量 | 说明 | 默认值 |
|---|---|---|
| `VITE_API_BASE_URL` | API 基础地址 | `/api` |

### 路径别名

`@` 映射至 `src/` 目录，在 TypeScript 和 Vite 中均已配置。

## 设计风格

采用暗色玻璃拟态 (Glassmorphism) 主题：

- **主色调** — 紫色 `#8b5cf6` + 青色 `#06b6d4` 双强调色
- **背景** — 深色渐变 + SVG 噪点纹理叠加
- **卡片** — 半透明毛玻璃效果 (`.glass`)
- **按钮** — 主按钮 (`.btn-primary`) / 幽灵按钮 (`.btn-ghost`) / 青色按钮 (`.btn-cyan`)
- **标签** — 默认标签 (`.tag`) / 青色标签 (`.tag-cyan`)
- **滚动条** — 自定义紫色强调色滚动条
- **动画** — 页面进入 (`pageIn`)、脉冲 (`pulse`)、加载旋转 (`spin`)

## 当前状态

项目处于 **前端 Demo 阶段**：

- 数据来源为 `src/mock/data.ts` 中的模拟数据，尚未对接真实 API
- 生成器使用 `setTimeout` 模拟生成过程，输出硬编码模板
- 登录/注册展示 Demo Toast，未调用真实接口
- `marked` 和 `highlight.js` 已安装但尚未在视图中使用，预留用于 Markdown 渲染与代码高亮
