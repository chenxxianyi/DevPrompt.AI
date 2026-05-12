import request from './request'
import type { ApiResponse, PaginatedData, TrialRequest } from '@/types'

export function getTrialRequests(params?: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<PaginatedData<TrialRequest>>>('/admin/trial-requests', { params })
}

export function updateTrialRequestStatus(id: number, data: { status: TrialRequest['status']; adminNote?: string }) {
  return request.put<ApiResponse<null>>(`/admin/trial-requests/${id}`, data)
}
