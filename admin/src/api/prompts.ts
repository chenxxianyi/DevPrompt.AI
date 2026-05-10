import request from './request'
import type { ApiResponse, PaginatedData, PromptTemplate } from '@/types'

export function getPrompts(params: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<PaginatedData<PromptTemplate>>>('/admin/prompts', { params })
}

export function createPrompt(data: Partial<PromptTemplate>) {
  return request.post<ApiResponse<PromptTemplate>>('/admin/prompts', data)
}

export function updatePrompt(id: number, data: Partial<PromptTemplate>) {
  return request.put<ApiResponse<PromptTemplate>>(`/admin/prompts/${id}`, data)
}

export function deletePrompt(id: number) {
  return request.delete<ApiResponse<null>>(`/admin/prompts/${id}`)
}
