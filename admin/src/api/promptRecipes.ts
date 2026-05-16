import request from './request'
import type { ApiResponse, PaginatedData, PromptRecipe } from '@/types'

export function getPromptRecipes(params?: { type?: string; targetTool?: string; status?: string }) {
  return request.get<ApiResponse<PaginatedData<PromptRecipe>>>('/admin/prompt-recipes', { params })
}

export function createPromptRecipe(data: Partial<PromptRecipe>) {
  return request.post<ApiResponse<PromptRecipe>>('/admin/prompt-recipes', data)
}

export function updatePromptRecipe(id: number, data: Partial<PromptRecipe>) {
  return request.put<ApiResponse<PromptRecipe>>(`/admin/prompt-recipes/${id}`, data)
}

export function deletePromptRecipe(id: number) {
  return request.delete<ApiResponse<null>>(`/admin/prompt-recipes/${id}`)
}

export function setDefaultPromptRecipe(id: number) {
  return request.put<ApiResponse<null>>(`/admin/prompt-recipes/${id}/set-default`)
}
