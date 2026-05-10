import request from './request'
import type { ApiResponse, DashboardStats } from '@/types'

export function getDashboardStats() {
  return request.get<ApiResponse<DashboardStats>>('/admin/dashboard')
}
