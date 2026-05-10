import request from './request'
import type { ApiResponse, PaginatedData, User } from '@/types'

export function getUsers(params: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<PaginatedData<User>>>('/admin/users', { params })
}
