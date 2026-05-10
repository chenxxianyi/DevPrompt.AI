import request from './request'
import type { ApiResponse, PaginatedData, PromptCategory } from '@/types'

export function getCategories() {
  return request.get<ApiResponse<PaginatedData<PromptCategory>>>('/admin/categories')
}

export function createCategory(data: Partial<PromptCategory>) {
  return request.post<ApiResponse<PromptCategory>>('/admin/categories', data)
}

export function updateCategory(id: number, data: Partial<PromptCategory>) {
  return request.put<ApiResponse<PromptCategory>>(`/admin/categories/${id}`, data)
}

export function deleteCategory(id: number) {
  return request.delete<ApiResponse<null>>(`/admin/categories/${id}`)
}
