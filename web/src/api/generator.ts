import request from './request'
import type { ProjectGeneratorParams, CursorRulesParams, ClaudeCodeParams, OptimizeParams, GeneratedPrompt, PaginatedData } from '@/types'

export function generateProject(data: ProjectGeneratorParams) {
  return request.post<{ data: GeneratedPrompt }>('/generator/project', data)
}

export function generateCursorRules(data: CursorRulesParams) {
  return request.post<{ data: GeneratedPrompt }>('/generator/cursor-rules', data)
}

export function generateClaudeCode(data: ClaudeCodeParams) {
  return request.post<{ data: GeneratedPrompt }>('/generator/claude-code', data)
}

export function optimizePrompt(data: OptimizeParams) {
  return request.post<{ data: GeneratedPrompt }>('/generator/optimize', data)
}

export function getHistory(params?: { page?: number; pageSize?: number }) {
  return request.get<{ data: PaginatedData<GeneratedPrompt> }>('/generator/history', { params })
}
