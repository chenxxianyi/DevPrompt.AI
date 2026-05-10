import request from './request'
import type { ApiResponse, PaginatedData, Order } from '@/types'

export function getOrders(params: { page?: number; pageSize?: number }) {
  return request.get<ApiResponse<PaginatedData<Order>>>('/admin/orders', { params })
}
