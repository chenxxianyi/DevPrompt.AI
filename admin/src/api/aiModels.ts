import request from './request'
import type { ApiResponse, PaginatedData, AIModel } from '@/types'

export function getAiModels() {
  return request.get<ApiResponse<PaginatedData<AIModel>>>('/admin/ai-models')
}

export function createAiModel(data: Partial<AIModel>) {
  return request.post<ApiResponse<AIModel>>('/admin/ai-models', data)
}

export function updateAiModel(id: number, data: Partial<AIModel>) {
  return request.put<ApiResponse<AIModel>>(`/admin/ai-models/${id}`, data)
}

export function deleteAiModel(id: number) {
  return request.delete<ApiResponse<null>>(`/admin/ai-models/${id}`)
}
