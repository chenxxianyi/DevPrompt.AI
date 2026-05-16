export interface User {
  id: number
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

export interface PromptCategory {
  id: number
  name: string
  slug: string
  description: string
  sort: number
  status: 'active' | 'disabled'
  count?: number
}

export interface PromptTemplate {
  id: number
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
  isLiked?: boolean
  isFavorited?: boolean
  category?: PromptCategory
  createdAt: string
  updatedAt: string
}

export interface GeneratedPrompt {
  id: number
  userId: number
  type: 'project' | 'cursor-rules' | 'claude-code' | 'optimize'
  title: string
  input: string | Record<string, any>
  output: string
  model: string
  provider: string
  tokens: number
  createdAt: string
}

export interface MembershipPlan {
  id: number
  name: string
  code: 'free' | 'pro' | 'team' | 'enterprise'
  price: number
  durationDays: number
  dailyLimit: number
  features: string[]
  status: 'active' | 'disabled'
}

export interface TrialRequestPayload {
  planCode: string
  contact?: string
  company?: string
  teamSize?: string
  useCase?: string
  message?: string
}

export interface LoginParams {
  email: string
  password: string
}

export interface RegisterParams {
  username: string
  email: string
  password: string
}

export type QualityMode = 'concise' | 'standard' | 'expert'
export type OutputFormat = 'markdown' | 'checklist' | 'json' | 'plain'

export interface QualityOptions {
  qualityMode?: QualityMode
  outputFormat?: OutputFormat
  includeAcceptanceCriteria?: boolean
  includeRiskCheck?: boolean
  includeTestPlan?: boolean
  includeDeploymentNotes?: boolean
}

export interface ProjectGeneratorParams {
  projectName: string
  projectType: string
  techStack: string[]
  features: string[]
  targetAiTool: string
  qualityMode?: QualityMode
  outputFormat?: OutputFormat
  includeAcceptanceCriteria?: boolean
  includeRiskCheck?: boolean
  includeTestPlan?: boolean
  includeDeploymentNotes?: boolean
}

export interface CursorRulesParams {
  language: string
  framework: string
  codeStyle: string
  rules: string[]
  qualityMode?: QualityMode
  outputFormat?: OutputFormat
  includeAcceptanceCriteria?: boolean
  includeRiskCheck?: boolean
  includeTestPlan?: boolean
  includeDeploymentNotes?: boolean
}

export interface ClaudeCodeParams {
  task: string
  context: string
  requirements: string[]
  qualityMode?: QualityMode
  outputFormat?: OutputFormat
  includeAcceptanceCriteria?: boolean
  includeRiskCheck?: boolean
  includeTestPlan?: boolean
  includeDeploymentNotes?: boolean
}

export interface OptimizeParams {
  rawPrompt: string
  targetTool: string
  optimizeLevel: 'basic' | 'professional' | 'expert'
  qualityMode?: QualityMode
  outputFormat?: OutputFormat
  includeAcceptanceCriteria?: boolean
  includeRiskCheck?: boolean
  includeTestPlan?: boolean
  includeDeploymentNotes?: boolean
}

export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

export interface ProjectType {
  id: number
  name: string
  value: string
  description: string
  sort: number
  status: 'active' | 'disabled'
}
