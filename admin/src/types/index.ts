// 统一 API 响应格式
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页数据结构
export interface PaginatedData<T = any> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// 用户
export interface User {
  id: number
  createdAt: string
  updatedAt: string
  username: string
  email: string
  avatar: string
  role: 'user' | 'admin'
  membershipLevel: 'free' | 'pro' | 'team' | 'enterprise'
  membershipExpiredAt: string | null
  dailyGenerateCount: number
  lastGenerateDate: string | null
  status: 'active' | 'disabled'
}

// Prompt 分类
export interface PromptCategory {
  id: number
  createdAt: string
  updatedAt: string
  name: string
  slug: string
  description: string
  sort: number
  status: 'active' | 'disabled'
}

// Prompt 模板
export interface PromptTemplate {
  id: number
  createdAt: string
  updatedAt: string
  categoryId: number
  title: string
  slug: string
  description: string
  content: string
  tags: string[]
  useCount: number
  likeCount: number
  favoriteCount: number
  status: 'active' | 'disabled'
  category?: PromptCategory
}

// AI 模型
export interface AIModel {
  id: number
  createdAt: string
  updatedAt: string
  provider: string
  modelName: string
  displayName: string
  isDefault: boolean
  status: 'active' | 'disabled'
  priority: number
  timeoutSeconds: number
}

// AI 调用日志
export interface AICallLog {
  id: number
  userId: number
  provider: string
  model: string
  requestType: string
  promptTokens: number
  completionTokens: number
  totalTokens: number
  status: 'success' | 'failed'
  errorMessage: string
  latencyMs: number
  createdAt: string
}

// 会员套餐
export interface MembershipPlan {
  id: number
  createdAt: string
  updatedAt: string
  name: string
  code: 'free' | 'pro' | 'team' | 'enterprise'
  price: number
  durationDays: number
  dailyLimit: number
  features: string[]
  status: 'active' | 'disabled'
}

// 订单
export interface Order {
  id: number
  createdAt: string
  updatedAt: string
  userId: number
  planId: number
  orderNo: string
  amount: number
  status: 'pending' | 'paid' | 'cancelled' | 'refunded'
  paidAt: string | null
}

// 仪表盘统计
export interface DashboardStats {
  userCount: number
  promptCount: number
  callLogCount: number
  planCount: number
  projectTypeCount: number
}

// 项目类型
export interface ProjectType {
  id: number
  createdAt: string
  updatedAt: string
  name: string
  value: string
  description: string
  sort: number
  status: 'active' | 'disabled'
}

// 登录
export interface LoginInput {
  email: string
  password: string
}

export interface LoginResult {
  token: string
  user: User
}

export interface TrialRequest {
  id: number
  createdAt: string
  updatedAt: string
  userId: number
  planCode: 'pro' | 'team' | 'enterprise'
  contact: string
  company: string
  teamSize: string
  useCase: string
  message: string
  status: 'pending' | 'contacted' | 'approved' | 'rejected'
  adminNote: string
}
