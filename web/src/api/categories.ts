import request from './request'
import type { ApiResponse, PromptCategory } from '@/types'

export function getCategories(params?: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<{ list: PromptCategory[]; total: number; page: number; pageSize: number }>>('/categories', { params })
}
