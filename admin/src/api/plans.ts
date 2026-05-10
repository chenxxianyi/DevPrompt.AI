import request from './request'
import type { ApiResponse, PaginatedData, MembershipPlan } from '@/types'

export function getPlans() {
  return request.get<ApiResponse<PaginatedData<MembershipPlan>>>('/admin/membership-plans')
}

export function createPlan(data: Partial<MembershipPlan>) {
  return request.post<ApiResponse<MembershipPlan>>('/admin/membership-plans', data)
}

export function updatePlan(id: number, data: Partial<MembershipPlan>) {
  return request.put<ApiResponse<MembershipPlan>>(`/admin/membership-plans/${id}`, data)
}

export function deletePlan(id: number) {
  return request.delete<ApiResponse<null>>(`/admin/membership-plans/${id}`)
}
