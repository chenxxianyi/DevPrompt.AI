import type { PromptCategory, PromptTemplate, MembershipPlan } from '@/types'

export const categories: PromptCategory[] = [
  { id: 0, name: '全部', slug: 'all', description: '所有模板', sort: 0, status: 'active', count: 9 },
  { id: 1, name: '前端开发', slug: 'frontend', description: '前端开发相关 Prompt', sort: 1, status: 'active', count: 3 },
  { id: 2, name: '后端开发', slug: 'backend', description: '后端开发相关 Prompt', sort: 2, status: 'active', count: 3 },
  { id: 3, name: 'DevOps', slug: 'devops', description: 'DevOps 相关 Prompt', sort: 3, status: 'active', count: 1 },
  { id: 4, name: '数据库', slug: 'database', description: '数据库相关 Prompt', sort: 4, status: 'active', count: 1 },
  { id: 5, name: 'AI/ML', slug: 'ai-ml', description: 'AI/ML 相关 Prompt', sort: 5, status: 'active', count: 1 },
  { id: 6, name: '移动端', slug: 'mobile', description: '移动端开发 Prompt', sort: 6, status: 'active', count: 1 },
]

export const templates: PromptTemplate[] = [
  {
    id: 1, categoryId: 1, title: 'Next.js 全栈项目 Prompt', slug: 'nextjs-fullstack',
    description: '适用于 Next.js 14 + TypeScript + Tailwind CSS 全栈项目的开发 Prompt，包含路由设计、API 层、认证和部署配置。',
    tags: ['Next.js', 'TypeScript', 'Tailwind'], useCount: 520, likeCount: 328, favoriteCount: 156,
    status: 'active', createdAt: '2025-01-15', updatedAt: '2025-01-20',
    content: `# Next.js 全栈项目开发 Prompt

## 项目概述
你是一个资深的全栈开发工程师，擅长使用 Next.js 14 + TypeScript + Tailwind CSS 构建现代 Web 应用。

## 技术栈
- Framework: Next.js 14 (App Router)
- Language: TypeScript
- Styling: Tailwind CSS + Shadcn UI
- State: Zustand
- Database: Prisma + PostgreSQL
- Auth: NextAuth.js

## 编码规范
1. 使用 App Router，所有页面放在 app/ 目录下
2. 组件使用 Server Components 优先，仅在需要交互时使用 Client Components
3. 使用 TypeScript 严格模式，所有 props 必须定义类型
4. API 路由使用 Route Handlers，返回统一响应格式
5. 错误处理使用 error.tsx 和 not-found.tsx

## 项目结构
\`\`\`
app/
  (auth)/login/page.tsx
  (dashboard)/dashboard/page.tsx
  api/
    auth/[...nextauth]/route.ts
    users/route.ts
components/
  ui/
  layout/
lib/
  request.ts
  auth.ts
types/
\`\`\`

## 要求
- 遵循 RESTful API 设计
- 实现 JWT 认证和权限控制
- 使用 React Hook Form 管理表单
- 响应式布局，支持移动端
- 代码注释使用中文`,
  },
  {
    id: 2, categoryId: 2, title: 'Golang Gin 微服务 Prompt', slug: 'golang-gin-microservice',
    description: '基于 Golang Gin 框架的微服务开发 Prompt，涵盖分层架构、中间件、数据库和 Redis 缓存。',
    tags: ['Golang', 'Gin', '微服务'], useCount: 410, likeCount: 256, favoriteCount: 98,
    status: 'active', createdAt: '2025-01-10', updatedAt: '2025-01-18',
    content: `# Golang Gin 微服务开发 Prompt

## 角色定义
你是一个精通 Golang 的后端架构师，擅长使用 Gin 框架构建高性能微服务。

## 技术栈
- Language: Go 1.21+
- Framework: Gin
- ORM: GORM
- Cache: Redis
- MQ: RabbitMQ
- Config: Viper
- Log: Zap

## 分层架构
\`\`\`
internal/
  api/        # HTTP Handler
  service/    # 业务逻辑
  repository/ # 数据访问
  model/      # 数据模型
  middleware/ # 中间件
  provider/   # 外部服务
  response/   # 统一响应
  utils/      # 工具函数
\`\`\`

## 编码规范
1. 接口返回统一格式: {code, message, data}
2. 错误使用自定义错误码，不暴露内部信息
3. 数据库操作使用 GORM，禁止裸 SQL
4. Redis Key 使用冒号分隔命名空间
5. 所有配置通过 Viper 读取，不硬编码
6. 日志使用结构化日志，包含 request_id`,
  },
  {
    id: 3, categoryId: 1, title: 'Cursor Rules 前端规范', slug: 'cursor-rules-frontend',
    description: '面向前端项目的 .cursorrules 配置文件生成模板，规范 AI 编码行为和代码风格。',
    tags: ['Cursor', '前端', '规范'], useCount: 680, likeCount: 412, favoriteCount: 203,
    status: 'active', createdAt: '2025-01-08', updatedAt: '2025-01-22',
    content: `# .cursorrules — 前端项目规范

## 项目信息
- 项目类型: 前端 Web 应用
- 框架: React + TypeScript
- 样式: Tailwind CSS

## 代码风格
- 使用函数式组件和 Hooks
- 使用 TypeScript 严格模式
- 组件文件使用 PascalCase 命名
- 工具函数使用 camelCase 命名
- 常量使用 UPPER_SNAKE_CASE

## 禁止事项
- 禁止使用 any 类型
- 禁止内联样式
- 禁止硬编码 API 地址
- 禁止在组件中直接调用 API
- 禁止使用 var 声明变量

## 组件规范
- 每个组件一个文件
- Props 必须定义 interface
- 使用 React.memo 优化渲染
- 事件处理使用 useCallback
- 副作用使用 useEffect 并清理`,
  },
  {
    id: 4, categoryId: 2, title: 'Claude Code 分层架构 Prompt', slug: 'claude-code-layered',
    description: '使用 Claude Code 生成分层架构代码的 Prompt 模板，适用于中大型后端项目。',
    tags: ['Claude Code', '架构', '后端'], useCount: 290, likeCount: 189, favoriteCount: 87,
    status: 'active', createdAt: '2025-01-12', updatedAt: '2025-01-19',
    content: `# Claude Code — 分层架构开发 Prompt

## 任务
开发一个基于分层架构的后端服务模块。

## 上下文
- 语言: Golang
- 框架: Gin
- 数据库: MySQL + Redis
- 认证: JWT

## 需求
1. 严格遵循 API → Service → Repository 三层架构
2. 每层只依赖下一层，不跨层调用
3. 使用接口定义层间契约，便于测试和替换
4. 统一错误处理，使用自定义错误码
5. 数据库操作使用事务保证一致性
6. 缓存使用 Redis，实现 Cache-Aside 模式

## 输出要求
- 先输出架构设计说明
- 再逐层输出代码
- 每个文件标注完整路径
- 包含单元测试示例`,
  },
  {
    id: 5, categoryId: 3, title: 'Docker Compose 部署 Prompt', slug: 'docker-compose-deploy',
    description: '生成 Docker Compose 多容器部署配置的 Prompt，包含 Web、API、数据库和 Nginx。',
    tags: ['Docker', 'DevOps', '部署'], useCount: 260, likeCount: 167, favoriteCount: 72,
    status: 'active', createdAt: '2025-01-05', updatedAt: '2025-01-16',
    content: `# Docker Compose 部署配置 Prompt

## 部署架构
- Nginx 反向代理 (80/443)
- Web 前端 (Node.js)
- API 后端 (Golang)
- MySQL 8 数据库
- Redis 缓存

## 要求
1. 所有服务使用 Alpine 基础镜像
2. 配置健康检查
3. 使用 named volumes 持久化数据
4. 环境变量通过 .env 文件管理
5. Nginx 配置 SSL 和 gzip
6. MySQL 配置字符集 utf8mb4
7. Redis 配置密码认证
8. 配置资源限制和重启策略`,
  },
  {
    id: 6, categoryId: 6, title: 'React Native 移动端 Prompt', slug: 'react-native-mobile',
    description: 'React Native 跨平台移动应用开发 Prompt，包含导航、状态管理和原生模块集成。',
    tags: ['React Native', '移动端', 'TypeScript'], useCount: 150, likeCount: 98, favoriteCount: 41,
    status: 'active', createdAt: '2025-01-14', updatedAt: '2025-01-21',
    content: `# React Native 移动端开发 Prompt

## 项目概述
使用 React Native + TypeScript 开发跨平台移动应用。

## 技术栈
- React Native 0.73
- TypeScript
- React Navigation 6
- Zustand 状态管理
- Axios 网络请求

## 编码规范
1. 使用函数式组件和 Hooks
2. 样式使用 StyleSheet.create
3. 导航使用类型安全的路由
4. 国际化支持 i18n
5. 无障碍 Accessibility 支持`,
  },
  {
    id: 7, categoryId: 2, title: 'Python FastAPI 后端 Prompt', slug: 'python-fastapi',
    description: '基于 Python FastAPI 框架的后端开发 Prompt，包含异步处理、Pydantic 验证和自动文档。',
    tags: ['Python', 'FastAPI', '异步'], useCount: 340, likeCount: 203, favoriteCount: 91,
    status: 'active', createdAt: '2025-01-11', updatedAt: '2025-01-23',
    content: `# Python FastAPI 后端开发 Prompt

## 角色定义
你是一个精通 Python 的后端工程师，擅长使用 FastAPI 构建高性能异步 API。

## 技术栈
- Python 3.11+
- FastAPI
- SQLAlchemy 2.0 (async)
- Pydantic v2
- Redis (aioredis)
- Celery 异步任务

## 编码规范
1. 使用 async/await 异步编程
2. Pydantic Model 做请求/响应验证
3. 依赖注入管理数据库连接
4. 统一异常处理中间件
5. API 版本化 /api/v1/
6. 自动生成 OpenAPI 文档`,
  },
  {
    id: 8, categoryId: 4, title: 'PostgreSQL 性能优化 Prompt', slug: 'postgresql-optimization',
    description: 'PostgreSQL 数据库性能优化的 Prompt 模板，涵盖索引策略、查询优化和配置调优。',
    tags: ['PostgreSQL', '性能', '优化'], useCount: 210, likeCount: 145, favoriteCount: 63,
    status: 'active', createdAt: '2025-01-09', updatedAt: '2025-01-17',
    content: `# PostgreSQL 性能优化 Prompt

## 任务
对 PostgreSQL 数据库进行全面的性能优化。

## 优化维度
1. 索引策略 — B-tree、GIN、BRIN 索引选择
2. 查询优化 — EXPLAIN ANALYZE 分析
3. 连接池 — PgBouncer 配置
4. 分区表 — 时间/范围分区
5. 配置调优 — shared_buffers, work_mem
6. 慢查询监控 — pg_stat_statements`,
  },
  {
    id: 9, categoryId: 5, title: 'LLM 应用开发 Prompt', slug: 'llm-app-dev',
    description: '大语言模型应用开发 Prompt，包含 RAG、Agent、Function Calling 和 Prompt Engineering。',
    tags: ['LLM', 'RAG', 'Agent'], useCount: 430, likeCount: 276, favoriteCount: 134,
    status: 'active', createdAt: '2025-01-13', updatedAt: '2025-01-24',
    content: `# LLM 应用开发 Prompt

## 项目概述
构建基于大语言模型的智能应用。

## 核心能力
1. RAG 检索增强生成
2. Agent 自主决策
3. Function Calling 工具调用
4. Prompt Chain 编排
5. 流式输出处理
6. Token 用量优化

## 技术选型
- LangChain / LlamaIndex
- Vector DB: Pinecone / Milvus
- Embedding: OpenAI / BGE
- Framework: FastAPI`,
  },
]

export const membershipPlans: MembershipPlan[] = [
  {
    id: 1, name: 'Free', code: 'free', price: 0, durationDays: 0, dailyLimit: 5,
    features: ['基础 Prompt 生成', '浏览模板库', '收藏模板', '生成历史', '一键复制'],
    status: 'active',
  },
  {
    id: 2, name: 'Pro', code: 'pro', price: 9.9, durationDays: 30, dailyLimit: 100,
    features: ['全部 Free 功能', '高级 Prompt 优化', '多 AI 工具支持', '优先生成速度', 'Markdown 渲染', '无广告体验'],
    status: 'active',
  },
  {
    id: 3, name: 'Team', code: 'team', price: 29.9, durationDays: 30, dailyLimit: 500,
    features: ['全部 Pro 功能', '团队协作空间', '模板共享', 'API 调用', '优先技术支持'],
    status: 'active',
  },
  {
    id: 4, name: 'Enterprise', code: 'enterprise', price: 0, durationDays: 0, dailyLimit: -1,
    features: ['全部 Team 功能', '私有部署', '自定义模型', 'SLA 保障', '专属客户经理'],
    status: 'active',
  },
]

export const techStacks = ['Next.js', 'Vue3', 'React', 'Golang', 'Python', 'MySQL', 'Redis', 'MongoDB', 'PostgreSQL', 'Docker', 'TypeScript', 'Rust']

export const aiTools = [
  { name: 'Cursor', color: '#06b6d4' },
  { name: 'Claude Code', color: '#8b5cf6' },
  { name: 'GPT', color: '#10b981' },
  { name: 'Gemini', color: '#f59e0b' },
  { name: 'DeepSeek', color: '#3b82f6' },
  { name: 'Qwen', color: '#f43f5e' },
]
