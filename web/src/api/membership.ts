import request from './request'
import type { ApiResponse, MembershipPlan } from '@/types'

export function getMembershipPlans() {
  return request.get<ApiResponse<MembershipPlan[]>>('/membership/plans')
}
