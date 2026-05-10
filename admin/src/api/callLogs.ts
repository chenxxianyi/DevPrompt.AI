import request from './request'
import type { ApiResponse, PaginatedData, AICallLog } from '@/types'

export function getCallLogs(params: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<PaginatedData<AICallLog>>>('/admin/ai-call-logs', { params })
}
