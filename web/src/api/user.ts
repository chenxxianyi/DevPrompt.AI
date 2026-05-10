import request from './request'

export interface GenerateStats {
  dailyUsed: number
  dailyLimit: number
  dailyRemaining: number
  membershipLevel: string
  membershipExpiredAt: string | null
}

export function getGenerateStats() {
  return request.get<{ data: GenerateStats }>('/user/generate-stats')
}
